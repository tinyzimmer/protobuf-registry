package remotecache

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/go-logr/glogr"
	"github.com/tinyzimmer/protobuf-registry/pkg/config"
	"github.com/tinyzimmer/protobuf-registry/pkg/types"
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
