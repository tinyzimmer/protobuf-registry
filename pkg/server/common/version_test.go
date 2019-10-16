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

package common

import (
	"testing"

	"github.com/tinyzimmer/protobuf-registry/pkg/protobuf"
	"github.com/tinyzimmer/protobuf-registry/pkg/util"
)

func TestGetVersionFromProtoSlice(t *testing.T) {

	protos := []*protobuf.Protobuf{
		{Version: util.StringPtr("0.0.1")},
		{Version: util.StringPtr("0.0.2")},
		{Version: util.StringPtr("0.1.0")},
	}

	proto, err := GetVersionFromProtoSlice(protos, "0.0.1")
	if err != nil {
		t.Fatal("Expected no error, got:", err)
	}
	if *proto.Version != "0.0.1" {
		t.Error("Got wrong version back, expected: 0.0.1, got:", *proto.Version)
	}

	_, err = GetVersionFromProtoSlice(protos, "0.0.5")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

}

func TestGetLatestVersion(t *testing.T) {
	protos := []*protobuf.Protobuf{
		{Version: util.StringPtr("invalid")},
		{Version: util.StringPtr("0.0.1")},
		{Version: util.StringPtr("0.0.2")},
		{Version: util.StringPtr("0.1.0")},
		{Version: util.StringPtr("invalid")},
	}
	latest := GetLatestVersion(protos)
	if *latest.Version != "0.1.0" {
		t.Error("Expected 0.1.0, to be the latest version, got:", *latest.Version)
	}

	protos = []*protobuf.Protobuf{
		{Version: util.StringPtr("0.1.0")},
	}
	latest = GetLatestVersion(protos)
	if *latest.Version != "0.1.0" {
		t.Error("Expected 0.1.0, to be the latest version, got:", *latest.Version)
	}
}

func TestSortVersions(t *testing.T) {
	protos := []*protobuf.Protobuf{
		{Version: util.StringPtr("invalid")},
		{Version: util.StringPtr("0.0.1")},
		{Version: util.StringPtr("0.0.2")},
		{Version: util.StringPtr("0.1.0")},
		{Version: util.StringPtr("invalid")},
	}
	expectedSorted := []*protobuf.Protobuf{
		{Version: util.StringPtr("0.1.0")},
		{Version: util.StringPtr("0.0.2")},
		{Version: util.StringPtr("0.0.1")},
		{Version: util.StringPtr("invalid")},
		{Version: util.StringPtr("invalid")},
	}

	sorted := SortVersions(protos)

	for i, x := range sorted {
		if *x.Version != *expectedSorted[i].Version {
			t.Error("Expected:", *expectedSorted[i].Version, "Got:", *x.Version)
		}
	}
}
