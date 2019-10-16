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

package common

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/tinyzimmer/protobuf-registry/pkg/config"
)

func testRouter(t *testing.T) *mux.Router {
	t.Helper()
	router := mux.NewRouter()
	router.HandleFunc("/test/{version}", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(GetVersion(r)))
	})
	router.HandleFunc("/{name}", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(GetName(r)))
	})
	return router
}

func TestGetName(t *testing.T) {
	router := testRouter(t)
	req, err := http.NewRequest("GET", "/test-name", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Body.String() != "test-name" {
		t.Error("Expected 'test-name', got:", rr.Body.String())
	}
}

func TestGetVersion(t *testing.T) {
	router := testRouter(t)
	req, err := http.NewRequest("GET", "/test/test-version", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Body.String() != "test-version" {
		t.Error("Expected 'test-version', got:", rr.Body.String())
	}
}

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

func TestUnmarshallInto(t *testing.T) {
	var testObj struct {
		Name    string
		Version string
	}
	var raw = []byte(`{"Name":"Test Name","Version":"Test Version"}`)
	if err := UnmarshallInto(bytes.NewReader(raw), &testObj); err != nil {
		t.Error("Expected to unmarshal object, got:", err)
	}
	if testObj.Name != "Test Name" {
		t.Error("Name was unmarshalled incorrectly, got:", testObj.Name)
	}
	if testObj.Version != "Test Version" {
		t.Error("Version was unmarshalled incorrectly, got:", testObj.Version)
	}

	// test empty body
	if err := UnmarshallInto(nil, &testObj); err == nil {
		t.Error("Expected error for no reader, got nil")
	}

	// test bad reader
	if err := UnmarshallInto(errReader(0), &testObj); err == nil {
		t.Error("Expected error for bad reader, got nil")
	}
}

func TestBadRequest(t *testing.T) {
	rr := httptest.NewRecorder()
	BadRequest(errors.New("bad"), rr)
	if rr.Code != http.StatusBadRequest {
		t.Error("Expected bad request, got:", rr.Code)
	}
}

func TestNotFound(t *testing.T) {
	rr := httptest.NewRecorder()
	NotFound(errors.New("bad"), rr)
	if rr.Code != http.StatusNotFound {
		t.Error("Expected not found, got:", rr.Code)
	}
}

func TestWriteJSONResponse(t *testing.T) {
	rr := httptest.NewRecorder()
	WriteJSONResponse(map[string]string{
		"key": "value",
	}, rr)
	if rr.Code != http.StatusOK {
		t.Error("Expected ok, got:", rr.Code)
	}

	rr = httptest.NewRecorder()
	WriteJSONResponse(make(chan int), rr)
	if rr.Code != http.StatusBadRequest {
		t.Error("Expected bad request, got:", rr.Code)
	}
}

func TestWriteRawResponse(t *testing.T) {
	config.SafeInit()

	rr := httptest.NewRecorder()
	WriteRawResponse([]byte("test"), rr)

	if rr.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Error("Got CORS header when CORS was disabled")
	}
	if rr.Code != http.StatusOK {
		t.Error("Expected ok, got:", rr.Code)
	}

	config.GlobalConfig.CORSEnabled = true
	rr = httptest.NewRecorder()
	WriteRawResponse([]byte("test"), rr)
	if rr.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("Expected wildcard origin CORS header")
	}

}

func TestServeFile(t *testing.T) {
	config.SafeInit()

	req, err := http.NewRequest("GET", "/test/path", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	ServeFile(rr, req, "test-file", bytes.NewReader([]byte("test-data")))

	if rr.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Error("Got CORS header when CORS was disabled")
	}
	if rr.Code != http.StatusOK {
		t.Error("Expected ok, got:", rr.Code)
	}

	config.GlobalConfig.CORSEnabled = true
	rr = httptest.NewRecorder()
	ServeFile(rr, req, "test-file", bytes.NewReader([]byte("test-data")))
	if rr.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("Expected wildcard origin CORS header")
	}
}
