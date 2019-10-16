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
	"net/http"

	"github.com/tinyzimmer/protobuf-registry/pkg/config"
	"github.com/tinyzimmer/protobuf-registry/pkg/protobuf"
	"github.com/tinyzimmer/protobuf-registry/pkg/server/common"
	"github.com/tinyzimmer/protobuf-registry/pkg/types"
)

func (api *apiServer) processReq(req *types.PostProtoRequest, force bool) (proto *protobuf.Protobuf, err error) {
	log.Info("Validating parameters...")
	// validate parameters
	if err = req.Validate(); err != nil {
		return
	}

	// Create a protobuf object from the request
	log.Info("Validating protocol specification...")
	proto = protobuf.NewFromRequest(req)
	if err = proto.SetRawFromBase64(req.Body); err != nil {
		return
	}

	// compile and set the descirptor set
	if err = proto.CompileDescriptorSet(); err != nil {
		return
	}

	log.Info("Registering new package", "proto", proto)
	// register to DB - object comes back with a generated ID if it doesn't exist.
	// This will return an error if an object with the same name and version exists
	// and force is false. If force is true, the object comes back with the ID of
	// the existing one which will cause it to be overwritten by the following
	// Storage() call.
	if proto, err = api.DB().StoreProtoVersion(proto, force); err != nil {
		return
	}

	// write raw proto to storage backend
	if err = api.Storage().StoreProtoPackage(proto); err != nil {
		return
	}

	return
}

func (api *apiServer) putProtoHandler(w http.ResponseWriter, r *http.Request) {

	var req *types.PostProtoRequest
	var err error

	log.Info("Unmarshaling new protobuf request...")
	// unmarshall the request
	if req, err = types.NewProtoReqFromReader(r.Body); err != nil {
		common.BadRequest(err, w)
		return
	}

	if proto, err := api.processReq(req, true); err != nil {
		common.BadRequest(err, w)
	} else {
		common.WriteJSONResponse(proto, w)
	}
}

func (api *apiServer) postProtoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		if config.GlobalConfig.CORSEnabled {
			w.Header().Set("Access-Control-Allow-Methods", "POST, PUT")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.WriteHeader(http.StatusOK)
			return
		}
		common.BadRequest(errors.New("CORS is not enabled, rejecting OPTIONS request"), w)
		return
	}

	var req *types.PostProtoRequest
	var err error

	log.Info("Unmarshaling new protobuf request...")
	// unmarshall the request
	if req, err = types.NewProtoReqFromReader(r.Body); err != nil {
		common.BadRequest(err, w)
		return
	}

	if proto, err := api.processReq(req, false); err != nil {
		common.BadRequest(err, w)
	} else {
		common.WriteJSONResponse(proto, w)
	}
}
