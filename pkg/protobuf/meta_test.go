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

import "testing"

func TestFileContents(t *testing.T) {
	proto := newTestProtoWithData(t)
	contents, err := proto.Contents("TestProtoMessage.proto")
	if err != nil {
		t.Error("Expected no error from test data, got:", err)
	}
	if string(contents) != testProtoMessageRawString {
		t.Error("Data was malformed on retrieval")
	}

	// test nested file
	if _, err := proto.Contents("details/TestProtoMessageDetails.proto"); err != nil {
		t.Error("Expected no error on file, got:", err)
	}

	if _, err := proto.Contents("NonExistFile.proto"); err == nil {
		t.Error("Expected error from non-existent file, got nil")
	}

	proto.SetRaw(nil)
	if _, err := proto.Contents("TestProtoMessage.proto"); err == nil {
		t.Error("Expected error from no raw contents, got nil")
	}
}
