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
	dbcommon "github.com/tinyzimmer/proto-registry/pkg/database/common"
	storagecommon "github.com/tinyzimmer/proto-registry/pkg/storage/common"
)

type ServerController struct {
	db      dbcommon.DBEngine
	storage storagecommon.Provider
}

func (s *ServerController) SetDBEngine(db dbcommon.DBEngine) {
	s.db = db
}

func (s *ServerController) SetStorageProvider(storage storagecommon.Provider) {
	s.storage = storage
}

func (s *ServerController) DB() dbcommon.DBEngine {
	return s.db
}

func (s *ServerController) Storage() storagecommon.Provider {
	return s.storage
}
