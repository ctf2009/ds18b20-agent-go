package main

import (
	"github.com/ctf2009/ds18b20-agent-go/internal/features/ds18b20"
	"github.com/ctf2009/ds18b20-agent-go/internal/features/web"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"log"
	"net/http"
	"strings"
)

func main() {
	//configuration, err := config.New()
	//if err != nil {
	//	log.Panicln("Configuration error", err)
	//}

	ds18b20.Init()
	web.Init()

	r := Router()
	log.Fatal(http.ListenAndServe(":8080", r))
}

func Router() *chi.Mux {
	router := chi.NewRouter()

	router.Use(
		render.SetContentType(render.ContentTypeJSON), // Set content-Type headers as application/json
		middleware.Logger,                             // Log API request calls
		middleware.DefaultCompress,                    // Compress results, mostly gzipping assets and json
		middleware.RedirectSlashes,                    // Redirect slashes to no slash URL versions
		middleware.Recoverer,                          // Recover from panics without crashing server
	)

	router.Mount("/api/ds18b20", ds18b20.Routes())

	// Serving Static Files
	fs := http.StripPrefix("/", http.FileServer(http.Dir("./public")))
	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	})

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.Replace(route, "/*/", "/", -1)
		log.Printf("Endpoint - %s %s\n", method, route)
		return nil
	}

	if err := chi.Walk(router, walkFunc); err != nil {
		log.Panicf("Logging err: %s\n", err.Error())
	}

	return router
}
