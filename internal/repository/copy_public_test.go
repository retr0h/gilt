// Copyright (c) 2018 John Dewey

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
// DEALINGS IN THE SOFTWARE.

package repository_test

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/retr0h/gilt/v2/internal/repository"
)

type CopyPublicTestSuite struct {
	suite.Suite

	appFs    afero.Fs
	cloneDir string
	dstDir   string
}

func (suite *CopyPublicTestSuite) NewTestCopyManager() repository.CopyManager {
	return repository.NewCopy(
		suite.appFs,
		slog.New(slog.NewTextHandler(os.Stdout, nil)),
	)
}

func (suite *CopyPublicTestSuite) SetupTest() {
	suite.appFs = afero.NewMemMapFs()
	suite.cloneDir = "/cloneDir"
	suite.dstDir = "/dstDir"
}

func (suite *CopyPublicTestSuite) TestCopyFileOk() {
	cm := suite.NewTestCopyManager()

	specs := []FileSpec{
		{
			appFs:   suite.appFs,
			srcDir:  filepath.Join(suite.cloneDir, "srcDir"),
			srcFile: filepath.Join(suite.cloneDir, "srcDir", "1.txt"),
		},
	}
	createFileSpecs(specs)

	assertFile := filepath.Join(suite.dstDir, "1.txt")
	err := cm.CopyFile(specs[0].srcFile, assertFile)
	assert.NoError(suite.T(), err)

	got, _ := afero.Exists(suite.appFs, assertFile)
	assert.True(suite.T(), got)
}

func (suite *CopyPublicTestSuite) TestCopyFileReturnsError() {
	cm := suite.NewTestCopyManager()

	assertFile := "/invalidSrc"
	err := cm.CopyFile("/invalidSrc", assertFile)
	assert.Error(suite.T(), err)

	got, _ := afero.Exists(suite.appFs, assertFile)
	assert.False(suite.T(), got)
}

func (suite *CopyPublicTestSuite) TestCopyFileReturnsErrorEPERM() {
	specs := []FileSpec{
		{
			appFs:   suite.appFs,
			srcDir:  filepath.Join(suite.cloneDir, "srcDir"),
			srcFile: filepath.Join(suite.cloneDir, "srcDir", "1.txt"),
		},
	}
	createFileSpecs(specs)

	assertFile := filepath.Join(suite.dstDir, "1.txt")
	// Replace the test FS with a read-only copy
	suite.appFs = afero.NewReadOnlyFs(suite.appFs)
	cm := suite.NewTestCopyManager()
	err := cm.CopyFile(specs[0].srcFile, assertFile)
	assert.Error(suite.T(), err)
}

func (suite *CopyPublicTestSuite) TestCopyDirOk() {
	cm := suite.NewTestCopyManager()

	specs := []FileSpec{
		{
			appFs:   suite.appFs,
			srcDir:  filepath.Join(suite.cloneDir, "srcDir"),
			srcFile: filepath.Join(suite.cloneDir, "srcDir", "1.txt"),
		},
	}
	createFileSpecs(specs)

	assertFile := filepath.Join(suite.dstDir, "1.txt")
	err := cm.CopyDir(specs[0].srcDir, suite.dstDir)
	assert.NoError(suite.T(), err)

	got, _ := afero.Exists(suite.appFs, assertFile)
	assert.True(suite.T(), got)
}

func (suite *CopyPublicTestSuite) TestCopyDirNestedOk() {
	cm := suite.NewTestCopyManager()

	specs := []FileSpec{
		{
			appFs:   suite.appFs,
			srcDir:  filepath.Join(suite.cloneDir, "srcDir", "subDir"),
			srcFile: filepath.Join(suite.cloneDir, "srcDir", "subDir", "1.txt"),
		},
	}
	createFileSpecs(specs)

	assertFile := filepath.Join(suite.dstDir, "subDir", "1.txt")
	err := cm.CopyDir(filepath.Join(suite.cloneDir, "srcDir"), suite.dstDir)
	assert.NoError(suite.T(), err)

	got, _ := afero.Exists(suite.appFs, assertFile)
	assert.True(suite.T(), got)
}

func (suite *CopyPublicTestSuite) TestCopyDirSymlinksOk() {
	suite.T().Skip("afero.MemMapFS does not support symlinks")
}

func (suite *CopyPublicTestSuite) TestCopyDirReturnsError() {
	cm := suite.NewTestCopyManager()

	assertDir := "/invalidDir"
	err := cm.CopyDir("/invalidSrc", assertDir)
	assert.Error(suite.T(), err)

	got, _ := afero.Exists(suite.appFs, assertDir)
	assert.False(suite.T(), got)
}

func (suite *CopyPublicTestSuite) TestCopyDirReturnsErrorEEXIST() {
	cm := suite.NewTestCopyManager()

	specs := []FileSpec{
		{
			appFs:   suite.appFs,
			srcDir:  filepath.Join(suite.cloneDir, "srcDir"),
			srcFile: filepath.Join(suite.cloneDir, "srcDir", "1.txt"),
		},
		{
			appFs:   suite.appFs,
			srcFile: filepath.Join(suite.dstDir, "1.txt"),
		},
	}
	createFileSpecs(specs)

	assertFile := filepath.Join(suite.dstDir, "1.txt")
	err := cm.CopyDir(specs[0].srcDir, suite.dstDir)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "destination already exists")

	got, _ := afero.Exists(suite.appFs, assertFile)
	assert.True(suite.T(), got)
}

func (suite *CopyPublicTestSuite) TestCopyDirReturnsErrorWhenSrcIsNotDir() {
	cm := suite.NewTestCopyManager()

	specs := []FileSpec{
		{
			appFs:   suite.appFs,
			srcDir:  filepath.Join(suite.cloneDir, "srcDir"),
			srcFile: filepath.Join(suite.cloneDir, "srcDir", "1.txt"),
		},
	}
	createFileSpecs(specs)

	assertFile := specs[0].srcFile
	err := cm.CopyDir(assertFile, suite.dstDir)
	assert.Error(suite.T(), err)
}

// In order for `go test` to run this suite, we need to create
// a normal test function and pass our suite to suite.Run.
func TestCopyPublicTestSuite(t *testing.T) {
	suite.Run(t, new(CopyPublicTestSuite))
}
