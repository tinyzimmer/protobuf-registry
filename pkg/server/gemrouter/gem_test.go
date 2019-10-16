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
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/tinyzimmer/protobuf-registry/pkg/config"
	"github.com/tinyzimmer/protobuf-registry/pkg/database"
	"github.com/tinyzimmer/protobuf-registry/pkg/protobuf"
	"github.com/tinyzimmer/protobuf-registry/pkg/server/common"
	"github.com/tinyzimmer/protobuf-registry/pkg/storage"
	"github.com/tinyzimmer/protobuf-registry/pkg/util"
)

func getController(t *testing.T) (*common.ServerController, func()) {
	t.Helper()
	config.SafeInit()
	config.GlobalConfig.CORSEnabled = true
	config.GlobalConfig.FileStoragePath, _ = ioutil.TempDir("", "")
	config.GlobalConfig.ProtocPath = "echo"
	config.GlobalConfig.ProtobufVersion = "3.6.1"
	ctrl := &common.ServerController{}
	ctrl.SetDBEngine(database.GetEngine(config.GlobalConfig))
	ctrl.SetStorageProvider(storage.GetProvider(config.GlobalConfig))
	return ctrl, func() { os.RemoveAll(config.GlobalConfig.FileStoragePath) }
}

func getServer(t *testing.T) (*gemServer, *mux.Router, func()) {
	t.Helper()
	ctrl, rm := getController(t)
	router := mux.NewRouter()
	RegisterRoutes(router, "/gem", ctrl)
	return &gemServer{ctrl: ctrl}, router, rm
}

func TestRoutes(t *testing.T) {
	srvr, router, rm := getServer(t)
	defer rm()

	req, err := http.NewRequest("GET", "/gem/specs.4.8.gz", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Error("Expected OK, got:", rr.Code)
	}

	req, err = http.NewRequest("GET", "/gem/quick/Marshal.4.8/test-proto-0.0.1.tar.gz.gemspec.rz", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Error("Expected bad request on non-exist proto, got:", rr.Code)
	}

	_, _ = srvr.ctrl.DB().StoreProtoVersion(&protobuf.Protobuf{
		Name:    util.StringPtr("test-proto"),
		Version: util.StringPtr("0.0.1"),
	}, false)

	req, err = http.NewRequest("GET", "/gem/quick/Marshal.4.8/test-proto-0.0.1.tar.gz.gemspec.rz", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Error("Expected ok, got:", rr.Code)
	}

	// test first with data
	req, err = http.NewRequest("GET", "/gem/specs.4.8.gz", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Error("Expected OK, got:", rr.Code)
	}
}
