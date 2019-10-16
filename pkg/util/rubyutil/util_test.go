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

package rubyutil

import (
	"testing"

	"github.com/tinyzimmer/protobuf-registry/pkg/protobuf"
	"github.com/tinyzimmer/protobuf-registry/pkg/util"
)

func TestRubyGemsList(t *testing.T) {
	pkgs := []*protobuf.Protobuf{
		{Name: util.StringPtr("test-proto"), Version: util.StringPtr("0.0.1")},
	}

	out, err := NewRubyGemsListFromPackages(pkgs)
	if err != nil {
		t.Error("Expected no error, got:", err)
	}
	if len(out) == 0 {
		t.Error("Expected output, got nil")
	}

	out, err = NewGemSpecFromPackage(pkgs[0])
	if err != nil {
		t.Error("Expected no error, got:", err)
	}
	if len(out) == 0 {
		t.Error("Expected output, got nil")
	}
}
