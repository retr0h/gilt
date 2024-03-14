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
	"errors"
	"io/fs"
	"log/slog"
	"os"
	"testing"

	"github.com/avfs/avfs"
	"github.com/avfs/avfs/vfs/memfs"
	"github.com/avfs/avfs/vfs/rofs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	failfs "github.com/retr0h/gilt/v2/internal/mocks/vfs"
	"github.com/retr0h/gilt/v2/internal/repository"
)

type CopyPublicTestSuite struct {
	suite.Suite

	appFs    avfs.VFS
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
	suite.appFs = memfs.New()
	suite.cloneDir = "/cloneDir"
	suite.dstDir = "/dstDir"
}

func (suite *CopyPublicTestSuite) TestCopyFileOk() {
	cm := suite.NewTestCopyManager()

	specs := []FileSpec{
		{
			appFs:   suite.appFs,
			srcDir:  suite.appFs.Join(suite.cloneDir, "srcDir"),
			srcFile: suite.appFs.Join(suite.cloneDir, "srcDir", "1.txt"),
		},
	}
	createFileSpecs(specs)

	assertFile := suite.appFs.Join(suite.dstDir, "1.txt")
	err := cm.CopyFile(specs[0].srcFile, assertFile)
	assert.NoError(suite.T(), err)

	got, _ := avfs.Exists(suite.appFs, assertFile)
	assert.True(suite.T(), got)
}

func (suite *CopyPublicTestSuite) TestCopyFileReturnsError() {
	cm := suite.NewTestCopyManager()

	assertFile := "/invalidSrc"
	err := cm.CopyFile("/invalidSrc", assertFile)
	assert.Error(suite.T(), err)

	got, _ := avfs.Exists(suite.appFs, assertFile)
	assert.False(suite.T(), got)
}

func (suite *CopyPublicTestSuite) TestCopyFileErrorStat() {
	specs := []FileSpec{
		{
			appFs:   suite.appFs,
			srcDir:  suite.appFs.Join(suite.cloneDir, "srcDir"),
			srcFile: suite.appFs.Join(suite.cloneDir, "srcDir", "1.txt"),
		},
	}
	createFileSpecs(specs)
	// Make Stat() calls fail
	suite.appFs = failfs.New(
		suite.appFs,
		map[string]interface{}{
			"Stat": func(string) (fs.FileInfo, error) { return nil, errors.New("FailFS!") },
		},
	)
	assertFile := suite.appFs.Join(suite.dstDir, "1.txt")
	cm := suite.NewTestCopyManager()
	err := cm.CopyFile(specs[0].srcFile, assertFile)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "FailFS!", err.Error())
}

func (suite *CopyPublicTestSuite) TestCopyFileReturnsErrorEPERM() {
	specs := []FileSpec{
		{
			appFs:   suite.appFs,
			srcDir:  suite.appFs.Join(suite.cloneDir, "srcDir"),
			srcFile: suite.appFs.Join(suite.cloneDir, "srcDir", "1.txt"),
		},
	}
	createFileSpecs(specs)

	assertFile := suite.appFs.Join(suite.dstDir, "1.txt")
	// Replace the test FS with a read-only copy
	suite.appFs = rofs.New(suite.appFs)
	cm := suite.NewTestCopyManager()
	err := cm.CopyFile(specs[0].srcFile, assertFile)
	assert.Error(suite.T(), err)
}

func (suite *CopyPublicTestSuite) TestCopyFileErrorSettingDestfilePerms() {
	specs := []FileSpec{
		{
			appFs:   suite.appFs,
			srcDir:  suite.appFs.Join(suite.cloneDir, "srcDir"),
			srcFile: suite.appFs.Join(suite.cloneDir, "srcDir", "1.txt"),
		},
	}
	createFileSpecs(specs)
	// Make Chmod() calls fail
	suite.appFs = failfs.New(
		suite.appFs,
		map[string]interface{}{
			"Chmod": func(string, fs.FileMode) error { return errors.New("FailFS!") },
		},
	)
	assertFile := suite.appFs.Join(suite.dstDir, "1.txt")
	cm := suite.NewTestCopyManager()
	err := cm.CopyFile(specs[0].srcFile, assertFile)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "FailFS!", err.Error())
}

