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
	"github.com/tinyzimmer/protobuf-registry/pkg/protobuf"
	"github.com/tinyzimmer/protobuf-registry/pkg/server/common"
	"github.com/tinyzimmer/protobuf-registry/pkg/util"
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

	log.Info("Serving download request for proto package", "name", name, "version", version, "language", language)
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
	} else if language == "descriptors" {
		common.ServeFile(w, r, fmt.Sprintf("%s-%s-descriptors.pb", *proto.Name, *proto.Version), proto.DescriptorReader())
		return
	}

	target, err := getGenerateTarget(language)
	if err != nil {
		common.BadRequest(err, w)
		return
	}

	prefix := *proto.Name
	if target == protobuf.GenerateTargetGo {
		prefix = ""
	}

	out, rm, err := proto.GenerateTo(target, prefix)
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

func getGenerateTarget(target string) (protobuf.GenerateTarget, error) {
	var out protobuf.GenerateTarget
	var err error
	switch target {
	case "cpp":
		out = protobuf.GenerateTargetCPP
	case "csharp":
		out = protobuf.GenerateTargetCSharp
	case "java":
		out = protobuf.GenerateTargetJava
	case "javanano":
		out = protobuf.GenerateTargetJavaNano
	case "js":
		out = protobuf.GenerateTargetJS
	case "objc":
		out = protobuf.GenerateTargetObjC
	case "php":
		out = protobuf.GenerateTargetPHP
	case "python":
		out = protobuf.GenerateTargetPython
	case "ruby":
		out = protobuf.GenerateTargetRuby
	case "go":
		out = protobuf.GenerateTargetGo
	default:
		err = errors.New("Unknown target")
	}
	return out, err
}
