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

package remotecache

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/tinyzimmer/protobuf-registry/pkg/config"
)

func setConfig(t *testing.T) (rm func()) {
	t.Helper()
	config.SafeInit()
	config.GlobalConfig.FileStoragePath, _ = ioutil.TempDir("", "")
	return func() { os.RemoveAll(config.GlobalConfig.FileStoragePath) }
}

func TestInitCache(t *testing.T) {
	rm := setConfig(t)
	defer rm()

	if err := InitCache(); err != nil {
		t.Error("Expected init with nothing to cache to have no error, got:", err)
	}

	if Cache() != cache {
		t.Error("Cache() should point to the global cache")
	}

	config.GlobalConfig.PreCachedRemotes = []string{"github.com/googleapis/api-common-protos"}
	if err := InitCache(); err != nil {
		t.Error("Expected no error while precaching googleapis")
	}

	if err := InitCache(); err != nil {
		t.Error("Expected no error while precaching pre-existing remote")
	}

	config.GlobalConfig.PreCachedRemotes = []string{"github.com/googleapis/i-dont-exist"}
	if err := InitCache(); err == nil {
		t.Error("Expected error trying to clone non-existant repo")
	}

	config.GlobalConfig.PreCachedRemotes = []string{"file://github\\bad/url/i-dont-exist\\"}
	if err := InitCache(); err == nil {
		t.Error("Expected error from invalid URL")
	}

	config.GlobalConfig.PreCachedRemotes = []string{"https://github/something/still/a/bad/url"}
	if err := InitCache(); err == nil {
		t.Error("Expected error from invalid URL")
	}
}

func TestGetRemotes(t *testing.T) {
	rm := setConfig(t)
	defer rm()

	if err := InitCache(); err != nil {
		t.Error("Expected init with nothing to cache to have no error, got:", err)
	}

	remotes, err := cache.GetAllRemotes()
	if err != nil {
		t.Fatal("Expected no error getting remotes, got:", err)
	}
	if len(remotes) != 0 {
		t.Error("Expected empty list of remotes, got:", remotes)
	}

	config.GlobalConfig.PreCachedRemotes = []string{"github.com/googleapis/api-common-protos"}
	if err := InitCache(); err != nil {
		t.Error("Expected no error while precaching googleapis")
	}

	remotes, err = cache.GetAllRemotes()
	if err != nil {
		t.Fatal("Expected no error getting remotes, got:", err)
	}
	if len(remotes) != 1 {
		t.Error("Expected list of 1 remote, got:", remotes)
	}
}
