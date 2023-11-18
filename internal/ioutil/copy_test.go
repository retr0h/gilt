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

package ioutil

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type IOUtilTestSuite struct {
	suite.Suite
}

func (suite *IOUtilTestSuite) TestCopyFile() {
	type spec struct {
		appFs   afero.Fs
		baseDir string
		srcFile string
		dstFile string
	}

	specs := []spec{
		{
			appFs:   afero.NewMemMapFs(),
			baseDir: "/base",
			srcFile: "/base/srcFile",
			dstFile: "/base/dstFile",
		},
	}

	for _, s := range specs {
		_ = s.appFs.MkdirAll(s.baseDir, 0o755)
		_, err := s.appFs.Create(s.srcFile)
		assert.NoError(suite.T(), err)

		err = CopyFile(s.appFs, s.srcFile, s.dstFile)
		assert.NoError(suite.T(), err)

		got, _ := afero.Exists(s.appFs, s.dstFile)
		assert.True(suite.T(), got)
	}
}

func (suite *IOUtilTestSuite) TestCopyFileReturnsError() {
	type spec struct {
		appFs   afero.Fs
		srcFile string
		dstFile string
	}

	specs := []spec{
		{appFs: afero.NewMemMapFs(), srcFile: "srcFile", dstFile: "/invalid/path"},
	}

	for _, s := range specs {
		err := CopyFile(s.appFs, s.srcFile, s.dstFile)
		assert.Error(suite.T(), err)

		got, _ := afero.Exists(s.appFs, s.dstFile)
		assert.False(suite.T(), got)
	}
}

func (suite *IOUtilTestSuite) TestCopyDir() {
	type spec struct {
		appFs   afero.Fs
		srcDir  string
		srcFile string
		dstDir  string
		dstFile string
	}

	specs := []spec{
		{
			appFs:   afero.NewMemMapFs(),
			srcDir:  "/src",
			srcFile: "/src/srcFile",
			dstDir:  "/dst",
			dstFile: "/dst/dstFile",
		},
	}

	for _, s := range specs {
		_ = s.appFs.MkdirAll(s.srcDir, 0o755)
		_, err := s.appFs.Create(s.srcFile)
		assert.NoError(suite.T(), err)

		err = CopyFile(s.appFs, s.srcDir, s.dstDir)
		assert.NoError(suite.T(), err)

		got, _ := afero.Exists(s.appFs, s.dstFile)
		assert.False(suite.T(), got)
	}
}

func (suite *IOUtilTestSuite) TestCopyDirReturnsError() {
	type spec struct {
		appFs  afero.Fs
		srcDir string
		dstDir string
	}

	specs := []spec{
		{appFs: afero.NewMemMapFs(), srcDir: "/src", dstDir: "/dst"},
	}

	for _, s := range specs {
		err := CopyDir(s.appFs, s.srcDir, s.dstDir)
		assert.Error(suite.T(), err)

		got, _ := afero.DirExists(s.appFs, s.dstDir)
		assert.False(suite.T(), got)
	}
}

// In order for `go test` to run this suite, we need to create
// a normal test function and pass our suite to suite.Run.
func TestIOUtilTestSuite(t *testing.T) {
	suite.Run(t, new(IOUtilTestSuite))
}
