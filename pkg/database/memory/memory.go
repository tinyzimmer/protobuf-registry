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

// a logger for this package
var log = glogr.New()

// Assert that memoryDatabase implements dbcomon.DBEngine
var _ dbcommon.DBEngine = &memoryDatabase{}

// memoryDatabase implements a DBEngine that stores its data in maps in-memory.
// It can optionally dump itself to the disk after CUD operations.
type memoryDatabase struct {
	persistToDisk bool
	storage       storagecommon.Provider

	mux       sync.Mutex
	protobufs map[string][]*protobuf.Protobuf
}

// NewEngine creates a new memoryDatabase from the given configuration
func NewEngine(conf *config.Config) *memoryDatabase {
	db := &memoryDatabase{
		persistToDisk: conf.PersistMemoryToDisk,
		protobufs:     make(map[string][]*protobuf.Protobuf),
	}
	if db.persistToDisk {
		db.storage = storage.GetProvider(conf)
	}
	return db
}

// Init checks if persistence is enabled, and if so attempts to load a pre-existing
// database from the storage provider
func (m *memoryDatabase) Init() error {
	if m.persistToDisk {
		if err := m.loadFromDisk(); err != nil {
			// Just log it and start a new one
			log.Error(err, "Failed to load memdb, will continue and attempt to create a new one")
		}
	}
	// Always return nil, returns error type just to satisfy the interface
	return nil
}

// GetAllProtoVersions returns the entire protobuf map in memory
func (m *memoryDatabase) GetAllProtoVersions() (map[string][]*protobuf.Protobuf, error) {
	return m.protobufs, nil
}

// RemoveProtoVersion removes a single proto version from the list of protobufs
// for a given name
func (m *memoryDatabase) RemoveProtoVersion(in *protobuf.Protobuf) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	if protos, ok := m.protobufs[*in.Name]; ok {
		newProtos := make([]*protobuf.Protobuf, len(protos))
		copy(newProtos, protos)
		for i, x := range protos {
			if *x.Version == *in.Version {
				newProtos = remove(newProtos, i)
			}
		}
		m.protobufs[*in.Name] = newProtos
	}
	if m.persistToDisk {
		if err := m.dumpToDisk(); err != nil {
			return err
		}
	}
	return nil
}

// RemoveAllVersionsForProto removes all versions for a protobuf package with
// the given name
func (m *memoryDatabase) RemoveAllVersionsForProto(name string) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	if _, ok := m.protobufs[name]; ok {
		delete(m.protobufs, name)
	}
	if m.persistToDisk {
		if err := m.dumpToDisk(); err != nil {
			return err
		}
	}
	return nil
}

// GetProtoVersions returns all versions for the package with the given name
func (m *memoryDatabase) GetProtoVersions(name string) ([]*protobuf.Protobuf, error) {
	if protos, ok := m.protobufs[name]; ok {
		return protos, nil
	}
	return nil, dbcommon.NewError(dbcommon.ProtobufNotExists{}, fmt.Errorf("No protobuf %s in registry", name))
}

// StoreProtoVersion writes a new proto version to the db, and optionally
// overwrites an existing one
func (m *memoryDatabase) StoreProtoVersion(proto *protobuf.Protobuf, force bool) (*protobuf.Protobuf, error) {
	m.mux.Lock()
	defer m.mux.Unlock()
	if proto.ID == nil || *proto.ID == "" {
		proto.ID = util.StringPtr(util.RandomString(32))
	}

	if existing, ok := m.protobufs[*proto.Name]; ok {
		for _, x := range existing {
			if *x.Version == *proto.Version {
				if force {
					// inherit the existing ID which will cause the storage interface
					// to overwrite the contents of the existing entry
					proto.ID = x.ID
				} else {
					// TODO - make checkable error
					return proto, fmt.Errorf("%s %s already exists", *proto.Name, *proto.Version)
				}
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
