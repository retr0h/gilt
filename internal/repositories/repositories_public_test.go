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

package repositories_test

import (
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/avfs/avfs"
	"github.com/avfs/avfs/vfs/memfs"
	"github.com/avfs/avfs/vfs/rofs"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/retr0h/gilt/v2/internal"
	"github.com/retr0h/gilt/v2/internal/mocks/exec"
	"github.com/retr0h/gilt/v2/internal/mocks/repository"
	"github.com/retr0h/gilt/v2/internal/repositories"
	"github.com/retr0h/gilt/v2/pkg/config"
)

type RepositoriesPublicTestSuite struct {
	suite.Suite

	ctrl     *gomock.Controller
	mockRepo *repository.MockRepositoryManager
	mockExec *exec.MockExecManager

	appFs            avfs.VFS
	dstDir           string
	giltDir          string
	gitURL           string
	gitVersion       string
	repoConfigDstDir []config.Repository
	SkipCommands     bool
	logger           *slog.Logger
}

func (suite *RepositoriesPublicTestSuite) NewTestRepositoriesManager(
	repoConfig []config.Repository,
) internal.RepositoriesManager {
	reposConfig := config.Repositories{
		Debug:        false,
		Parallel:     true,
		SkipCommands: suite.SkipCommands,
		GiltFile:     "Giltfile.yaml",
		GiltDir:      suite.giltDir,
		Repositories: repoConfig,
	}

	return repositories.New(
		suite.appFs,
		reposConfig,
		suite.mockRepo,
		suite.mockExec,
		suite.logger,
	)
}

func (suite *RepositoriesPublicTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockRepo = repository.NewMockRepositoryManager(suite.ctrl)
	suite.mockExec = exec.NewMockExecManager(suite.ctrl)

	suite.appFs = memfs.New()
	suite.dstDir = "/dstDir"
	suite.giltDir = "/giltDir"
	suite.gitURL = "https://example.com/user/repo.git"
	suite.gitVersion = "abc1234"
	suite.repoConfigDstDir = []config.Repository{
		{
			Git:     suite.gitURL,
			Version: suite.gitVersion,
			DstDir:  suite.dstDir,
		},
	}
	suite.SkipCommands = false
	suite.logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
}

func (suite *RepositoriesPublicTestSuite) TearDownTest() {
	defer suite.ctrl.Finish()
}

func (suite *RepositoriesPublicTestSuite) TestOverlayOkWhenDstDir() {
	repos := suite.NewTestRepositoriesManager(suite.repoConfigDstDir)
	expected := suite.appFs.Join(suite.giltDir, "cache")

	suite.mockRepo.EXPECT().
		Clone(suite.repoConfigDstDir[0], expected).
		Return(expected, nil)
	suite.mockRepo.EXPECT().
		Worktree(suite.repoConfigDstDir[0], expected, suite.dstDir).
		Return(nil)

	err := repos.Overlay()
	assert.NoError(suite.T(), err)
}

func (suite *RepositoriesPublicTestSuite) TestOverlayErrorWhenDstDir() {
	repos := suite.NewTestRepositoriesManager(suite.repoConfigDstDir)
	expected := suite.appFs.Join(suite.giltDir, "cache")
	errors := errors.New("tests error")

	suite.mockRepo.EXPECT().
		Clone(suite.repoConfigDstDir[0], expected).
		Return(expected, nil)
	suite.mockRepo.EXPECT().
		Worktree(suite.repoConfigDstDir[0], expected, suite.dstDir).
		Return(errors)

	err := repos.Overlay()
	assert.Error(suite.T(), err)
}

func (suite *RepositoriesPublicTestSuite) TestOverlayCacheDirCreateError() {
	// Replace the test FS with a read-only copy
	suite.appFs = rofs.New(suite.appFs)
	repos := suite.NewTestRepositoriesManager(suite.repoConfigDstDir)
	err := repos.Overlay()
	assert.Error(suite.T(), err)
}

func (suite *RepositoriesPublicTestSuite) TestOverlayDstDirExists() {
	repos := suite.NewTestRepositoriesManager(suite.repoConfigDstDir)

	suite.mockRepo.EXPECT().Clone(gomock.Any(), gomock.Any()).Return("", nil)
	suite.mockRepo.EXPECT().Worktree(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	_ = suite.appFs.MkdirAll(suite.dstDir, 0o755)
	err := repos.Overlay()
	assert.NoError(suite.T(), err)
}

func (suite *RepositoriesPublicTestSuite) TestOverlayErrorRemovingDstDir() {
	suite.mockRepo.EXPECT().Clone(gomock.Any(), gomock.Any()).Return("", nil)
	_ = suite.appFs.MkdirAll(suite.appFs.Join(suite.giltDir, "cache"), 0o700)
	_ = suite.appFs.MkdirAll(suite.dstDir, 0o755)
	// Replace the test FS with a read-only copy
	suite.appFs = rofs.New(suite.appFs)
	repos := suite.NewTestRepositoriesManager(suite.repoConfigDstDir)
	err := repos.Overlay()
	assert.Error(suite.T(), err)
}

func (suite *RepositoriesPublicTestSuite) TestOverlayReturnsErrorWhenDstDirDeleteFails() {
	suite.T().Skip("implement")
}

func (suite *RepositoriesPublicTestSuite) TestOverlayReturnsErrorWhenCloneErrors() {
	repos := suite.NewTestRepositoriesManager(suite.repoConfigDstDir)

	errors := errors.New("tests error")
	suite.mockRepo.EXPECT().Clone(gomock.Any(), gomock.Any()).Return("", errors)

	err := repos.Overlay()
	assert.Error(suite.T(), err)
}

func (suite *RepositoriesPublicTestSuite) TestOverlayOkWhenCopySources() {
	repoConfig := []config.Repository{
		{
			Git:     suite.gitURL,
			Version: suite.gitVersion,
			Sources: []config.Source{
				{
					Src:    "srcDir",
					DstDir: suite.dstDir,
				},
			},
		},
	}
	repos := suite.NewTestRepositoriesManager(repoConfig)

	suite.mockRepo.EXPECT().Clone(gomock.Any(), gomock.Any()).Return("", nil)
	suite.mockExec.EXPECT().
		RunInTempDir(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ string, _ string, fn func(string) error) error {
			if fn != nil {
				return fn("stub")
			}
			return nil
		})
	suite.mockRepo.EXPECT().Worktree(repoConfig[0], gomock.Any(), gomock.Any()).Return(nil)
	suite.mockRepo.EXPECT().CopySources(repoConfig[0], gomock.Any()).Return(nil)

	err := repos.Overlay()
	assert.NoError(suite.T(), err)
}

