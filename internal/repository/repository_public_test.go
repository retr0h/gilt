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
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/retr0h/go-gilt/internal"
	"github.com/retr0h/go-gilt/internal/git"
	"github.com/retr0h/go-gilt/internal/repository"
	"github.com/retr0h/go-gilt/pkg/config"
)

type RepositoryPublicTestSuite struct {
	suite.Suite

	ctrl            *gomock.Controller
	mockGit         *git.MockGitManager
	mockCopyManager *repository.MockCopyManager

	appFs    afero.Fs
	cloneDir string
	dstDir   string
	gitURL   string
	gitSHA   string
	gitTag   string
	logger   *slog.Logger
}

func (suite *RepositoryPublicTestSuite) NewRepositoryManager() internal.RepositoryManager {
	return repository.New(
		suite.appFs,
		suite.mockCopyManager,
		suite.mockGit,
		suite.logger,
	)
}

func (suite *RepositoryPublicTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockGit = git.NewMockGitManager(suite.ctrl)
	suite.mockCopyManager = repository.NewMockCopyManager(suite.ctrl)
	defer suite.ctrl.Finish()

	suite.appFs = afero.NewMemMapFs()
	suite.cloneDir = "/cloneDir"
	suite.dstDir = "/dstDir"
	suite.gitURL = "https://example.com/user/repo.git"
	suite.gitSHA = "abc123"
	suite.gitTag = "v1.1"
	suite.logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
}

func (suite *RepositoryPublicTestSuite) TestCloneOk() {
	repo := suite.NewRepositoryManager()

	c := config.Repository{
		Git: suite.gitURL,
		SHA: suite.gitSHA,
	}
	targetDir := filepath.Join(suite.cloneDir, filepath.Base(c.Git))

	gomock.InOrder(
		suite.mockGit.EXPECT().Clone(suite.gitURL, targetDir).Return(nil),
	)

	_, err := repo.Clone(c, suite.cloneDir)
	assert.NoError(suite.T(), err)
}

func (suite *RepositoryPublicTestSuite) TestCloneReturnsErrorWhenCloneErrors() {
	repo := suite.NewRepositoryManager()

	c := config.Repository{
		Git: suite.gitURL,
		SHA: suite.gitSHA,
	}

	errors := errors.New("tests error")
	gomock.InOrder(
		suite.mockGit.EXPECT().Clone(gomock.Any(), gomock.Any()).Return(errors),
	)

	_, err := repo.Clone(c, suite.cloneDir)
	assert.Error(suite.T(), err)
}

func (suite *RepositoryPublicTestSuite) TestCloneDoesNotCloneWhenCloneDirExists() {
	repo := suite.NewRepositoryManager()

	c := config.Repository{
		Git: suite.gitURL,
		SHA: suite.gitSHA,
	}
	targetDir := filepath.Join(suite.cloneDir, filepath.Base(c.Git))

	_ = suite.appFs.MkdirAll(targetDir, 0o755)

	_, err := repo.Clone(c, suite.cloneDir)
	assert.NoError(suite.T(), err)
}

func (suite *RepositoryPublicTestSuite) TestCopySourcesOkWhenSourceIsDirAndDstDirDoesNotExist() {
	repo := suite.NewRepositoryManager()
	specs := []FileSpec{
		{
			appFs:   suite.appFs,
			srcDir:  filepath.Join(suite.cloneDir, "subDir"),
			srcFile: filepath.Join(suite.cloneDir, "subDir", "1.txt"),
		},
	}
	createFileSpecs(specs)

	c := config.Repository{
		Git: suite.gitURL,
		SHA: suite.gitSHA,
		Sources: []config.Source{
			{
				Src:    filepath.Base(specs[0].srcDir),
				DstDir: suite.dstDir,
			},
		},
	}

	suite.mockCopyManager.EXPECT().
		CopyDir(filepath.Join(suite.cloneDir, c.Sources[0].Src), c.Sources[0].DstDir).
		Return(nil)

	err := repo.CopySources(c, suite.cloneDir)
	assert.NoError(suite.T(), err)
}

