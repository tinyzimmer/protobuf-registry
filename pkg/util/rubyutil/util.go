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
	"compress/gzip"
	"compress/zlib"

	"github.com/samcday/rmarsh"
	"github.com/tinyzimmer/protobuf-registry/pkg/protobuf"
)

func NewRubyGemsListFromPackages(in []*protobuf.Protobuf) (out []byte, err error) {
	var buf bytes.Buffer
	gzw := gzip.NewWriter(&buf)
	gen := rmarsh.NewGenerator(gzw)

	if err = gen.StartArray(len(in)); err != nil {
		return
	}
	for _, x := range in {
		if err = addPackage(gen, *x.Name, *x.Version); err != nil {
			return
		}
	}
	if err = gen.EndArray(); err != nil {
		return
	}

	if err = gzw.Close(); err != nil {
		return
	}

	out = buf.Bytes()
	return
}

func NewGemSpecFromPackage(in *protobuf.Protobuf) (out []byte, err error) {
	var buf bytes.Buffer
	zw := zlib.NewWriter(&buf)
	gen := rmarsh.NewGenerator(zw)

	objDef, err := newGemSpec(*in.Name, *in.Version)
	if err != nil {
		return nil, err
	}

	if err = gen.UserDefinedObject("Gem::Specification", objDef); err != nil {
		return nil, err
	}

	zw.Close()

	out = buf.Bytes()
	return
}
