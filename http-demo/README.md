# Yokai HTTP Demo

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go version](https://img.shields.io/badge/Go-1.24-blue)](https://go.dev/)

> HTTP REST API demo application, based on
> the [Yokai](https://github.com/ankorstore/yokai) Go framework.

<!-- TOC -->
* [Overview](#overview)
  * [Layout](#layout)
  * [Makefile](#makefile)
* [Usage](#usage)
  * [Start the application](#start-the-application)
  * [Available endpoints](#available-endpoints)
  * [Authentication](#authentication)
<!-- TOC -->

## Overview

This demo application is a simple REST API (CRUD) to manage [gophers](https://go.dev/blog/gopher).

It provides:

- a [Yokai](https://github.com/ankorstore/yokai) application container, with the [HTTP server](https://ankorstore.github.io/yokai/modules/fxhttpserver/) and [SQL](https://ankorstore.github.io/yokai/modules/fxsql/) modules to offer the gophers API
- a [MySQL](https://www.mysql.com/) container to store the gophers
- a [Jaeger](https://www.jaegertracing.io/) container to collect the application traces

### Layout

This demo application is following the [recommended project layout](https://go.dev/doc/modules/layout#server-project):

- `cmd/`: entry points
- `configs/`: configuration files
- `db/`:
  - `migrations/`: database migrations
  - `seeds/`: database seeds
- `internal/`:
  - `api/`: HTTP API
    - `handler/`: HTTP handlers
    - `middleware/`: HTTP middlewares
  - `domain/`: domain
    - `model.go`: gophers model 
    - `repository.go`: gophers repository 
    - `service.go`: gophers service 
  - `bootstrap.go`: bootstrap
  - `register.go`: dependencies registration
  - `router.go`: routing registration
- `templates/`: HTML templates

### Makefile

This demo application provides a `Makefile`:

```
make up      # start the docker compose stack
make down    # stop the docker compose stack
make logs    # stream the docker compose stack logs
make fresh   # refresh the docker compose stack
make migrate # run database migrations
make test    # run tests
make lint    # run linter
```

## Usage

### Start the application

To start the application, simply run:

```shell
make fresh
```

After a short moment, the application will offer:

- [http://localhost:8080](http://localhost:8080): application dashboard
- [http://localhost:8081](http://localhost:8081): application core dashboard
- [http://localhost:16686](http://localhost:16686): jaeger UI

### Available endpoints

On [http://localhost:8080](http://localhost:8080), you can use:

| Route                   | Description      | Type     |
|-------------------------|------------------|----------|
| `[GET] /`               | Dashboard        | template |
| `[GET] /gophers`        | List all gophers | REST     |
| `[POST] /gophers`       | Create a gopher  | REST     |
| `[GET] /gophers/:id`    | Get a gopher     | REST     |
| `[DELETE] /gophers/:id` | Delete a gopher  | REST     |

### Authentication

This demo application provides an example [authentication middleware](internal/api/middleware/authentication.go).

You can enable authentication in the application [configuration file](configs/config.yaml) with `config.authentication.enabled=true`.

If enabled, you need to provide the secret configured in `config.authentication.secret` as request `Authorization` header.