func (suite *CopyPublicTestSuite) TestCopyFileErrorCopy() {
	// Make the file unreadable
	suite.appFs = failfs.New(
		suite.appFs,
		map[string]interface{}{
			"file.Read": func([]byte) (int, error) { return 0, errors.New("FailFS!") },
		},
	)
	cm := suite.NewTestCopyManager()
	specs := []FileSpec{
		{
			appFs:   suite.appFs,
			srcDir:  suite.appFs.Join(suite.cloneDir, "srcDir"),
			srcFile: suite.appFs.Join(suite.cloneDir, "srcDir", "1.txt"),
		},
	}
	createFileSpecs(specs)
	assertFile := suite.appFs.Join(suite.dstDir, "1.txt")
	err := cm.CopyFile(specs[0].srcFile, assertFile)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "FailFS!", err.Error())
}

func (suite *CopyPublicTestSuite) TestCopyFileErrorSync() {
	specs := []FileSpec{
		{
			appFs:   suite.appFs,
			srcDir:  suite.appFs.Join(suite.cloneDir, "srcDir"),
			srcFile: suite.appFs.Join(suite.cloneDir, "srcDir", "1.txt"),
		},
	}
	createFileSpecs(specs)
	// Make Sync() calls fail
	suite.appFs = failfs.New(
		suite.appFs,
		map[string]interface{}{
			"file.Sync": func() error { return errors.New("FailFS!") },
		},
	)
	assertFile := suite.appFs.Join(suite.dstDir, "1.txt")
	cm := suite.NewTestCopyManager()
	err := cm.CopyFile(specs[0].srcFile, assertFile)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "FailFS!", err.Error())
}

func (suite *CopyPublicTestSuite) TestCopyFileErrorFinalizingDestfilePerms() {
	specs := []FileSpec{
		{
			appFs:   suite.appFs,
			srcDir:  suite.appFs.Join(suite.cloneDir, "srcDir"),
			srcFile: suite.appFs.Join(suite.cloneDir, "srcDir", "1.txt"),
		},
	}
	createFileSpecs(specs)
	// Make the second Chmod() call fail
	suite.appFs = failfs.New(
		suite.appFs,
		map[string]interface{}{
			"Chmod": func(name string, mode fs.FileMode) error {
				if mode == 0o600 {
					return nil
				}
				return errors.New("FailFS!")
			},
		},
	)
	assertFile := suite.appFs.Join(suite.dstDir, "1.txt")
	cm := suite.NewTestCopyManager()
	err := cm.CopyFile(specs[0].srcFile, assertFile)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "FailFS!", err.Error())
}

func (suite *CopyPublicTestSuite) TestCopyDirOk() {
	cm := suite.NewTestCopyManager()

	specs := []FileSpec{
		{
			appFs:   suite.appFs,
			srcDir:  suite.appFs.Join(suite.cloneDir, "srcDir"),
			srcFile: suite.appFs.Join(suite.cloneDir, "srcDir", "1.txt"),
		},
	}
	createFileSpecs(specs)

	assertFile := suite.appFs.Join(suite.dstDir, "1.txt")
	err := cm.CopyDir(specs[0].srcDir, suite.dstDir)
	assert.NoError(suite.T(), err)

	got, _ := avfs.Exists(suite.appFs, assertFile)
	assert.True(suite.T(), got)
}

func (suite *CopyPublicTestSuite) TestCopyDirNestedOk() {
	cm := suite.NewTestCopyManager()

	specs := []FileSpec{
		{
			appFs:   suite.appFs,
			srcDir:  suite.appFs.Join(suite.cloneDir, "srcDir", "subDir"),
			srcFile: suite.appFs.Join(suite.cloneDir, "srcDir", "subDir", "1.txt"),
		},
	}
	createFileSpecs(specs)

	assertFile := suite.appFs.Join(suite.dstDir, "subDir", "1.txt")
	err := cm.CopyDir(suite.appFs.Join(suite.cloneDir, "srcDir"), suite.dstDir)
	assert.NoError(suite.T(), err)

	got, _ := avfs.Exists(suite.appFs, assertFile)
	assert.True(suite.T(), got)
}

