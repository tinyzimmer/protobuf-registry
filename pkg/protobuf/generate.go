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

	"github.com/tinyzimmer/protobuf-registry/pkg/config"
)

type GenerateTarget int

const (
	_ GenerateTarget = iota
	GenerateTargetCPP
	GenerateTargetCSharp
	GenerateTargetJava
	GenerateTargetJavaNano
	GenerateTargetJS
	GenerateTargetObjC
	GenerateTargetPHP
	GenerateTargetPython
	GenerateTargetRuby
	GenerateTargetGo
)

func getTargetArg(target GenerateTarget) string {
	switch target {
	case GenerateTargetCPP:
		return "--cpp_out"
	case GenerateTargetCSharp:
		return "--csharp_out"
	case GenerateTargetJava:
		return "--java_out"
	case GenerateTargetJavaNano:
		return "--javanano_out"
	case GenerateTargetJS:
		return "--js_out"
	case GenerateTargetObjC:
		return "--objc_out"
	case GenerateTargetPHP:
		return "--php_out"
	case GenerateTargetPython:
		return "--python_out"
	case GenerateTargetRuby:
		return "--ruby_out"
	case GenerateTargetGo:
		return "--go_out"
	default:
		return ""
	}
}

func (p Protobuf) GenerateTo(target GenerateTarget, prefix string) (tempOut string, rm func(), err error) {

	rawPath, descriptorSet, tempFiles, err := p.newTempFilesFromRaw(true)
	if err != nil {
		return "", nil, err
	}
	defer os.RemoveAll(rawPath)

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
		fmt.Sprintf("--descriptor_set_in=%s", descriptorSet),
		fmt.Sprintf("%s=%s", getTargetArg(target), out),
	}
	if target == GenerateTargetGo {
		args = append(args, fmt.Sprintf("--plugin=%s", config.GlobalConfig.ProtocGenGoPath))
	}
	args = append(args, tempFilesToStrings(tempFiles, rawPath+"/")...)

	cmdout, err := exec.CommandContext(ctx,
		config.GlobalConfig.ProtocPath,
		args...,
	).CombinedOutput()
	if err != nil {
		rm()
		return "", nil, fmt.Errorf("failed to generate code for protocol spec: %s", string(cmdout))
	}
	return tempOut, rm, nil
}
