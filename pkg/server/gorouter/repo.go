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

package gorouter

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

type gitRouter struct {
	repos map[string]*repository
}

type repository struct {
	repo    *git.Repository
	gitPath string
	repoURL string
}

func (g *gitRouter) getRepo(name string) *repository {
	if repo, ok := g.repos[name]; ok {
		return repo
	}
	return nil
}

func (g *gitRouter) newRepoFromPath(r *http.Request, protoName, pkgName string) (srv *repository, err error) {
	log.Info("Initializing new in-memory repository for package", "name", protoName, "pkgName", pkgName)
	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		return
	}
	repoPath := filepath.Join(tempDir, protoName)

	repo, err := git.PlainInit(repoPath, true)
	if err != nil {
		return
	}

	g.repos[protoName] = &repository{
		repo:    repo,
		gitPath: repoPath,
		repoURL: fmt.Sprintf("%s://%s/golang/git/%s/%s", getGitScheme(r), r.Host, protoName, pkgName),
	}

	return g.repos[protoName], nil
}

func (g *gitRouter) addToRepo(dir string, repo *repository) error {
	fs := memfs.New()
	storer := memory.NewStorage()
	memrepo, err := git.Init(storer, fs)
	if err != nil {
		return err
	}
	if _, err := memrepo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{fmt.Sprintf("file://%s", repo.gitPath)},
	}); err != nil {
		return err
	}
	wt, err := memrepo.Worktree()
	if err != nil {
		return err
	}
	if err := billyCopyDir(dir, ".", wt.Filesystem); err != nil {
		return err
	}
	if err := wt.AddGlob("*"); err != nil {
		return err
	}
	if _, err := wt.Commit("initial commit", &git.CommitOptions{
		All: true,
		Author: &object.Signature{
			Name:  "proto-registry",
			Email: "proto@registry.com",
			When:  time.Now(),
		},
	}); err != nil {
		return err
	}
	if err := memrepo.Push(&git.PushOptions{}); err != nil {
		return err
	}
	return nil
}

func billyCopyDir(src, dst string, fs billy.Filesystem) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = fs.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}
	if err == nil {
		return fmt.Errorf("destination already exists")
	}

	err = fs.MkdirAll(dst, si.Mode())
	if err != nil {
		return
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = billyCopyDir(srcPath, dstPath, fs)
			if err != nil {
				return
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = billyCopyFile(srcPath, dstPath, fs)
			if err != nil {
				return
			}
		}
	}
	return
}

func billyCopyFile(src, dst string, fs billy.Filesystem) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := fs.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	return
}
