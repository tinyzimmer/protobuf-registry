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
	"strings"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc"
)

// ProtobufDescriptors is a more human readable representation of raw
// file descriptors
type ProtobufDescriptors struct {
	Messages    []*ProtobufMessage `json:"messages"`
	SourceFiles []string           `json:"sourceFiles"`

	// Will add more as they are needed
	JavaPackages []*string `json:"javaPackages,omitempty"`
	GoPackages   []*string `json:"goPackages,omitempty"`
}

// ProtoMessage is a more human-readable representation of raw message
// descriptors
type ProtobufMessage struct {
	Name   string            `json:"name"`
	Fields map[string]string `json:"fields"`
}

// Descriptors returns the human-readable representation of the raw
// file descriptors for this object
func (p *Protobuf) Descriptors() (*ProtobufDescriptors, error) {
	descriptors, err := p.GetDescriptors()
	if err != nil {
		return nil, err
	}
	//defer os.RemoveAll(tempPath)
	out := &ProtobufDescriptors{
		Messages:     make([]*ProtobufMessage, 0),
		SourceFiles:  make([]string, 0),
		JavaPackages: make([]*string, 0),
		GoPackages:   make([]*string, 0),
	}
	for _, x := range descriptors {
		out = appendPkgsFromDescriptor(out, x)
		out.SourceFiles = append(out.SourceFiles, strings.TrimPrefix(x.GetName(), "/"))
		for _, msg := range x.GetMessageTypes() {
			out.Messages = append(out.Messages, protoMessageFromDescriptor(msg))
		}
	}
	return out, nil
}

// GetDescriptors returns the raw file descriptors for a protobuf object
// This is primarily a helper for functions that return more human readable
// formats
func (p *Protobuf) GetDescriptors() (map[string]*desc.FileDescriptor, error) {
	descriptorSet := new(descriptor.FileDescriptorSet)
	if err := descriptorSet.XXX_Unmarshal(p.DescriptorBytes()); err != nil {
		return nil, err
	}
	return desc.CreateFileDescriptorsFromSet(descriptorSet)
}
