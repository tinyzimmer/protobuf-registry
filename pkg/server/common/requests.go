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
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/tinyzimmer/proto-registry/pkg/config"
)

func GetName(r *http.Request) string {
	vars := mux.Vars(r)
	return vars["name"]
}

func GetVersion(r *http.Request) string {
	vars := mux.Vars(r)
	return vars["version"]
}

func UnmarshallInto(body io.Reader, obj interface{}) error {
	raw, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(raw, obj)
	return err
}

func BadRequest(err error, w http.ResponseWriter) {
	res := &serverError{ErrMsg: fmt.Sprintf("Bad request: %s", err.Error())}
	http.Error(w, res.JSON(), http.StatusBadRequest)
}

func NotFound(err error, w http.ResponseWriter) {
	res := &serverError{ErrMsg: fmt.Sprintf("Not Found: %s", err.Error())}
	http.Error(w, res.JSON(), http.StatusNotFound)
}

func WriteJSONResponse(res interface{}, w http.ResponseWriter) {
	out, _ := json.MarshalIndent(res, "", "    ")
	WriteRawResponse(out, w)
}

func WriteRawResponse(out []byte, w http.ResponseWriter) {
	if config.GlobalConfig.CORSEnabled {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}
	// w.WriteHeader(http.StatusOK)
	if _, err := w.Write(append(out, "\n"...)); err != nil {
		BadRequest(err, w)
	}
}

func ServeFile(w http.ResponseWriter, r *http.Request, name string, rdr io.ReadSeeker) {
	if config.GlobalConfig.CORSEnabled {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", name))
	http.ServeContent(w, r, name, time.Now(), rdr)
}
