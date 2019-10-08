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
	"io"
	"net/http"

	"github.com/tinyzimmer/proto-registry/pkg/config"
	"github.com/tinyzimmer/proto-registry/pkg/protobuf"
	"github.com/tinyzimmer/proto-registry/pkg/server/common"
	"github.com/tinyzimmer/proto-registry/pkg/util"
)

func validatePost(req *PostProtoRequest) (*PostProtoRequest, error) {
	if req.Name == nil || req.Body == nil {
		return req, errors.New("'name' and 'body' are required")
	} else if req.Version == nil {
		req.Version = util.StringPtr("0.0.1")
	}
	return req, nil
}

type PostProtoRequest struct {
	ID      *string `json:"id,omitempty"`
	Name    *string `json:"name,omitempty"`
	Body    *string `json:"body,omitempty"`
	Version *string `json:"version,omitempty"`
}

func unmarshalProtoRequest(rdr io.ReadCloser) (*PostProtoRequest, error) {
	defer rdr.Close()
	var req PostProtoRequest
	err := common.UnmarshallInto(rdr, &req)
	return &req, err
}

func (api *apiServer) postProtoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		if config.GlobalConfig.CORSEnabled {
			w.Header().Set("Access-Control-Allow-Methods", "POST")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.WriteHeader(http.StatusOK)
			return
		}
		common.BadRequest(errors.New("CORS is not enabled, rejecting OPTIONS request"), w)
		return
	}

	var req *PostProtoRequest
	var err error

	// unmarshall the request
	if req, err = unmarshalProtoRequest(r.Body); err != nil {
		common.BadRequest(err, w)
		return
	}

	// validate parameters
	if req, err = validatePost(req); err != nil {
		common.BadRequest(err, w)
		return
	}

	// Create a protobuf object from the request
	proto := protobuf.New(req.ID, req.Name, req.Version)
	if err := proto.SetRawFromBase64(req.Body); err != nil {
		common.BadRequest(err, w)
		return
	}

	// make sure it compiles
	if err := proto.Compile(); err != nil {
		common.BadRequest(err, w)
		return
	}

	// register to DB - object comes back with a generated ID if it doesn't exist
	// this will return an error if an object with the same name and version exists
	if proto, err = api.DB().StoreProtoVersion(proto); err != nil {
		common.BadRequest(err, w)
		return
	}

	// write raw proto to storage backend
	if err = api.Storage().StoreProtoPackage(proto); err != nil {
		common.BadRequest(err, w)
		return
	}

	common.WriteJSONResponse(proto, w)
}
