$(document).ready(function () {
    $(".dropdown-trigger").dropdown();

    let ds18b20Manager = new Ds18b20Manager()

    $("main").append(ds18b20Manager.render())
});

class Ds18b20Manager {

    constructor() {
        this.probesMap = {}
        this.container = this.generateContainer()

        // Initial Loading Graphic
        this.container.append(
            $("<div>", {id: "probe-loader", class: "center-align"}).append(generateLoader()))

        this.container.append(
            $("<div>", {id: "probe-summary", class: "row"})
        )
    }

    generateContainer() {
        let container = $("<div>", {id: "probe-container", class: "container"})
        //container.append("<h4 class='blue-text text-accent-3 center'>Probes<h4>")
        return container
    }

    render() {
        console.log("Rendering")
        this.update()
        return this.container
    }

    update() {
        $.ajax({
            url: "/api/ds18b20/",
            context: this,
            success: function (data) {

                // Remove Loader
                $("#probe-loader").remove();

                data.forEach(probe => {
                    console.log("Probe:" + probe)
                    if (this.probesMap[probe.id]) {
                        console.log("Update")
                    } else {
                        // We havent seen this Probe Before
                        let ds18b20 = new Ds18b20(probe)
                        $("#probe-summary").append(ds18b20.render())
                    }
                })
            }
        })
    }
}

class Ds18b20 {
    constructor(data) {
        this.id = data.id
        this.label = data.label ? data.label : data.id
        this.temp = data.temperature
        this.unit = "C"
    }

    update(data) {
        this.temp = data.temperature
    }

    render() {

        let current_datetime = new Date()
        let formatted_date = current_datetime.getFullYear()
            + "-" + (current_datetime.getMonth() + 1)
            + "-" + current_datetime.getDate() + " "
            + current_datetime.getHours() + ":"
            + current_datetime.getMinutes() + ":"
            + current_datetime.getSeconds()

        return $(`
        <div class="col s12 m7 offset-m2 probe-card">
            <div class="card indigo darken-3">
                <div class="card-content white-text">
                    <span class="card-title activator blue-text text-accent-3">${this.label}<i class="material-icons right">more_vert</i></span>
                    <p><span class="supplimentary blue-text text-accent-3">Id: </span>${this.id}</p>
                    <div class="temperature">
                        ${this.temp}\xB0${this.unit}
                    </div>
                    <div class="last-updated">
                        <span class="blue-text text-accent-3">Last Updated: </span>${formatted_date}
                    </div>
                </div>
                <div class="card-reveal light-blue lighten-5">
                    <span class="card-title indigo-text text-darken-3">Update Label<i class="material-icons right">close</i></span>
                    <form action="/api/ds18b20/${this.id}" method="post">
                        <div class="row">
                            <div class="col s8">
                                <input type="hidden" name="probeId" value="${this.id}">
                                <input placeholder="${this.label ? this.label : this.id}" name="label" type="text" class="validate"> 
                                <label for="${this.id}-update">Probe Name</label>
                            </div>
                            <div class="col s4 center">
                                <button class="btn waves-effect waves-light blue update" type="submit">Update</button>
                            </div>
                        </div>
                    </form>                  
                </div>
            </div>
        </div>`)
    }
}

// TODO: Move to Utils
function generateLoader() {
    return $(`
    <div class="preloader-wrapper big active">
    <div class="spinner-layer spinner-blue-only">
      <div class="circle-clipper left">
        <div class="circle"></div>
      </div><div class="gap-patch">
        <div class="circle"></div>
      </div><div class="circle-clipper right">
        <div class="circle"></div>
      </div>
    </div>
  </div>`)
}