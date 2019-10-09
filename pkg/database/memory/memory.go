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
	"fmt"
	"sync"
	"time"

	"github.com/go-logr/glogr"
	"github.com/tinyzimmer/protobuf-registry/pkg/config"
	dbcommon "github.com/tinyzimmer/protobuf-registry/pkg/database/common"
	"github.com/tinyzimmer/protobuf-registry/pkg/protobuf"
	"github.com/tinyzimmer/protobuf-registry/pkg/storage"
	storagecommon "github.com/tinyzimmer/protobuf-registry/pkg/storage/common"
	"github.com/tinyzimmer/protobuf-registry/pkg/util"
)

var log = glogr.New()

type memoryDatabase struct {
	persistToDisk bool
	storage       storagecommon.Provider

	mux       sync.Mutex
	protoBufs map[string][]*protobuf.Protobuf
}

func NewEngine(conf *config.Config) dbcommon.DBEngine {
	db := &memoryDatabase{
		persistToDisk: conf.PersistMemoryToDisk,
		protoBufs:     make(map[string][]*protobuf.Protobuf),
	}
	if db.persistToDisk {
		db.storage = storage.GetProvider(conf)
	}
	return db
}

func (m *memoryDatabase) Init() error {
	if m.persistToDisk {
		if err := m.loadFromDisk(); err != nil {
			log.Error(err, "Failed to load memdb, will continue and attempt to create a new one")
		}
	}
	return nil
}

func (m *memoryDatabase) GetAllProtoVersions() (map[string][]*protobuf.Protobuf, error) {
	return m.protoBufs, nil
}

func (m *memoryDatabase) RemoveProtoVersion(in *protobuf.Protobuf) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	if protos, ok := m.protoBufs[*in.Name]; ok {
		newProtos := make([]*protobuf.Protobuf, len(protos))
		copy(newProtos, protos)
		for i, x := range protos {
			if *x.Version == *in.Version {
				newProtos = remove(newProtos, i)
			}
		}
		m.protoBufs[*in.Name] = newProtos
	}
	if m.persistToDisk {
		if err := m.dumpToDisk(); err != nil {
			return err
		}
	}
	return nil
}

func (m *memoryDatabase) RemoveAllVersionsForProto(name string) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	if _, ok := m.protoBufs[name]; ok {
		delete(m.protoBufs, name)
	}
	if m.persistToDisk {
		if err := m.dumpToDisk(); err != nil {
			return err
		}
	}
	return nil
}

func (m *memoryDatabase) GetProtoVersions(name string) ([]*protobuf.Protobuf, error) {
	if protos, ok := m.protoBufs[name]; ok {
		return protos, nil
	}
	return nil, dbcommon.NewError(&dbcommon.ProtobufNotExists{}, fmt.Errorf("No protobuf %s in registry", name))
}

func (m *memoryDatabase) StoreProtoVersion(proto *protobuf.Protobuf) (*protobuf.Protobuf, error) {
	m.mux.Lock()
	defer m.mux.Unlock()
	if proto.ID == nil || *proto.ID == "" {
		proto.ID = util.StringPtr(util.RandomString(32))
	}
	if existing, ok := m.protoBufs[*proto.Name]; ok {
		for _, x := range existing {
			if *x.Version == *proto.Version {
				return proto, fmt.Errorf("%s %s already exists", *proto.Name, *proto.Version)
			}
		}
	}
	proto.LastUpdated = time.Now().UTC()
	m.addProtobuf(proto)
	if m.persistToDisk {
		if err := m.dumpToDisk(); err != nil {
			return proto, err
		}
	}
	return proto, nil
}
