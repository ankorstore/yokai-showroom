# Yokai Worker Demo

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go version](https://img.shields.io/badge/Go-1.24-blue)](https://go.dev/)

> Demo application working with [Pub/Sub](https://cloud.google.com/pubsub), based on
> the [Yokai](https://github.com/ankorstore/yokai) Go framework.

<!-- TOC -->
* [Overview](#overview)
  * [Layout](#layout)
  * [Makefile](#makefile)
* [Usage](#usage)
  * [Start the application](#start-the-application)
  * [Message publication](#message-publication)
  * [Message subscription](#message-subscription)
<!-- TOC -->

## Overview

This demo application is a simple subscriber to [Pub/Sub](https://cloud.google.com/pubsub).

It provides:

- a [Yokai](https://github.com/ankorstore/yokai) application container, with the [worker](https://ankorstore.github.io/yokai/modules/fxworker/) module to offer a worker subscribing to Pub/Sub (using the [fxgcppubsub](https://github.com/ankorstore/yokai-contrib/tree/main/fxgcppubsub) contrib module)
- a [Pub/Sub emulator](https://cloud.google.com/pubsub) container, with preconfigured topic and subscription
- a [Pub/Sub emulator UI](https://github.com/echocode-io/gcp-pubsub-emulator-ui) container, preconfigured to work with the emulator container
- a [Jaeger](https://www.jaegertracing.io/) container to collect the application traces

### Layout

This demo application is following the [recommended project layout](https://go.dev/doc/modules/layout#server-project):

- `cmd/`: entry points
- `configs/`: configuration files
- `internal/`:
  - `worker/`: workers
  - `bootstrap.go`: bootstrap
  - `register.go`: dependencies registration

### Makefile

This demo application provides a `Makefile`:

```
make up     # start the docker compose stack
make down   # stop the docker compose stack
make logs   # stream the docker compose stack logs
make fresh  # refresh the docker compose stack
make test   # run tests
make lint   # run linter
```

## Usage

### Start the application

To start the application, simply run:

```shell
make fresh
```

After a short moment, the application will offer:

- [http://localhost:8081](http://localhost:8081): application core dashboard
- [http://localhost:8680](http://localhost:8680): pub/sub emulator UI
- [http://localhost:16686](http://localhost:16686): jaeger UI

### Message publication

You can use the Pub/Sub emulator UI to publish a message to the preconfigured topic:

[http://localhost:8680/project/demo-project/topic/demo-topic](http://localhost:8680/project/demo-project/topic/demo-topic)

### Message subscription

Check your application logs by running:

```shell
make logs
```

You will see the [SubscribeWorker](internal/worker/subscribe.go) subscribed to Pub/Sub in action, logging the received
messages.
