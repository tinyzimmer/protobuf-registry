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
	"archive/zip"
	"bytes"
	"errors"
	"os"

	"github.com/tinyzimmer/proto-registry/pkg/util"
)

func (p *Protobuf) rawZipFiles() ([]*zip.File, error) {
	if p.Raw() == nil {
		return nil, errors.New("raw zip is nil, need to call p.SetRaw() or p.SetRawFromBase64()")
	}
	size := len(p.Raw())
	reader, err := zip.NewReader(bytes.NewReader(p.Raw()), int64(size))
	if err != nil {
		return nil, err
	}
	return reader.File, nil
}

func (p *Protobuf) newTempFilesFromRaw() (path string, filenames []string, rm func(), err error) {
	var files []*zip.File
	if files, err = p.rawZipFiles(); err != nil {
		return
	}
	if path, filenames, err = util.WriteZipFilesToTempDir(files); err != nil {
		return
	}
	rm = func() {
		if err := os.RemoveAll(path); err != nil {
			log.Error(err, "Failed to remove tempdir")
		}
	}
	return
}
