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
	"github.com/gorilla/mux"
	"github.com/tinyzimmer/proto-registry/pkg/server/common"
)

type apiServer struct {
	ctrl   *common.ServerController
	router *mux.Router
}

func RegisterRoutes(router *mux.Router, path string, ctrl *common.ServerController) {
	api := &apiServer{ctrl, router}

	apiRouter := router.PathPrefix(path).Subrouter()

	apiRouter.HandleFunc("",
		api.walkRouter).
		Methods("GET")

	apiRouter.HandleFunc("/config",
		api.getConfigHandler).
		Methods("GET")

	apiRouter.HandleFunc("/proto",
		api.postProtoHandler).
		Methods("OPTIONS", "POST")

	apiRouter.HandleFunc("/proto",
		api.getAllProtoHandler).
		Methods("GET")

	apiRouter.HandleFunc("/proto/{name}",
		api.getProtoHandler).
		Methods("GET")
	apiRouter.HandleFunc("/proto/{name}",
		api.deleteAllProtoVersionsHandler).
		Methods("OPTIONS", "DELETE")

	apiRouter.HandleFunc("/proto/{name}/{version}",
		api.getProtoVersionMetaHandler).
		Methods("GET")
	apiRouter.HandleFunc("/proto/{name}/{version}",
		api.deleteProtoVersionHandler).
		Methods("OPTIONS", "DELETE")

	apiRouter.HandleFunc("/proto/{name}/{version}/{language}",
		api.downloadProtoHandler).
		Methods("GET")

	apiRouter.PathPrefix("/proto/{name}/{version}/raw/{filename}").HandlerFunc(
		api.getRawProtoFile).
		Methods("GET")

	apiRouter.PathPrefix("/proto/{name}/{version}/meta/{filename}").HandlerFunc(
		api.getMetaForProtoFile).
		Methods("GET")

}
