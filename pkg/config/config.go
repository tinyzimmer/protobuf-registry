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

package config

import (
	"encoding/json"
	"os/exec"
	"runtime"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

const (
	envPrefix = "proto_registry"
)

var GlobalConfig *Config

func init() {
	var err error
	if GlobalConfig, err = newConfig(); err != nil {
		panic(err)
	}
}

type Config struct {
	// Server Settings
	BindAddress  string `envconfig:"bind_address" default:"0.0.0.0:8080" json:"bind_address"`
	ReadTimeout  int    `envconfig:"read_timeout" default:"15" json:"read_timeout"`
	WriteTimeout int    `envconfig:"write_timeout" default:"15" json:"write_timeout"`

	// Protobuf Compilation Settings
	CompileTimeout int    `envconfig:"compile_timeout" default:"10" json:"compile_timeout"`
	ProtocPath     string `envconfig:"protoc_path" default:"/usr/bin/protoc" json:"protoc_path"`

	// Database Settings
	DatabaseDriver string `envconfig:"database_driver" default:"memory" json:"database_driver"`

	// Storage Settings
	StorageDriver string `envconfig:"storage_driver" default:"file" json:"storage_driver"`

	// File Storage settings
	FileStoragePath string `envconfig:"file_storage_path" default:"/data" json:"file_storage_path"`

	// Memory database settings
	PersistMemoryToDisk bool `envconfig:"persist_memory" default:"false" json:"persist_memory"`

	// Pre-populate cache with remote dependencies
	PreCachedRemotes []string `envconfig:"pre_cached_remotes" json:"pre_cached_remotes"`

	// UI Settings
	RedirectNotFoundToUI bool `envconfig:"ui_redirect_all" default:"true" json:"ui_redirect_all"`
	CORSEnabled          bool `envconfig:"enable_cors" default:"false" json:"cors_enabled"`

	// protobuf version as detected at boot
	ProtobufVersion string `ignored:"true" json:"protobuf_version"`

	// Compile environment
	GoVersion   string `ignored:"true" json:"go_version"`
	GitCommit   string `ignored:"true" json:"git_commit"`
	CompileDate string `ignored:"true" json:"compile_date"`
}

func (c *Config) JSON() []byte {
	out, _ := json.MarshalIndent(c, "", "  ")
	return append(out, "\n"...)
}

func newConfig() (*Config, error) {
	c := &Config{}
	err := envconfig.Process(envPrefix, c)
	if err != nil {
		return nil, err
	}
	out, err := exec.Command(c.ProtocPath, "--version").CombinedOutput()
	if err != nil {
		return nil, err
	}
	c.ProtobufVersion = strings.TrimSpace(string(out))
	c.GoVersion = runtime.Version()
	return c, nil
}
