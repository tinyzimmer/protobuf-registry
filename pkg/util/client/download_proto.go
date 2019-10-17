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

// TODO: Validate format with switch
func (r *registryClient) DownloadProtoPackage(name, version, format string) ([]byte, error) {
	return r.DownloadProtoPackageWithContext(context.Background(), name, version, format)
}

func (r *registryClient) DownloadProtoPackageWithContext(ctx context.Context, name, version, format string) ([]byte, error) {
	path := fmt.Sprintf("api/proto/%s/%s/%s", name, version, format)
	req, err := r.newRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	return r.doRaw(req)
}
