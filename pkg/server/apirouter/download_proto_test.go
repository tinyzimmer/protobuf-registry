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

func TestDownloadProto(t *testing.T) {
	srvr, rm := getServer(t)
	defer rm()
	addTestDataToServer(t, srvr)

	tt := []struct {
		routeNameVar    string
		routeVersionVar string
		routeTargetVar  string
		shouldPass      bool
	}{
		{testProtoName, testProtoVersion, "raw", true},
		{testProtoName, testProtoVersion, "descriptors", true},
		{testProtoName, testProtoVersion, "go", true},
		{testProtoName, testProtoVersion, "csharp", true},
		{testProtoName, testProtoVersion, "java", true},
		{testProtoName, testProtoVersion, "js", true},
		{testProtoName, testProtoVersion, "objc", true},
		{testProtoName, testProtoVersion, "php", true},
		{testProtoName, testProtoVersion, "python", true},
		{testProtoName, testProtoVersion, "ruby", true},
		{testProtoName, testProtoVersion, "cpp", true},
		{testProtoName, testProtoVersion, "tinyzimmerlang", false},
		{testProtoName, "0.0.2", "go", false},
		{"not-exists", "0.0.1", "go", false},
	}

	for _, x := range tt {
		path := fmt.Sprintf("/api/proto/%s/%s/%s", x.routeNameVar, x.routeVersionVar, x.routeTargetVar)
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
