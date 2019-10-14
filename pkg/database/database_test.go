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

package database

import (
	"reflect"
	"testing"

	"github.com/tinyzimmer/protobuf-registry/pkg/config"
)

func TestGetEngine(t *testing.T) {
	config := &config.Config{DatabaseDriver: dbEngineMemory}
	engine := GetEngine(config)
	if reflect.TypeOf(engine).String() != "*memory.memoryDatabase" {
		t.Error("Expected memory database got:", reflect.TypeOf(engine).Name())
	}
	config.DatabaseDriver = "not-exists"
	engine = GetEngine(config)
	// should still return the default (and currently only) driver
	if reflect.TypeOf(engine).String() != "*memory.memoryDatabase" {
		t.Error("Expected memory database got:", reflect.TypeOf(engine).Name())
	}
}
