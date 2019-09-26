[![Build Status](https://travis-ci.com/gapsquare/goevent.svg?branch=master)](https://travis-ci.com/gapsquare/goevent)
[![Coverage Status](https://coveralls.io/repos/github/gapsquare/goevent/badge.svg?branch=master&service=github)](https://coveralls.io/github/gapsquare/goevent?branch=master&service=github)


# GoEvent

GoEvent is a simple implementation of eventing toolkit for Go

**NOTE: GoEvent is used in production systems but the API is not final!**


GoEvent defines a set of interfaces to work with events which can be published and handled over a Event Bus. 

It also define a **local** and a Google PubSub **gpubsub** event bus

Inspired by the following libraries/examples:

- https://github.com/looplab/eventhorizon


Suggestions are welcome!

# Usage

See the example folder for a few examples to get you started.


# Messaging drivers

These are the drivers for messaging, currently only publishers.

### Local / in memory

Fully synchrounos. Useful for testing/experimentation.

### GCP Cloud Pub/Sub

Experimental driver.

## Development

To develop GoEvent you need to have Docker and Docker Compose installed.

To start all needed services and run all tests, simply run make:

```bash
make
```

To manually run the services and stop them:

```bash
make services
make stop
```

When the services are running testing can be done either locally or with Docker:

```bash
make test
make test_docker
go test ./...
```

The difference between `make test` and `go test ./...` is that `make test` also prints coverage info.

# Get Involved

- Check out the [contribution guidelines](CONTRIBUTING.md)

# License
GoEvent licensed under Apache License 2.0

http://www.apache.org/licenses/LICENSE-2.0
