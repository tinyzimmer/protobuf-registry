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

func (r *registryClient) GetCachedRemotes() ([]string, error) {
	return r.GetCachedRemotesWithContext(context.Background())
}

func (r *registryClient) PutCachedRemote(req *types.PutRemoteRequest) error {
	return r.PutCachedRemoteWithContext(context.Background(), req)
}

func (r *registryClient) GetCachedRemotesWithContext(ctx context.Context) ([]string, error) {
	req, err := r.newRequest(ctx, "GET", "api/remotes", nil)
	if err != nil {
		return nil, err
	}
	var res []string
	return res, r.doInto(req, &res)
}

func (r *registryClient) PutCachedRemoteWithContext(ctx context.Context, req *types.PutRemoteRequest) error {
	hreq, err := r.newRequest(ctx, "PUT", "api/remotes", req)
	if err != nil {
		return err
	}
	_, err = r.doRaw(hreq)
	return err
}
