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
	"errors"
	"testing"
)

type randomError struct{ error }

func TestProtobufNotExistsError(t *testing.T) {
	err := NewError(ProtobufNotExists{}, errors.New("not exists"))
	if !IsProtobufNotExists(err) {
		t.Error("Expected protobuf not exists error to be true, got false")
	}

	err = NewError(randomError{}, errors.New("random"))
	if IsProtobufNotExists(err) {
		t.Error("Expected protobuf not exists error to be false, got true")
	}
}
