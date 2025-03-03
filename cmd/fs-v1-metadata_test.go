// Copyright (c) 2015-2021 MinIO, Inc.
//
// This file is part of MinIO Object Storage stack
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package cmd

import (
	"reflect"
	"testing"
)

// Tests ToObjectInfo function.
func TestFSV1MetadataObjInfo(t *testing.T) {
	fsMeta := newFSMetaV1()
	objInfo := fsMeta.ToObjectInfo("testbucket", "testobject", nil)
	if objInfo.Size != 0 {
		t.Fatal("Unexpected object info value for Size", objInfo.Size)
	}
	if !objInfo.ModTime.Equal(timeSentinel) {
		t.Fatal("Unexpected object info value for ModTime ", objInfo.ModTime)
	}
	if objInfo.IsDir {
		t.Fatal("Unexpected object info value for IsDir", objInfo.IsDir)
	}
	if !objInfo.Expires.IsZero() {
		t.Fatal("Unexpected object info value for Expires ", objInfo.Expires)
	}
}

func TestFSChecksumV1MarshalJSON(t *testing.T) {
	var cs FSChecksumInfoV1

	testCases := []struct {
		checksum       FSChecksumInfoV1
		expectedResult string
	}{
		{cs, `{"algorithm":"","blocksize":0,"hashes":null}`},
		{FSChecksumInfoV1{Algorithm: "highwayhash", Blocksize: 500}, `{"algorithm":"highwayhash","blocksize":500,"hashes":null}`},
		{FSChecksumInfoV1{Algorithm: "highwayhash", Blocksize: 10, Hashes: [][]byte{[]byte("hello")}}, `{"algorithm":"highwayhash","blocksize":10,"hashes":["68656c6c6f"]}`},
	}

	for _, testCase := range testCases {
		data, _ := testCase.checksum.MarshalJSON()
		if testCase.expectedResult != string(data) {
			t.Fatalf("expected: %v, got: %v", testCase.expectedResult, string(data))
		}
	}
}

func TestFSChecksumV1UnMarshalJSON(t *testing.T) {
	var cs FSChecksumInfoV1

	testCases := []struct {
		data           []byte
		expectedResult FSChecksumInfoV1
	}{
		{[]byte(`{"algorithm":"","blocksize":0,"hashes":null}`), cs},
		{[]byte(`{"algorithm":"highwayhash","blocksize":500,"hashes":null}`), FSChecksumInfoV1{Algorithm: "highwayhash", Blocksize: 500}},
		{[]byte(`{"algorithm":"highwayhash","blocksize":10,"hashes":["68656c6c6f"]}`), FSChecksumInfoV1{Algorithm: "highwayhash", Blocksize: 10, Hashes: [][]byte{[]byte("hello")}}},
	}

	for _, testCase := range testCases {
		err := (&cs).UnmarshalJSON(testCase.data)
		if err != nil {
			t.Fatal("Unexpected error during checksum unmarshalling ", err)
		}
		if !reflect.DeepEqual(testCase.expectedResult, cs) {
			t.Fatalf("expected: %v, got: %v", testCase.expectedResult, cs)
		}
	}
}
