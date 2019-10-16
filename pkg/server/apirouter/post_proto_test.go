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
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tinyzimmer/protobuf-registry/pkg/config"
	"github.com/tinyzimmer/protobuf-registry/pkg/protobuf"
	"github.com/tinyzimmer/protobuf-registry/pkg/types"
)

func TestPostProtoHandler(t *testing.T) {
	srvr, rm := getServer(t)
	defer rm()
	addTestDataToServer(t, srvr)

	// test OPTIONS
	req, err := http.NewRequest("OPTIONS", "/api/proto", nil)
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

	// test invalid body
	req, err = http.NewRequest("POST", "/api/proto", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	srvr.router.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Error("Expected bad request, got:", rr.Code)
	}

	// test bad compiler
	config.GlobalConfig.ProtocPath = "/not/exist/exec"
	r := types.PostProtoRequest{
		Name:    "some-proto",
		Version: "0.0.1",
		Body:    protobuf.TestProtoZip,
	}
	out, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	req, err = http.NewRequest("POST", "/api/proto", bytes.NewReader(out))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	srvr.router.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Error("Expected bad request, got:", rr.Code)
	}
	config.GlobalConfig.ProtocPath = "echo"

	tt := []struct {
		routeNameVar    string
		routeVersionVar string
		routeBodyVar    string
		shouldPass      bool
	}{
		// already exists
		{testProtoName, "0.0.1", protobuf.TestProtoZip, false},
		// new version
		{testProtoName, "0.0.2", protobuf.TestProtoZip, true},
		// invalid body
		{testProtoName, "0.0.3", "bad", false},
		// invalid name
		{"", "0.0.2", protobuf.TestProtoZip, false},
	}

	for _, x := range tt {
		path := fmt.Sprintf("/api/proto")
		r := types.PostProtoRequest{
			Name:    x.routeNameVar,
			Version: x.routeVersionVar,
			Body:    x.routeBodyVar,
		}
		out, err := json.Marshal(r)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("POST", path, bytes.NewReader(out))
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

func TestPutProtoHandler(t *testing.T) {
	srvr, rm := getServer(t)
	defer rm()
	addTestDataToServer(t, srvr)

	// test invalid body
	req, err := http.NewRequest("PUT", "/api/proto", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	srvr.router.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Error("Expected bad request, got:", rr.Code)
	}

	tt := []struct {
		routeNameVar    string
		routeVersionVar string
		routeBodyVar    string
		shouldPass      bool
	}{
		// already exists - put should overwrite
		{testProtoName, "0.0.1", protobuf.TestProtoZip, true},
		// new version
		{testProtoName, "0.0.2", protobuf.TestProtoZip, true},
		// invalid body
		{testProtoName, "0.0.3", "bad", false},
		// invalid name
		{"", "0.0.2", protobuf.TestProtoZip, false},
	}

	for _, x := range tt {
		path := fmt.Sprintf("/api/proto")
		r := types.PostProtoRequest{
			Name:    x.routeNameVar,
			Version: x.routeVersionVar,
			Body:    x.routeBodyVar,
		}
		out, err := json.Marshal(r)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("PUT", path, bytes.NewReader(out))
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
