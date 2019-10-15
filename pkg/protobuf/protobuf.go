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

package protobuf

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/go-logr/glogr"
	"github.com/tinyzimmer/protobuf-registry/pkg/types"
)

var log = glogr.New()

// Protobuf represents a single versioned protobuf package. It implements methods
// for utility functions across different packages. Storage and Database providers
// are responsible for setting, storing, and retrieving the raw data.
type Protobuf struct {
	ID           *string                  `json:"id"`
	Name         *string                  `json:"name"`
	Version      *string                  `json:"version"`
	LastUpdated  time.Time                `json:"lastUpdated"`
	Dependencies []*types.ProtoDependency `json:"dependencies"`
	// raw zip bytes
	raw []byte
	// raw descriptor set bytes
	descriptor []byte
}

// NewFromRequest converts a PostProtoRequest to a bare protobuf object
func NewFromRequest(req *types.PostProtoRequest) *Protobuf {
	return &Protobuf{
		ID:           &req.ID,
		Name:         &req.Name,
		Version:      &req.Version,
		Dependencies: req.RemoteDepends,
	}
}

// SetRawFromBase64 sets the raw protobuf data from a base64 string
func (p *Protobuf) SetRawFromBase64(body string) error {
	var raw []byte
	var err error
	if raw, err = base64.StdEncoding.DecodeString(body); err != nil {
		return fmt.Errorf("Could not decode base64: %s", err.Error())
	}
	p.SetRaw(raw)
	return nil
}

// SetRaw sets the raw protobuf data from a byte slice
func (p *Protobuf) SetRaw(raw []byte) {
	p.raw = raw
}

// SetDescritptor sets the raw descriptor set from a byte slice
func (p *Protobuf) SetDescriptor(raw []byte) {
	p.descriptor = raw
}

// Raw returns the raw zip bytes of the protobuf object
func (p *Protobuf) Raw() []byte {
	return p.raw
}

// DescriptorBytes returns the descriptor set bytes for the protobuf object
func (p *Protobuf) DescriptorBytes() []byte {
	return p.descriptor
}

// RawFilename returns a zip filename for the protobuf object
func (p *Protobuf) RawFilename() string {
	return fmt.Sprintf("%s-%s.zip", *p.Name, *p.Version)
}

// RawReader returns a ReadSeeker for the raw zip data - useful for serving in http requests
func (p *Protobuf) RawReader() io.ReadSeeker {
	return bytes.NewReader(p.raw)
}

// DescriptorReader returns a ReadSeeker for the raw desciptor set
func (p *Protobuf) DescriptorReader() io.ReadSeeker {
	return bytes.NewReader(p.descriptor)
}

// SHA256 computes the sha256zum for the protobuf zip contents
func (p *Protobuf) SHA256() (string, error) {
	if p.Raw() == nil {
		return "", errors.New("raw zip is nil, need to call p.SetRaw() or p.SetRawFromBase64()")
	}
	h := sha256.New()
	if _, err := h.Write(p.Raw()); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
