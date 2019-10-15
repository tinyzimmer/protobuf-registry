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
	"io/ioutil"
	"os"
	"testing"

	"github.com/tinyzimmer/protobuf-registry/pkg/config"
	dbcommon "github.com/tinyzimmer/protobuf-registry/pkg/database/common"
	"github.com/tinyzimmer/protobuf-registry/pkg/protobuf"
	"github.com/tinyzimmer/protobuf-registry/pkg/util"
)

// newTempDir returns a new directory and clean function for testing
func newTempDir(t *testing.T) (string, func()) {
	t.Helper()
	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal("Could not allocate a temp directory")
	}
	return tempDir, func() { os.RemoveAll(tempDir) }
}

// newConfWithTempStorage returns a new config with a temp storage path
func newConfWithTempStorage(t *testing.T) (conf *config.Config, clean func()) {
	t.Helper()
	tempDir, clean := newTempDir(t)
	return &config.Config{FileStoragePath: tempDir}, clean
}

// newPlainEngine returns a memoryDatabase with an empty configuration
func newPlainEngine(t *testing.T) (db *memoryDatabase) {
	t.Helper()
	conf := &config.Config{}
	return NewEngine(conf)
}

// newEngineWithPersistence returns a memoryDatabase with persistence
func newEngineWithPersistence(t *testing.T) (db *memoryDatabase, clean func()) {
	t.Helper()
	conf, rm := newConfWithTempStorage(t)
	conf.PersistMemoryToDisk = true
	return NewEngine(conf), rm
}

// newEngineWithBadPersistence returns a memoryDatabase with unwriteable persistence
func newEngineWithBadPersistence(t *testing.T) (db *memoryDatabase) {
	t.Helper()
	conf := &config.Config{}
	conf.FileStoragePath = "/not/exists/at/all"
	conf.PersistMemoryToDisk = true
	return NewEngine(conf)
}

func TestNewEngine(t *testing.T) {
	conf := &config.Config{}
	engine := NewEngine(conf)
	if len(engine.protobufs) > 0 {
		t.Error("Expected empty map for protobuf storage, got:", engine.protobufs)
	}
	conf.PersistMemoryToDisk = true
	engine = NewEngine(conf)
	if !engine.persistToDisk {
		t.Error("Expected to set memory persistence to true, got false")
	} else if engine.storage == nil {
		t.Error("Expected to setup a storage interface, got nil")
	}
}

func TestInit(t *testing.T) {
	conf, clean := newConfWithTempStorage(t)
	defer clean()

	// Check that inits return no errors
	engine := NewEngine(conf)
	if err := engine.Init(); err != nil {
		t.Error("Expected no error on plain init, got:", err)
	}
	conf.PersistMemoryToDisk = true
	engine = NewEngine(conf)
	if err := engine.Init(); err != nil {
		t.Error("Expected no error on init with persistence and new db, got:", err)
	}

	// set a dummy value to the db and dump to disk
	engine.protobufs["test"] = make([]*protobuf.Protobuf, 0)
	if err := engine.dumpToDisk(); err != nil {
		t.Fatal("Got error attempting to dump dummy data to disk")
	}

	// create a new engine
	engine = NewEngine(conf)

	// assert that we have an empty db first
	if len(engine.protobufs) > 0 {
		t.Error("Expected new engine with empty db, got:", engine.protobufs)
	}

	if err := engine.Init(); err != nil {
		t.Error("Expected no error during init, got:", err)
	}

	// our dummy data should be there now
	if _, ok := engine.protobufs["test"]; !ok {
		t.Error("Expected to load dummy data from list, did not appear to happen")
	}
}

func TestGetAllProtoVersions(t *testing.T) {
	engine := newPlainEngine(t)

	// test empty retrieval
	if protos, err := engine.GetAllProtoVersions(); err != nil {
		t.Error("Expected no error on GetAllProtoVersions(), got:", err)
	} else if len(protos) > 0 {
		t.Error("Expected empty protos map, got:", protos)
	}

	// add some dummy data
	engine.protobufs["test"] = make([]*protobuf.Protobuf, 0)

	// test object retrieval
	if protos, err := engine.GetAllProtoVersions(); err != nil {
		t.Error("Expected no error on GetAllProtoVersions(), got:", err)
	} else if len(protos) != 1 {
		t.Error("Expected one proto entry in map, got:", protos)
	} else if _, ok := protos["test"]; !ok {
		t.Error("Expected 'test' key in proto map, was not there")
	}
}

