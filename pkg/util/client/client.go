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
	"net/http"
	"strings"

	"github.com/tinyzimmer/protobuf-registry/pkg/config"
	"github.com/tinyzimmer/protobuf-registry/pkg/protobuf"
	"github.com/tinyzimmer/protobuf-registry/pkg/types"
)

// RegistryClient is the exported interface for client-side protobuf-registry
// requests and operations.
type RegistryClient interface {

	// GetServerConfig returns the configuration of the server
	GetServerConfig() (*config.Config, error)
	GetServerConfigWithContext(context.Context) (*config.Config, error)

	// UploadProtoPackage uploads a new protocol buffer spec to the server.
	// If overwrite is true, it will replace an existing package with the same
	// name and version.
	UploadProtoPackage(req *types.PostProtoRequest, overwrite bool) (*Protobuf, error)
	UploadProtoPackageWithContext(ctx context.Context, req *types.PostProtoRequest, overwrite bool) (*Protobuf, error)

	// ListProtoPackages returns all the protocol packages on the server. The
	// protobuf object methods returned by this call are not safe to use. They are the
	// same interface the server uses and will later be seperated from this request.
	// To get a protobuf object with use-able methods, use GetProtoPackageVersions
	ListProtoPackages() (*types.ListProtoResponse, error)
	ListProtoPackagesWithContext(context.Context) (*types.ListProtoResponse, error)

	// GetProtoPackageVersions will retrieve a list of all the protocol buffer
	// packages with the provided name
	GetProtoPackageVersions(name string) ([]*Protobuf, error)
	GetProtoPackageVersionsWithContext(ctx context.Context, name string) ([]*Protobuf, error)

	// GetProtoPackage retrieves details about the package with the given name
	// and version
	GetProtoPackage(name, version string) (*protobuf.ProtobufDescriptors, error)
	GetProtoPackageWithContext(ctx context.Context, name, version string) (*protobuf.ProtobufDescriptors, error)

	// DeleteAllProtoPackageVersions deletes all versions for the given protocol
	// name
	DeleteAllProtoPackageVersions(name string) error
	DeleteAllProtoPackageVersionsWithContext(ctx context.Context, name string) error

	// DeleteProtoPackage deletes the package with the given name and version
	DeleteProtoPackage(name, version string) error
	DeleteProtoPackageWithContext(ctx context.Context, name, version string) error

	// DownloadProtoPackage retrieves the raw bytes for the package name and version
	// in the provided format.
	// TODO: Will use an enum for the format. Also note the following formats used:
	// - raw: .zip
	// - descriptors: marshaled proto descriptor set (equivalent of --descriptor_set_out)
	// - language (codegen): .tar.gz
	// Will likely make these more consistent
	DownloadProtoPackage(name, version, format string) ([]byte, error)
	DownloadProtoPackageWithContext(ctx context.Context, name, version, format string) ([]byte, error)

	// GetFileContents retrieves the contents of the file in the given package
	// name and version
	GetFileContents(pkgName, pkgVersion, filename string) (*types.GetFileContentsResponse, error)
	GetFileContentsWithContext(ctx context.Context, pkgName, pkgVersion, filename string) (*types.GetFileContentsResponse, error)

	// GetFileDocs returns documentation for the file in the given package name
	// and version
	GetFileDocs(pkgName, pkgVersion, filename string) (map[string]interface{}, error)
	GetFileDocsWithContext(ctx context.Context, pkgName, pkgVersion, filename string) (map[string]interface{}, error)

	// GetCachedRemotes retrieves the list of currently cached remote repositories
	// on the server
	GetCachedRemotes() ([]string, error)
	GetCachedRemotesWithContext(ctx context.Context) ([]string, error)

	// PutCachedRemote ensures a cached remote repository on the server -
	// This call is idempotent.
	PutCachedRemote(req *types.PutRemoteRequest) error
	PutCachedRemoteWithContext(ctx context.Context, req *types.PutRemoteRequest) error
}

// Protobuf struct is a client-side implementation of the API's protobuf struct -
// Methods in the server interface will be replaced here with client-side API calls
type Protobuf struct {
	ID           *string                     `json:"id"`
	Name         *string                     `json:"name"`
	Version      *string                     `json:"version"`
	Dependencies []*protobuf.ProtoDependency `json:"dependencies"`
	client       *registryClient
}

// registryClient implements RegistryClient
type registryClient struct {
	baseURL    string
	httpclient *http.Client
}

// New returns a new RegistryClient using the provided URL
func New(url string) RegistryClient {
	return &registryClient{
		baseURL:    strings.TrimSuffix(url, "/"),
		httpclient: &http.Client{},
	}
}