func (suite *CopyPublicTestSuite) TestCopyDirSymlinksOk() {
	cm := suite.NewTestCopyManager()

	specs := []FileSpec{
		{
			appFs:   suite.appFs,
			srcDir:  suite.appFs.Join(suite.cloneDir, "srcDir", "subDir"),
			srcFile: suite.appFs.Join(suite.cloneDir, "srcDir", "subDir", "1.txt"),
		},
	}
	createFileSpecs(specs)
	_ = suite.appFs.Symlink(
		suite.appFs.Join(suite.cloneDir, "srcDir"),
		suite.appFs.Join(suite.cloneDir, "srcLink"),
	)

	assertFile := suite.appFs.Join(suite.dstDir, "subDir", "1.txt")
	err := cm.CopyDir(suite.appFs.Join(suite.cloneDir, "srcLink"), suite.dstDir)
	assert.NoError(suite.T(), err)

	got, _ := avfs.Exists(suite.appFs, assertFile)
	assert.True(suite.T(), got)
}

func (suite *CopyPublicTestSuite) TestCopyDirReturnsError() {
	cm := suite.NewTestCopyManager()

	assertDir := "/invalidDir"
	err := cm.CopyDir("/invalidSrc", assertDir)
	assert.Error(suite.T(), err)

	got, _ := avfs.Exists(suite.appFs, assertDir)
	assert.False(suite.T(), got)
}

func (suite *CopyPublicTestSuite) TestCopyDirReturnsErrorEEXIST() {
	cm := suite.NewTestCopyManager()

	specs := []FileSpec{
		{
			appFs:   suite.appFs,
			srcDir:  suite.appFs.Join(suite.cloneDir, "srcDir"),
			srcFile: suite.appFs.Join(suite.cloneDir, "srcDir", "1.txt"),
		},
		{
			appFs:   suite.appFs,
			srcDir:  suite.dstDir,
			srcFile: suite.appFs.Join(suite.dstDir, "1.txt"),
		},
	}
	createFileSpecs(specs)

	assertFile := suite.appFs.Join(suite.dstDir, "1.txt")
	err := cm.CopyDir(specs[0].srcDir, suite.dstDir)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "destination already exists")

	got, _ := avfs.Exists(suite.appFs, assertFile)
	assert.True(suite.T(), got)
}

func (suite *CopyPublicTestSuite) TestCopyDirReturnsErrorWhenSrcIsNotDir() {
	cm := suite.NewTestCopyManager()

	specs := []FileSpec{
		{
			appFs:   suite.appFs,
			srcDir:  suite.appFs.Join(suite.cloneDir, "srcDir"),
			srcFile: suite.appFs.Join(suite.cloneDir, "srcDir", "1.txt"),
		},
	}
	createFileSpecs(specs)

	assertFile := specs[0].srcFile
	err := cm.CopyDir(assertFile, suite.dstDir)
	assert.Error(suite.T(), err)
}

func (suite *CopyPublicTestSuite) TestCopyDirReturnsErrorCreatingDestDir() {
	specs := []FileSpec{
		{
			appFs:   suite.appFs,
			srcDir:  suite.appFs.Join(suite.cloneDir, "srcDir"),
			srcFile: suite.appFs.Join(suite.cloneDir, "srcDir", "1.txt"),
		},
	}
	createFileSpecs(specs)
	suite.appFs = failfs.New(
		suite.appFs,
		map[string]interface{}{
			"MkdirAll": func(string, fs.FileMode) error { return errors.New("FailFS!") },
		},
	)
	cm := suite.NewTestCopyManager()
	err := cm.CopyDir(specs[0].srcDir, suite.dstDir)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "FailFS!", err.Error())
}

