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

import (
	"math/rand"
	"strings"
	"time"
)

var randStringChars []rune = []rune("abcdefghijklmnopqrstuvwxyz" + "0123456789")

func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(randStringChars[rand.Intn(len(randStringChars))])
	}
	return b.String()
}

func StringPtr(str string) *string {
	return &str
}

func StringSliceContains(sl []string, s string) bool {
	for _, x := range sl {
		if x == s {
			return true
		}
	}
	return false
}

func StringPtrSliceContains(sl []*string, s *string) bool {
	for _, x := range sl {
		if *x == *s {
			return true
		}
	}
	return false
}
