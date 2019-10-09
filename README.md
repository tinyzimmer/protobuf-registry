# protobuf-registry

[![GPL License][gpl-img]][gpl]
![Build][build-img]
![Test][test-img]
[![Docker][docker-img]][docker]

A repository, package manager, and file viewer for Protocol Buffers.

## Description

This project is a repository for protocol buffer packages.
Groups of `.proto` files can be kept in the registry and versioned appropriately.
You can then retrieve details on your packages as well as generated code and descriptor sets with simple API calls.
There is also a Web UI available for visualizing your packages and other ad-hoc operations.

Where feasible, I want to support package manager discovery for various languages.
This functionality already works for `pip`, however it is obviously a challenge for other ones.
For example, with `maven`, I'd rather not include a Java compiler with the build image, but you can still download a ready-to-package directory.

For the pip discovery, if you uploaded a package called `my-app-protocol` with versions `0.0.1` and `0.0.2`:


```bash
# Assuming the registry is running locally on 8080

# Would install version 0.0.2
$> pip install --index-url http://localhost:8080/pip my-app-protocol
# Install version 0.0.1
$> pip install --index-url http://localhost:8080/pip my-app-protocol==0.0.1
```

This should be possible for _most_ interpreted languages (ruby marshaling in golang is a bitch). I'm sure I can do it with `go get` also. However, at time of writing, only pip has it.

This project was primarily just me being bored and wanting to build something, but also was an opportunity for me to finally start learning some modern front-end technologies while making the UI. However, if I decide to keep building on it, I eventually want to incorporate elements from [`prototool`](https://github.com/uber/prototool) as well. Since it's also in go it would be easy to include some of their functionality (e.g. linting).

Heck, maybe even add Avro support later.

## Running

You can build and run the image locally with `make run` or there is a docker image available at `tinyzimmer/protobuf-registry`.

Using docker:

```bash
$> docker run -p 8080:8080 tinyzimmer/protobuf-registry
```

You can then visit the UI at http://localhost:8080.
There are docs for the API in the UI, but you could also just pull them with curl using:

```bash
$> curl http://localhost:8080/api
```

**Note:** This runs the registry _without_ persistence.
To enable persistence run with the following flags:

```bash
$> docker run \
    -p 8080:8080 \
    -v "`pwd`/data:/data" \
    -e PROTO_REGISTRY_PERSIST_MEMORY=true \
    tinyzimmer/protobuf-registry
```

The data directory will also hold the cache of remote repositories that are referenced by protobuf packages.
`POST` operations may take a while if they rely on large repositories for imported definitions that are not yet cached.
You can enforce a cache of certain repsitories by setting `PROTO_REGISTRY_PRE_CACHED_REMOTES` in the environment.
For example, `PROTO_REGISTRY_PRE_CACHED_REMOTES="github.com/googleapis/googleapis"`.

### Configuration

The image can be configured via environment variables.
Options are limited right now, but it is setup in a way to easily add new interfaces for different backends.

| Name                              | Default           |Description                                                                       |
|:---------------------------------:|:-----------------:|:---------------------------------------------------------------------------------|
|`PROTO_REGISTRY_BIND_ADDRESS`      |`0.0.0.0:8080`     |The address and port to bind to.                                                  |
|`PROTO_REGISTRY_READ_TIMEOUT`      | `15`              |Read timeout for API/UI `HTTP` requests.                                          |
|`PROTO_REGISTRY_WRITE_TIMEOUT`     | `15`              |Write timeout for API/UI `HTTP` requests.                                         |
|`PROTO_REGISTRY_COMPILE_TIMEOUT`   | `10`              |Timeout for `protoc` invocations.                                                 |
|`PROTO_REGISTRY_PROTOC_PATH`       | `/usr/bin/protoc` |Path to the `protoc` executable. Leave unchanged in docker image.                 |
|`PROTO_REGISTRY_DATABASE_DRIVER`   | `memory`          |Driver to use for database operations, only `memory` currently.                   |
|`PROTO_REGISTRY_STORAGE_DRIVER`    | `file`            |Driver to use for file storage operations, only `file` currently.                 |
|`PROTO_REGISTRY_FILE_STORAGE_PATH` | `/data`           |Path to file storage when using file storage driver.                              |
|`PROTO_REGISTRY_PERSIST_MEMORY`    | `false`           |Persist the in-memory database to disk after write operations.                    |
|`PROTO_REGISTRY_PRE_CACHED_REMOTES`| `[]`              |A comma-separated list of remote git repositories to pre-cache for compilations.  |
|`PROTO_REGISTRY_UI_REDIRECT_ALL`   | `true`            |Redirect all unknown routes to the UI. Useful to turn off for discovery debugging.|
|`PROTO_REGISTRY_ENABLE_CORS`       | `false`           |Enable CORS headers for API requests.                                             |

## Development

Coming soon

#  

#### TODO

 - [ ] unit tests
 - [ ] dev docs
 - [ ] validateOnly/linting on `POST /api/proto`
 - [ ] day theme/night theme
 - [ ] proto version compile options - or maybe worker nodes
   - trying to keep the image small


[build-img]: https://github.com/tinyzimmer/protobuf-registry/workflows/Build/badge.svg
[test-img]: https://github.com/tinyzimmer/protobuf-registry/workflows/Test/badge.svg
[gpl-img]: https://img.shields.io/badge/license-GPL-blue
[gpl]: https://github.com/tinyzimmer/protobuf-registry/blob/master/COPYING
[docker-img]: https://img.shields.io/badge/docker%20build-automated-066da5
[docker]: https://hub.docker.com/r/tinyzimmer/protobuf-registry
