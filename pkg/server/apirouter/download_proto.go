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

package apirouter

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tinyzimmer/proto-registry/pkg/protobuf"
	"github.com/tinyzimmer/proto-registry/pkg/server/common"
	"github.com/tinyzimmer/proto-registry/pkg/util"
)

func getLanguage(r *http.Request) string {
	vars := mux.Vars(r)
	return vars["language"]
}

func (api *apiServer) downloadProtoHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	name := common.GetName(r)
	version := common.GetVersion(r)
	language := getLanguage(r)
	var protos []*protobuf.Protobuf
	if protos, err = api.DB().GetProtoVersions(name); err != nil {
		common.BadRequest(err, w)
		return
	}
	var proto *protobuf.Protobuf
	if proto, err = common.GetVersionFromProtoSlice(protos, version); err != nil {
		common.BadRequest(err, w)
		return
	}
	if proto, err = api.Storage().GetRawProto(proto); err != nil {
		common.BadRequest(err, w)
		return
	}

	if language == "raw" {
		common.ServeFile(w, r, proto.RawFilename(), proto.RawReader())
		return
	}

	target, err := getCompileTarget(language)
	if err != nil {
		common.BadRequest(err, w)
		return
	}

	out, rm, err := proto.CompileTo(target, *proto.Name)
	if err != nil {
		common.BadRequest(err, w)
		return
	}
	defer rm()

	filename := fmt.Sprintf("%s-%s-%s.tar.gz", *proto.Name, language, *proto.Version)
	archive, err := util.NewTarGZArchive(out)
	if err != nil {
		common.BadRequest(err, w)
		return
	}

	common.ServeFile(w, r, filename, bytes.NewReader(archive))
}

func getCompileTarget(target string) (protobuf.CompileTarget, error) {
	var out protobuf.CompileTarget
	var err error
	switch target {
	case "cpp":
		out = protobuf.CompileTargetCPP
	case "csharp":
		out = protobuf.CompileTargetCSharp
	case "java":
		out = protobuf.CompileTargetJava
	case "javanano":
		out = protobuf.CompileTargetJavaNano
	case "js":
		out = protobuf.CompileTargetJS
	case "objc":
		out = protobuf.CompileTargetObjC
	case "php":
		out = protobuf.CompileTargetPHP
	case "python":
		out = protobuf.CompileTargetPython
	case "ruby":
		out = protobuf.CompileTargetRuby
	default:
		err = errors.New("Unknown target")
	}
	return out, err
}
