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

package piprouter

import (
	"github.com/gorilla/mux"
	"github.com/tinyzimmer/protobuf-registry/pkg/server/common"
)

type pipServer struct {
	ctrl *common.ServerController
}

func RegisterRoutes(router *mux.Router, path string, ctrl *common.ServerController) {
	pip := &pipServer{ctrl}
	pipRouter := router.PathPrefix(path).Subrouter()
	pipRouter.HandleFunc("/{name}/", pip.getPipVersionsHandler).Methods("GET")
	pipRouter.HandleFunc("/download/{name}", pip.getPipDownloadHandler).Methods("GET")
}
