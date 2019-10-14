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
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"syscall"

	"github.com/tinyzimmer/protobuf-registry/pkg/server/common"
)

func (g *gitRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	path := strings.TrimPrefix(r.RequestURI, "/golang/git/")

	split := strings.Split(path, "/")
	if len(split) < 4 {
		common.BadRequest(fmt.Errorf("Incomplete path: %s", path), w)
		return
	}

	pkgRef := split[0]
	repo, ok := g.repos[pkgRef]
	if !ok {
		common.BadRequest(fmt.Errorf("No reference to %s available", pkgRef), w)
		return
	}

	switch split[len(split)-1] {
	case "refs?service=git-upload-pack":
		advertiseRefs(w, r, repo)

	case "git-upload-pack":
		servePack(w, r, repo)
	}
}

func servePack(w http.ResponseWriter, r *http.Request, repo *repository) {
	if r.Method != http.MethodPost {
		common.BadRequest(errors.New("Only POST allowed for this method"), w)
		return
	}
	if r.Header.Get("Content-Type") != "application/x-git-upload-pack-request" {
		err := fmt.Errorf("unexpected Content-Type: %v", r.Header.Get("Content-Type"))
		common.BadRequest(err, w)
		return
	}
	cmd := exec.CommandContext(r.Context(), "git-upload-pack", "--stateless-rpc", repo.gitPath)
	cmd.Stdin = r.Body
	var buf bytes.Buffer
	cmd.Stdout = &buf
	err := cmd.Start()
	if err != nil {
		http.Error(w, fmt.Errorf("could not start command: %v", err).Error(), http.StatusInternalServerError)
		return
	}
	err = cmd.Wait()
	if ee, _ := err.(*exec.ExitError); ee != nil && ee.Sys().(syscall.WaitStatus).ExitStatus() == 128 {
		log.Info("The remote end hung up")
	} else if err != nil {
		log.Error(err, "git-upload-pack command failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/x-git-upload-pack-result")
	_, err = io.Copy(w, &buf)
	if err != nil {
		log.Error(err, "failed to write response to buffer")
	}
}

func advertiseRefs(w http.ResponseWriter, r *http.Request, repo *repository) {
	defer r.Body.Close()
	if r.Method != http.MethodGet {
		common.BadRequest(errors.New("Only GET allowed for this method"), w)
		return
	}
	cmd := exec.CommandContext(r.Context(), "git-upload-pack", "--advertise-refs", repo.gitPath)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	err := cmd.Start()
	if err != nil {
		http.Error(w, fmt.Errorf("could not start command: %v", err).Error(), http.StatusInternalServerError)
		return
	}
	err = cmd.Wait()
	if err != nil {
		log.Error(err, "git-upload-pack command failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/x-git-upload-pack-advertisement")
	_, err = io.WriteString(w, "001e# service=git-upload-pack\n0000")
	if err != nil {
		log.Error(err, "failed to write upload-pack header")
		return
	}
	_, err = io.Copy(w, &buf)
	if err != nil {
		log.Error(err, "failed to write response to buffer")
	}
}
