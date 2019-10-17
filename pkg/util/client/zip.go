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

package client

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// DirToUploadBody takes a directory and creates an upload body that can be
// included in a PostProtoRequest. All files that don't end in .proto are
// ignored.
func DirToUploadBody(dir string) (string, error) {
	body, err := ZipDir(dir)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(body), nil
}

// ZipDir zips all .proto files in a directory and it's sub-directories. The
// zip depth starts at one level below the directory provided as the argument.
func ZipDir(dir string) ([]byte, error) {
	var buf bytes.Buffer
	zwr := zip.NewWriter(&buf)
	if err := addDirToZip(zwr, dir, filepath.Base(dir)); err != nil {
		zwr.Close()
		return nil, err
	}
	if err := zwr.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func addDirToZip(zwr *zip.Writer, dir string, stripPrefix string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() {
			if err := addDirToZip(zwr, filepath.Join(dir, file.Name()), stripPrefix); err != nil {
				return err
			}
		}
		if err := addFileToZip(zwr, filepath.Join(dir, file.Name()), stripPrefix); err != nil {
			return err
		}
	}
	return nil
}

func addFileToZip(zwr *zip.Writer, filepath string, stripPrefix string) error {
	if !strings.HasSuffix(filepath, ".proto") {
		return nil
	}
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	info, err := file.Stat()
	if err != nil {
		return err
	}
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = strings.TrimPrefix(filepath, stripPrefix+"/")
	header.Method = zip.Deflate

	wrtr, err := zwr.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(wrtr, file)
	return err
}
