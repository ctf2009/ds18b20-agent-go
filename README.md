# ds18b20-agent-go

A golang agent that is able to query ds18b20 temperature probes.

The agent was built with the Raspberry Pi in mind, however it is entirely possible this could work on other platforms with minimal changes
The agent provides both a UI and REST interface in order to show the temperatures of the connected probes

### Assumptions

The app currently expects device information to exist in a single root folder.
For Raspberry Pi's the root folder is usually `/sys/devices/w1_bus_master1` but it it possible to override this via configuration

Regardless of the root folder used, the following are expected to exist

`w1_master_slaves` - This is a file that lists the ds18b20 probes

`28-XXXXXXXX` - One or more of these directories which represent each ds18b20 probe

Each ds18b20 probe directory should contain a `w1_slave` file which is used to obtain the temperatues
An example of the expected layout is show in the `test/data/probes` directory in this repository

### Running 

You can run the app straight away using `go run cmd/ds18b20-service/main.go 

By default, this will pick up the `agent-config.yaml` configuration file in the `test` directory in this repository.
This will point the application to the test probes in `test/data/probes` and allow you to access the UI at `localhost:8080`

In order to run the agent with real probes, you can create a new `agent-config.yaml` on the root of the repository.
This will allow you specify a `DS18B20_ROOT` location such as `sys/devices/w1_bus_master1`

If you choose to build using `go build cmd/ds18b20-service/main.go` and try to run the resulting binary, 
simply ensure your `agent-config.yaml` file is placed alongside the binary for it to be picked up

The `STORE_DIR` property specifies a location where the application stores its data, 
including a probe map file which allows the id of the probe to be mapped to an easy to read label which is updatable via the UI or REST interface


### Docker

The project includes two Dockerfiles:

`Dockerfile-amd64` - This allows you to run locally on a compatible machine. 
This is used by the `docker-compose.yml` in the `test` directory and is mainly used for testing

`Dockerfile-arm` - This allows the agent to be built for ALL Raspberry Pi's, 
including the zero which runs using arm32/v6 architecture

It is possible to build both of these and push to dockerhub with a multi architecture manifest as follows:

```
docker build -t <USERNAME>/ds18b20-agent-go:amd64-<VERSION> -f Dockerfile-amd64 .
docker build -t <USERNAME>/ds18b20-agent-go:arm-<VERSION> -f Dockerfile-arm .

docker push <USERNAME>/ds18b20-agent-go:amd64-<VERSION>
docker push <USERNAME>/ds18b20-agent-go:arm-<VERSION>

docker manifest create <USERNAME>/ds18b20-agent-go:<VERSION> <USERNAME>/ds18b20-agent-go:amd64-<VERSION> <USERNAME>/ds18b20-agent-go:arm-<VERSION>
docker manifest annotate <USERNAME>/ds18b20-agent-go:<VERSION> <USERNAME>/ds18b20-agent-go:arm-<VERSION> --arch arm
docker manifest push <USERNAME>/ds18b20-agent-go:<VERSION>
```

To pull down the image is as simple as: 

`docker pull <USERNAME>/ds18b20-agent-go:<VERSION>`

The correct image will be pulled down depending on the architecture required

#### Running the Image

An example of how to run the image is shown as follows:

```
docker container run -d -p 8080:8080 \
-v /sys/devices/w1_bus_master1:/app/data/probes:rw  \
-v /data/ds18b20/store:/app/data/store:rw <USERNAME>/ds18b20-agent-go:<VERSION>
```

### UI Screenshots

![UI Screenshot](screenshots/ui-screenshot.png?raw=true "UI Screenshot")