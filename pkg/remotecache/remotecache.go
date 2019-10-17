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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/go-logr/glogr"
	"github.com/tinyzimmer/protobuf-registry/pkg/config"
)

var cache *RemoteCache

var log = glogr.New()

func InitCache() error {
	log.Info("Initializing remote dependency cache, POST operations using remote dependencies will hang until this completes")
	cache = newCache()
	if err := os.MkdirAll(cache.cacheRoot, 0700); err != nil {
		return err
	}
	for _, x := range config.GlobalConfig.PreCachedRemotes {
		log.Info(fmt.Sprintf("Updating remote dependency cache for: %s", x))
		if _, err := cache.GetGitDependency(x, "", "master"); err != nil {
			return err
		} else {
			log.Info(fmt.Sprintf("Fetched remote dependency: %s", x))
		}
	}
	log.Info("Finished initializing remote dependency cache.")
	return nil
}

type RemoteCache struct {
	cacheRoot string
	mux       sync.Mutex
}

func Cache() *RemoteCache { return cache }

func newCache() *RemoteCache {
	return &RemoteCache{
		cacheRoot: filepath.Join(config.GlobalConfig.FileStoragePath, "cache"),
	}
}

func (c *RemoteCache) GetAllRemotes() (remotes []string, err error) {
	var dirs []os.FileInfo
	remotes = make([]string, 0)
	dirs, err = ioutil.ReadDir(c.cacheRoot)
	if err != nil {
		return
	}
	for _, x := range dirs {
		if x.IsDir() {
			var dirRemotes []string
			dirRemotes, err = c.enumerateHostDir(filepath.Join(c.cacheRoot, x.Name()))
			if err != nil {
				return
			}
			remotes = append(remotes, dirRemotes...)
		}
	}
	return
}

func (c *RemoteCache) enumerateHostDir(dir string) (remotes []string, err error) {
	var dirs []os.FileInfo
	remotes = make([]string, 0)
	dirs, err = ioutil.ReadDir(dir)
	if err != nil {
		return
	}
	for _, x := range dirs {
		if x.IsDir() {
			var dirRemotes []string
			dirRemotes, err = c.enumerateGroupDir(filepath.Join(dir, x.Name()))
			if err != nil {
				return
			}
			remotes = append(remotes, dirRemotes...)
		}
	}
	return
}

func (c *RemoteCache) enumerateGroupDir(dir string) (remotes []string, err error) {
	var dirs []os.FileInfo
	remotes = make([]string, 0)
	dirs, err = ioutil.ReadDir(dir)
	if err != nil {
		return
	}
	for _, x := range dirs {
		if x.IsDir() {
			remotes = append(remotes,
				strings.Replace(filepath.Join(dir, x.Name()), c.cacheRoot+"/", "", 1),
			)
		}
	}
	return
}
