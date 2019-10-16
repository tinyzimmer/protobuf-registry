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

import "testing"

func TestServerError(t *testing.T) {
	err := &ServerError{ErrMsg: "test error"}
	if err.Error() != "test error" {
		t.Error("Error message came back malformed, got:", err.Error())
	}
	o := err.JSON()
	if len(o) == 0 {
		t.Error("Got back empty string for JSON method")
	}
}
