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
	"flag"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/go-logr/glogr"
	"github.com/jessevdk/go-flags"
	"github.com/tinyzimmer/protobuf-registry/pkg/util"
)

var GlobalConfig *Config
var log = glogr.New()

func Init(ignoreProtoc bool) error {
	var err error
	if GlobalConfig, err = newConfig(ignoreProtoc); err != nil {
		return err
	}
	return nil
}

// SafeInit is used from unit tests to ignore protoc
func SafeInit() {
	_ = Init(true)
}

type Config struct {
	// Server Settings
	BindAddress  string `long:"bind-address" env:"BIND_ADDRESS" default:"0.0.0.0:8080" json:"bind_address" description:"The address and port for the server to listen on"`
	ReadTimeout  int    `long:"read-timeout" env:"READ_TIMEOUT" default:"15" json:"read_timeout" description:"The read timeout in seconds for web requests"`
	WriteTimeout int    `long:"write-timeout" env:"WRITE_TIMEOUT" default:"15" json:"write_timeout" description:"The write timeout in seconds for web requests"`

	// Protobuf Compilation Settings
	CompileTimeout  int    `long:"compile-timeout" env:"COMPILE_TIMEOUT" default:"10" json:"compile_timeout" description:"The timeout in seconds for protoc operations"`
	ProtocPath      string `long:"protoc-path" env:"PROTOC_PATH" default:"/usr/bin/protoc" json:"protoc_path" description:"The path to the protoc executable"`
	ProtocGenGoPath string `long:"protoc-gen-go-path" env:"PROTOC_GEN_GO_PATH" default:"/opt/proto-registry/bin/protoc-gen-go" json:"protoc_gen_go_path" description:"The path to the protoc-gen-go plugin"`
	// Database Settings
	DatabaseDriver string `long:"database-driver" env:"DATABASE_DRIVER" default:"memory" json:"database_driver" description:"The database driver to use"`

	// Storage Settings
	StorageDriver string `long:"storage-driver" env:"STORAGE_DRIVER" default:"file" json:"storage_driver" description:"The storage driver to use"`

	// File Storage settings
	FileStoragePath string `long:"file-storage-path" env:"FILE_STORAGE_PATH" default:"/opt/proto-registry/data" json:"file_storage_path" description:"The filepath used by the file storage driver and remote cache"`

	// Memory database settings
	PersistMemoryToDisk bool `long:"persist-memory" env:"PERSIST_MEMORY" json:"persist_memory" description:"Whether to persist the memory database to disk after write/delete operations"`

	// Pre-populate cache with remote dependencies
	PreCachedRemotes []string `long:"pre-cached-remotes" env-delim:"," env:"PRE_CACHED_REMOTES" json:"pre_cached_remotes" description:"Remote repositories to pre-cache as protoc imports"`

	// UI Settings
	RedirectNotFoundToUI bool `long:"ui-redirect-all" env:"UI_REDIRECT_ALL" json:"ui_redirect_all" description:"Redirect all unhandled requests to the UI instead of the default 404 handler"`
	CORSEnabled          bool `long:"enable-cors" env:"ENABLE_CORS" json:"cors_enabled" description:"Enable CORS headers for requests"`

	// protobuf version as detected at boot
	ProtobufVersion string `no-flag:"" json:"protobuf_version"`

	// Compile environment - add a short t to trick unit test parsing
	GoVersion   string `hidden:"true" short:"t" json:"go_version"`
	GitCommit   string `hidden:"true" json:"git_commit"`
	CompileDate string `hidden:"true" json:"compile_date"`
}

// JSON returns the json bytes for the config object
func (c *Config) JSON() []byte {
	out, _ := json.MarshalIndent(c, "", "  ")
	return append(out, "\n"...)
}

// newConfig parses the environment and cli flags for both go-flags and glogr
func newConfig(ignoreProtoc bool) (*Config, error) {
	c := &Config{}
	parser := flags.NewParser(c, flags.Default|flags.PassAfterNonOption)
	if _, err := parser.Parse(); err != nil {
		return nil, err
	}
	out, err := exec.Command(c.ProtocPath, "--version").CombinedOutput()
	if err != nil && !ignoreProtoc {
		// hack for --help somehow getting set over the default protoc path
		if !util.StringSliceContains(os.Args, "--help") {
			log.Error(err, "No protoc at the provided path")
		}
		return nil, err
	}
	c.ProtobufVersion = strings.TrimSpace(string(out))
	c.GoVersion = runtime.Version()
	// handle flags for logging
	handleLogFlags()
	return c, nil
}

// handleLogFlags wipes out previous flags so we can just configure glogr -
// yes this is rather hacky to just be able to use a logging library
func handleLogFlags() {
	os.Args = []string{os.Args[0]}
	if err := flag.Set("alsologtostderr", "true"); err != nil {
		panic(err)
	}
	flag.Parse()
}
