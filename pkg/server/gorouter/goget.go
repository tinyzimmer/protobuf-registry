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
	"html/template"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-logr/glogr"
	"github.com/tinyzimmer/protobuf-registry/pkg/protobuf"
	"github.com/tinyzimmer/protobuf-registry/pkg/server/common"
)

var log = glogr.New()

var funcMap = template.FuncMap{
	"url": func(s string) template.URL {
		return template.URL(s)
	},
}

var getMetaTemplString = `
<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
<meta name="go-import" content="{{ .RequestedPath }} git {{ .GitPath | url }}">
</head>
<body>
Nothing to see here
</body>
</html>`

var getMetaTemplate = template.Must(template.New("go_get").Funcs(funcMap).Parse(getMetaTemplString))

func (gosrv *goServer) getPackage(w http.ResponseWriter, r *http.Request) {
	stripArg := strings.Replace(r.RequestURI, "?go-get=1", "", 1)
	pkgName := fmt.Sprintf("%s%s", r.Host, stripArg)

	pkgs, err := gosrv.ctrl.DB().GetAllProtoVersions()
	if err != nil {
		common.BadRequest(err, w)
		return
	}

	pkgMatch, err := gosrv.getLatestMatchingGoPackage(pkgs, pkgName)
	if err != nil {
		common.BadRequest(err, w)
		return
	}

	repo := gosrv.gitRouter.getRepo(*pkgMatch.Name)
	if repo == nil {
		log.Info("Generating go code for package", "name", pkgMatch.Name)
		out, rm, err := pkgMatch.GenerateTo(protobuf.GenerateTargetGo, "")
		if err != nil {
			common.BadRequest(err, w)
			return
		}
		codeDir := filepath.Join(out, pkgName)
		defer rm()
		repo, err = gosrv.gitRouter.newRepoFromPath(r, *pkgMatch.Name, pkgName)
		if err != nil {
			common.BadRequest(err, w)
			return
		}
		if err := gosrv.gitRouter.addToRepo(codeDir, repo); err != nil {
			common.BadRequest(err, w)
			return
		}
	}

	if err := getMetaTemplate.Execute(w, map[string]string{
		"RequestedPath": pkgName,
		"GitPath":       repo.repoURL,
	}); err != nil {
		common.BadRequest(err, w)
	}

}

func getGitScheme(r *http.Request) string {
	var gitScheme string
	if r.TLS == nil {
		gitScheme = "http"
	} else {
		gitScheme = "https"
	}
	return gitScheme
}

func (gosrv *goServer) getLatestMatchingGoPackage(pkgs map[string][]*protobuf.Protobuf, goPkg string) (*protobuf.Protobuf, error) {
	for _, protos := range pkgs {
		sorted := common.SortVersions(protos)
		latest := sorted[0]
		var err error
		if latest, err = gosrv.ctrl.Storage().GetRawProto(latest); err != nil {
			return nil, err
		}
		desc, err := latest.Descriptors()
		if err != nil {
			return nil, err
		}
		if len(desc.GoPackages) > 0 {
			for _, pkg := range desc.GoPackages {
				log.Info("has go package", "name", latest.Name, "package", *pkg)
				if *pkg == goPkg {
					log.Info("Have go package match", "name", latest.Name)
					return latest, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("Could not match %s to any packages", goPkg)
}
