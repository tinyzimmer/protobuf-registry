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
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/tinyzimmer/protobuf-registry/pkg/util"
)

// rawZipFiles returns raw zip readers for the zip data
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

// newTempFilesFromRaw writes the raw files, and optionally descriptor set, and returns
// the path to the files and a map of directories to their children attributes
func (p *Protobuf) newTempFilesFromRaw(withDescriptor bool) (rootPath, descriptorPath string, filesInfo map[string][]os.FileInfo, err error) {
	if withDescriptor {
		if p.DescriptorBytes() == nil {
			err = errors.New("raw descriptor set is nil, need to call p.SetDescriptor() or load from storage")
			return
		}
	}
	var files []*zip.File
	if files, err = p.rawZipFiles(); err != nil {
		return
	}
	if rootPath, filesInfo, err = util.WriteZipFilesToTempDir(files); err != nil {
		return
	}
	if withDescriptor {
		descriptorPath = filepath.Join(rootPath, "descriptor.pb")
		err = ioutil.WriteFile(descriptorPath, p.DescriptorBytes(), 0600)
	}
	return
}
