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
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/tinyzimmer/protobuf-registry/pkg/config"
)

type CatchAllHandler struct {
	http.Handler
}

func (c *CatchAllHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if config.GlobalConfig.RedirectNotFoundToUI {
		http.Redirect(w, r, "/ui", http.StatusSeeOther)
		return
	}
	fmt.Printf("%+v", r)
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}
	fmt.Println(string(body))
	http.Error(w, "not found", http.StatusNotFound)
}

func NewCatchAllHandler() *CatchAllHandler {
	return &CatchAllHandler{}
}