func TestRemoveProtoVersion(t *testing.T) {
	// use persistence to cover the dump functionality
	engine, clean := newEngineWithPersistence(t)
	defer clean()
	engine.protobufs["test"] = []*protobuf.Protobuf{
		{Version: util.StringPtr("0.0.1")},
		{Version: util.StringPtr("0.0.2")},
	}

	if err := engine.RemoveProtoVersion(&protobuf.Protobuf{
		Name:    util.StringPtr("test"),
		Version: util.StringPtr("0.0.1"),
	}); err != nil {
		t.Error("Expected no error, got:", err)
	}

	protos, ok := engine.protobufs["test"]
	if !ok {
		t.Fatal("The entire test slice dissapeared!")
	} else if len(protos) != 1 {
		t.Error("Expected only one protobuf left in slice")
	}

	// test error on remove attempt with unwriteable persistence
	engine = newEngineWithBadPersistence(t)
	engine.protobufs["test"] = []*protobuf.Protobuf{
		{Version: util.StringPtr("0.0.1")},
		{Version: util.StringPtr("0.0.2")},
	}
	if err := engine.RemoveProtoVersion(&protobuf.Protobuf{
		Name:    util.StringPtr("test"),
		Version: util.StringPtr("0.0.1"),
	}); err == nil {
		t.Error("Expected error attempting to write db to disk, got nil")
	}
}

func TestRemoveAllVersionsForProto(t *testing.T) {
	engine, clean := newEngineWithPersistence(t)
	defer clean()
	engine.protobufs["test"] = []*protobuf.Protobuf{
		{Version: util.StringPtr("0.0.1")},
		{Version: util.StringPtr("0.0.2")},
	}

	if err := engine.RemoveAllVersionsForProto("test"); err != nil {
		t.Error("Expected no error, got:", err)
	}

	if _, ok := engine.protobufs["test"]; ok {
		t.Error("Expected 'test' slice to be gone, it still exists")
	}

	engine = newEngineWithBadPersistence(t)
	engine.protobufs["test"] = []*protobuf.Protobuf{
		{Version: util.StringPtr("0.0.1")},
		{Version: util.StringPtr("0.0.2")},
	}
	if err := engine.RemoveAllVersionsForProto("test"); err == nil {
		t.Error("Expected error attempting to write db to disk, got nil")
	}
}

func TestGetProtoVersions(t *testing.T) {
	engine := newPlainEngine(t)
	engine.protobufs["test"] = []*protobuf.Protobuf{
		{Version: util.StringPtr("0.0.1")},
		{Version: util.StringPtr("0.0.2")},
	}

	// test exists
	protos, err := engine.GetProtoVersions("test")
	if err != nil {
		t.Fatal("Expected no error, got:", err)
	}
	if len(protos) != 2 {
		t.Error("Expected proto slice of length 2, got:", protos)
	}

	// test non-exists
	_, err = engine.GetProtoVersions("not-exists")
	if err == nil {
		t.Error("Expected error, got nil")
	} else if !dbcommon.IsProtobufNotExists(err) {
		// error checker should work on the returned error
		t.Error("Expected protobuf not exists error, got:", err)
	}
}

func TestStoreProtoVersion(t *testing.T) {
	engine, clean := newEngineWithPersistence(t)
	defer clean()

	inProto := &protobuf.Protobuf{
		Name:    util.StringPtr("test"),
		Version: util.StringPtr("0.0.1"),
	}

	if proto, err := engine.StoreProtoVersion(inProto, false); err != nil {
		t.Fatal("Expected no error storing new proto, got:", err)
	} else if proto.Name != inProto.Name {
		t.Error("Expected same name ptr as input, got:", proto.Name)
	} else if proto.Version != inProto.Version {
		t.Error("Expected same version ptr as input, got:", proto.Version)
	}

	// same proto again without force should fail
	if _, err := engine.StoreProtoVersion(inProto, false); err == nil {
		t.Error("Expecter error on existing proto, got nil")
	}

	// with force should be the same as always
	if proto, err := engine.StoreProtoVersion(inProto, true); err != nil {
		t.Fatal("Expected no error storing new proto, got:", err)
	} else if proto.Name != inProto.Name {
		t.Error("Expected same name ptr as input, got:", proto.Name)
	} else if proto.Version != inProto.Version {
		t.Error("Expected same version ptr as input, got:", proto.Version)
	}

	// test bad persistence
	engine = newEngineWithBadPersistence(t)
	if _, err := engine.StoreProtoVersion(inProto, false); err == nil {
		t.Error("Expected error from bad persistence, got nil")
	}
}
