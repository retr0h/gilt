//go:build integration
// +build integration

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
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/retr0h/go-gilt/internal/git"
	"github.com/retr0h/go-gilt/internal/repositories"
	"github.com/retr0h/go-gilt/internal/repository"
	helper "github.com/retr0h/go-gilt/internal/testing"
)

type RepositoryIntegrationTestSuite struct {
	suite.Suite
	r  repository.Repository
	rr repositories.Repositories
	g  *git.Git
}

func (suite *RepositoryIntegrationTestSuite) unmarshalYAML(data []byte) error {
	return helper.UnmarshalYAML([]byte(data), &suite.rr.Repositories)
}

func (suite *RepositoryIntegrationTestSuite) SetupTest() {
	suite.rr = repositories.Repositories{}
	suite.g = git.NewGit(suite.rr.Debug)
}

func (suite *RepositoryIntegrationTestSuite) TearDownTest() {
	helper.RemoveTempDirectory(suite.r.GiltDir)
}

func (suite *RepositoryIntegrationTestSuite) TestCopySourcesHasErrorWhenDstDirDoesNotExist() {
	data := `
- git: https://github.com/lorin/openstack-ansible-modules.git
  version: 2677cc3
  sources:
    - src: "*_manage"
      dstDir: invalid/path
`
	err := suite.unmarshalYAML([]byte(data))
	assert.NoError(suite.T(), err)

	r := suite.rr.Repositories[0]
	r.GiltDir = helper.CreateTempDirectory()
	err = suite.g.Clone(r)
	assert.NoError(suite.T(), err)

	err = r.CopySources()
	assert.Error(suite.T(), err)
}

func (suite *RepositoryIntegrationTestSuite) TestCopySourcesHasErrorWhenFileCopyFails() {
	tempDir := helper.CreateTempDirectory()
	dstDir := filepath.Join(tempDir, "library")
	data := fmt.Sprintf(`
- git: https://github.com/lorin/openstack-ansible-modules.git
  version: 2677cc3
  sources:
    - src: cinder_manage
      dstDir: %s
`, dstDir)
	err := suite.unmarshalYAML([]byte(data))
	assert.NoError(suite.T(), err)

	r := suite.rr.Repositories[0]
	r.GiltDir = tempDir

	err = suite.g.Clone(r)
	assert.NoError(suite.T(), err)

	os.Mkdir(dstDir, 0o755)

	originalCopyFile := repository.CopyFile
	repository.CopyFile = func(src string, dst string) error {
		return errors.New("Failed to copy file")
	}
	defer func() { repository.CopyFile = originalCopyFile }()

	err = r.CopySources()
	assert.Error(suite.T(), err)
}

func (suite *RepositoryIntegrationTestSuite) TestCopySourcesHasErrorWhenDstFileDoesNotExist() {
	data := `
- git: https://github.com/lorin/openstack-ansible-modules.git
  version: 2677cc3
  sources:
    - src: cinder_manage
      dstFile: invalid/path
`
	err := suite.unmarshalYAML([]byte(data))
	assert.NoError(suite.T(), err)

	r := suite.rr.Repositories[0]
	r.GiltDir = helper.CreateTempDirectory()

	err = suite.g.Clone(r)
	assert.NoError(suite.T(), err)

	err = r.CopySources()
	assert.Error(suite.T(), err)
}

func (suite *RepositoryIntegrationTestSuite) TestCopySourcesCopiesFile() {
	tempDir := helper.CreateTempDirectory()
	dstFile := filepath.Join(tempDir, "cinder_manage")
	dstDir := filepath.Join(tempDir, "library")
	dstDirFile := filepath.Join(tempDir, "library", "glance_manage")
	data := fmt.Sprintf(`
- git: https://github.com/lorin/openstack-ansible-modules.git
  version: 2677cc3
  sources:
    - src: cinder_manage
      dstFile: %s
    - src: glance_manage
      dstDir: %s
`, dstFile, dstDir)
	err := suite.unmarshalYAML([]byte(data))
	assert.NoError(suite.T(), err)

	r := suite.rr.Repositories[0]
	r.GiltDir = tempDir

	os.Mkdir(dstDir, 0o755)
	err = suite.g.Clone(r)
	assert.NoError(suite.T(), err)

	err = r.CopySources()
	assert.NoError(suite.T(), err)
	assert.FileExistsf(suite.T(), dstFile, "File does not exist")
	assert.FileExistsf(suite.T(), dstDirFile, "File does not exist")
}

func (suite *RepositoryIntegrationTestSuite) TestCopySourcesHasErrorWhenDirExistsAndDirCopyFails() {
	tempDir := helper.CreateTempDirectory()
	dstDir := filepath.Join(tempDir, "tests")
	data := fmt.Sprintf(`
- git: https://github.com/lorin/openstack-ansible-modules.git
  version: 2677cc3
  sources:
    - src: tests
      dstDir: %s
`, dstDir)
	err := suite.unmarshalYAML([]byte(data))
	assert.NoError(suite.T(), err)

	r := suite.rr.Repositories[0]
	r.GiltDir = tempDir

	err = suite.g.Clone(r)
	assert.NoError(suite.T(), err)

	originalCopyDir := repository.CopyDir
	repository.CopyDir = func(src string, dst string) error {
		return errors.New("Failed to copy dir")
	}
	defer func() { repository.CopyDir = originalCopyDir }()

	err = r.CopySources()
	assert.Error(suite.T(), err)
}

func (suite *RepositoryIntegrationTestSuite) TestCopySourcesCopiesDir() {
	tempDir := helper.CreateTempDirectory()
	dstDir := filepath.Join(tempDir, "tests")
	data := fmt.Sprintf(`
- git: https://github.com/lorin/openstack-ansible-modules.git
  version: 2677cc3
  sources:
    - src: tests
      dstDir: %s
`, dstDir)
	err := suite.unmarshalYAML([]byte(data))
	assert.NoError(suite.T(), err)

	r := suite.rr.Repositories[0]
	r.GiltDir = tempDir
	os.Mkdir(dstDir, 0o755) // execute the dstDir cleanup code prior to copy.

	err = suite.g.Clone(r)
	assert.NoError(suite.T(), err)

	err = r.CopySources()
	assert.NoError(suite.T(), err)
	assert.DirExistsf(suite.T(), dstDir, "Dir does not exist")
}

// In order for `go test` to run this suite, we need to create
// a normal test function and pass our suite to suite.Run.
func TestRepositoryIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryIntegrationTestSuite))
}
