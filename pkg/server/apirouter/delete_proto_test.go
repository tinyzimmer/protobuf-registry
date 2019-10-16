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
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tinyzimmer/protobuf-registry/pkg/config"
	"github.com/tinyzimmer/protobuf-registry/pkg/protobuf"
	"github.com/tinyzimmer/protobuf-registry/pkg/util"
)

func TestDeleteAllProtoVersionsHandler(t *testing.T) {
	srvr, rm := getServer(t)
	defer rm()

	// test OPTIONS
	req, err := http.NewRequest("OPTIONS", "/api/proto/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	srvr.router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Error("Expected no error on OPTIONS request with CORS enabled")
	}
	config.GlobalConfig.CORSEnabled = false
	rr = httptest.NewRecorder()
	srvr.router.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Error("Expected bad request on OPTIONS request with CORS disabled got:", rr.Code)
	}
	config.GlobalConfig.CORSEnabled = true

	// put some test data in the db
	_, err = srvr.DB().StoreProtoVersion(&protobuf.Protobuf{
		Name:    util.StringPtr("test"),
		Version: util.StringPtr("0.0.1"),
	}, false)
	if err != nil {
		t.Fatal(err)
	}

	tt := []struct {
		routeVar   string
		shouldPass bool
	}{
		{"test", true},
		{"not-exists", false},
	}

	for _, x := range tt {
		path := fmt.Sprintf("/api/proto/%s", x.routeVar)
		req, err := http.NewRequest("DELETE", path, nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		srvr.router.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK && x.shouldPass {
			t.Error("Handler should have passed on routeVar:", x.routeVar)
		}
		if rr.Code == http.StatusOK && !x.shouldPass {
			t.Error("Handler should have failed on routeVar:", x.routeVar)
		}
	}
}

func TestDeleteProtoVersionHandler(t *testing.T) {
	srvr, rm := getServer(t)
	defer rm()

	// test OPTIONS
	req, err := http.NewRequest("OPTIONS", "/api/proto/test/0.0.1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	srvr.router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Error("Expected no error on OPTIONS request with CORS enabled")
	}
	config.GlobalConfig.CORSEnabled = false
	rr = httptest.NewRecorder()
	srvr.router.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Error("Expected bad request on OPTIONS request with CORS disabled got:", rr.Code)
	}
	config.GlobalConfig.CORSEnabled = true

	// put some test data in the db
	_, err = srvr.DB().StoreProtoVersion(&protobuf.Protobuf{
		Name:    util.StringPtr("test"),
		Version: util.StringPtr("0.0.1"),
	}, false)
	if err != nil {
		t.Fatal(err)
	}

	tt := []struct {
		routeNameVar    string
		routeVersionVar string
		shouldPass      bool
	}{
		{"test", "0.0.1", true},
		{"test", "0.0.2", false},
		{"not-exists", "0.0.1", false},
	}

	for _, x := range tt {
		path := fmt.Sprintf("/api/proto/%s/%s", x.routeNameVar, x.routeVersionVar)
		req, err := http.NewRequest("DELETE", path, nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		srvr.router.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK && x.shouldPass {
			t.Error("Handler should have passed on routeVars:", x)
		}
		if rr.Code == http.StatusOK && !x.shouldPass {
			t.Error("Handler should have failed on routeVars:", x)
		}
	}
}
