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

package util

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestTarGZArchive(t *testing.T) {
	dir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(dir)
	if err := ioutil.WriteFile(
		filepath.Join(dir, "test.txt"),
		[]byte("hello world"),
		0600,
	); err != nil {
		t.Fatal("Failed to write test file to dir")
	}

	out, err := NewTarGZArchive(dir)
	if err != nil {
		t.Error("Expected no error, got:", err)
	}
	if len(out) == 0 {
		t.Error("Got empty response")
	}
}