func (suite *RepositoryPublicTestSuite) TestCopySourcesReturnsErrorWhenSourceIsDirAndDstDirDoesNotExistAndCopyDirErrors() {
	repo := suite.NewRepositoryManager()
	specs := []FileSpec{
		{
			appFs:   suite.appFs,
			srcDir:  filepath.Join(suite.cloneDir, "subDir"),
			srcFile: filepath.Join(suite.cloneDir, "subDir", "1.txt"),
		},
	}
	createFileSpecs(specs)

	c := config.Repository{
		Git: suite.gitURL,
		SHA: suite.gitSHA,
		Sources: []config.Source{
			{
				Src:    filepath.Base(specs[0].srcDir),
				DstDir: suite.dstDir,
			},
		},
	}

	errors := errors.New("tests error")
	suite.mockCopyManager.EXPECT().CopyDir(gomock.Any(), gomock.Any()).Return(errors)

	err := repo.CopySources(c, suite.cloneDir)
	assert.Error(suite.T(), err)
}

func (suite *RepositoryPublicTestSuite) TestCopySourcesOkWhenSourceIsDirAndDstDirExists() {
	repo := suite.NewRepositoryManager()
	specs := []FileSpec{
		{
			appFs:   suite.appFs,
			srcDir:  filepath.Join(suite.cloneDir, "subDir"),
			srcFile: filepath.Join(suite.cloneDir, "subDir", "1.txt"),
		},
	}
	createFileSpecs(specs)

	c := config.Repository{
		Git: suite.gitURL,
		SHA: suite.gitSHA,
		Sources: []config.Source{
			{
				Src:    filepath.Base(specs[0].srcDir),
				DstDir: suite.dstDir,
			},
		},
	}

	suite.mockCopyManager.EXPECT().
		CopyDir(filepath.Join(suite.cloneDir, c.Sources[0].Src), c.Sources[0].DstDir).
		Return(nil)

	// create dstDir
	_ = suite.appFs.MkdirAll(suite.dstDir, 0o755)
	err := repo.CopySources(c, suite.cloneDir)
	assert.NoError(suite.T(), err)
}

func (suite *RepositoryPublicTestSuite) TestCopySourcesOkWhenSourceIsFilesAndDstDir() {
	repo := suite.NewRepositoryManager()
	specs := []FileSpec{
		{
			appFs:  suite.appFs,
			srcDir: filepath.Join(suite.cloneDir, "subDir"),
			srcFiles: []string{
				filepath.Join(suite.cloneDir, "subDir", "1.txt"),
				filepath.Join(suite.cloneDir, "subDir", "cinder_manage"),
				filepath.Join(suite.cloneDir, "subDir", "nova_manage"),
				filepath.Join(suite.cloneDir, "subDir", "glance_manage"),
			},
		},
		{
			appFs:  suite.appFs,
			srcDir: suite.cloneDir,
			srcFiles: []string{
				filepath.Join(suite.cloneDir, "1.txt"),
				filepath.Join(suite.cloneDir, "cinder_manage"),
				filepath.Join(suite.cloneDir, "nova_manage"),
				filepath.Join(suite.cloneDir, "glance_manage"),
			},
		},
	}
	createFileSpecs(specs)

	c := config.Repository{
		Git: suite.gitURL,
		SHA: suite.gitSHA,
		Sources: []config.Source{
			{
				Src:    "subDir/*_manage",
				DstDir: suite.dstDir,
			},
			{
				Src:    "*_manage",
				DstDir: suite.dstDir,
			},
		},
	}

	suite.mockCopyManager.EXPECT().
		CopyFile(filepath.Join(suite.cloneDir, "subDir", "cinder_manage"), filepath.Join(suite.dstDir, "cinder_manage")).
		Return(nil)
	suite.mockCopyManager.EXPECT().
		CopyFile(filepath.Join(suite.cloneDir, "subDir", "glance_manage"), filepath.Join(suite.dstDir, "glance_manage")).
		Return(nil)
	suite.mockCopyManager.EXPECT().
		CopyFile(filepath.Join(suite.cloneDir, "subDir", "nova_manage"), filepath.Join(suite.dstDir, "nova_manage")).
		Return(nil)

	suite.mockCopyManager.EXPECT().
		CopyFile(filepath.Join(suite.cloneDir, "cinder_manage"), filepath.Join(suite.dstDir, "cinder_manage")).
		Return(nil)
	suite.mockCopyManager.EXPECT().
		CopyFile(filepath.Join(suite.cloneDir, "glance_manage"), filepath.Join(suite.dstDir, "glance_manage")).
		Return(nil)
	suite.mockCopyManager.EXPECT().
		CopyFile(filepath.Join(suite.cloneDir, "nova_manage"), filepath.Join(suite.dstDir, "nova_manage")).
		Return(nil)

	err := repo.CopySources(c, suite.cloneDir)
	assert.NoError(suite.T(), err)
}

