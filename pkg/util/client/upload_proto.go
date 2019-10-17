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

package client

import (
	"context"

	"github.com/tinyzimmer/protobuf-registry/pkg/types"
)

// UploadProtoPackage uploads a new protocol buffer spec with the given
// request parameters.
func (r *registryClient) UploadProtoPackage(req *types.PostProtoRequest, overwrite bool) (*Protobuf, error) {
	return r.UploadProtoPackageWithContext(context.Background(), req, overwrite)
}

func (r *registryClient) UploadProtoPackageWithContext(ctx context.Context, req *types.PostProtoRequest, overwrite bool) (*Protobuf, error) {
	var method string
	if overwrite {
		method = "PUT"
	} else {
		method = "POST"
	}
	hreq, err := r.newRequest(ctx, method, "api/proto", req)
	if err != nil {
		return nil, err
	}
	var proto Protobuf
	return &proto, r.doInto(hreq, &proto)
}
