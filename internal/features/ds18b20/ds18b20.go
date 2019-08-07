package ds18b20

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ctf2009/ds18b20-agent-go/internal/agent"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
)

const DEVICES_FILE = "w1_master_slaves"
const DEVICE_FILE = "w1_slave"

const PROBE_MAP_FILE = "probe-map.json"

var (
	probeMap = make(map[string]string)
	probes   = make(map[string]bool)

	ds18b20Root string
	devicesFile string
	storeDir    string
)

type ds18b20 struct {
	Id          string  `json:"id"`
	Label       string  `json:"label"`
	Temperature float64 `json:"temperature"`
	Error       string  `json:"error"`
}

func Init(config *agent.Config) {
	log.Println("Initialising DS18b20 Probes")

	ds18b20Root = config.DS18B20_ROOT
	devicesFile = path.Join(ds18b20Root, DEVICES_FILE)
	storeDir = config.STORE_DIR

	w1MasterSlaveData, err := ioutil.ReadFile(devicesFile)
	if err != nil {
		log.Panicf("Unable to locate devices file %s", devicesFile)
	}

	sensors := strings.Split(string(bytes.TrimSpace(w1MasterSlaveData)), "\n")

	if len(sensors) != 0 {
		log.Printf("Found %d probes\n", len(sensors))
		for _, sensor := range sensors {
			probes[sensor] = true
			log.Printf("Setting up Probe with name: %s\n", sensor)
		}
	} else {
		log.Println("No Probes Found")
	}

	if _, err := os.Stat(storeDir); os.IsNotExist(err) {
		os.Mkdir(storeDir, os.ModePerm)
	} else {
		loadProbeMapFile()
	}
}

func loadProbeMapFile() {
	probeMapFilePath := path.Join(storeDir, PROBE_MAP_FILE)

	if _, err := os.Stat(probeMapFilePath); os.IsNotExist(err) {
		log.Printf("%s file not found\n", PROBE_MAP_FILE)
	} else {
		log.Printf("Loading %s\n", PROBE_MAP_FILE)
		if jsonFile, err := os.Open(probeMapFilePath); err != nil {
			log.Printf("Error reading %s: %s", PROBE_MAP_FILE, err)
		} else {
			byteValue, _ := ioutil.ReadAll(jsonFile)
			if err := json.Unmarshal([]byte(byteValue), &probeMap); err != nil {
				log.Printf("Unable to parse %s: %s", PROBE_MAP_FILE, err)
			}
		}
	}
}

func updateProbeMapFile() {
	probeMapFilePath := path.Join(storeDir, PROBE_MAP_FILE)

	if _, err := os.Stat(storeDir); os.IsNotExist(err) {
		os.Mkdir(storeDir, os.ModePerm)
	}
	jsonString, _ := json.Marshal(probeMap)
	if err := ioutil.WriteFile(probeMapFilePath, jsonString, 0644); err != nil {
		fmt.Printf("Failed to write updated %s due to error %s", probeMapFilePath, err)
	}
}

func getAllDs18b20(w http.ResponseWriter, r *http.Request) {
	results := make([] *ds18b20, 0)

	if len(probes) != 0 {
		var wg sync.WaitGroup
		ds18b20Chan := make(chan *ds18b20, len(probes))

		for probeId := range probes {
			wg.Add(1)
			go func(probeId string) {
				ds18b20Chan <- getDs18b20(probeId)
				wg.Done()
			}(probeId)
		}

		wg.Wait()
		close(ds18b20Chan)

		for ds18b20 := range ds18b20Chan {
			results = append(results, ds18b20)
		}
	}
	render.JSON(w, r, results)
}

func getDs18b20(probeId string) *ds18b20 {
	probe := &ds18b20{
		Id:    probeId,
		Label: probeMap[probeId],
	}

	data, err := ioutil.ReadFile(path.Join(ds18b20Root, probeId, DEVICE_FILE))
	if err != nil {
		fmt.Println(err)
		probe.Temperature = -1
	}

	raw := string(bytes.TrimSpace(data))
	temperatureIndex := strings.LastIndex(raw, "t=")

	if temperatureIndex == -1 {
		probe.Temperature = -1
	} else {
		if c, err := strconv.ParseFloat(raw[temperatureIndex+2:], 64); err != nil {
			probe.Error = err.Error()
			probe.Temperature = -1
		} else {
			probe.Temperature = math.Round((c/1000)*100) / 100
		}
	}

	return probe
}

// Currently this parses a http form and processes the update
// This should be done front end in future
func updateDs18b20ById(w http.ResponseWriter, r *http.Request) {

	contentType := r.Header.Get("Content-Type")

	if len(contentType) == 0 {
		http.Error(w, "Missing Content-Type header", http.StatusBadRequest)
	}

	switch contentType {
	case "application/json":
		decoder := json.NewDecoder(r.Body)
		jsonMap := make(map[string]string)
		err := decoder.Decode(&jsonMap)

		if err != nil {
			fmt.Println(err)
			http.Error(w, "Unable to Parse JSON", http.StatusBadRequest)
		}

		probeId := chi.URLParam(r, "id")
		label := jsonMap["label"]

		if len(probeId) != 0 && len(label) != 0 {
			fmt.Println("Updating Label: " + label + " " + probeId)
			w.WriteHeader(http.StatusNoContent)
		} else {
			w.WriteHeader(http.StatusNotModified)
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err == nil {
			probeId, label := r.FormValue("probeId"), r.FormValue("label")
			if len(probeId) != 0 && len(label) != 0 {
				probeMap[probeId] = label
				updateProbeMapFile()
			}
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		http.Error(w, "Unsupported Content Type", http.StatusBadRequest)
	}
}

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/", getAllDs18b20)

	// Currently handle everything server side
	router.Post("/{id}", updateDs18b20ById)

	return router
}
