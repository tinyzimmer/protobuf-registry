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
	"errors"
	"io/ioutil"
	"strings"
	"testing"
)

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

func (errReader) Close() error {
	return errors.New("test error")
}

var rawBasic = `
{
  "name": "test-proto",
  "body": "test-body"
}
`

var rawWithVersion = `
{
  "name": "test-proto",
  "body": "test-body",
  "version": "0.0.1"
}
`

var rawWithBadBody = `
{
  "name": "test-proto",
  "version": "0.0.1"
}
`

var rawWithBadName = `
{
  "body": "test-body",
  "version": "0.0.1"
}
`

var rawWithRemoteDepends = `
{
  "name": "test-proto",
  "body": "test-body",
  "version": "0.0.1",
  "remoteDeps": [
    {
      "url": "github.com/googleapis/api-common-protos"
    }
  ]
}
`

var rawWithBadRemoteDepends = `
{
  "name": "test-proto",
  "body": "test-body",
  "version": "0.0.1",
  "remoteDeps": [
    {
      "revision": "master"
    }
  ]
}
`

func TestPostProtoReq(t *testing.T) {

	if _, err := NewProtoReqFromReader(nil); err == nil {
		t.Error("Expected error from no body, got nil")
	}

	if _, err := NewProtoReqFromReader(errReader(0)); err == nil {
		t.Error("Expected error from bad reader, got nil")
	}

	tt := []struct {
		body           string
		shouldValidate bool
	}{
		{rawBasic, true},
		{rawWithVersion, true},
		{rawWithBadBody, false},
		{rawWithBadName, false},
		{rawWithRemoteDepends, true},
		{rawWithBadRemoteDepends, false},
	}

	for _, x := range tt {
		req, err := NewProtoReqFromReader(ioutil.NopCloser(strings.NewReader(x.body)))
		if err != nil {
			t.Fatal(err)
		}
		err = req.Validate()
		if err == nil && !x.shouldValidate {
			t.Error("Expected body to not pass validation:", x.body)
		}
		if err != nil && x.shouldValidate {
			t.Error("Expected body to pass validation:", x.body)
		}
	}

}
