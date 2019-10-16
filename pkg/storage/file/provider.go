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

	"github.com/go-logr/glogr"
	"github.com/tinyzimmer/protobuf-registry/pkg/config"
	"github.com/tinyzimmer/protobuf-registry/pkg/protobuf"
	"github.com/tinyzimmer/protobuf-registry/pkg/storage/common"
)

const zipFileName = "proto.zip"
const descriptorFileName = "descriptors.pb"

var log = glogr.New()

type fileProvider struct {
	common.Provider

	conf *config.Config
}

func NewProvider(conf *config.Config) common.Provider {
	return &fileProvider{conf: conf}
}

func (f *fileProvider) GetRawProto(in *protobuf.Protobuf) (*protobuf.Protobuf, error) {
	path := filepath.Join(f.protoRoot(), *in.ID, zipFileName)
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return in, err
	}
	in.SetRaw(raw)

	descriptorPath := filepath.Join(f.protoRoot(), *in.ID, descriptorFileName)
	dRaw, err := ioutil.ReadFile(descriptorPath)
	if err != nil {
		return in, err
	}
	in.SetDescriptor(dRaw)
	return in, nil
}

func (f *fileProvider) StoreProtoPackage(in *protobuf.Protobuf) error {
	path := filepath.Join(f.protoRoot(), *in.ID, zipFileName)
	if err := f.writeRawToDisk(path, in.Raw()); err != nil {
		return err
	}
	descriptorPath := filepath.Join(f.protoRoot(), *in.ID, descriptorFileName)
	return f.writeRawToDisk(descriptorPath, in.DescriptorBytes())
}

func (f *fileProvider) DeleteProtoPackage(in *protobuf.Protobuf) error {
	path := filepath.Join(f.protoRoot(), *in.ID)
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
