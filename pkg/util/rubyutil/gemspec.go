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
	"bytes"

	"github.com/samcday/rmarsh"
)

func newGemSpec(name, version string) (out string, err error) {
	var buf bytes.Buffer
	gen := rmarsh.NewGenerator(&buf)

	if err = gen.StartArray(6); err != nil {
		return
	}

	// generic version info - ruby version
	if err = rawStrIVar(gen, "2.4.8"); err != nil {
		return
	}

	// some number
	if err = gen.Fixnum(4); err != nil {
		return
	}

	// name info
	if err = newNameInfo(gen, name); err != nil {
		return
	}

	// version info
	if err = newVersionInfo(gen, version, false); err != nil {
		return
	}

	// modtime info
	if err = newTimeObject(gen); err != nil {
		return
	}

	// summary
	if err = rawStrIVar(gen, "foo syummary"); err != nil {
		return
	}

	if err = gen.EndArray(); err != nil {
		return
	}

	out = buf.String()
	return
}

func newMarshalledTime() (out string, err error) {
	var buf bytes.Buffer
	gen := rmarsh.NewGenerator(&buf)
	if err = gen.String("2011-10-05"); err != nil {
		return
	}
	out = buf.String()
	return
}

func newTimeObject(gen *rmarsh.Generator) (err error) {
	t, err := newMarshalledTime()
	if err != nil {
		return
	}

	if err = gen.StartIVar(1); err != nil {
		return
	}

	if err = gen.UserDefinedObject("Time", t); err != nil {
		return
	}

	if err = gen.Symbol("submicro"); err != nil {
		return
	}

	if err = gen.String(""); err != nil {
		return
	}

	if err = gen.EndIVar(); err != nil {
		return
	}

	return
}
