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
	"fmt"
	"html/template"

	dbcommon "github.com/tinyzimmer/protobuf-registry/pkg/database/common"
	"github.com/tinyzimmer/protobuf-registry/pkg/protobuf"
	storagecommon "github.com/tinyzimmer/protobuf-registry/pkg/storage/common"
)

var funcMap = template.FuncMap{
	"url": func(s string) template.URL {
		return template.URL(s)
	},
}

// TODO: persistent storage of compiled artifacts and use sha256
// "{{ $.Host|url }}/{{ .Filename }}#sha256={{ .SHA256 }}"
var getPipTemplateString = `
<html>
  <head>
    <title>Links for {{ .PackageName }}</title>
  </head>
  <body>
    <h1>Links for {{ .PackageName }}</h1>
    {{- range .Packages }}
    <a href="{{ $.Host|url }}/{{ .Filename }}">{{ .Filename }}</a><br/>
    {{- end }}
  </body>
</html>
`
var getPipTemplate = template.Must(template.New("get_pip").Funcs(funcMap).Parse(getPipTemplateString))

var filenameTemplate = "%s-%s.tar.gz"

func (pip *pipServer) DB() dbcommon.DBEngine {
	return pip.ctrl.DB()
}

func (pip *pipServer) Storage() storagecommon.Provider {
	return pip.ctrl.Storage()
}

type PipTemplateOptions struct {
	Host        string
	PackageName string
	Packages    []*PipPackage
}

type PipPackage struct {
	//SHA256   string
	Filename string
}

func (pip *pipServer) protosToTemplateOpts(pkg string, downloadURL string, protos []*protobuf.Protobuf) *PipTemplateOptions {
	opts := &PipTemplateOptions{
		Host:        downloadURL,
		PackageName: pkg,
		Packages:    make([]*PipPackage, 0),
	}

	for _, proto := range protos {
		//var err error
		// if proto, err = pip.Storage().GetRawProto(proto); err != nil {
		// 	return nil, err
		// }
		// sha256, err := proto.SHA256()
		// if err != nil {
		// 	return nil, err
		// }
		pipPkg := &PipPackage{
			// SHA256:   sha256,
			Filename: fmt.Sprintf(filenameTemplate, *proto.Name, *proto.Version),
		}
		opts.Packages = append(opts.Packages, pipPkg)
	}
	return opts
}
