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

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	servercommon "github.com/tinyzimmer/protobuf-registry/pkg/server/common"
)

// newRequest returns a request object with the given path and body
func (r *registryClient) newRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	var err error
	var out []byte
	// if an object was provided for body, marshal it into json
	if body != nil {
		out, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}
	path = fmt.Sprintf("%s/%s", r.baseURL, path)
	return http.NewRequestWithContext(ctx, method, path, bytes.NewReader(out))
}

// doInto performs the provided request and unmarshals the response into the
// given interface
func (r *registryClient) doInto(req *http.Request, obj interface{}) error {
	// do the request with our httpclient
	body, err := r.doRaw(req)
	if err != nil {
		return err
	}
	// unmarshal the response into the provided interface
	return json.Unmarshal(body, obj)
}

func (r *registryClient) doRaw(req *http.Request) ([]byte, error) {
	res, err := r.httpclient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	// read the response body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		// attempt to unmarshal the body into a server error
		var err servercommon.ServerError
		if merr := json.Unmarshal(body, &err); merr != nil {
			// if it doesn't unmarshal - return a generic error
			return nil, fmt.Errorf("Could not parse error message from server\nbody: %s\nerror: %s", string(body), merr.Error())
		}
		return nil, &err
	}
	return body, nil
}
