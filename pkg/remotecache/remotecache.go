package remotecache

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/go-logr/glogr"
	"github.com/tinyzimmer/protobuf-registry/pkg/config"
	"github.com/tinyzimmer/protobuf-registry/pkg/types"
	"gopkg.in/src-d/go-git.v4"
)

var cache *RemoteCache

var log = glogr.New()

func InitCache() error {
	log.Info("Initializing remote dependency cache, POST operations using remote dependencies will hang until this completes")
	cache = newCache()
	for _, x := range config.GlobalConfig.PreCachedRemotes {
		log.Info(fmt.Sprintf("Updating remote dependency cache for: %s", x))
		if _, err := cache.GetGitDependency(&types.ProtoDependency{
			URL:      x,
			Revision: "master",
		}); err != nil {
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
	err = gdep.Checkout()
	return
}
