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

func TestCopyDir(t *testing.T) {
	dir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(dir)
	if err := os.MkdirAll(filepath.Join(dir, "data"), 0700); err != nil {
		t.Fatal("Failed to make nested directory")
	}
	if err := ioutil.WriteFile(
		filepath.Join(dir, "data", "test.txt"),
		[]byte("hello world"),
		0600,
	); err != nil {
		t.Fatal("Failed to write test file to dir")
	}

	outDir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(outDir)

	if err := CopyDir(dir, filepath.Join(outDir, "out")); err != nil {
		t.Error("Expected no error, got:", err)
	}

	if err := CopyDir("/non/exist/path", "/non/exist/path"); err == nil {
		t.Error("Expected error on non-existent paths, got nil")
	}
}
