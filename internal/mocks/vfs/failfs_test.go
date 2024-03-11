// Copyright (c) 2023 Nicolas Simonds

//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//  	http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.
//

package failfs_test

import (
	"testing"

	"github.com/avfs/avfs"
	"github.com/avfs/avfs/test"
	"github.com/avfs/avfs/vfs/memfs"
	failfs "github.com/retr0h/gilt/v2/internal/mocks/vfs"
)

var (
	// Tests that failfs.failfs struct implements avfs.VFS interface.
	_ avfs.VFS = &failfs.FailFS{}

	// Tests that failfs.FailFile struct implements avfs.File interface.
	_ avfs.File = &failfs.FailFile{}
)

func initTest(t *testing.T) *test.Suite {
	vfsSetup := memfs.New()
	vfs := failfs.New(vfsSetup, nil)

	ts := test.NewSuiteFS(t, vfsSetup, vfs)

	return ts
}

func TestFailFS(t *testing.T) {
	ts := initTest(t)
	ts.TestVFSAll(t)
}

func TestFailFSConfig(t *testing.T) {
	vfsWrite := memfs.New()
	vfs := failfs.New(vfsWrite, nil)

	wantFeatures := vfs.Features() &^ avfs.FeatIdentityMgr
	if vfs.Features() != wantFeatures {
		t.Errorf("Features : want Features to be %s, got %s", wantFeatures, vfs.Features())
	}

	name := vfs.Name()
	if name != "" {
		t.Errorf("Name : want name to be empty, got %v", name)
	}

	osType := vfs.OSType()
	if osType != vfsWrite.OSType() {
		t.Errorf("OSType : want os type to be %v, got %v", vfsWrite.OSType(), osType)
	}
}
