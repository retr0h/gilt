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
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/retr0h/go-gilt/internal"
	"github.com/retr0h/go-gilt/internal/mocks/exec"
	"github.com/retr0h/go-gilt/internal/mocks/repository"
	"github.com/retr0h/go-gilt/internal/repositories"
	"github.com/retr0h/go-gilt/pkg/config"
)

type RepositoriesPublicTestSuite struct {
	suite.Suite

	ctrl     *gomock.Controller
	mockRepo *repository.MockRepositoryManager
	mockExec *exec.MockExecManager

	appFs            afero.Fs
	dstDir           string
	giltDir          string
	gitURL           string
	gitSHA           string
	repoConfigDstDir []config.Repository
	logger           *slog.Logger
}

func (suite *RepositoriesPublicTestSuite) NewTestRepositoriesManager(
	repoConfig []config.Repository,
) internal.RepositoriesManager {
	reposConfig := config.Repositories{
		Debug:        false,
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
	defer suite.ctrl.Finish()

	suite.appFs = afero.NewMemMapFs()
	suite.dstDir = "/dstDir"
	suite.giltDir = "/giltDir"
	suite.gitURL = "https://example.com/user/repo.git"
	suite.gitSHA = "abc1234"
	suite.repoConfigDstDir = []config.Repository{
		{
			Git:    suite.gitURL,
			SHA:    suite.gitSHA,
			DstDir: suite.dstDir,
		},
	}

	suite.logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
}

func (suite *RepositoriesPublicTestSuite) TestOverlayOkWhenDstDir() {
	repos := suite.NewTestRepositoriesManager(suite.repoConfigDstDir)
	expected := filepath.Join(suite.giltDir, "cache")

	suite.mockRepo.EXPECT().
		Clone(suite.repoConfigDstDir[0], filepath.Join(suite.giltDir, "cache")).
		Return(expected, nil)
	suite.mockRepo.EXPECT().
		Worktree(suite.repoConfigDstDir[0], filepath.Join(suite.giltDir, "cache"), suite.dstDir).
		Return(nil)

	err := repos.Overlay()
	assert.NoError(suite.T(), err)
}

func (suite *RepositoriesPublicTestSuite) TestOverlayDstDirExists() {
	repos := suite.NewTestRepositoriesManager(suite.repoConfigDstDir)

	suite.mockRepo.EXPECT().Clone(gomock.Any(), gomock.Any()).Return("", nil)
	suite.mockRepo.EXPECT().Worktree(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	_ = suite.appFs.MkdirAll(suite.dstDir, 0o755)
	err := repos.Overlay()
	assert.NoError(suite.T(), err)
}

func (suite *RepositoriesPublicTestSuite) TestOverlayReturnsErrorWhenDstDirDeleteFails() {
	suite.T().Skip("implement")
}

func (suite *RepositoriesPublicTestSuite) TestOverlayReturnsErrorWhenCloneErrors() {
	repos := suite.NewTestRepositoriesManager(suite.repoConfigDstDir)

	errors := errors.New("tests error")
	suite.mockRepo.EXPECT().Clone(gomock.Any(), gomock.Any()).Return("", errors)
	suite.mockRepo.EXPECT().Worktree(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	err := repos.Overlay()
	assert.Error(suite.T(), err)
}

func (suite *RepositoriesPublicTestSuite) TestOverlayOkWhenCopySources() {
	repoConfig := []config.Repository{
		{
			Git: suite.gitURL,
			SHA: suite.gitSHA,
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
	suite.mockExec.EXPECT().RunInTempDir(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	suite.mockRepo.EXPECT().Worktree(repoConfig[0], gomock.Any(), gomock.Any()).Return(nil)
	suite.mockRepo.EXPECT().CopySources(repoConfig[0], gomock.Any()).Return(nil)

	err := repos.Overlay()
	assert.NoError(suite.T(), err)
}

func (suite *RepositoriesPublicTestSuite) TestOverlayReturnsErrorWhenCopySourcesErrors() {
	repoConfig := []config.Repository{
		{
			Git: suite.gitURL,
			SHA: suite.gitSHA,
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
	suite.mockExec.EXPECT().RunInTempDir(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	suite.mockRepo.EXPECT().Worktree(repoConfig[0], gomock.Any(), gomock.Any()).Return(nil)
	suite.mockRepo.EXPECT().CopySources(gomock.Any(), gomock.Any()).Return(errors)

	err := repos.Overlay()
	assert.Error(suite.T(), err)
}

func (suite *RepositoriesPublicTestSuite) TestOverlayOkWhenCommands() {
	repoConfig := []config.Repository{
		{
			Git:    suite.gitURL,
			SHA:    suite.gitSHA,
			DstDir: suite.dstDir,
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
	suite.mockExec.EXPECT().RunCmd("touch", []string{"/tmp/foo"}).Return(nil)

	err := repos.Overlay()
	assert.NoError(suite.T(), err)
}

func (suite *RepositoriesPublicTestSuite) TestOverlayReturnsErrorWhenCommandErrors() {
	repoConfig := []config.Repository{
		{
			Git:    suite.gitURL,
			SHA:    suite.gitSHA,
			DstDir: suite.dstDir,
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
	suite.mockExec.EXPECT().RunCmd(gomock.Any(), gomock.Any()).Return(errors)

	err := repos.Overlay()
	assert.Error(suite.T(), err)
}

// In order for `go test` to run this suite, we need to create
// a normal test function and pass our suite to suite.Run.
func TestRepositoriesPublicTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoriesPublicTestSuite))
}
