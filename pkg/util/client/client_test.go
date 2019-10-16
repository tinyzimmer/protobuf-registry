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

import "testing"

// Placeholder tests for coverage

func TestClient(t *testing.T) {
	client := New("test-url")
	client.GetServerConfig()
	client.UploadProtoPackage(nil, false)
	client.UploadProtoPackageFromDir("", false)
	client.GetProtoPackageVersions("test")
	client.DeleteAllProtoPackageVersions("test")
	client.GetProtoPackage("test", "0.0.1")
	client.DeleteProtoPackage("test", "0.0.1")
	client.DownloadProtoPackage("test", "0.0.1", "raw")
	client.GetFileContents("test", "0.0.1", "testmessage.proto")
	client.GetFileDocs("test", "0.0.1", "testmessage.proto")
	client.GetCachedRemotes()
	client.PutCachedRemote("test-remote")
}
