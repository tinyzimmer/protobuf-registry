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

package protobuf

import (
	"encoding/base64"
	"io/ioutil"
	"testing"

	"github.com/tinyzimmer/protobuf-registry/pkg/types"
	"github.com/tinyzimmer/protobuf-registry/pkg/util"
)

func newTestProto(t *testing.T) *Protobuf {
	t.Helper()
	return &Protobuf{
		Name:    util.StringPtr("test-proto"),
		Version: util.StringPtr("0.0.1"),
	}
}

func newTestProtoWithData(t *testing.T) *Protobuf {
	t.Helper()
	proto := newTestProto(t)
	if err := proto.SetRawFromBase64(testProtoZip); err != nil {
		t.Fatal(err)
	}
	desc, err := base64.StdEncoding.DecodeString(descriptorSetBase64)
	if err != nil {
		t.Fatal(err)
	}
	proto.SetDescriptor(desc)
	return proto
}

func TestNewFromRequest(t *testing.T) {
	req := &types.PostProtoRequest{
		Name:    "test-proto",
		Version: "0.0.1",
	}
	proto := NewFromRequest(req)
	if *proto.Name != req.Name {
		t.Error("Name was malformed, got:", *proto.Name)
	} else if *proto.Version != req.Version {
		t.Error("Version was malformed, got:", *proto.Version)
	}
}

func TestSetRaw(t *testing.T) {
	proto := newTestProto(t)
	if err := proto.SetRawFromBase64(testProtoZip); err != nil {
		t.Error("Expected no error on valid base64, got:", err)
	}
	if err := proto.SetRawFromBase64("non-valid base64"); err == nil {
		t.Error("Expected error from invalid base64, got nil")
	}

	base64Raw := base64.StdEncoding.EncodeToString(proto.Raw())
	if base64Raw != testProtoZip {
		t.Error("Data was malformed between decoding and re-encoding")
	}
}

func TestDescriptorSetter(t *testing.T) {
	proto := newTestProto(t)
	proto.SetDescriptor([]byte("some data"))
	if string(proto.DescriptorBytes()) != "some data" {
		t.Error("Descriptor data was malformed after setting")
	}
}

func TestRawFilename(t *testing.T) {
	proto := newTestProto(t)
	if proto.RawFilename() != "test-proto-0.0.1.zip" {
		t.Error("Raw filename was not expected, got:", proto.RawFilename())
	}
}

func TestReaders(t *testing.T) {
	proto := newTestProtoWithData(t)
	body, err := ioutil.ReadAll(proto.RawReader())
	if err != nil {
		t.Error("Could not read body from proto reader")
	}
	base64Raw := base64.StdEncoding.EncodeToString(body)
	if base64Raw != testProtoZip {
		t.Error("Data was malformed after going through reader")
	}

	proto.SetDescriptor([]byte("some data"))
	body, err = ioutil.ReadAll(proto.DescriptorReader())
	if err != nil {
		t.Error("Could not read body from proto descriptor reader")
	}
	if string(body) != "some data" {
		t.Error("Data was malformed after going through reader")
	}
}

func TestSHA256(t *testing.T) {
	proto := newTestProto(t)
	if _, err := proto.SHA256(); err == nil {
		t.Error("Expected error from no sha data, got nil")
	}

	proto = newTestProtoWithData(t)
	if sha, err := proto.SHA256(); err != nil {
		t.Error("Expected no error getting sha256, got:", err)
	} else if sha != testProtoSha {
		t.Error("SHA256SUM was not expected:", testProtoSha, "got:", sha)
	}
}
