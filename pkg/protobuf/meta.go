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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/jhump/protoreflect/desc"
	"github.com/tinyzimmer/protobuf-registry/pkg/util"
)

func appendPkgsFromDescriptor(p *ProtobufDescriptors, f *desc.FileDescriptor) (o *ProtobufDescriptors) {
	opts := f.GetFileOptions()
	if opts == nil {
		return p
	}
	if opts.JavaPackage != nil {
		if !util.StringPtrSliceContains(p.JavaPackages, opts.JavaPackage) {
			p.JavaPackages = append(p.JavaPackages, opts.JavaPackage)
		}
	}
	if opts.GoPackage != nil {
		if !util.StringPtrSliceContains(p.GoPackages, opts.GoPackage) {
			p.GoPackages = append(p.GoPackages, opts.GoPackage)
		}
	}
	return p
}

func protoMessageFromDescriptor(msg *desc.MessageDescriptor) *ProtobufMessage {
	return &ProtobufMessage{
		Name:   msg.GetName(),
		Fields: parseMessageFields(msg.GetFields()),
	}
}

func parseMessageFields(fields []*desc.FieldDescriptor) map[string]string {
	fieldData := make(map[string]string)
	for _, field := range fields {
		msgType := field.GetMessageType()
		if msgType != nil {
			fieldData[field.GetName()] = fmt.Sprintf("%s:%s", field.GetType().String(), msgType.GetFullyQualifiedName())
		} else {
			fieldData[field.GetName()] = field.GetType().String()
		}
	}
	return fieldData
}

func (p *Protobuf) Contents(filename string) ([]byte, error) {
	// write raw proto to temp files
	tempPath, _, filesInfo, err := p.newTempFilesFromRaw(false)
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tempPath)

	filePath, err := getFilePathFromZipFiles(tempPath, filesInfo, filename)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadFile(filePath)
}

func getFilePathFromZipFiles(tempPath string, filesInfo map[string][]os.FileInfo, filename string) (string, error) {
	var file string
	var err error
	for dir, files := range filesInfo {
		rawDir := strings.TrimPrefix(strings.Replace(dir, tempPath, "", 1), "/")
		for _, x := range files {
			if len(strings.Split(filename, "/")) == 1 {
				if x.Name() == filename {
					file = filepath.Join(dir, x.Name())
					break
				}
			}
			if rawDir == strings.TrimPrefix(filepath.Dir(filename), "/") && strings.TrimPrefix(x.Name(), "/") == filepath.Base(filename) {
				file = filepath.Join(dir, x.Name())
				break
			}
		}
	}
	if file == "" {
		err = fmt.Errorf("No file %s in this protobuf package", filename)
	}
	return file, err
}
