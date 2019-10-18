# protobuf-registry

[![GPL License][gpl-img]][gpl]
![Build][build-img]
![Test][test-img]
[![Coverage Status](https://coveralls.io/repos/github/tinyzimmer/protobuf-registry/badge.svg?branch=master)](https://coveralls.io/github/tinyzimmer/protobuf-registry?branch=master)
[![Docker][docker-img]][docker]

A repository, package manager, and file viewer for Protocol Buffers.

## Description

This project is a repository for protocol buffer packages.
Groups of `.proto` files can be kept in the registry and versioned appropriately.
You can then retrieve details on your packages as well as generated code and descriptor sets with simple API calls.
There is also a Web UI available for visualizing your packages, their documentation, and other ad-hoc operations.

Where feasible, I want to support package manager discovery for various languages.
This functionality already works for `pip` and `go`, however it is obviously a challenge for other ones.
For example, with `maven`, I'd rather not include a Java compiler with the build image, but you can still download a ready-to-package directory.

### Uploading Packages

You can upload packages via the UI or the API.
Via the API it is done by sending a POST or PUT request to `/api/proto`.
A PUT request will overwrite an existing package of the same name and version, but a POST attempt to overwrite will fail.

An example self-contained package

```json
{
  "name": "package_name",
  "version": "package_version (default: 0.0.1)",
  "body": "<base64_encoded_zip>"
}
```

Package with remote imports

```json
{
  "name": "package_name",
  "version": "package_version (default: 0.0.1)",
  "body": "<base64_encoded_zip>",
  "remoteDeps": [
    {
      "url": "github.com/googleapis/api-common-protos",
      "revision": "master"
    }
  ]
}
```

The above example uses the Google common protocol buffer files, but it is not required to import those explicitly. They are automatically included in the package if they are needed.

### Pip discovery

For the pip discovery, if you uploaded a package called `my-app-protocol` with versions `0.0.1` and `0.0.2`:


```bash
# Assuming the registry is running locally on 8080

# Would install version 0.0.2
$> pip install --index-url http://localhost:8080/pip my-app-protocol
# Install version 0.0.1
$> pip install --index-url http://localhost:8080/pip my-app-protocol==0.0.1
```

### Gocode discovery

For the go discovery, it works based off the `go_package` option supplied in the `.proto` files.
This currently only works for the latest version of the package in the registry.
The option value must also correlate to, or at least resolve to, wherever you host the registry.

For example with a registry running locally at `http://registry.localhost` and a protocol spec (of any "package name") using:

```proto
option go_package = "registry.localhost/golang/test-protobuf";
```

You could then use `go get` to fetch the gocode like this:

```bash
$> go get -insecure registry.localhost/golang/test-protobuf

go: finding registry.localhost/golang/test-protobuf latest
go: downloading registry.localhost/golang/test-protobuf v0.0.0-20191014153456-93e6948efcdf
go: extracting registry.localhost/golang/test-protobuf v0.0.0-20191014153456-93e6948efcdf
go: finding google.golang.org/genproto latest
go: downloading google.golang.org/genproto v0.0.0-20191009194640-548a555dbc03
go: extracting google.golang.org/genproto v0.0.0-20191009194640-548a555dbc03
```

This should be possible for _most_ interpreted languages (ruby marshaling in golang is a bitch).

This project was primarily just me being bored and wanting to build something, but also was an opportunity for me to finally start learning some modern front-end technologies while making the UI. However, if I decide to keep building on it, I eventually want to incorporate elements from [`prototool`](https://github.com/uber/prototool) as well. Since it's also in go it would be easy to include some of their functionality (e.g. linting).

## Running

If you have a Kubernetes cluster and want to just try it out on that, I made a quick and dirty helm chart.
It doesn't support all the configuration options but for the most part you can just:

```
$> helm install --name proto-registry chart/
```

Refer to the `values.yaml` for the available options for now.

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
    tinyzimmer/protobuf-registry --persist-memory
```

The data directory will also hold the cache of remote repositories that are referenced by protobuf packages.
`POST` operations may take a while if they rely on large repositories for imported definitions that are not yet cached.
You can enforce a cache of certain repsitories by setting `PRE_CACHED_REMOTES` in the environment.
For example, `PRE_CACHED_REMOTES="github.com/googleapis/api-common-protos"`.

### Test data

If you don't have any `proto` files handy and want to see the UI with some data in it, then you can use the test data inside this repository.

If you have `go` installed and the server running at http://localhost:8080 then you can just run `make test_data`.
If you don't have `go` installed, just look at the `hack/add_test_data.sh` script and swap out the parts that are commented with the line that invokes `util.go`.
I used `util.go` to try out the client interface, but I left the commands that just use raw `zip` and `curl` calls in case those are easier.

I may make a full CLI from the client interfaces in `pkg/util/client` at some point.

### Configuration

The image can be configured via environment variables or on the command-line.
Options are limited right now, but it is setup in a way to easily add new interfaces for different backends.

| EnvVar             | Command-Line         |  Default                 |  Description                                                                       |
|:------------------:|:--------------------:|:------------------------:|:-----------------------------------------------------------------------------------|
|`BIND_ADDRESS`      |`--bind-address`      |`0.0.0.0:8080`             |The address and port to bind to.                                                   |
|`READ_TIMEOUT`      |`--read-timeout`      | `15`                      |Read timeout for API/UI `HTTP` requests.                                           |
|`WRITE_TIMEOUT`     |`--write-timeout`     | `15`                      |Write timeout for API/UI `HTTP` requests.                                          |
|`COMPILE_TIMEOUT`   |`--compile-timeout`   | `10`                      |Timeout for `protoc` invocations. Only applies to codegen now.                     |
|`PROTOC_PATH`       |`--protoc-path`       | `/usr/bin/protoc`         |Path to the `protoc` executable. Leave unchanged in docker image. Used for codegen.|
|`PROTOC_GEN_GO_PATH`|`--protoc-gen-go-path`|`/opt/proto-registry/bin/protoc-gen-go` |The path to a compiled `protoc-gen-go` plugin until I can get gocode generation to work independently via its exported interfaces|
|`DATABASE_DRIVER`   |`--database-driver`   | `memory`                  |Driver to use for database operations, only `memory` currently.                    |
|`STORAGE_DRIVER`    |`--storage-driver`    | `file`                    |Driver to use for file storage operations, only `file` currently.                  |
|`FILE_STORAGE_PATH` |`--file-storage-path` | `/opt/proto-registry/data`|Path to file storage when using file storage driver.                               |
|`PERSIST_MEMORY`    |`--persist-memory`    | `false`                   |Persist the in-memory database to disk after write operations.                     |
|`PRE_CACHED_REMOTES`|`--pre-cached-remotes`| `[]`                      |A comma-separated list of remote git repositories to pre-cache for compilations.   |
|`UI_REDIRECT_ALL`   |`--ui-redirect-all`   | `false`                   |Redirect all unknown routes to the UI. Useful to turn off for discovery debugging. |
|`ENABLE_CORS`       |`--enable-cors`       | `false`                   |Enable CORS headers for API requests.                                              |

## Development

Coming soon

#  

#### TODO

 - [ ] dev docs
 - [ ] fetch protoc as required only for codegen, and switch back to scratch image
 - [ ] validateOnly/linting on `POST /api/proto`
 - [ ] day theme/night theme
 - [ ] proto version codegen options - or maybe worker nodes
   - trying to keep the image small

#### Special Thanks

 - https://github.com/jhump/protoreflect - Protoreflect is an incredible library for dynamically parsing proto files and deserves special mention.
 - https://github.com/pseudomuto/protoc-gen-doc - Great for generating documentation directly from .proto files in any format you desire.

[build-img]: https://github.com/tinyzimmer/protobuf-registry/workflows/Build/badge.svg
[test-img]: https://github.com/tinyzimmer/protobuf-registry/workflows/Test/badge.svg
[gpl-img]: https://img.shields.io/badge/license-GPL-blue
[gpl]: https://github.com/tinyzimmer/protobuf-registry/blob/master/COPYING
[docker-img]: https://img.shields.io/badge/docker%20build-automated-066da5
[docker]: https://hub.docker.com/r/tinyzimmer/protobuf-registry
