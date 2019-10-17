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
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/tinyzimmer/protobuf-registry/pkg/remotecache"
)

func accessor(filename string) (io.ReadCloser, error) {
	log.Info(filename)
	return os.Open(filename)
}

func (p *Protobuf) CompileToDescriptorSet() error {
	parser := &protoparse.Parser{
		ImportPaths:           make([]string, 0),
		InferImportPaths:      true,
		IncludeSourceCodeInfo: true,
		Accessor:              protoparse.FileAccessor(accessor),
	}
	tempPath, _, tempFiles, err := p.newTempFilesFromRaw(false)
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempPath)
	parser.ImportPaths = append(parser.ImportPaths, tempPath)
	// always import common protos
	dep, err := remotecache.Cache().GetGitDependency(remotecache.APICommonProtos, "", remotecache.DefaultBranch)
	if err != nil {
		return err
	}
	parser.ImportPaths = append(parser.ImportPaths, dep.Dir())
	if len(p.Dependencies) > 0 {
		for _, remoteDep := range p.Dependencies {
			dep, err := remotecache.Cache().GetGitDependency(remoteDep.URL, remoteDep.Path, remoteDep.Revision)
			if err != nil {
				return err
			}
			parser.ImportPaths = append(parser.ImportPaths, dep.Dir())
		}
	}
	resolved, err := protoparse.ResolveFilenames(parser.ImportPaths, tempFilesToStrings(tempFiles, "")...)
	if err != nil {
		return err
	}
	fds, err := parser.ParseFiles(resolved...)
	if err != nil {
		return err
	}
	set := desc.ToFileDescriptorSet(fds...)
	p.descriptor, err = proto.Marshal(set)
	return err
}

func tempFilesToStrings(in map[string][]os.FileInfo, trimPrefix string) []string {
	out := make([]string, 0)
	for dir, files := range in {
		for _, file := range files {
			if !file.IsDir() {
				fpath := filepath.Join(dir, file.Name())
				if trimPrefix != "" {
					fpath = strings.TrimPrefix(fpath, trimPrefix)
				}
				out = append(out, fpath)
			}
		}
	}
	return out
}
