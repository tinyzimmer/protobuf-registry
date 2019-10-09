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
	"strconv"
	"strings"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	docreq "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/jhump/protoreflect/desc"
	docgen "github.com/pseudomuto/protoc-gen-doc"
	"github.com/tinyzimmer/proto-registry/pkg/config"
)

func appendPkgsFromDescriptor(p *ProtobufDescriptors, f *desc.FileDescriptor) (o *ProtobufDescriptors) {
	opts := f.GetFileOptions()
	if opts == nil {
		return p
	}
	if opts.JavaPackage != nil {
		if !contains(p.JavaPackages, opts.JavaPackage) {
			p.JavaPackages = append(p.JavaPackages, opts.JavaPackage)
		}
	}
	if opts.GoPackage != nil {
		if !contains(p.GoPackages, opts.GoPackage) {
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

func intPtr(s string) *int32 {
	i, _ := strconv.Atoi(s)
	i32 := int32(i)
	return &i32
}

func parseProtocVersion() *docreq.Version {
	spl := strings.Split(config.GlobalConfig.ProtobufVersion, " ")
	vers := spl[len(spl)-1]
	versSplit := strings.Split(vers, ".")
	return &docreq.Version{
		Major: intPtr(versSplit[0]),
		Minor: intPtr(versSplit[1]),
		Patch: intPtr(versSplit[2]),
	}
}

func (p *Protobuf) DocJSON(filename string) ([]byte, error) {
	// write raw proto to temp files
	descriptors, err := p.GetDescriptors()
	if err != nil {
		return nil, err
	}
	var desc *desc.FileDescriptor
	var rawDescriptors []*descriptor.FileDescriptorProto
	for _, x := range descriptors {
		rawDescriptors = append(rawDescriptors, x.AsFileDescriptorProto())
		if x.GetName() == filename || x.GetName() == strings.TrimPrefix(filename, "/") {
			desc = x
			break
		}
	}
	if desc == nil {
		return nil, fmt.Errorf("No file %s in this protobuf package", filename)
	}
	plugin := docgen.Plugin{}
	param := "json,docs.json"
	res, err := plugin.Generate(&docreq.CodeGeneratorRequest{
		FileToGenerate:  []string{desc.GetName()},
		ProtoFile:       rawDescriptors,
		CompilerVersion: parseProtocVersion(),
		Parameter:       &param,
	})
	if err != nil {
		return nil, err
	} else if len(res.File) == 0 {
		return nil, fmt.Errorf("No documentation returned from the plugin")
	}
	content := *res.File[0].Content
	return []byte(content), nil
}

func (p *Protobuf) Contents(filename string) ([]byte, error) {
	// write raw proto to temp files
	tempPath, filesInfo, remove, err := p.newTempFilesFromRaw()
	if err != nil {
		return nil, err
	}
	defer remove()

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

func contains(slice []*string, s *string) bool {
	for _, x := range slice {
		if *x == *s {
			return true
		}
	}
	return false
}
