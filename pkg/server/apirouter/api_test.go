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
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/tinyzimmer/protobuf-registry/pkg/config"
	"github.com/tinyzimmer/protobuf-registry/pkg/database"
	"github.com/tinyzimmer/protobuf-registry/pkg/protobuf"
	"github.com/tinyzimmer/protobuf-registry/pkg/remotecache"
	"github.com/tinyzimmer/protobuf-registry/pkg/server/common"
	"github.com/tinyzimmer/protobuf-registry/pkg/storage"
	"github.com/tinyzimmer/protobuf-registry/pkg/types"
)

var testProtoName = "test-proto"
var testProtoVersion = "0.0.1"

func getController(t *testing.T) (*common.ServerController, func()) {
	t.Helper()
	config.SafeInit()
	config.GlobalConfig.CORSEnabled = true
	config.GlobalConfig.FileStoragePath, _ = ioutil.TempDir("", "")
	config.GlobalConfig.ProtocPath = "echo"
	config.GlobalConfig.ProtobufVersion = "3.6.1"
	if err := remotecache.InitCache(); err != nil {
		t.Fatal(err)
	}
	ctrl := &common.ServerController{}
	ctrl.SetDBEngine(database.GetEngine(config.GlobalConfig))
	ctrl.SetStorageProvider(storage.GetProvider(config.GlobalConfig))
	return ctrl, func() { os.RemoveAll(config.GlobalConfig.FileStoragePath) }
}

func getServer(t *testing.T) (*apiServer, func()) {
	ctrl, rm := getController(t)
	router := mux.NewRouter()
	RegisterRoutes(router, "/api", ctrl)
	return &apiServer{ctrl: ctrl, router: router}, rm
}

func addTestDataToServer(t *testing.T, srvr *apiServer) {
	t.Helper()
	r := &types.PostProtoRequest{
		Name:    testProtoName,
		Version: testProtoVersion,
		Body:    protobuf.TestProtoZip,
		RemoteDepends: []*protobuf.ProtoDependency{
			{URL: "github.com/googleapis/api-common-protos"},
		},
	}
	o, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/api/proto", bytes.NewReader(o))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	srvr.router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatal("Failed to post test data", rr.Body.String())
	}
	var proto protobuf.Protobuf
	err = json.Unmarshal(rr.Body.Bytes(), &proto)
	if err != nil {
		t.Fatal(err)
	}
	// overwrite storage with valid descriptor set and re-add zip to object
	descSet, err := base64.StdEncoding.DecodeString(protobuf.TestProtoDescriptorSet)
	if err != nil {
		t.Fatal(err)
	}
	if err := proto.SetRawFromBase64(protobuf.TestProtoZip); err != nil {
		t.Fatal(err)
	}
	proto.SetDescriptor(descSet)
	err = srvr.ctrl.Storage().StoreProtoPackage(&proto)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRegisterRoutes(t *testing.T) {
	router := mux.NewRouter()
	ctrl, rm := getController(t)
	defer rm()
	RegisterRoutes(router, "/api", ctrl)
}
