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
	"os"
	"path/filepath"

	"github.com/tinyzimmer/protobuf-registry/pkg/types"
	"gopkg.in/src-d/go-git.v4"
)

func (c *RemoteCache) GetGitDependency(dep *types.ProtoDependency) (gdep *GitDependency, err error) {
	c.mux.Lock()
	defer c.mux.Unlock()
	cloneURL, subPath, err := resolveURL(dep.URL)
	if err != nil {
		return
	}
	path := filepath.Join(c.cacheRoot, cloneURL.Path)
	// check if we already have a cached clone
	if _, err = os.Stat(path); err == nil {
		gdep = &GitDependency{
			LocalPath:    path,
			LocalSubPath: subPath,
			Revision:     dep.Revision,
			ImportPath:   dep.Path,
		}
		err = gdep.Checkout()
		return
	}
	cloneOpts := &git.CloneOptions{
		URL: cloneURL.String(),
	}
	log.Info(fmt.Sprintf("Cloning %s", cloneOpts.URL))
	_, err = git.PlainClone(path, false, cloneOpts)
	if err != nil {
		return
	}
	gdep = &GitDependency{
		LocalPath:    path,
		LocalSubPath: subPath,
		Revision:     dep.Revision,
		ImportPath:   dep.Path,
	}
	log.Info(fmt.Sprintf("Checking out %s of %s", dep.Revision, dep.URL))
	err = gdep.Checkout()
	return
}