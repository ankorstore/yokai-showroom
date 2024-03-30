# Yokai gRPC Demo

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go version](https://img.shields.io/badge/Go-1.22-blue)](https://go.dev/)

> gRPC API demo application, based on
> the [Yokai](https://github.com/ankorstore/yokai) Go framework.

<!-- TOC -->
* [Overview](#overview)
  * [Layout](#layout)
  * [Makefile](#makefile)
* [Usage](#usage)
  * [Start the application](#start-the-application)
  * [Available services](#available-services)
  * [Authentication](#authentication)
  * [Examples](#examples)
<!-- TOC -->

## Overview

This demo application is a simple gRPC API offering a [text transformation service](proto/example.proto).

It provides:

- a [Yokai](https://github.com/ankorstore/yokai) application container, with the [gRPC server](https://ankorstore.github.io/yokai/modules/fxgrpcserver/) module to offer the gRPC API
- a [Jaeger](https://www.jaegertracing.io/) container to collect the application traces

### Layout

This demo application is following the [recommended project layout](https://go.dev/doc/modules/layout):

- `cmd/`: entry points
- `configs/`: configuration files
- `internal/`:
  - `interceptor/`: gRPC interceptors
  - `service/`: gRPC services
  - `bootstrap.go`: bootstrap (modules, lifecycles, etc)
  - `services.go`: dependency injection
- `proto/`: protobuf definition and stubs

### Makefile

This demo application provides a `Makefile`:

```
make up     # start the docker compose stack
make down   # stop the docker compose stack
make logs   # stream the docker compose stack logs
make fresh  # refresh the docker compose stack
make stubs  # generate gRPC stubs with protoc
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

- `localhost:50051`: application gRPC server
- [http://localhost:8081](http://localhost:8081): application core dashboard
- [http://localhost:16686](http://localhost:16686): jaeger UI

### Available services

This demo application provides a [TransformTextService](proto/example.proto), with the following `RPCs`:

| RPC                     | Type      | Description                                                  |
|-------------------------|-----------|--------------------------------------------------------------|
| `TransformText`         | unary     | Transforms a given text using a given transformer            |
| `TransformAndSplitText` | streaming | Transforms and splits a given text using a given transformer |

If no `Transformer` is provided, the transformation configured in `config.transform.default` will be applied.

If you update the [proto definition](proto/example.proto), you can run `make stubs` to regenerate the stubs.

This demo application also provides [reflection](https://ankorstore.github.io/yokai/modules/fxgrpcserver/#reflection) and [health check](https://ankorstore.github.io/yokai/modules/fxgrpcserver/#health-check) services.

### Authentication

This demo application provides example [authentication interceptors](internal/interceptor/authentication.go).

You can enable authentication in the application [configuration file](configs/config.yaml) with `config.authentication.enabled=true`.

If enabled, you need to provide the secret configured in `config.authentication.secret` as context `authorization` metadata.

### Examples

Usage examples with [grpcurl](https://github.com/fullstorydev/grpcurl):

- with `TransformTextService/TransformText`:

```shell
grpcurl -plaintext -d '{"text":"abc","transformer":"TRANSFORMER_UPPERCASE"}' localhost:50051 example.TransformTextService/TransformText
{
  "text": "ABC"
}
```

- with `TransformTextService/TransformAndSplitText`:

```shell
grpcurl -plaintext -d '{"text":"ABC DEF","transformer":"TRANSFORMER_LOWERCASE"}' localhost:50051 example.TransformTextService/TransformAndSplitText
{
  "text": "abc"
}
{
  "text": "def"
}
```

You can use any gRPC clients, for example [Postman](https://learning.postman.com/docs/sending-requests/grpc/grpc-request-interface/) or [Evans](https://github.com/ktr0731/evans).
