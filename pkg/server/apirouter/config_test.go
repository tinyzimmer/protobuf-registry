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
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/tinyzimmer/protobuf-registry/pkg/config"
)

func TestGetConfigHandler(t *testing.T) {
	srvr, rm := getServer(t)
	defer rm()
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/api/config", nil)
	if err != nil {
		t.Fatal(err)
	}

	srvr.router.ServeHTTP(rr, req)

	var conf config.Config
	err = json.Unmarshal(rr.Body.Bytes(), &conf)
	if err != nil {
		t.Error("Did not get valid JSON back")
	}
	if !reflect.DeepEqual(conf, *config.GlobalConfig) {
		t.Error("Config was malformed in response")
	}
}
