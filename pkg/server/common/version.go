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

package common

import (
	"fmt"
	"sort"

	version "github.com/hashicorp/go-version"
	"github.com/tinyzimmer/protobuf-registry/pkg/protobuf"
)

func GetVersionFromProtoSlice(protos []*protobuf.Protobuf, version string) (*protobuf.Protobuf, error) {
	var protoInstance *protobuf.Protobuf
	for _, x := range protos {
		if *x.Version == version {
			protoInstance = x
		}
	}
	if protoInstance == nil {
		return nil, fmt.Errorf("version %s not found", version)
	}
	return protoInstance, nil
}

func GetLatestVersion(protos []*protobuf.Protobuf) (latest *protobuf.Protobuf) {
	// if there is only one, return it
	if len(protos) == 1 {
		latest = protos[0]
		return
	}
	var latestVersion *version.Version
	for _, proto := range protos {
		// if latest isn't set yet - set it to this instance and continue
		if latest == nil {
			var err error
			if latestVersion, err = version.NewVersion(*proto.Version); err != nil {
				continue
			}
			latest = proto
			continue
		}
		// check if this instance is newer than currently set latest
		if v, err := version.NewVersion(*proto.Version); err != nil {
			continue
		} else {
			if v.GreaterThan(latestVersion) {
				latestVersion = v
				latest = proto
			}
		}
	}

	return
}

func SortVersions(protos []*protobuf.Protobuf) (sorted []*protobuf.Protobuf) {

	// make slices to hold raw versions and invalid objects
	invalid := make([]*protobuf.Protobuf, 0)
	versions := make([]*version.Version, 0)

	// populate the versions and invalid slices
	for _, proto := range protos {
		v, err := version.NewVersion(*proto.Version)
		if err != nil {
			invalid = append(invalid, proto)
		} else {
			versions = append(versions, v)
		}
	}

	// sort the versions slice in descending order
	sort.Sort(sort.Reverse(version.Collection(versions)))

	// make a new slice for sorted results
	sorted = make([]*protobuf.Protobuf, 0)

	// rebuild proto slice from sorted version slice
	for _, v := range versions {
		for _, x := range protos {
			if v.Original() == *x.Version {
				sorted = append(sorted, x)
			}
		}
	}

	// add invalid versions to the end of the list
	sorted = append(sorted, invalid...)
	return
}
