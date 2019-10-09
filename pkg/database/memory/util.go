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

package memory

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io/ioutil"

	"github.com/tinyzimmer/protobuf-registry/pkg/protobuf"
)

func (m *memoryDatabase) dumpToDisk() error {
	var buf bytes.Buffer
	gzw := gzip.NewWriter(&buf)
	out, err := json.Marshal(m.protoBufs)
	if err != nil {
		return err
	}
	if _, err := gzw.Write(out); err != nil {
		return err
	}
	gzw.Close()
	return m.storage.StoreRawFile("mem.db", buf.Bytes())
}

func (m *memoryDatabase) loadFromDisk() error {
	data, err := m.storage.GetRawFile("mem.db")
	if err != nil {
		return err
	}
	gzr, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return err
	}
	raw, err := ioutil.ReadAll(gzr)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(raw, &m.protoBufs); err != nil {
		return err
	}
	return nil
}

func (m *memoryDatabase) addProtobuf(proto *protobuf.Protobuf) {
	if _, ok := m.protoBufs[*proto.Name]; !ok {
		m.protoBufs[*proto.Name] = []*protobuf.Protobuf{proto}
	} else {
		m.protoBufs[*proto.Name] = append(m.protoBufs[*proto.Name], proto)
	}
}

func remove(s []*protobuf.Protobuf, i int) []*protobuf.Protobuf {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}
