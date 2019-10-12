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
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/tinyzimmer/protobuf-registry/pkg/config"
	"github.com/tinyzimmer/protobuf-registry/pkg/remotecache"
)

func (p *Protobuf) CompileDescriptorSet() error {
	var importPaths []string
	if len(p.Dependencies) > 0 {
		for _, remoteDep := range p.Dependencies {
			dep, err := remotecache.Cache().GetGitDependency(remoteDep)
			if err != nil {
				return err
			}
			importPaths = append(importPaths, dep.Dir())
		}
	}
	tempPath, _, tempFiles, err := p.newTempFilesFromRaw(false)
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempPath)
	importPaths = append(importPaths, tempPath)
	tempOut, err := ioutil.TempDir("", "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempOut)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.GlobalConfig.CompileTimeout)*time.Second)
	defer cancel()

	args := make([]string, 0)
	for _, x := range importPaths {
		args = append(args, fmt.Sprintf("-I=%s", x))
	}
	args = append(args, []string{
		"--include_imports",
		"--include_source_info",
		fmt.Sprintf("--descriptor_set_out=%s", filepath.Join(tempOut, "descriptor.pb")),
	}...)
	args = append(args, tempFilesToStrings(tempFiles, "")...)
	out, err := exec.CommandContext(ctx,
		config.GlobalConfig.ProtocPath,
		args...,
	).CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to compile protocol spec: %s", string(out))
	}
	p.descriptor, err = ioutil.ReadFile(filepath.Join(tempOut, "descriptor.pb"))
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
