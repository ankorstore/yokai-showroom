# Yokai HTTP Demo

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go version](https://img.shields.io/badge/Go-1.20-blue)](https://go.dev/)

> HTTP REST API demo application, based on
> the [Yokai](https://github.com/ankorstore/yokai) Go framework.

<!-- TOC -->
* [Overview](#overview)
* [Usage](#usage)
  * [Start the application](#start-the-application)
  * [Available endpoints](#available-endpoints)
  * [Authentication](#authentication)
* [Contents](#contents)
  * [Layout](#layout)
  * [Makefile](#makefile)
<!-- TOC -->

## Overview

This demo application is a simple REST API (CRUD) to manage [gophers](https://go.dev/blog/gopher).

It provides:

- a demo [Yokai](https://github.com/ankorstore/yokai) application container, with the [fxhttpserver](https://github.com/ankorstore/yokai/tree/main/fxhttpserver) module to offer the REST API
- a [MySQL](https://www.mysql.com/) container to store the gophers

See the [Yokai documentation](https://ankorstore.github.io/yokai) for more details.

## Usage

### Start the application

To start the application, simply run:

```shell
make fresh
```

After a short moment, the application will offer:

- [http://localhost:8080](http://localhost:8080): application dashboard
- [http://localhost:8081](http://localhost:8081): application core dashboard

### Available endpoints

| Route                   | Description      | Type     |
|-------------------------|------------------|----------|
| `[GET] /`               | Dashboard        | template |
| `[GET] /gophers`        | List all gophers | REST     |
| `[POST] /gophers`       | Create a gopher  | REST     |
| `[GET] /gophers/:id`    | Get a gopher     | REST     |
| `[PATCH] /gophers/:id`  | Update a gopher  | REST     |
| `[DELETE] /gophers/:id` | Delete a gopher  | REST     |

### Authentication

This demo application provides an example [authentication middleware](internal/middleware/authentication.go).

You can enable authentication in the application [configuration file](configs/config.yaml) with `config.authentication.enabled=true`.

If enabled, you need to provide the secret configured in `config.authentication.secret` as request `Authorization` header.

## Contents

### Layout

This template is following the [standard go project layout](https://github.com/golang-standards/project-layout):

- `cmd/`: entry points
- `configs/`: configuration files
- `internal/`:
  - `handler/`: HTTP handlers
  - `middleware/`: HTTP middlewares
  - `model/`: models
  - `service/`: services
  - `bootstrap.go`: bootstrap (modules, lifecycles, etc)
  - `routing.go`: routing
  - `services.go`: dependency injection

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
