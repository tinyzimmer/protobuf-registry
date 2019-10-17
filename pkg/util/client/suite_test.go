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

package client

import (
	"reflect"
	"testing"

	"github.com/tinyzimmer/protobuf-registry/pkg/config"
	"github.com/tinyzimmer/protobuf-registry/pkg/protobuf"
	"github.com/tinyzimmer/protobuf-registry/pkg/server"
	"github.com/tinyzimmer/protobuf-registry/pkg/types"
)

// Placeholder tests for coverage

func doUpload(t *testing.T, client RegistryClient, version string) {
	t.Helper()
	proto, err := client.UploadProtoPackage(&types.PostProtoRequest{
		Name:    "test-proto",
		Version: version,
		Body:    protobuf.TestProtoZip,
		RemoteDepends: []*protobuf.ProtoDependency{
			{URL: "github.com/googleapis/api-common-protos"},
		},
	}, false)
	if err != nil {
		t.Fatal(err)
	}
	if *proto.Name != "test-proto" {
		t.Error("Name was malformed on retrieval")
	}
	if version != "" && *proto.Version != version {
		t.Error("Expected proto to come back with default version")
	}
}

func TestClient(t *testing.T) {
	// get a test server and defer cleanup
	srvr, close, err := server.NewTestServer()
	if err != nil {
		t.Fatal(err)
	}
	defer close()

	// open a client pointing at the test server
	client := New(srvr.URL)

	// Test Get Server Config
	conf, err := client.GetServerConfig()
	if err != nil {
		t.Fatal(err)
	}
	// Config should be the same as the current GlobalConfig
	if !reflect.DeepEqual(conf, config.GlobalConfig) {
		t.Error("Did not get a valid config object back")
	}

	// Test list empty packages
	pkgs, err := client.ListProtoPackages()
	if err != nil {
		t.Fatal(err)
	}
	if len(pkgs.Items) != 0 {
		t.Error("Expected empty list of packages back, got:", pkgs.Items)
	}

	// Upload new proto package version
	doUpload(t, client, "")

	// Test list with single package
	pkgs, err = client.ListProtoPackages()
	if err != nil {
		t.Fatal(err)
	}
	if len(pkgs.Items) != 1 {
		t.Error("Expected empty list of packages back, got:", pkgs.Items)
	}

	// get list of protos with previously uploaded name
	protos, err := client.GetProtoPackageVersions("test-proto")
	if err != nil {
		t.Fatal(err)
	}
	if len(protos) != 1 {
		t.Error("Expected only one version, got:", protos)
	}

	// test single version retrieval
	_, err = client.GetProtoPackage("test-proto", "0.0.1")
	if err != nil {
		t.Error("Expected to retrieve package version")
	}
	_, err = client.GetProtoPackage("test-proto", "0.0.2")
	if err == nil {
		t.Error("Retrieved non-existant version")
	}

	if err := client.DeleteAllProtoPackageVersions("test-proto"); err != nil {
		t.Fatal(err)
	}

	if err := client.DeleteProtoPackage("test-proto", "0.0.1"); err == nil {
		t.Error("Expected to not find already deleted package")
	}

	// re-upload package
	doUpload(t, client, "")

	// test download raw
	out, err := client.DownloadProtoPackage("test-proto", "0.0.1", "raw")
	if err != nil {
		t.Fatal(err)
	}
	if len(out) == 0 {
		t.Error("Expected to get back zip bytes")
	}

	// test download descriptors
	out, err = client.DownloadProtoPackage("test-proto", "0.0.1", "descriptors")
	if err != nil {
		t.Fatal(err)
	}
	if len(out) == 0 {
		t.Error("Expected to get back zip bytes")
	}

	// test get file contents
	content, err := client.GetFileContents("test-proto", "0.0.1", "TestProtoMessage.proto")
	if err != nil {
		t.Fatal(err)
	}
	if content.Content == "" {
		t.Error("Got back empty contents")
	}

	_, err = client.GetFileDocs("test-proto", "0.0.1", "TestProtoMessage.proto")
	if err != nil {
		t.Fatal(err)
	}

	remotes, err := client.GetCachedRemotes()
	if err != nil {
		t.Fatal(err)
	}
	if len(remotes) != 1 {
		t.Error("Expected list of one remote")
	}

	if err := client.PutCachedRemote(&types.PutRemoteRequest{URL: "github.com/googleapis/api-common-protos"}); err != nil {
		t.Fatal(err)
	}
}
