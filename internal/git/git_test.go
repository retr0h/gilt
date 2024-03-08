// Copyright (c) 2023 John Dewey

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

package git

import (
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/avfs/avfs/vfs/memfs"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/retr0h/gilt/v2/internal"
	"github.com/retr0h/gilt/v2/internal/mocks/exec"
)

type GitManagerTestSuite struct {
	suite.Suite

	ctrl     *gomock.Controller
	mockExec *exec.MockExecManager

	gitURL     string
	gitVersion string
	cloneDir   string
	dstDir     string

	gm internal.GitManager
}

func (suite *GitManagerTestSuite) NewTestGitManager() internal.GitManager {
	return New(
		memfs.New(),
		suite.mockExec,
		slog.New(slog.NewTextHandler(os.Stdout, nil)),
	)
}

func (suite *GitManagerTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockExec = exec.NewMockExecManager(suite.ctrl)

	suite.gitURL = "https://example.com/user/repo.git"
	suite.gitVersion = "abc123"
	suite.cloneDir = "/cloneDir"
	suite.dstDir = "/dstDir"

	suite.gm = suite.NewTestGitManager()
}

func (suite *GitManagerTestSuite) TearDownTest() {
	defer suite.ctrl.Finish()
}

func (suite *GitManagerTestSuite) TestWorktreeErrorWhenAbsErrors() {
	originalAbsFn := AbsFn
	AbsFn = func(g *Git, _ string) (string, error) {
		return "", fmt.Errorf("failed to get abs path")
	}
	defer func() { AbsFn = originalAbsFn }()

	err := suite.gm.Worktree(suite.cloneDir, suite.gitVersion, suite.dstDir)
	assert.Error(suite.T(), err)
}

// In order for `go test` to run this suite, we need to create
// a normal test function and pass our suite to suite.Run.
func TestGitManagerTestSuite(t *testing.T) {
	suite.Run(t, new(GitManagerTestSuite))
}
