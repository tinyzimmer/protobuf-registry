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
	"archive/zip"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteZipFilesToTempDir(t *testing.T) {

	// set up a tempdir
	tempDir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(tempDir)
	tempFile := filepath.Join(tempDir, "test.txt")
	err := ioutil.WriteFile(tempFile, []byte("test string"), 0644)
	if err != nil {
		t.Fatal("Failed to write test file")
	}

	// write a new zip file to an in-memory buffer
	var buf bytes.Buffer
	wrt := zip.NewWriter(&buf)

	file, _ := os.Open(tempFile)
	info, _ := file.Stat()
	header, _ := zip.FileInfoHeader(info)

	// make it a nested header - will cause the dir check to happen
	header.Name = filepath.Join("test-dir", "test.txt")

	// add the tempfile to the zipfile
	writer, err := wrt.CreateHeader(header)
	if err != nil {
		t.Fatal("Failed to write zip header")
	}
	_, err = io.Copy(writer, file)
	if err != nil {
		t.Fatal("Failed to write file to buffer")
	}

	// close and flush the zip buffer
	wrt.Close()
	zipBytes := buf.Bytes()

	// open a reader on the zip
	zipReader, err := zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))
	if err != nil {
		t.Fatal("Failed to open a reader on the zip buffer")
	}

	// pass the reader through the function
	path, _, err := WriteZipFilesToTempDir(zipReader.File)
	if err != nil {
		t.Error("Expected to write zip to temp dir, got error:", err)
	}
	defer os.RemoveAll(path)

	// the root of the generated dir should have one dir entry
	files, _ := ioutil.ReadDir(path)
	if len(files) != 1 {
		t.Error("Expected temp path to have one entry, got:", len(files))
	} else if !files[0].IsDir() {
		t.Error("Expected first entry to be dir, got something else")
	}

	// the file contents of the tempfile should match the new one
	file, err = os.Open(filepath.Join(path, files[0].Name(), "test.txt"))
	if err != nil {
		t.Error("Could not open unzipped test file")
	}
	body, err := ioutil.ReadAll(file)
	if err != nil {
		t.Error("Could not read the unzipped file")
	} else if string(body) != "test string" {
		t.Error("Body was malformed when unzipped, got:", string(body))
	}
}
