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

package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-logr/glogr"
	"github.com/tinyzimmer/protobuf-registry/pkg/config"
	"github.com/tinyzimmer/protobuf-registry/pkg/remotecache"
	"github.com/tinyzimmer/protobuf-registry/pkg/server"
)

var CompileDate string
var GitCommit string

var log = glogr.New()
var catchSignals = []os.Signal{syscall.SIGINT, syscall.SIGTERM}

func main() {

	if err := config.Init(false); err != nil {
		os.Exit(1)
	}

	setupLog := log.WithName("setup")

	// log the configuration
	setupLog.Info(string(config.GlobalConfig.JSON()))

	config.GlobalConfig.CompileDate = CompileDate
	config.GlobalConfig.GitCommit = GitCommit

	// start the webserver
	setupLog.Info("Starting proto-registry server")
	srvr, err := server.New()
	if err != nil {
		setupLog.Error(err, "Failed to initialize the proto-registry server")
		os.Exit(1)
	}

	// kick off the webserver in a goroutine
	go func() {
		if err := srvr.ListenAndServe(); err != nil {
			setupLog.Error(err, "Fatal error serving web interface")
			os.Exit(1)
		}
	}()

	// initialize the cache
	if err := remotecache.InitCache(); err != nil {
		setupLog.Error(err, "Failed to initialize remote dependency cache")
		os.Exit(1)
	}

	// make a signal channel
	c := make(chan os.Signal, 1)
	// catch SIGINT and SIGTERM
	signal.Notify(c, catchSignals...)

	// block until we receive our signal.
	<-c

	// setup a shutdown logger
	shutdownLog := log.WithName("shutdown")
	shutdownLog.Info("Shutting down proto-registry")
	// create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 15)
	defer cancel()
	// shutdown the web server
	if err := srvr.Shutdown(ctx); err != nil {
		shutdownLog.Error(err, "Failed to shutdown web server")
		os.Exit(1)
	}

	os.Exit(0)
}