func (suite *CopyPublicTestSuite) TestCopyDirReturnsErrorReadingSrcDir() {
	specs := []FileSpec{
		{
			appFs:   suite.appFs,
			srcDir:  suite.appFs.Join(suite.cloneDir, "srcDir"),
			srcFile: suite.appFs.Join(suite.cloneDir, "srcDir", "1.txt"),
		},
	}
	createFileSpecs(specs)
	suite.appFs = failfs.New(
		suite.appFs,
		map[string]interface{}{
			"ReadDir": func(string) ([]fs.DirEntry, error) { return nil, errors.New("FailFS!") },
		},
	)
	cm := suite.NewTestCopyManager()
	err := cm.CopyDir(specs[0].srcDir, suite.dstDir)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "FailFS!", err.Error())
}

func (suite *CopyPublicTestSuite) TestCopyDirReturnsErrorCheckingSrcFile() {
	specs := []FileSpec{
		{
			appFs:   suite.appFs,
			srcDir:  suite.appFs.Join(suite.cloneDir, "srcDir"),
			srcFile: suite.appFs.Join(suite.cloneDir, "srcDir", "1.txt"),
		},
	}
	createFileSpecs(specs)
	// FailFS cannot fall-through, so preload the responses it is expected to give
	srcDirStat, _ := suite.appFs.Stat(specs[0].srcDir)
	suite.appFs = failfs.New(
		suite.appFs,
		map[string]interface{}{
			"Stat": func(path string) (fs.FileInfo, error) {
				if path == specs[0].srcDir {
					return srcDirStat, nil
				}
				return nil, errors.New("FailFS!")
			},
		},
	)
	cm := suite.NewTestCopyManager()
	err := cm.CopyDir(specs[0].srcDir, suite.dstDir)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "FailFS!", err.Error())
}

func (suite *CopyPublicTestSuite) TestCopyDirNestedReturnsErrorOnCopyDir() {
	specs := []FileSpec{
		{
			appFs:   suite.appFs,
			srcDir:  suite.appFs.Join(suite.cloneDir, "srcDir", "subDir"),
			srcFile: suite.appFs.Join(suite.cloneDir, "srcDir", "subDir", "1.txt"),
		},
		{
			appFs:  suite.appFs,
			srcDir: suite.dstDir,
		},
	}
	createFileSpecs(specs)
	dstDirHandle, _ := suite.appFs.Open(suite.dstDir)
	suite.appFs = failfs.New(
		suite.appFs,
		map[string]interface{}{
			"Open": func(path string) (avfs.File, error) {
				if path == suite.dstDir {
					return dstDirHandle, errors.New("FailFS!")
				}
				return dstDirHandle, nil
			},
		},
	)
	cm := suite.NewTestCopyManager()
	err := cm.CopyDir(suite.appFs.Join(suite.cloneDir, "srcDir"), suite.dstDir)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "destination already exists", err.Error())
}

func (suite *CopyPublicTestSuite) TestCopyDirNestedReturnsErrorOnCopyFile() {
	specs := []FileSpec{
		{
			appFs:   suite.appFs,
			srcDir:  suite.appFs.Join(suite.cloneDir, "srcDir", "subDir"),
			srcFile: suite.appFs.Join(suite.cloneDir, "srcDir", "subDir", "1.txt"),
		},
	}
	createFileSpecs(specs)
	suite.appFs = failfs.New(
		suite.appFs,
		map[string]interface{}{
			"file.Read": func([]byte) (int, error) { return 0, errors.New("FailFS!") },
		},
	)
	cm := suite.NewTestCopyManager()
	err := cm.CopyDir(suite.appFs.Join(suite.cloneDir, "srcDir"), suite.dstDir)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "FailFS!", err.Error())
}

// In order for `go test` to run this suite, we need to create
// a normal test function and pass our suite to suite.Run.
func TestCopyPublicTestSuite(t *testing.T) {
	suite.Run(t, new(CopyPublicTestSuite))
}
