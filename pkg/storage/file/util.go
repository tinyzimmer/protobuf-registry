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

package file

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func (f *fileProvider) root() string {
	return f.conf.FileStoragePath
}

func (f *fileProvider) rawFilePath() string {
	return filepath.Join(f.root(), "raw")
}

func (f *fileProvider) mkdirAll(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0700); err != nil {
			return err
		}
	}
	return nil
}

func (f *fileProvider) writeRawToDisk(path string, data []byte) error {
	if err := f.mkdirAll(filepath.Dir(path)); err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0644)
}

func (f *fileProvider) deleteFromDisk(path string) error {
	// If the path does not exist, RemoveAll returns nil (no error)
	return os.RemoveAll(path)
}
