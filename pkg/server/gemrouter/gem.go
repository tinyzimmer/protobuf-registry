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

package gemrouter

import (
	"bytes"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tinyzimmer/proto-registry/pkg/protobuf"
	"github.com/tinyzimmer/proto-registry/pkg/server/common"
	"github.com/tinyzimmer/proto-registry/pkg/util/rubyutil"
)

type gemServer struct {
	ctrl *common.ServerController
}

func RegisterRoutes(router *mux.Router, path string, ctrl *common.ServerController) {
	gem := &gemServer{ctrl}
	gemRouter := router.PathPrefix(path).Subrouter()
	gemRouter.HandleFunc("/specs.4.8.gz", gem.getSpecsHandler).Methods("GET")
	gemRouter.HandleFunc("/quick/Marshal.4.8/{name}.gemspec.rz", gem.getPackageSpecHandler).Methods("GET")
}

func (gem *gemServer) getSpecsHandler(w http.ResponseWriter, r *http.Request) {
	protos, err := gem.ctrl.DB().GetAllProtoVersions()
	if err != nil {
		common.BadRequest(err, w)
		return
	}
	spec, err := rubyutil.NewRubyGemsListFromPackages(squashProtoMap(protos))
	if err != nil {
		common.BadRequest(err, w)
	}

	common.ServeFile(w, r, "specs.4.8.gz", bytes.NewReader(spec))
}

func (gem *gemServer) getPackageSpecHandler(w http.ResponseWriter, r *http.Request) {
	filename := common.GetName(r)
	pkg, version := common.ParseNameVersionExtString(filename, ".tar.gz")
	protos, err := gem.ctrl.DB().GetProtoVersions(pkg)
	if err != nil {
		common.BadRequest(err, w)
		return
	}
	var proto *protobuf.Protobuf
	if proto, err = common.GetVersionFromProtoSlice(protos, version); err != nil {
		common.BadRequest(err, w)
		return
	}

	out, err := rubyutil.NewGemSpecFromPackage(proto)
	if err != nil {
		common.BadRequest(err, w)
		return
	}

	common.ServeFile(w, r, filename, bytes.NewReader(out))
}

func squashProtoMap(in map[string][]*protobuf.Protobuf) []*protobuf.Protobuf {
	out := make([]*protobuf.Protobuf, 0)
	for _, x := range in {
		out = append(out, x...)
	}
	return out
}
