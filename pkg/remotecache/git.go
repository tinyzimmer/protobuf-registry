package remotecache

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/tinyzimmer/protobuf-registry/pkg/types"
	"github.com/tinyzimmer/protobuf-registry/pkg/util"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

type GitDependency struct {
	LocalPath    string
	LocalSubPath string
	Revision     string
	ImportPath   string
}

func (dep *GitDependency) InjectToPath(path string) (err error) {
	localFullPath := filepath.Join(dep.LocalPath, dep.LocalSubPath)
	if dep.ImportPath == "" {
		dep.ImportPath = dep.LocalSubPath
	}
	importFullPath := filepath.Join(path, dep.ImportPath)
	log.Info(fmt.Sprintf("Injecting %s into import path %s", localFullPath, importFullPath))
	return util.CopyDir(localFullPath, importFullPath)
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

func (dep *GitDependency) Checkout() error {
	repo, err := git.PlainOpen(dep.LocalPath)
	if err != nil {
		return err
	}
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}
	if err := worktree.Pull(&git.PullOptions{}); err != nil {
		if !strings.Contains(err.Error(), "already up-to-date") {
			return err
		}
	}
	h, err := repo.ResolveRevision(plumbing.Revision(dep.Revision))
	if err != nil {
		return err
	}
	return worktree.Checkout(&git.CheckoutOptions{
		Hash: *h,
	})
}

func resolveURL(in string) (u *url.URL, subPath string, err error) {
	u, err = url.Parse(in)
	if err != nil {
		return
	}
	if u.Scheme == "" {
		u.Scheme = "https"
	}
	split := strings.Split(u.Path, "/")
	if len(split) >= 4 {
		u.Path = strings.Join(split[0:3], "/")
		subPath = strings.Join(split[3:len(split)-1], "/")
	}
	return
}
