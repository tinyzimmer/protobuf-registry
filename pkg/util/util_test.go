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

package util

import "testing"

func TestRandomString(t *testing.T) {
	randStr := RandomString(16)
	if len(randStr) != 16 {
		t.Error("Expected string of 16 characters, got:", len(randStr))
	}
}

func TestStrPointer(t *testing.T) {
	strPtr := StringPtr("test string")
	if *strPtr != "test string" {
		t.Error("Expected pointer to reference 'test string', got:", *strPtr)
	}
}

func TestStringSliceContains(t *testing.T) {
	sl := []string{"1", "2", "3"}
	if !StringSliceContains(sl, "1") {
		t.Error("Expected slice contains 1 to be true, got false")
	}
	if StringSliceContains(sl, "4") {
		t.Error("Expected slice contains 4 to be false, got true")
	}
}

func TestStringPtrSliceContains(t *testing.T) {
	sl := []*string{StringPtr("1"), StringPtr("2")}
	if !StringPtrSliceContains(sl, StringPtr("1")) {
		t.Error("Expected slice contains pointer to 1 to be true, got false")
	}
	if StringPtrSliceContains(sl, StringPtr("4")) {
		t.Error("Expected slice contains pointer to 4 to be false, got true")
	}
}