func (suite *RepositoryPublicTestSuite) TestCopySourcesReturnsErrorWhenSourceIsFilesAndDstDirAndCopyFileErrors() {
	repo := suite.NewRepositoryManager()
	specs := []FileSpec{
		{
			appFs:  suite.appFs,
			srcDir: filepath.Join(suite.cloneDir, "subDir"),
			srcFiles: []string{
				filepath.Join(suite.cloneDir, "subDir", "1.txt"),
			},
		},
	}
	createFileSpecs(specs)

	c := config.Repository{
		Git: suite.gitURL,
		SHA: suite.gitSHA,
		Sources: []config.Source{
			{
				Src:    "subDir/*.txt",
				DstDir: suite.dstDir,
			},
		},
	}

	errors := errors.New("tests error")
	suite.mockCopyManager.EXPECT().CopyFile(gomock.Any(), gomock.Any()).Return(errors)

	err := repo.CopySources(c, suite.cloneDir)
	assert.Error(suite.T(), err)
}

func (suite *RepositoryPublicTestSuite) TestCopySourcesOkWhenSourceIsFileAndDstFile() {
	repo := suite.NewRepositoryManager()
	specs := []FileSpec{
		{
			appFs:   suite.appFs,
			srcDir:  filepath.Join(suite.cloneDir),
			srcFile: filepath.Join(suite.cloneDir, "1.txt"),
		},
	}
	createFileSpecs(specs)

	c := config.Repository{
		Git: suite.gitURL,
		SHA: suite.gitSHA,
		Sources: []config.Source{
			{
				Src:     "1.txt",
				DstFile: filepath.Join(suite.dstDir, "1.txt"),
			},
		},
	}

	suite.mockCopyManager.EXPECT().
		CopyFile(filepath.Join(suite.cloneDir, "1.txt"), filepath.Join(suite.dstDir, "1.txt")).
		Return(nil)

	err := repo.CopySources(c, suite.cloneDir)
	assert.NoError(suite.T(), err)
}

func (suite *RepositoryPublicTestSuite) TestCopySourcesReturnsErrorWhenSourceIsFileAndDstFileAndCopyFileErrors() {
	repo := suite.NewRepositoryManager()
	specs := []FileSpec{
		{
			appFs:   suite.appFs,
			srcDir:  filepath.Join(suite.cloneDir),
			srcFile: filepath.Join(suite.cloneDir, "1.txt"),
		},
	}
	createFileSpecs(specs)

	c := config.Repository{
		Git: suite.gitURL,
		SHA: suite.gitSHA,
		Sources: []config.Source{
			{
				Src:     "1.txt",
				DstFile: filepath.Join(suite.dstDir, "1.txt"),
			},
		},
	}

	errors := errors.New("tests error")
	suite.mockCopyManager.EXPECT().CopyFile(gomock.Any(), gomock.Any()).Return(errors)

	err := repo.CopySources(c, suite.cloneDir)
	assert.Error(suite.T(), err)
}

// In order for `go test` to run this suite, we need to create
// a normal test function and pass our suite to suite.Run.
func TestRepositoryPublicTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryPublicTestSuite))
}
