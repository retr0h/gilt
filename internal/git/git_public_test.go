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

package git_test

import (
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/avfs/avfs"
	"github.com/avfs/avfs/vfs/failfs"
	"github.com/avfs/avfs/vfs/memfs"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/retr0h/gilt/v2/internal"
	"github.com/retr0h/gilt/v2/internal/git"
	"github.com/retr0h/gilt/v2/internal/mocks/exec"
)

type GitManagerPublicTestSuite struct {
	suite.Suite

	ctrl     *gomock.Controller
	mockExec *exec.MockExecManager
	appFs    avfs.VFS

	gitURL     string
	gitVersion string
	origin     string
	cloneDir   string
	dstDir     string

	gm internal.GitManager
}

func (suite *GitManagerPublicTestSuite) NewTestGitManager() internal.GitManager {
	return git.New(
		suite.appFs,
		suite.mockExec,
		slog.New(slog.NewTextHandler(os.Stdout, nil)),
	)
}

func (suite *GitManagerPublicTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockExec = exec.NewMockExecManager(suite.ctrl)
	suite.appFs = memfs.New()

	suite.gitURL = "https://example.com/user/repo.git"
	suite.gitVersion = "abc123"
	suite.origin = "gilt"
	suite.cloneDir = "/cloneDir"
	suite.dstDir = "/dstDir"

	suite.gm = suite.NewTestGitManager()
}

func (suite *GitManagerPublicTestSuite) TearDownTest() {
	defer suite.ctrl.Finish()
}

func (suite *GitManagerPublicTestSuite) TestCloneOk() {
	suite.T().Skip("Skipping until we can mock go-git properly")
	suite.mockExec.EXPECT().
		RunCmd("git", []string{
			"-c", "clone.defaultRemoteName=" + suite.origin,
			"clone", "--bare", "--filter=blob:none", suite.gitURL, suite.cloneDir,
		}).
		Return("", nil)
	suite.mockExec.EXPECT().
		RunCmdInDir("git", []string{"remote", "rename", "origin", suite.origin}, suite.cloneDir).
		Return("", nil)

	err := suite.gm.Clone(suite.gitURL, suite.origin, suite.cloneDir)
	assert.NoError(suite.T(), err)
}

func (suite *GitManagerPublicTestSuite) TestCloneReturnsError() {
	suite.T().Skip("Skipping until we can mock go-git properly")
	errors := errors.New("tests error")
	suite.mockExec.EXPECT().RunCmd(gomock.Any(), gomock.Any()).Return("", errors)
	// `git remote rename` is not called if the clone throws errors
	suite.mockExec.EXPECT().RunCmdInDir("git", gomock.Any(), suite.cloneDir).Times(0)

	err := suite.gm.Clone(suite.gitURL, suite.origin, suite.cloneDir)
	assert.Error(suite.T(), err)
}

func (suite *GitManagerPublicTestSuite) TestWorktreeOk() {
	suite.mockExec.EXPECT().
		RunCmdInDir("git", []string{"worktree", "add", "--force", suite.dstDir, suite.gitVersion}, suite.cloneDir).
		Return("", nil)
	suite.mockExec.EXPECT().
		RunCmdInDir("git", []string{"worktree", "prune", "--verbose"}, suite.cloneDir).
		Return("", nil)
	err := suite.gm.Worktree(suite.cloneDir, suite.gitVersion, suite.dstDir)
	assert.NoError(suite.T(), err)
}

func (suite *GitManagerPublicTestSuite) TestWorktreeError() {
	errors := errors.New("tests error")
	suite.mockExec.EXPECT().
		RunCmdInDir("git", []string{"worktree", "add", "--force", suite.dstDir, suite.gitVersion}, suite.cloneDir).
		Return("", errors)
	err := suite.gm.Worktree(suite.cloneDir, suite.gitVersion, suite.dstDir)
	assert.Error(suite.T(), err)
}

func (suite *GitManagerPublicTestSuite) TestWorktreeErrorWhenAbsErrors() {
	// Make Abs() calls fail
	vfs := failfs.New(suite.appFs)
	_ = vfs.SetFailFunc(func(_ avfs.VFSBase, fn avfs.FnVFS, _ *failfs.FailParam) error {
		if fn == avfs.FnAbs {
			return errors.New("FailFS!")
		}
		return nil
	})
	suite.appFs = vfs

	gm := suite.NewTestGitManager()

	err := gm.Worktree(suite.cloneDir, suite.gitVersion, suite.dstDir)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "FailFS!", err.Error())
}

func (suite *GitManagerPublicTestSuite) TestUpdateOk() {
	suite.T().Skip("Skipping until we can mock go-git properly")
	suite.mockExec.EXPECT().
		RunCmdInDir("git", []string{"fetch", "--tags", "--force", suite.origin, "+refs/heads/*:refs/heads/*"}, suite.cloneDir).
		Return("", nil)
	err := suite.gm.Update(suite.origin, suite.cloneDir)
	assert.NoError(suite.T(), err)
}

func (suite *GitManagerPublicTestSuite) TestUpdateError() {
	suite.T().Skip("Skipping until we can mock go-git properly")
	errors := errors.New("tests error")
	suite.mockExec.EXPECT().
		RunCmdInDir("git", []string{"fetch", "--tags", "--force", suite.origin, "+refs/heads/*:refs/heads/*"}, suite.cloneDir).
		Return("", errors)
	err := suite.gm.Update(suite.origin, suite.cloneDir)
	assert.Error(suite.T(), err)
}

func (suite *GitManagerPublicTestSuite) TestRemoteOk() {
	suite.T().Skip("Skipping until we can mock go-git properly")
	suite.mockExec.EXPECT().RunCmdInDir("git", []string{"remote"}, suite.cloneDir).Return("", nil)
	_, err := suite.gm.RemoteExists(suite.cloneDir, suite.origin)
	assert.NoError(suite.T(), err)
}

func (suite *GitManagerPublicTestSuite) TestRemoteError() {
	suite.T().Skip("Skipping until we can mock go-git properly")
	errors := errors.New("tests error")
	suite.mockExec.EXPECT().
		RunCmdInDir("git", []string{"remote"}, suite.cloneDir).
		Return("", errors)
	_, err := suite.gm.RemoteExists(suite.cloneDir, suite.origin)
	assert.Error(suite.T(), err)
}

// In order for `go test` to run this suite, we need to create
// a normal test function and pass our suite to suite.Run.
func TestGitManagerPublicTestSuite(t *testing.T) {
	suite.Run(t, new(GitManagerPublicTestSuite))
}
