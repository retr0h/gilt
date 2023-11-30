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

package repositories

import (
	"log/slog"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/retr0h/go-gilt/internal/exec"
	"github.com/retr0h/go-gilt/internal/repository"
	"github.com/retr0h/go-gilt/pkg/config"
)

type RepositoriesTestSuite struct {
	suite.Suite

	ctrl     *gomock.Controller
	mockRepo *repository.MockRepositoryManager
	mockExec *exec.MockExecManager

	appFs   afero.Fs
	giltDir string
	logger  *slog.Logger
}

func (suite *RepositoriesTestSuite) NewTestRepositories(
	giltDir string,
) *Repositories {
	reposConfig := config.Repositories{
		Debug:        false,
		GiltFile:     "Giltfile.yaml",
		GiltDir:      giltDir,
		Repositories: []config.Repository{},
	}

	return New(
		suite.appFs,
		reposConfig,
		suite.mockRepo,
		suite.mockExec,
		suite.logger,
	)
}

func (suite *RepositoriesTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockRepo = repository.NewMockRepositoryManager(suite.ctrl)
	suite.mockExec = exec.NewMockExecManager(suite.ctrl)
	defer suite.ctrl.Finish()

	suite.appFs = afero.NewMemMapFs()
	suite.giltDir = "/giltDir"

	suite.logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
}

func (suite *RepositoriesTestSuite) TestgetCloneDirOk() {
	repos := suite.NewTestRepositories(suite.giltDir)

	got := repos.getCloneDir(
		suite.giltDir,
		config.Repository{
			Git: "https://example.com/user/repo2.git",
			SHA: "abc123",
		},
	)
	assert.Equal(suite.T(), "/giltDir/https---example.com-user-repo2.git-abc123", got)
}

func (suite *RepositoriesTestSuite) TestgetCloneDirOkByTag() {
	repos := suite.NewTestRepositories(suite.giltDir)

	got := repos.getCloneDir(
		suite.giltDir,
		config.Repository{
			Git: "https://example.com/user/repo2.git",
			Tag: "v1.1",
		},
	)
	assert.Equal(suite.T(), "/giltDir/https---example.com-user-repo2.git-v1.1", got)
}

func (suite *RepositoriesTestSuite) TestgetCloneHashOk() {
	repos := suite.NewTestRepositories(suite.giltDir)

	got := repos.getCloneHash(
		config.Repository{
			Git: "https://example.com/user/repo2.git",
			SHA: "abc123",
		},
	)
	assert.Equal(suite.T(), "https---example.com-user-repo2.git-abc123", got)
}

func (suite *RepositoriesTestSuite) TestgetCacheDir() {
	repos := suite.NewTestRepositories(suite.giltDir)

	expectedDir := "/giltDir/cache"
	got, err := repos.getCacheDir()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedDir, got)

	exists, err := afero.Exists(suite.appFs, expectedDir)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), exists)
}

// In order for `go test` to run this suite, we need to create
// a normal test function and pass our suite to suite.Run.
func TestRepositoriesTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoriesTestSuite))
}
