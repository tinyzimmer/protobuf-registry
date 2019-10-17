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
	"fmt"

	"github.com/tinyzimmer/protobuf-registry/pkg/protobuf"
	"github.com/tinyzimmer/protobuf-registry/pkg/types"
)

func (r *registryClient) GetProtoPackageVersions(name string) ([]*Protobuf, error) {
	return r.GetProtoPackageVersionsWithContext(context.Background(), name)
}

func (r *registryClient) GetProtoPackage(name, version string) (*protobuf.ProtobufDescriptors, error) {
	return r.GetProtoPackageWithContext(context.Background(), name, version)
}

func (r *registryClient) GetFileContents(pkgName, pkgVersion, filename string) (*types.GetFileContentsResponse, error) {
	return r.GetFileContentsWithContext(context.Background(), pkgName, pkgVersion, filename)
}

func (r *registryClient) GetFileDocs(pkgName, pkgVersion, filename string) (map[string]interface{}, error) {
	return r.GetFileDocsWithContext(context.Background(), pkgName, pkgVersion, filename)
}

func (r *registryClient) GetProtoPackageVersionsWithContext(ctx context.Context, name string) ([]*Protobuf, error) {
	path := fmt.Sprintf("api/proto/%s", name)
	req, err := r.newRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	var res []*Protobuf
	return res, r.doInto(req, &res)
}

func (r *registryClient) GetProtoPackageWithContext(ctx context.Context, name, version string) (*protobuf.ProtobufDescriptors, error) {
	path := fmt.Sprintf("api/proto/%s/%s", name, version)
	req, err := r.newRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	var res protobuf.ProtobufDescriptors
	return &res, r.doInto(req, &res)
}

func (r *registryClient) GetFileContentsWithContext(ctx context.Context, pkgName, pkgVersion, filename string) (*types.GetFileContentsResponse, error) {
	path := fmt.Sprintf("api/proto/%s/%s/raw/%s", pkgName, pkgVersion, filename)
	req, err := r.newRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	var res types.GetFileContentsResponse
	return &res, r.doInto(req, &res)
}

func (r *registryClient) GetFileDocsWithContext(ctx context.Context, pkgName, pkgVersion, filename string) (map[string]interface{}, error) {
	path := fmt.Sprintf("api/proto/%s/%s/meta/%s", pkgName, pkgVersion, filename)
	req, err := r.newRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	var res map[string]interface{}
	return res, r.doInto(req, &res)
}
