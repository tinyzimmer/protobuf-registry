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
	"bytes"
	"fmt"
	"net/http"

	"github.com/tinyzimmer/protobuf-registry/pkg/protobuf"
	"github.com/tinyzimmer/protobuf-registry/pkg/server/common"
	"github.com/tinyzimmer/protobuf-registry/pkg/util"
)

func (pip *pipServer) getPipVersionsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	pkg := common.GetName(r)
	if pkg == "pip" {
		w.WriteHeader(http.StatusOK)
		return
	}

	protos, err := pip.DB().GetProtoVersions(pkg)
	if err != nil {
		common.BadRequest(err, w)
		return
	}

	downloadURL := fmt.Sprintf("%s%s/pip/download", r.URL.Scheme, r.URL.Host)
	if err := getPipTemplate.Execute(w, pip.protosToTemplateOpts(pkg, downloadURL, protos)); err != nil {
		common.BadRequest(err, w)
	}
}

func (pip *pipServer) getPipDownloadHandler(w http.ResponseWriter, r *http.Request) {
	filename := common.GetName(r)
	pkg, version := common.ParseNameVersionExtString(filename, ".tar.gz")
	protos, err := pip.DB().GetProtoVersions(pkg)
	if err != nil {
		common.BadRequest(err, w)
		return
	}
	var proto *protobuf.Protobuf
	if proto, err = common.GetVersionFromProtoSlice(protos, version); err != nil {
		common.BadRequest(err, w)
		return
	}
	if proto, err = pip.Storage().GetRawProto(proto); err != nil {
		common.BadRequest(err, w)
		return
	}
	path, rm, err := proto.GenerateTo(protobuf.GenerateTargetPython, *proto.Name)
	if err != nil {
		common.BadRequest(err, w)
		return
	}
	defer rm()
	if err := addSetupPy(proto, path); err != nil {
		common.BadRequest(err, w)
	}
	archive, err := util.NewTarGZArchive(path)
	if err != nil {
		common.BadRequest(err, w)
		return
	}
	rdr := bytes.NewReader(archive)
	common.ServeFile(w, r, filename, rdr)
}
