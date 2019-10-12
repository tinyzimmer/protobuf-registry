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

package mvnrouter

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/tinyzimmer/protobuf-registry/pkg/config"
	"github.com/tinyzimmer/protobuf-registry/pkg/protobuf"
	"github.com/tinyzimmer/protobuf-registry/pkg/server/common"
	"github.com/tinyzimmer/protobuf-registry/pkg/util"
)

type mavenServer struct {
	ctrl *common.ServerController
}

func RegisterRoutes(router *mux.Router, path string, ctrl *common.ServerController) {
	mvn := &mavenServer{ctrl}
	mvnRouter := router.PathPrefix(path).Subrouter()
	mvnRouter.HandleFunc("/{name}/{version}", mvn.getMavenPkgBoilerplate).Methods("GET")
}

var pomXMLTemplateString = `
<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
  xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd">
  <modelVersion>4.0.0</modelVersion>

  <groupId>{{ .JavaPackage }}</groupId>
  <artifactId>{{ .Name }}</artifactId>
  <version>{{ .Version }}-RELEASE</version>

  <properties>
    <maven.compiler.source>1.7</maven.compiler.source>
    <maven.compiler.target>1.7</maven.compiler.target>
  </properties>

  <dependencies>
    <dependency>
      <groupId>com.google.protobuf</groupId>
      <artifactId>protobuf-java</artifactId>
      <version>{{ .ProtocVersion }}</version>
    </dependency>
  </dependencies>
</project>
`
var pomXMLTemplate = template.Must(template.New("pom.xml").Parse(pomXMLTemplateString))

func getProtocVersionStr() string {
	split := strings.Split(config.GlobalConfig.ProtobufVersion, " ")
	return split[len(split)-1]
}

func getPathPrefix(in *protobuf.Protobuf) string {
	return filepath.Join(*in.Name, "src", "main", "java")
}

func addPomXML(pkgs []*string, proto *protobuf.Protobuf, path string) error {
	if len(pkgs) > 1 {
		return errors.New("Can only do single java-package maven archives")
	} else if len(pkgs) == 0 {
		return errors.New("Cannot build maven archive for protobufs without option java_package")
	}
	file, err := os.Create(filepath.Join(path, "pom.xml"))
	if err != nil {
		return err
	}
	defer file.Close()
	return pomXMLTemplate.Execute(file, map[string]interface{}{
		"Name":          *proto.Name,
		"Version":       *proto.Version,
		"JavaPackage":   *pkgs[0],
		"ProtocVersion": getProtocVersionStr(),
	})
}

func (mvn *mavenServer) getMavenPkgBoilerplate(w http.ResponseWriter, r *http.Request) {
	name := common.GetName(r)
	version := common.GetVersion(r)
	protos, err := mvn.ctrl.DB().GetProtoVersions(name)
	if err != nil {
		common.BadRequest(err, w)
		return
	}
	var proto *protobuf.Protobuf
	if proto, err = common.GetVersionFromProtoSlice(protos, version); err != nil {
		common.BadRequest(err, w)
		return
	}
	if proto, err = mvn.ctrl.Storage().GetRawProto(proto); err != nil {
		common.BadRequest(err, w)
		return
	}
	descriptors, err := proto.Descriptors()
	if err != nil {
		common.BadRequest(err, w)
		return
	}
	path, rm, err := proto.GenerateTo(protobuf.GenerateTargetJava, getPathPrefix(proto))
	if err != nil {
		common.BadRequest(err, w)
		return
	}
	defer rm()
	if err := addPomXML(descriptors.JavaPackages, proto, filepath.Join(path, *proto.Name)); err != nil {
		common.BadRequest(err, w)
	}
	archive, err := util.NewTarGZArchive(path)
	if err != nil {
		common.BadRequest(err, w)
		return
	}
	filename := fmt.Sprintf("%s-%s-maven.tar.gz", name, version)
	common.ServeFile(w, r, filename, bytes.NewReader(archive))
}
