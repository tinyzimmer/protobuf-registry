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

package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-logr/glogr"
	"github.com/gorilla/mux"
	"github.com/tinyzimmer/protobuf-registry/pkg/config"
	"github.com/tinyzimmer/protobuf-registry/pkg/database"
	"github.com/tinyzimmer/protobuf-registry/pkg/server/apirouter"
	"github.com/tinyzimmer/protobuf-registry/pkg/server/common"
	"github.com/tinyzimmer/protobuf-registry/pkg/server/gemrouter"
	"github.com/tinyzimmer/protobuf-registry/pkg/server/gorouter"
	"github.com/tinyzimmer/protobuf-registry/pkg/server/mvnrouter"
	"github.com/tinyzimmer/protobuf-registry/pkg/server/piprouter"
	"github.com/tinyzimmer/protobuf-registry/pkg/storage"
)

var log = glogr.New()

type CoreServer interface {
	ListenAndServe() error
	Shutdown(context.Context) error
}

type coreServer struct {
	CoreServer

	srvr *http.Server
}

func New() (CoreServer, error) {
	var ctrl *common.ServerController
	var err error
	srvr := &coreServer{}
	if ctrl, err = srvr.InitController(config.GlobalConfig); err != nil {
		return nil, err
	}
	srvr.srvr = &http.Server{
		Handler:      srvr.configureRouter(ctrl),
		Addr:         config.GlobalConfig.BindAddress,
		WriteTimeout: time.Duration(config.GlobalConfig.WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(config.GlobalConfig.ReadTimeout) * time.Second,
	}

	return srvr, nil
}

func (c *coreServer) InitController(conf *config.Config) (*common.ServerController, error) {
	// get storage and db interfaces
	ctrl := &common.ServerController{}
	ctrl.SetDBEngine(database.GetEngine(conf))
	ctrl.SetStorageProvider(storage.GetProvider(conf))
	// run init on the db interface
	err := ctrl.DB().Init()
	return ctrl, err
}

func (c *coreServer) configureRouter(ctrl *common.ServerController) *mux.Router {
	// new base router
	router := mux.NewRouter()
	router.HandleFunc("/healthz", healthz).Methods("GET")
	// UI router
	router.PathPrefix("/ui").Handler(http.StripPrefix("/ui", http.FileServer(http.Dir("./static"))))
	// main API
	apirouter.RegisterRoutes(router, "/api", ctrl)
	// pkg manager routers
	piprouter.RegisterRoutes(router, "/pip", ctrl)
	gorouter.RegisterRoutes(router, "/golang", ctrl)
	mvnrouter.RegisterRoutes(router, "/mvn", ctrl)
	gemrouter.RegisterRoutes(router, "/gem", ctrl)
	// catch-all - for debugging package discovery
	router.PathPrefix("/").Handler(common.NewCatchAllHandler())
	router.Use(loggingMiddleware)
	return router
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info(fmt.Sprintf("%s %s", r.Method, r.RequestURI))
		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

func (c *coreServer) ListenAndServe() error {
	return c.srvr.ListenAndServe()
}

func (c *coreServer) Shutdown(ctx context.Context) error {
	return c.srvr.Shutdown(ctx)
}

func healthz(w http.ResponseWriter, r *http.Request) {
	common.WriteJSONResponse(map[string]string{
		"status": "ok",
	}, w)
}
