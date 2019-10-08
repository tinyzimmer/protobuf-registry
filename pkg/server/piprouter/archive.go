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

package piprouter

import (
	"html/template"
	"os"
	"path/filepath"

	"github.com/tinyzimmer/proto-registry/pkg/protobuf"
)

var setupPyTemplateString = `
"""Distribution and setup script for {{ .Name }}"""

from setuptools import setup

version = "{{ .Version }}"

setup(
    name="{{ .Name }}",
    version=version,
    install_requires=[
        "protobuf"
    ],
    packages=[
        "{{ .Name }}",
    ]
)

`
var setupPyTemplate = template.Must(template.New("setup.py").Funcs(funcMap).Parse(setupPyTemplateString))

func addSetupPy(proto *protobuf.Protobuf, path string) error {
	file, err := os.Create(filepath.Join(path, "setup.py"))
	if err != nil {
		return err
	}
	defer file.Close()
	return setupPyTemplate.Execute(file, proto)
}
