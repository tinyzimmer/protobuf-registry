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
	"testing"

	"github.com/tinyzimmer/protobuf-registry/pkg/config"
)

func TestDocJSON(t *testing.T) {
	proto := newTestProtoWithData(t)
	config.SafeInit()
	config.GlobalConfig.ProtobufVersion = "3.6.1"
	out, err := proto.DocJSON("TestProtoMessage.proto")
	if err != nil {
		t.Fatal("Expected no error, got:", err)
	}
	if out == nil {
		t.Error("Expected JSON response, got nil")
	}

	_, err = proto.DocJSON("non-exist")
	if err == nil {
		t.Error("Expected error for non-exist file, got nil")
	}

	proto.SetDescriptor(nil)
	if _, err = proto.DocJSON("TestProtoMessage.proto"); err == nil {
		t.Error("Expected error from no descriptor set, got nil")
	}
}
