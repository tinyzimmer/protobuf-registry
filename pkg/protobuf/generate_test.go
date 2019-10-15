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

func TestGenerateTo(t *testing.T) {
	proto := newTestProtoWithData(t)
	os.Setenv("IGNORE_PROTOC", "true")
	_ = config.Init()
	// set protoc executable to echo
	config.GlobalConfig.ProtocPath = "echo"

	proto.SetDescriptor(nil)
	if _, _, err := proto.GenerateTo(GenerateTargetGo, "prefix"); err == nil {
		t.Error("Expected error from no descriptors, got:", err)
	}

	proto.SetDescriptor([]byte("some data"))
	for _, x := range []GenerateTarget{
		GenerateTargetCPP,
		GenerateTargetCSharp,
		GenerateTargetJava,
		GenerateTargetJavaNano,
		GenerateTargetJS,
		GenerateTargetObjC,
		GenerateTargetPHP,
		GenerateTargetPython,
		GenerateTargetRuby,
		GenerateTargetGo,
	} {
		_, rm, err := proto.GenerateTo(x, "")
		if err != nil {
			t.Fatal("Expected no error, got:", err)
		}
		rm()
	}

	config.GlobalConfig.ProtocPath = "/non/existant/exec"
	if _, _, err := proto.GenerateTo(GenerateTargetGo, "prefix"); err == nil {
		t.Error("Expected error from bad executable, got:", err)
	}

	os.Unsetenv("IGNORE_PROTOC")
}
