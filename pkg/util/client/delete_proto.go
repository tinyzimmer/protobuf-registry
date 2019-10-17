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
)

func (r *registryClient) DeleteAllProtoPackageVersions(name string) error {
	return r.DeleteAllProtoPackageVersionsWithContext(context.Background(), name)
}

func (r *registryClient) DeleteAllProtoPackageVersionsWithContext(ctx context.Context, name string) error {
	path := fmt.Sprintf("api/proto/%s", name)
	req, err := r.newRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return err
	}
	_, err = r.doRaw(req)
	return err
}

func (r *registryClient) DeleteProtoPackage(name, version string) error {
	return r.DeleteProtoPackageWithContext(context.Background(), name, version)
}

func (r *registryClient) DeleteProtoPackageWithContext(ctx context.Context, name, version string) error {
	path := fmt.Sprintf("api/proto/%s/%s", name, version)
	req, err := r.newRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return err
	}
	_, err = r.doRaw(req)
	return err
}
