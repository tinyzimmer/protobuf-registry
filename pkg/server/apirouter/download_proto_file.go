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
	"fmt"
	"net/http"
	"strings"

	"github.com/tinyzimmer/protobuf-registry/pkg/protobuf"
	"github.com/tinyzimmer/protobuf-registry/pkg/server/common"
)

func getFileVars(r *http.Request, splitStr string) (name, version, filename string) {
	name = common.GetName(r)
	version = common.GetVersion(r)
	filename = strings.Replace(r.URL.Path, fmt.Sprintf("api/proto/%s/%s/%s/", name, version, splitStr), "", 1)
	return name, version, filename
}

func (api *apiServer) getRawProtoFile(w http.ResponseWriter, r *http.Request) {
	var err error
	name, version, filename := getFileVars(r, "raw")
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
	out, err := proto.Contents(filename)
	if err != nil {
		common.BadRequest(err, w)
		return
	}
	common.WriteJSONResponse(map[string]string{
		"content": string(out),
	}, w)
}

func (api *apiServer) getMetaForProtoFile(w http.ResponseWriter, r *http.Request) {
	var err error
	name, version, filename := getFileVars(r, "meta")
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
	out, err := proto.DocJSON(filename)
	if err != nil {
		common.BadRequest(err, w)
		return
	}

	common.WriteRawResponse(out, w)
}