func (suite *RepositoriesPublicTestSuite) TestOverlayReturnsErrorWhenCopySourcesErrors() {
	repoConfig := []config.Repository{
		{
			Git:     suite.gitURL,
			Version: suite.gitVersion,
			Sources: []config.Source{
				{
					Src:    "srcDir",
					DstDir: suite.dstDir,
				},
			},
		},
	}
	repos := suite.NewTestRepositoriesManager(repoConfig)
	errors := errors.New("tests error")

	suite.mockRepo.EXPECT().Clone(gomock.Any(), gomock.Any()).Return("", nil)
	suite.mockExec.EXPECT().
		RunInTempDir(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ string, _ string, fn func(string) error) error {
			if fn != nil {
				return fn("stub")
			}
			return nil
		})
	suite.mockRepo.EXPECT().Worktree(repoConfig[0], gomock.Any(), gomock.Any()).Return(nil)
	suite.mockRepo.EXPECT().CopySources(gomock.Any(), gomock.Any()).Return(errors)

	err := repos.Overlay()
	assert.Error(suite.T(), err)
}

func (suite *RepositoriesPublicTestSuite) TestOverlayErrorCreatingCopySourcesWorktree() {
	repoConfig := []config.Repository{
		{
			Git:     suite.gitURL,
			Version: suite.gitVersion,
			Sources: []config.Source{
				{
					Src:    "srcDir",
					DstDir: suite.dstDir,
				},
			},
		},
	}
	repos := suite.NewTestRepositoriesManager(repoConfig)
	errors := errors.New("tests error")

	suite.mockRepo.EXPECT().Clone(gomock.Any(), gomock.Any()).Return("", nil)
	suite.mockExec.EXPECT().
		RunInTempDir(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ string, _ string, fn func(string) error) error {
			if fn != nil {
				return fn("stub")
			}
			return nil
		})
	suite.mockRepo.EXPECT().Worktree(repoConfig[0], gomock.Any(), gomock.Any()).Return(errors)

	err := repos.Overlay()
	assert.Error(suite.T(), err)
}

func (suite *RepositoriesPublicTestSuite) TestOverlayOkWhenCommands() {
	repoConfig := []config.Repository{
		{
			Git:     suite.gitURL,
			Version: suite.gitVersion,
			DstDir:  suite.dstDir,
			Commands: []config.Command{
				{
					Cmd:  "touch",
					Args: []string{"/tmp/foo"},
				},
			},
		},
	}

	repos := suite.NewTestRepositoriesManager(repoConfig)

	suite.mockRepo.EXPECT().Clone(gomock.Any(), gomock.Any()).Return("", nil)
	suite.mockRepo.EXPECT().Worktree(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	suite.mockExec.EXPECT().RunCmd("touch", []string{"/tmp/foo"}).Return("", nil)

	err := repos.Overlay()
	assert.NoError(suite.T(), err)
}

func (suite *RepositoriesPublicTestSuite) TestOverlayReturnsErrorWhenCommandErrors() {
	repoConfig := []config.Repository{
		{
			Git:     suite.gitURL,
			Version: suite.gitVersion,
			DstDir:  suite.dstDir,
			Commands: []config.Command{
				{
					Cmd:  "touch",
					Args: []string{"/tmp/foo"},
				},
			},
		},
	}

	repos := suite.NewTestRepositoriesManager(repoConfig)

	suite.mockRepo.EXPECT().Clone(gomock.Any(), gomock.Any()).Return("", nil)
	suite.mockRepo.EXPECT().Worktree(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	errors := errors.New("tests error")
	suite.mockExec.EXPECT().RunCmd(gomock.Any(), gomock.Any()).Return("", errors)

	err := repos.Overlay()
	assert.Error(suite.T(), err)
}

func (suite *RepositoriesPublicTestSuite) TestOverlaySkipsCommands() {
	suite.SkipCommands = true
	repoConfig := []config.Repository{
		{
			Git:     suite.gitURL,
			Version: suite.gitVersion,
			DstDir:  suite.dstDir,
			Commands: []config.Command{
				{
					Cmd:  "touch",
					Args: []string{"/tmp/foo"},
				},
			},
		},
	}

	repos := suite.NewTestRepositoriesManager(repoConfig)

	suite.mockRepo.EXPECT().Clone(gomock.Any(), gomock.Any()).Return("", nil)
	suite.mockRepo.EXPECT().Worktree(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	// Explicitly check that RunCmd is never called
	suite.mockExec.EXPECT().RunCmd(gomock.Any(), gomock.Any()).Times(0)

	err := repos.Overlay()
	assert.NoError(suite.T(), err)
}

// In order for `go test` to run this suite, we need to create
// a normal test function and pass our suite to suite.Run.
func TestRepositoriesPublicTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoriesPublicTestSuite))
}
