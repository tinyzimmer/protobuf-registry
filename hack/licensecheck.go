package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var licenseHeader = `// Copyright © 2019 tinyzimmer
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
`

func main() {
	missing := make([]string, 0)
	for _, dir := range []string{"pkg/", "cmd/"} {
		if err := filepath.Walk(dir, func(path string, info os.FileInfo, e error) error {
			if e != nil {
				return e
			}
			if info.IsDir() {
				return nil
			}
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			if !strings.HasPrefix(string(data), licenseHeader) {
				missing = append(missing, path)
			}
			return nil
		}); err != nil {
			panic(err)
		}
	}

	if len(missing) > 0 {
		for _, x := range missing {
			log.Printf("%s is missing the license header", x)
		}
		os.Exit(1)
	} else {
		log.Println("All source files contain the license header")
	}
}
