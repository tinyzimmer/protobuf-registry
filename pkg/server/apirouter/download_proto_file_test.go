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
)

func TestGetRawProtoFile(t *testing.T) {
	srvr, rm := getServer(t)
	defer rm()
	addTestDataToServer(t, srvr)

	tt := []struct {
		routeNameVar     string
		routeVersionVar  string
		routeFilenameVar string
		shouldPass       bool
	}{
		{testProtoName, testProtoVersion, "TestProtoMessage.proto", true},
		{testProtoName, testProtoVersion, "non-exist-file.proto", false},
		{testProtoName, "0.0.2", "TestProtoMessage.proto", false},
		{"not-exists", "0.0.1", "somemessage.proto", false},
	}

	for _, x := range tt {
		path := fmt.Sprintf("/api/proto/%s/%s/raw/%s", x.routeNameVar, x.routeVersionVar, x.routeFilenameVar)
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		srvr.router.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK && x.shouldPass {
			t.Error("Handler should have passed on routeVars:", x, "res:", rr.Body.String())
		}
		if rr.Code == http.StatusOK && !x.shouldPass {
			t.Error("Handler should have failed on routeVars:", x, "res:", rr.Body.String())
		}
	}
}

func TestGetMetaForProtoFile(t *testing.T) {
	srvr, rm := getServer(t)
	defer rm()
	addTestDataToServer(t, srvr)

	tt := []struct {
		routeNameVar     string
		routeVersionVar  string
		routeFilenameVar string
		shouldPass       bool
	}{
		{testProtoName, testProtoVersion, "TestProtoMessage.proto", true},
		{testProtoName, testProtoVersion, "non-exist-file.proto", false},
		{testProtoName, "0.0.2", "TestProtoMessage.proto", false},
		{"not-exists", "0.0.1", "somemessage.proto", false},
	}

	for _, x := range tt {
		path := fmt.Sprintf("/api/proto/%s/%s/meta/%s", x.routeNameVar, x.routeVersionVar, x.routeFilenameVar)
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		srvr.router.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK && x.shouldPass {
			t.Error("Handler should have passed on routeVars:", x, "res:", rr.Body.String())
		}
		if rr.Code == http.StatusOK && !x.shouldPass {
			t.Error("Handler should have failed on routeVars:", x, "res:", rr.Body.String())
		}
	}
}
