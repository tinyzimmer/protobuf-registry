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

package config

import (
	"os"
	"reflect"
	"testing"
)

var defaults = map[string]interface{}{
	"BindAddress":          "0.0.0.0:8080",
	"ReadTimeout":          15,
	"WriteTimeout":         15,
	"CompileTimeout":       10,
	"ProtocPath":           "/usr/bin/protoc",
	"ProtocGenGoPath":      "/opt/proto-registry/bin/protoc-gen-go",
	"DatabaseDriver":       "memory",
	"StorageDriver":        "file",
	"FileStoragePath":      "/opt/proto-registry/data",
	"PersistMemoryToDisk":  false,
	"PreCachedRemotes":     []string{},
	"RedirectNotFoundToUI": true,
	"CORSEnabled":          false,
}

func TestInit(t *testing.T) {
	os.Setenv("PROTO_REGISTRY_PROTOC_PATH", "/not/exists")
	if err := Init(); err == nil {
		t.Error("Expected error no protoc, got nil")
	}
	os.Setenv("PROTO_REGISTRY_IGNORE_PROTOC", "true")
	if err := Init(); err != nil {
		t.Error("Expected to ignore protoc, got:", err)
	}
	os.Unsetenv("PROTO_REGISTRY_PROTOC_PATH")

	jsonbytes := GlobalConfig.JSON()
	if len(jsonbytes) == 0 {
		t.Error("Expected json response, got empty byte slice")
	}

}

func TestBadConfig(t *testing.T) {
	os.Setenv("PROTO_REGISTRY_WRITE_TIMEOUT", "notint")
	if _, err := newConfig(); err == nil {
		t.Error("Expected error got nil")
	}
	os.Unsetenv("PROTO_REGISTRY_WRITE_TIMEOUT")
}

func TestDefaults(t *testing.T) {
	os.Setenv("PROTO_REGISTRY_IGNORE_PROTOC", "true")
	c, err := newConfig()
	if err != nil {
		t.Error("Expected no error on new config, got:", err)
		return
	}
	for k, v := range defaults {
		ftype := reflect.TypeOf(v).Name()
		switch ftype {
		case "string":
			val := getStrField(c, k)
			if val != v {
				t.Error("Field:", k, "Expected:", v, "got:", val)
			}
		case "int":
			val := getIntField(c, k)
			if val != v {
				t.Error("Field:", k, "Expected:", v, "got:", val)
			}
		case "bool":
			val := getBoolField(c, k)
			if val != v {
				t.Error("Field:", k, "Expected:", v, "got:", val)
			}
		default:
			// TODO: Test string slices
		}
	}
}

func getStrField(c *Config, field string) string {
	r := reflect.ValueOf(c)
	f := reflect.Indirect(r).FieldByName(field)
	return f.String()
}

func getIntField(c *Config, field string) int {
	r := reflect.ValueOf(c)
	f := reflect.Indirect(r).FieldByName(field)
	return int(f.Int())
}

func getBoolField(c *Config, field string) bool {
	r := reflect.ValueOf(c)
	f := reflect.Indirect(r).FieldByName(field)
	return f.Bool()
}
