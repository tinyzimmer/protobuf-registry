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
	"os"
	"testing"

	"github.com/tinyzimmer/protobuf-registry/pkg/config"
)

func TestCompileDescriptorSet(t *testing.T) {
	proto := newTestProtoWithData(t)
	_ = config.Init()
	config.GlobalConfig.ProtocPath = "echo"
	if err := proto.CompileDescriptorSet(); !os.IsNotExist(err) {
		t.Error("Expected only error from final read attempt got:", err)
	}
	config.GlobalConfig.ProtocPath = "bad-non-exist"
	if err := proto.CompileDescriptorSet(); err == nil {
		t.Error("Expected error from bad executable, got nil")
	}
	proto.SetRaw(nil)
	if err := proto.CompileDescriptorSet(); err == nil {
		t.Error("Expected error from no data, got nil")
	}
}
