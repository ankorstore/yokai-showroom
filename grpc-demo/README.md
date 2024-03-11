# Yokai gRPC Demo

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go version](https://img.shields.io/badge/Go-1.20-blue)](https://go.dev/)

> gRPC API demo application, based on
> the [Yokai](https://github.com/ankorstore/yokai) Go framework.

<!-- TOC -->
* [Overview](#overview)
* [Usage](#usage)
  * [Start the application](#start-the-application)
  * [Available services](#available-services)
  * [Authentication](#authentication)
  * [Examples](#examples)
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

- `localhost:50051`: application gRPC server
- [http://localhost:8081](http://localhost:8081): application core dashboard

### Available services

This demo application provides a [TransformTextService](proto/transform.proto), with the following `RPCs`:

| RPC                     | Type      | Description                                                  |
|-------------------------|-----------|--------------------------------------------------------------|
| `TransformText`         | unary     | Transforms a given text using a given transformer            |
| `TransformAndSplitText` | streaming | Transforms and splits a given text using a given transformer |

This demo application also provides [reflection](https://ankorstore.github.io/yokai/modules/fxgrpcserver/#reflection) and [health check ](https://ankorstore.github.io/yokai/modules/fxgrpcserver/#health-check) services.

### Authentication

This demo application provides example [authentication interceptors](internal/interceptor/authentication.go).

You can enable authentication in the application [configuration file](configs/config.yaml) with `config.authentication.enabled=true`.

If enabled, you need to provide the secret configured in `config.authentication.secret` as request `authorization` metadata.

### Examples

Usage examples with [grpcurl](https://github.com/fullstorydev/grpcurl):

- with `TransformTextService/TransformText`:

```shell
grpcurl -plaintext -d '{"text":"abc","transformer":"TRANSFORMER_UPPERCASE"}' localhost:50051 transform.TransformTextService/TransformText
{
  "text": "ABC"
}
```

- with `TransformTextService/TransformAndSplitText`:

```shell
grpcurl -plaintext -d '{"text":"ABC DEF","transformer":"TRANSFORMER_LOWERCASE"}' localhost:50051 transform.TransformTextService/TransformAndSplitText
{
  "text": "abc"
}
{
  "text": "def"
}
```

You can use any gRPC clients, such as [Postman](https://learning.postman.com/docs/sending-requests/grpc/grpc-request-interface/) or [Evans](https://github.com/ktr0731/evans).

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
make up        # start the docker compose stack
make down      # stop the docker compose stack
make logs      # stream the docker compose stack logs
make fresh     # refresh the docker compose stack
make protogen  # generate gRPC stubs with protoc
make test      # run tests
make lint      # run linter
```
