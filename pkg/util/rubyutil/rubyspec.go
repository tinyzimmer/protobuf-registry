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

package rubyutil

import (
	"github.com/samcday/rmarsh"
)

func addPackage(gen *rmarsh.Generator, name, version string) (err error) {

	// Each package is an array of length 2 with Ivar representing package name/meta
	// And a user marshalled version information object
	if err = gen.StartArray(2); err != nil {
		return
	}

	if err = newNameInfo(gen, name); err != nil {
		return
	}

	if err = newVersionInfo(gen, version, true); err != nil {
		return
	}

	if err = gen.EndArray(); err != nil {
		return
	}
	return nil
}

func newNameInfo(gen *rmarsh.Generator, name string) (err error) {
	return rawStrIVar(gen, name)
}

func newVersionInfo(gen *rmarsh.Generator, version string, withRuby bool) (err error) {
	// VERSION OBJECT
	if err = gen.StartUserMarshalled("Gem::Version"); err != nil {
		return
	}

	var length int
	if withRuby {
		length = 2
	} else {
		length = 1
	}

	if err = gen.StartArray(length); err != nil {
		return
	}

	if err = rawStrIVar(gen, version); err != nil {
		return
	}

	if withRuby {
		if err = rawStrIVar(gen, "ruby"); err != nil {
			return
		}
	}

	if err = gen.EndArray(); err != nil {
		return
	}

	if err = gen.EndUserMarshalled(); err != nil {
		return
	}

	return nil
}
