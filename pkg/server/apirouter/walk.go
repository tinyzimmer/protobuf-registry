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
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tinyzimmer/proto-registry/pkg/server/common"
)

type WalkResponse struct {
	Routes []*Route `json:"routes"`
}

func newWalkResponse() *WalkResponse {
	return &WalkResponse{
		Routes: make([]*Route, 0),
	}
}

type Route struct {
	Path        string `json:"path"`
	PathRegexp  string `json:"pathRegex"`
	Method      string `json:"method"`
	Description string `json:"description"`
}

func (api *apiServer) walkRouter(w http.ResponseWriter, r *http.Request) {
	res := newWalkResponse()
	err := api.router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err != nil {
			return err
		} else if pathTemplate == "/ui" {
			return nil
		}
		pathRegex, err := route.GetPathRegexp()
		if err != nil {
			return err
		}
		methods, err := route.GetMethods()
		if err != nil {
			return nil
		}
		method := getNonOptionsMethod(methods)
		path := &Route{
			Path:        pathTemplate,
			PathRegexp:  pathRegex,
			Method:      method,
			Description: GetDoc(pathTemplate, method),
		}
		res.Routes = append(res.Routes, path)
		return nil
	})

	if err != nil {
		common.BadRequest(err, w)
		return
	}
	common.WriteJSONResponse(res, w)
}

func getNonOptionsMethod(methods []string) string {
	for _, x := range methods {
		if x != "OPTIONS" {
			return x
		}
	}
	return ""
}
