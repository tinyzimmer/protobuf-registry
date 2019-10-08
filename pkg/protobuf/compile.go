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
	"time"

	"github.com/tinyzimmer/proto-registry/pkg/config"
)

func (p *Protobuf) CompileDescriptorSet() ([]byte, error) {
	tempPath, tempFiles, remove, err := p.newTempFilesFromRaw()
	if err != nil {
		return nil, err
	}
	defer remove()
	tempOut, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tempOut)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.GlobalConfig.CompileTimeout)*time.Second)
	defer cancel()
	args := []string{
		fmt.Sprintf("-I=%s", tempPath),
		"--include_imports",
		fmt.Sprintf("--descriptor_set_out=%s", filepath.Join(tempOut, "descriptor.pb")),
	}
	args = append(args, tempFilesToStrings(tempFiles)...)
	out, err := exec.CommandContext(ctx,
		config.GlobalConfig.ProtocPath,
		args...,
	).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to compile protocol spec: %s", string(out))
	}
	return ioutil.ReadFile(filepath.Join(tempOut, "descriptor.pb"))
}

type CompileTarget int

const (
	_ CompileTarget = iota
	CompileTargetCPP
	CompileTargetCSharp
	CompileTargetJava
	CompileTargetJavaNano
	CompileTargetJS
	CompileTargetObjC
	CompileTargetPHP
	CompileTargetPython
	CompileTargetRuby
)

func getTargetArg(target CompileTarget) string {
	switch target {
	case CompileTargetCPP:
		return "--cpp_out"
	case CompileTargetCSharp:
		return "--csharp_out"
	case CompileTargetJava:
		return "--java_out"
	case CompileTargetJavaNano:
		return "--javanano_out"
	case CompileTargetJS:
		return "--js_out"
	case CompileTargetObjC:
		return "--objc_out"
	case CompileTargetPHP:
		return "--php_out"
	case CompileTargetPython:
		return "--python_out"
	case CompileTargetRuby:
		return "--ruby_out"
	default:
		return ""
	}
}

func (p Protobuf) CompileTo(target CompileTarget, prefix string) (tempOut string, rm func(), err error) {
	rawPath, tempFiles, remove, err := p.newTempFilesFromRaw()
	if err != nil {
		return "", nil, err
	}
	defer remove()

	tempOut, err = ioutil.TempDir("", "")
	if err != nil {
		return "", nil, err
	}
	rm = func() {
		if err := os.RemoveAll(tempOut); err != nil {
			log.Error(err, "Failed to remove tempdir")
		}
	}

	var out string
	if prefix != "" {
		out = filepath.Join(tempOut, prefix)
		if err = os.MkdirAll(out, 0700); err != nil {
			rm()
			return "", nil, err
		}
	} else {
		out = tempOut
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.GlobalConfig.CompileTimeout)*time.Second)
	defer cancel()
	args := []string{
		fmt.Sprintf("-I=%s", rawPath),
		fmt.Sprintf("%s=%s", getTargetArg(target), out),
	}
	args = append(args, tempFilesToStrings(tempFiles)...)

	cmdout, err := exec.CommandContext(ctx,
		config.GlobalConfig.ProtocPath,
		args...,
	).CombinedOutput()
	if err != nil {
		rm()
		return "", nil, fmt.Errorf("failed to compile protocol spec: %s", string(cmdout))
	}
	return tempOut, rm, nil
}

func tempFilesToStrings(in map[string][]os.FileInfo) []string {
	out := make([]string, 0)
	for dir, files := range in {
		for _, file := range files {
			if !file.IsDir() {
				out = append(out, filepath.Join(dir, file.Name()))
			}
		}
	}
	return out
}
