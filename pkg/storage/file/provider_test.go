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

package file

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/tinyzimmer/protobuf-registry/pkg/config"
	"github.com/tinyzimmer/protobuf-registry/pkg/protobuf"
	"github.com/tinyzimmer/protobuf-registry/pkg/util"
)

func getTestProvider(t *testing.T) (*fileProvider, func()) {
	t.Helper()
	config.SafeInit()
	config.GlobalConfig.FileStoragePath, _ = ioutil.TempDir("", "")
	return NewProvider(config.GlobalConfig), func() { os.RemoveAll(config.GlobalConfig.FileStoragePath) }
}

func TestStoreAndGetRawProto(t *testing.T) {
	f, rm := getTestProvider(t)
	defer rm()

	proto := &protobuf.Protobuf{ID: util.StringPtr("test-id")}
	proto.SetRaw([]byte("test raw data"))
	proto.SetDescriptor([]byte("test descriptor data"))

	if err := f.StoreProtoPackage(proto); err != nil {
		t.Error("Expected no error, got:", err)
	}

	proto.SetRaw(nil)
	proto.SetDescriptor(nil)
	if err := f.StoreProtoPackage(proto); err == nil {
		t.Error("Expected error, got nil")
	}

	proto, err := f.GetRawProto(proto)
	if err != nil {
		t.Fatal("Expected no error on retrieving proto, got:", err)
	}
	if string(proto.Raw()) != "test raw data" {
		t.Error("Data was malformed on retrieval, got:", string(proto.Raw()))
	}

	if string(proto.DescriptorBytes()) != "test descriptor data" {
		t.Error("Descriptor data was malformed on retrieval, got:", string(proto.DescriptorBytes()))
	}

	descPath := filepath.Join(f.protoRoot(), *proto.ID, descriptorFileName)
	if err := os.RemoveAll(descPath); err != nil {
		t.Fatal("Failed to remove descriptor file")
	}
	proto, err = f.GetRawProto(proto)
	if err == nil {
		t.Fatal("Expected error on descriptor set gone, got nil")
	}

	if err := f.DeleteProtoPackage(proto); err != nil {
		t.Error("Expected no error deleting proto package, got:", err)
	}

	_, err = f.GetRawProto(proto)
	if err == nil {
		t.Fatal("Expected error on non-existant proto, got nil")
	}
}

func TestRawFiles(t *testing.T) {
	f, rm := getTestProvider(t)
	defer rm()

	if err := f.StoreRawFile("test-file", []byte("test-data")); err != nil {
		t.Error("Expected no error storing raw file, got:", err)
	}

	data, err := f.GetRawFile("test-file")
	if err != nil {
		t.Fatal("Expected to retrieve raw file, got:", err)
	}
	if string(data) != "test-data" {
		t.Error("Data was malformed on retrieval, got:", string(data))
	}

	if err := f.DeleteRawFile("test-file"); err != nil {
		t.Error("Expected no error deleting raw file, got:", err)
	}

	if _, err = f.GetRawFile("test-file"); err == nil {
		t.Error("Expected error getting deleted file, got nil")
	}
}
