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
	"net/http"

	dbcommon "github.com/tinyzimmer/protobuf-registry/pkg/database/common"
	"github.com/tinyzimmer/protobuf-registry/pkg/protobuf"
	"github.com/tinyzimmer/protobuf-registry/pkg/server/common"
	"github.com/tinyzimmer/protobuf-registry/pkg/types"
)

func (api *apiServer) getAllProtoHandler(w http.ResponseWriter, r *http.Request) {
	var protos map[string][]*protobuf.Protobuf
	var err error

	log.Info("Fetching list of all protobuf packages...")
	if protos, err = api.DB().GetAllProtoVersions(); err != nil {
		common.BadRequest(err, w)
		return
	}
	out := &types.ListProtoResponse{
		Items: make([]*types.ProtoMeta, 0),
	}
	for name, protos := range protos {
		// safety check that this isn't an empty slice
		if len(protos) == 0 {
			continue
		}
		// sort the protobuf versions
		protos = common.SortVersions(protos)
		out.Items = append(out.Items, &types.ProtoMeta{
			Name:           name,
			Versions:       protos,
			Latest:         *protos[0].Version,
			LatestUploaded: protos[0].LastUpdated,
		})
	}
	common.WriteJSONResponse(out.Items, w)
}

func (api *apiServer) getProtoHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	name := common.GetName(r)
	var protos []*protobuf.Protobuf

	log.Info("Fetching versions for protobuf package", "name", name)
	if protos, err = api.DB().GetProtoVersions(name); err != nil {
		if dbcommon.IsProtobufNotExists(err) {
			common.NotFound(err, w)
			return
		}
		common.BadRequest(err, w)
		return
	}
	protos = common.SortVersions(protos)
	common.WriteJSONResponse(protos, w)
}

func (api *apiServer) getProtoVersionMetaHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	name := common.GetName(r)
	version := common.GetVersion(r)
	var protos []*protobuf.Protobuf

	log.Info("Fetching meta for version of proto package", "name", name, "version", version)
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
		common.NotFound(err, w)
		return
	}
	log.Info("Fetching descriptor sets from storage")
	if proto, err = api.Storage().GetRawProto(proto); err != nil {
		common.BadRequest(err, w)
		return
	}
	log.Info("Gathering protobuf descriptors")
	out, err := proto.Descriptors()
	if err != nil {
		common.BadRequest(err, w)
		return
	}
	common.WriteJSONResponse(out, w)
}
