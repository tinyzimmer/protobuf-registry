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
)

func TestRawZipGetter(t *testing.T) {
	proto := newTestProtoWithData(t)

	// test all is well
	files, err := proto.rawZipFiles()
	if err != nil {
		t.Error("Expected no error, got nil")
	} else if len(files) != 7 {
		t.Error("Expected 7 files in the zip, got:", len(files))
	}

	// no data
	proto.SetRaw(nil)
	if _, err := proto.rawZipFiles(); err == nil {
		t.Error("Expected error from no zip, got nil")
	}

	// invalid zip
	proto.SetRaw([]byte("invalid zip data"))
	if _, err := proto.rawZipFiles(); err == nil {
		t.Error("Expected error from bad zip data, got nil")
	}
}

func TestTempFilesFromRaw(t *testing.T) {
	proto := newTestProtoWithData(t)
	proto.SetDescriptor([]byte("some data"))
	path, descPath, filesInfo, err := proto.newTempFilesFromRaw(true)
	if err != nil {
		t.Error("Expected no error, got:", err)
	}
	if descPath == "" {
		t.Error("Expected to get the path to a descriptor set back as well")
	}
	if len(filesInfo) != 3 {
		t.Error("Expected 3 sub directories")
	}
	defer os.RemoveAll(path)

	proto.SetDescriptor(nil)
	if _, _, _, err := proto.newTempFilesFromRaw(true); err == nil {
		t.Error("Expected error from no descriptor bytes, got nil")
	}

	proto.SetRaw(nil)
	if _, _, _, err := proto.newTempFilesFromRaw(false); err == nil {
		t.Error("Expected error from no zip bytes, got nil")
	}
}
