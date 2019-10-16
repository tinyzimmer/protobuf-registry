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

package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
)

const defaultVersion = "0.0.1"
const defaultRevision = "master"

type PostProtoRequest struct {
	ID            string             `json:"id,omitempty"`
	Name          string             `json:"name,omitempty"`
	Body          string             `json:"body,omitempty"`
	Version       string             `json:"version,omitempty"`
	RemoteDepends []*ProtoDependency `json:"remoteDeps,omitempty"`
}

type PutRemoteRequest struct {
	URL string `json:"url"`
}

type ProtoDependency struct {
	URL      string   `json:"url,omitempty"`
	Revision string   `json:"revision,omitempty"`
	Path     string   `json:"path,omitempty"`
	Ignores  []string `json:"ignores,omitempty"`
}

func NewProtoReqFromReader(rdr io.ReadCloser) (*PostProtoRequest, error) {
	if rdr == nil {
		return nil, errors.New("Got nil body")
	}
	defer rdr.Close()
	req := PostProtoRequest{
		Version:       defaultVersion,
		RemoteDepends: make([]*ProtoDependency, 0),
	}
	body, err := ioutil.ReadAll(rdr)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &req)
	return &req, err
}

func (req *PostProtoRequest) Validate() error {
	if req.Name == "" || req.Body == "" {
		return errors.New("'name' and 'body' are required")
	}
	if len(req.RemoteDepends) > 0 {
		for _, depend := range req.RemoteDepends {
			if depend.URL == "" {
				return fmt.Errorf("Remote dependency URL cannot be blank")
			}
			if depend.Revision == "" {
				depend.Revision = defaultRevision
			}
		}
	}
	return nil
}
