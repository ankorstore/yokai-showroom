# Yokai gRPC Demo

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go version](https://img.shields.io/badge/Go-1.20-blue)](https://go.dev/)

> gRPC REST API demo application, based on
> the [Yokai](https://github.com/ankorstore/yokai) Go framework.

<!-- TOC -->
* [Overview](#overview)
* [Usage](#usage)
  * [Start the application](#start-the-application)
  * [Available endpoints](#available-endpoints)
* [Contents](#contents)
  * [Layout](#layout)
  * [Makefile](#makefile)
<!-- TOC -->

## Overview

This demo application is a simple gRPC API offering a [text transformation service](proto/transform.proto).

It provides a demo [Yokai](https://github.com/ankorstore/yokai) application container, with the [fxgrpcserver](https://github.com/ankorstore/yokai/tree/main/fxgrpcserver) module to offer the gRPC API

See the [Yokai documentation](https://ankorstore.github.io/yokai) for more details.

## Usage

### Start the application

To start the application, simply run:

```shell
make fresh
```

After a short moment, the application will offer:

- localhost:50051: gRPC service
- [http://localhost:8081](http://localhost:8081): application core dashboard

### Available service

This demo application provides a [TransformTextService](proto/transform.proto), with the following `RPCs`:

| RPC                     | Type      | Description                                                  |
|-------------------------|-----------|--------------------------------------------------------------|
| `TransformText`         | unary     | Transforms a given text using a given transformer            |
| `TransformAndSplitText` | streaming | Transforms and splits a given text using a given transformer |

## Contents

### Layout

This demo application is following the [standard go project layout](https://github.com/golang-standards/project-layout):

- `cmd/`: entry points
- `configs/`: configuration files
- `internal/`:
  - `interceptor/`: gRPC interceptors
  - `service/`: gRPC services
  - `bootstrap.go`: bootstrap (modules, lifecycles, etc)
  - `services.go`: dependency injection

### Makefile

This demo application provides a `Makefile`:

```
make up     # start the docker compose stack
make down   # stop the docker compose stack
make logs   # stream the docker compose stack logs
make fresh  # refresh the docker compose stack
make proto  # generate gRPC stubs with protoc
make test   # run tests
make lint   # run linter
```
