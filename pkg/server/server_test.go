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
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/tinyzimmer/protobuf-registry/pkg/config"
)

type fakeHandler struct {
	called bool
}

func (f *fakeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.called = true
}

func getServer(t *testing.T) CoreServer {
	t.Helper()
	os.Setenv("IGNORE_PROTOC", "true")
	_ = config.Init()
	srvr, err := New()
	if err != nil {
		t.Fatal("Got error building server")
	}
	os.Unsetenv("IGNORE_PROTOC")
	return srvr
}

func TestNew(t *testing.T) {
	srvr := getServer(t)
	if srvr == nil {
		t.Error("Expected a server interface, got nil")
	}
}

func TestMiddleware(t *testing.T) {
	req, err := http.NewRequest("GET", "/api", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	fk := &fakeHandler{}
	handler := loggingMiddleware(fk)
	handler.ServeHTTP(rr, req)
	if !fk.called {
		t.Error("Expected fake handler to be called after middleware")
	}
}

func TestShutdown(t *testing.T) {
	srvr := getServer(t)
	if err := srvr.Shutdown(context.Background()); err != nil {
		t.Error("Failed to shutdown server")
	}
}

func TestHealthz(t *testing.T) {
	req, err := http.NewRequest("GET", "/api", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(healthz)
	handler.ServeHTTP(rr, req)
	if !strings.Contains(rr.Body.String(), "ok") {
		t.Error("Expected response with ok, got:", rr.Body.String())
	}
}
