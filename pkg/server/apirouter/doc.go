// Copyright © 2019 tinyzimmer
//
// This file is part of protobuf-registry
//
// protobuf-registry is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// protobuf-registry is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with protobuf-registry.  If not, see <https://www.gnu.org/licenses/>.

package apirouter

var routeDocumentation = map[string]map[string]string{

	"/api": {
		"GET": `Retrieve route information for the entire API`,
	},

	"/api/config": {
		"GET": `Retrieve the current server configuration`,
	},

	"/api/proto": {
		"GET": `Retrieve all protobuf specs currently in the registry`,
		"PUT": `Upload new protobuf spec, overwriting an existing one with the same name and version if it exists.

See POST /api/proto for more details`,
		"POST": `Upload new protobuf spec to the registry

Example Payloads
----------------

Self-contained package

{
  "name": "package_name",
  "version": "package_version (default: 0.0.1)",
  "body": "<base64_encoded_zip>"
}

Package with remote imports

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

Note that packages with remote imports will seek to fetch them mid-flight if they are not already cached via /api/remotes or the environment.`,
	},

	"/api/proto/{name}": {
		"GET":    `Get list of all versions for the spec {name}`,
		"DELETE": `Delete all versions for the spec {name}`,
	},

	"/api/proto/{name}/{version}": {
		"GET":    `Get details for version {version} of spec {name}`,
		"DELETE": `Delete version {version} of spec {name}`,
	},

	"/api/proto/{name}/{version}/{language}": {
		"GET": `Download version {version} of spec {name} in language {language}

Example
-------
> curl -JLO http://protoregistry.example.com/api/proto/my-app-protocol/0.0.1/descriptors

{language} options
------------------
raw (zip archive containing raw .proto files)
descriptors (descriptor set)
cpp
csharp
java
javanano (not functional)
js
objc
php
python
ruby`,
	},

	"/api/proto/{name}/{version}/raw/{filename}": {
		"GET": `Retrieve the raw contents of file {filename} from version {version} of package {name}

{filename} can be a full path with slashes (/) for a nested file in a protobuf package. However, a path to a directory will currently return an error.

This may be adapted to return a list of filenames when a directory is provided.`,
	},

	"/api/proto/{name}/{version}/meta/{filename}": {
		"GET": `Retrieves documentation for the data within {filename} of package {name} version {version}.`,
	},

	"/api/remotes": {
		"GET": "Retrieve a list of the currently cached remotes.",
		"PUT": `Ensure a remote is cached in the server.

Example payload
---------------
{
  "url": "github.com/googleapis/api-common-protos"
}
`,
	},

	"/pip/{name}/": {
		"GET": `Used for pip discovery - returns a list of available packages for spec {name} as parsed by pip.

Usage
-----
pip install --extra-index-url http://protoregistry.example.com/pip my-app-protocol`,
	},

	"/pip/download/{name}": {
		"GET": `Used for pip discovery - downloads the package/version specified by {name}`,
	},

	"/mvn/{name}/{version}": {
		"GET": `Unfornately, I don't feel like adding a java compiler to this image, but the API can still build a ready-to-go directory for packaging.

Perhaps later, since the server could be run outside docker or in an alternate image with a java installation, this route could be adapted for repo discovery as well.

Usage
-----
> curl http://protoregistry.example.com/mvn/example-proto-package/0.0.1 | tar xzf -
> cd example-proto-package && mvn package
`,
	},

	"/gem/specs.4.8.gz": {
		"GET": `Used for rubygems discovery - returns a ruby marshalled spec of available packages

** This functionality is not yet complete **

Usage
-----
gem install -s http://protoregistry.example.com/gem my-app-protocol`,
	},

	"/gem/quick/Marshal.4.8/{name}.gemspec.rz": {
		"GET": `(NOT WORKING) Used for rubygems discovery - returns a Gem::Specification for a given package`,
	},
}

func GetDoc(path, method string) string {
	if docObj, ok := routeDocumentation[path]; ok {
		if doc, ok := docObj[method]; ok {
			return doc
		}
	}
	return "No documentation for this route"
}
