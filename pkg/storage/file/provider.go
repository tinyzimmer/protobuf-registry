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

	"github.com/tinyzimmer/proto-registry/pkg/config"
	"github.com/tinyzimmer/proto-registry/pkg/protobuf"
	"github.com/tinyzimmer/proto-registry/pkg/storage/common"
)

const zipFileName = "proto.zip"

type fileProvider struct {
	common.Provider

	conf *config.Config
}

func NewProvider(conf *config.Config) common.Provider {
	return &fileProvider{conf: conf}
}

func (f *fileProvider) GetRawProto(in *protobuf.Protobuf) (*protobuf.Protobuf, error) {
	path := filepath.Join(f.root(), *in.ID, zipFileName)
	var file *os.File
	var err error
	if file, err = os.Open(path); err != nil {
		return in, err
	}
	defer file.Close()
	raw, err := ioutil.ReadAll(file)
	if err != nil {
		return in, err
	}
	in.SetRaw(raw)
	return in, nil
}

func (f *fileProvider) StoreProtoPackage(in *protobuf.Protobuf) error {
	path := filepath.Join(f.root(), *in.ID, zipFileName)
	return f.writeRawToDisk(path, in.Raw())
}

func (f *fileProvider) DeleteProtoPackage(in *protobuf.Protobuf) error {
	path := filepath.Join(f.root(), *in.ID)
	return f.deleteFromDisk(path)
}

func (f *fileProvider) GetRawFile(name string) ([]byte, error) {
	path := filepath.Join(f.rawFilePath(), name)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ioutil.ReadAll(file)
}

func (f *fileProvider) StoreRawFile(name string, data []byte) error {
	path := filepath.Join(f.rawFilePath(), name)
	return f.writeRawToDisk(path, data)
}

func (f *fileProvider) DeleteRawFile(name string) error {
	path := filepath.Join(f.rawFilePath(), name)
	return f.deleteFromDisk(path)
}
