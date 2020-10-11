# Basic demo

This program demonstrates the basic functionality of `agoradb`.

The focus in this demo is on showcasing:

* [Prerequisites](#prerequisites)
    * [Install cli](#agoradb-cli)
    * [A running cluster](#a-running-agoradb-cluster)
* [Setting up a client](#client-setup)
* [Migrations, changing the schema](#migrations-changing-the-schema)
* [Generate the client api](#generate-the-client-api)
* [Basic usage of the client](#basic-usage-of-the-client)

## Prerequisites

#### agoradb-cli
Installing the `agoradb` cli is only required for the client api generation:
```
go get -u github.com/featme-inc/agoradb/cli/
```

#### A running agoradb cluster
Make sure to have a running `agoradb` cluster. Check:

```bash
docker-compose up
``` 

## Client setup

Setting up a client [main.go](main.go)

## Migrations, changing the schema

perform migrations, changing, defining the schema [migrate.go](migrate.go)

## Generate the client API

This section deals with the generation of the client api for the schema. Also note that the we have generated 
the client and it is checked in with the demo, so you don't need to worry about it. 

One of the simplest way to directly embed the following line in you main.go:

```go
//go:generate agoradb generate client --addr=loclhost:5750 --auth=BearerToken --grpc --output_dir=basic basic
func main() {
 ...
}
```
The options provided for `--addr` and `--output_dir` are the default values. 

## Basic usage of the client

This section concentrates on the basic usage of generated library as seen in [demo.go](demo.go).