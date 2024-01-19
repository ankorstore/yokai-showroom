# Yokai Worker Demo

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go version](https://img.shields.io/badge/Go-1.20-blue)](https://go.dev/)

> Demo application working with [Pub/Sub](https://cloud.google.com/pubsub), based on
> the [Yokai](https://github.com/ankorstore/yokai) Go framework.

<!-- TOC -->

* [Overview](#overview)
* [Usage](#usage)
	* [Start the application](#start-the-application)
	* [Message publication](#message-publication)
	* [Message subscription](#message-subscription)
* [Contents](#contents)
	* [Layout](#layout)
	* [Makefile](#makefile)

<!-- TOC -->

## Overview

This demo provides:

- a demo [Yokai](https://github.com/ankorstore/yokai) application, with:
	- the [fxhttpserver](https://github.com/ankorstore/yokai/tree/main/fxhttpserver) module to offer a Pub/Sub
	  publication endpoint
	- the [fxworker](https://github.com/ankorstore/yokai/tree/main/fxworker) module to offer a worker running a Pub/Sub
	  subscriber
- a [Pub/Sub](https://cloud.google.com/pubsub) emulator

See the [Yokai documentation](https://ankorstore.github.io/yokai) for more details.

## Usage

### Start the application

To start the application, simply run:

```shell
make fresh
```

After a short moment, the application will offer:

- [http://localhost:8080](http://localhost:8080): application http server
- [http://localhost:8081](http://localhost:8081): application core dashboard

### Message publication

You can use the [[GET] /publish](http://localhost:8080/publish) endpoint to publish a message to the Pub/Sub emulator.

By default the current datetime will be sent, but you can use the `message` query parameter to send a custom message
with, for example, [[GET] /publish?message=test](http://localhost:8080/publish?message=test).

### Message subscription

Check your application logs by running:

```shell
make logs
```

You will see the [SubscribeWorker](internal/worker/subscribe.go) subscribed to Pub/Sub in action, logging the received
messages.

## Contents

### Layout

This template is following the [standard go project layout](https://github.com/golang-standards/project-layout):

- `cmd/`: entry points
- `configs/`: configuration files
- `internal/`:
	- `handler/`: HTTP handlers
	- `service/`: services
	- `worker/`: workers
	- `bootstrap.go`: bootstrap (modules, lifecycles, etc)
	- `routing.go`: routing
	- `services.go`: dependency injection
- `modules/`: application modules
	- `fxpubsub/`: Pub/Sub module

### Makefile

This template provides a `Makefile`:

```
make up     # start the docker compose stack
make down   # stop the docker compose stack
make logs   # stream the docker compose stack logs
make fresh  # refresh the docker compose stack
make test   # run tests
make lint   # run linter
```
