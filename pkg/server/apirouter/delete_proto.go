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
	"errors"
	"fmt"
	"net/http"

	"github.com/tinyzimmer/protobuf-registry/pkg/config"
	dbcommon "github.com/tinyzimmer/protobuf-registry/pkg/database/common"
	"github.com/tinyzimmer/protobuf-registry/pkg/protobuf"
	"github.com/tinyzimmer/protobuf-registry/pkg/server/common"
)

func (api *apiServer) deleteAllProtoVersionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		if config.GlobalConfig.CORSEnabled {
			w.Header().Set("Access-Control-Allow-Methods", "DELETE")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.WriteHeader(http.StatusOK)
			return
		}
		common.BadRequest(errors.New("CORS is not enabled, rejecting OPTIONS request"), w)
		return
	}
	var err error
	name := common.GetName(r)
	var protos []*protobuf.Protobuf
	if protos, err = api.DB().GetProtoVersions(name); err != nil {
		if dbcommon.IsProtobufNotExists(err) {
			common.NotFound(err, w)
			return
		}
		common.BadRequest(err, w)
		return
	}
	if err := api.DB().RemoveAllVersionsForProto(name); err != nil {
		common.BadRequest(err, w)
		return
	}
	for _, proto := range protos {
		if err := api.Storage().DeleteProtoPackage(proto); err != nil {
			common.BadRequest(err, w)
			return
		}
	}
	common.WriteJSONResponse(map[string]string{
		"result": fmt.Sprintf("All versions for %s deleted", name),
	}, w)
}

func (api *apiServer) deleteProtoVersionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		if config.GlobalConfig.CORSEnabled {
			w.Header().Set("Access-Control-Allow-Methods", "DELETE")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.WriteHeader(http.StatusOK)
			return
		}
		common.BadRequest(errors.New("CORS is not enabled, rejecting OPTIONS request"), w)
		return
	}
	var err error
	name := common.GetName(r)
	version := common.GetVersion(r)
	var protos []*protobuf.Protobuf
	if protos, err = api.DB().GetProtoVersions(name); err != nil {
		if dbcommon.IsProtobufNotExists(err) {
			common.NotFound(err, w)
			return
		}
		common.BadRequest(err, w)
		return
	}
	var proto *protobuf.Protobuf
	if proto, err = common.GetVersionFromProtoSlice(protos, version); err != nil {
		common.BadRequest(err, w)
		return
	}
	if err := api.DB().RemoveProtoVersion(proto); err != nil {
		common.BadRequest(err, w)
		return
	}
	if err := api.Storage().DeleteProtoPackage(proto); err != nil {
		common.BadRequest(err, w)
		return
	}
	common.WriteJSONResponse(map[string]string{
		"result": fmt.Sprintf("Deleted %s version %s", name, version),
	}, w)
}
