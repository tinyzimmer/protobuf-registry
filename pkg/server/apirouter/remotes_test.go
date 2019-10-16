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

	"github.com/tinyzimmer/protobuf-registry/pkg/remotecache"
	"github.com/tinyzimmer/protobuf-registry/pkg/types"
)

func TestGetRemotes(t *testing.T) {
	srvr, rm := getServer(t)
	defer rm()

	if err := remotecache.InitCache(); err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/api/remotes", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	srvr.router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Error("Expected ok response, got:", rr.Code)
	}
}

func TestPutRemotes(t *testing.T) {
	srvr, rm := getServer(t)
	defer rm()

	if err := remotecache.InitCache(); err != nil {
		t.Fatal(err)
	}

	// test invalid body
	req, err := http.NewRequest("PUT", "/api/remotes", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	srvr.router.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Error("Expected bad request, got:", rr.Code)
	}

	tt := []struct {
		routeURLVar string
		shouldPass  bool
	}{
		{"github.com/googleapis/api-common-protos", true},
		{"github.com/nononononono/nononononononno", false},
	}

	for _, x := range tt {
		path := fmt.Sprintf("/api/remotes")
		r := types.ProtoDependency{
			URL: x.routeURLVar,
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
