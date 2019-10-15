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
	"net/http"

	"github.com/tinyzimmer/protobuf-registry/pkg/config"
	"github.com/tinyzimmer/protobuf-registry/pkg/protobuf"
	"github.com/tinyzimmer/protobuf-registry/pkg/types"
)

// RegistryClient is the exported interface for client-side protobuf-registry
// requests and operations.
type RegistryClient interface {
	// GetServerConfig returns the configuration of the server
	GetServerConfig() (*config.Config, error)
	UploadProtoPackage(*types.PostProtoRequest, bool) (*Protobuf, error)
	UploadProtoPackageFromDir(dir string, force bool) (*Protobuf, error)
	GetProtoPackageVersions(name string) ([]*Protobuf, error)
	DeleteAllProtoPackageVersions(name string) error
	GetProtoPackage(name, version string) (*protobuf.ProtobufDescriptors, error)
	DeleteProtoPackage(name, version string) error
	DownloadProtoPackage(name, version, format string) ([]byte, error)
	GetFileContents(pkgName, pkgVersion, filename string) ([]byte, error)
	GetFileDocs(pkgName, pkgVersion, filename string) (map[string]interface{}, error)
	GetCachedRemotes() ([]string, error)
	PutCachedRemote(remote string) error
}

// Protobuf struct is a client-side implementation of the protobuf objects -
// Methods in the server interface will be replaced with client-side API calls
type Protobuf struct {
	ID      *string
	Name    *string
	Version *string

	client *registryClient
}

// registryClient implements RegistryClient
type registryClient struct {
	Server string

	httpclient *http.Client
}

// New returns a new RegistryClient using the provided URL
func New(url string) RegistryClient {
	return &registryClient{Server: url}
}

// TODO //

func (r *registryClient) GetServerConfig() (*config.Config, error) { return nil, nil }
func (r *registryClient) UploadProtoPackage(*types.PostProtoRequest, bool) (*Protobuf, error) {
	return nil, nil
}
func (r *registryClient) UploadProtoPackageFromDir(dir string, force bool) (*Protobuf, error) {
	return nil, nil
}
func (r *registryClient) GetProtoPackageVersions(name string) ([]*Protobuf, error) { return nil, nil }
func (r *registryClient) DeleteAllProtoPackageVersions(name string) error          { return nil }
func (r *registryClient) GetProtoPackage(name, version string) (*protobuf.ProtobufDescriptors, error) {
	return nil, nil
}
func (r *registryClient) DeleteProtoPackage(name, version string) error { return nil }
func (r *registryClient) DownloadProtoPackage(name, version, format string) ([]byte, error) {
	return nil, nil
}
func (r *registryClient) GetFileContents(pkgName, pkgVersion, filename string) ([]byte, error) {
	return nil, nil
}
func (r *registryClient) GetFileDocs(pkgName, pkgVersion, filename string) (map[string]interface{}, error) {
	return nil, nil
}
func (r *registryClient) GetCachedRemotes() ([]string, error) { return nil, nil }
func (r *registryClient) PutCachedRemote(remote string) error { return nil }
