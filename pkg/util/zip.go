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

package util

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func WriteZipFilesToTempDir(files []*zip.File) (path string, dirs map[string][]os.FileInfo, err error) {
	path, err = ioutil.TempDir("", "")
	if err != nil {
		return
	}
	dirs = make(map[string][]os.FileInfo)
	for _, f := range files {
		fpath := filepath.Join(path, f.Name)

		dirName := filepath.Dir(fpath)
		if _, ok := dirs[dirName]; !ok {
			dirs[dirName] = make([]os.FileInfo, 0)
		}

		dirs[dirName] = append(dirs[dirName], f.FileInfo())

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(path)+string(os.PathSeparator)) {
			return path, dirs, fmt.Errorf("%s: illegal file path", fpath)
		}

		if f.FileInfo().IsDir() {
			// Make Folder
			if err := os.MkdirAll(fpath, os.ModePerm); err != nil {
				return path, dirs, err
			}
			continue
		}
		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return path, dirs, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return path, dirs, err
		}

		rc, err := f.Open()
		if err != nil {
			return path, dirs, err
		}

		if _, err = io.Copy(outFile, rc); err != nil {
			return path, dirs, err
		}

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()
	}
	return
}
