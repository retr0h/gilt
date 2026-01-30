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
	"path/filepath"
	"testing"

	"github.com/avfs/avfs"
	"github.com/avfs/avfs/vfs/failfs"
	"github.com/avfs/avfs/vfs/memfs"
	"github.com/go-git/go-billy/v6"
	gogit "github.com/go-git/go-git/v6"
	gogitconfig "github.com/go-git/go-git/v6/config"
	gogitobject "github.com/go-git/go-git/v6/plumbing/object"
	gogitstorage "github.com/go-git/go-git/v6/storage"
	gogitmemfs "github.com/go-git/go-git/v6/storage/memory"
	gogitworktree "github.com/go-git/go-git/v6/x/plumbing/worktree"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/retr0h/gilt/v2/internal"
	"github.com/retr0h/gilt/v2/internal/git"
)

type GitManagerPublicTestSuite struct {
	suite.Suite

	ctrl  *gomock.Controller
	appFs avfs.VFS

	gitURL     string
	gitVersion string
	origin     string
	cloneDir   string
	dstDir     string

	gitClone    func(string, *gogit.CloneOptions) (*gogit.Repository, error)
	gitOpen     func(string) (*gogit.Repository, error)
	gitWorktree func(gogitstorage.Storer) (git.Worktree, error)

	gm internal.GitManager
}

// The go-git fixture in these tests is a memfs-backed repository,
// that we can populate as needed for "happy path" testing.  To inject
// errors, we'll take advantage of the fact that most go-git operations
// load the repo config as a first step, so by making that fail fast,
// we can avoid having to mock out a bunch of go-git internals.
type failStorer struct {
	gogitstorage.Storer
	err error
}

var errConfig = errors.New("Config error")

func (s *failStorer) Config() (*gogitconfig.Config, error) {
	return nil, s.err
}

// A simple set of stubs for the Worktree interface, these can be fleshed out
// if we ever need them to do anything more than "absolutely nothing"
type mockWorktree struct{}

func (w *mockWorktree) Add(
	_ billy.Filesystem,
	_ string,
	_ ...gogitworktree.Option,
) error {
	return nil
}

func (w *mockWorktree) Remove(_ string) error { return nil }

func (suite *GitManagerPublicTestSuite) NewTestGitManager() internal.GitManager {
	return git.NewWithOverrides(
		suite.appFs,
		slog.New(slog.NewTextHandler(os.Stdout, nil)),
		suite.gitClone,
		suite.gitOpen,
		suite.gitWorktree,
	)
}

func (suite *GitManagerPublicTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.appFs = memfs.New()

	suite.gitURL = "https://example.com/user/repo.git"
	suite.gitVersion = "abc123"
	suite.origin = "gilt"
	suite.cloneDir = "/cloneDir"
	suite.dstDir = "/dstDir"

	suite.gitClone = nil
	suite.gitOpen = nil
	suite.gitWorktree = nil

	suite.gm = suite.NewTestGitManager()
}

func (suite *GitManagerPublicTestSuite) TearDownTest() {
	defer suite.ctrl.Finish()
}

// TestInTheWorstMannerPossible uses the real go-git library to clone a real
// repository and create a real worktree in a tempdir, so that we can set
// breakpoints, inspect go-git internals, and all the other goodies that
// the integration tests cannot easily do.
func (suite *GitManagerPublicTestSuite) TestInTheWorstMannerPossible() {
	suite.T().Skip("In case of emergency, break glass")
	suite.gitClone = gogit.PlainClone
	suite.gitOpen = gogit.PlainOpen
	suite.gitWorktree = git.NewWorktree
	suite.gm = suite.NewTestGitManager()

	tmpDir, err := filepath.EvalSymlinks(suite.T().TempDir())
	URL := `https://github.com/lorin/openstack-ansible-modules.git`

	assert.NoError(suite.T(), err)

	cloneDir := suite.appFs.Join(tmpDir, "clone")
	dstDir := suite.appFs.Join(tmpDir, "dst")

	// Clone URL into tmpDir
	err = suite.gm.Clone(URL, "gilt", cloneDir)
	assert.NoError(suite.T(), err)

	// Create worktree at HEAD in dstDir
	err = suite.gm.Worktree(cloneDir, "2677cc3", dstDir)
	assert.NoError(suite.T(), err)
}

// NOTE(nic): the behavior of this test might change when go-git
// v6 goes final; in the current alpha release, the memfs storer
// does not support worktrees, but either way, this is a test that
// the wrapper, wraps
func (suite *GitManagerPublicTestSuite) TestNewWorktree() {
	wt, err := git.NewWorktree(gogitmemfs.NewStorage())
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), wt)
}

func (suite *GitManagerPublicTestSuite) TestClone() {
	var repo *gogit.Repository

	suite.gitClone = func(path string, _ *gogit.CloneOptions) (*gogit.Repository, error) {
		assert.Equal(suite.T(), suite.cloneDir, path)
		repo, _ = gogit.Init(gogitmemfs.NewStorage(), nil)
		return repo, nil
	}
	suite.gm = suite.NewTestGitManager()

	err := suite.gm.Clone(suite.gitURL, suite.origin, suite.cloneDir)
	assert.NoError(suite.T(), err)

	// Check that the remote workarounds were applied correctly.  Can likely drop
	// this once go-git supports blobless clones properly
	cfg, err := repo.Storer.Config()
	assert.NoError(suite.T(), err)
	opts := cfg.Raw.Section("remote").Subsection(suite.origin).Options
	assert.Equal(suite.T(), "true", opts.Get("promisor"))
	// NOTE(nic): use this assertion instead when blobless clone support is finished
	// assert.Equal(suite.T(), "blob:none", opts.Get("partialclonefilter"))
	assert.Equal(suite.T(), "", opts.Get("partialclonefilter"))
}

func (suite *GitManagerPublicTestSuite) TestCloneErrorOnOpen() {
	suite.gitClone = func(path string, _ *gogit.CloneOptions) (*gogit.Repository, error) {
		assert.Equal(suite.T(), suite.cloneDir, path)
		return nil, errors.New("gitOpen error")
	}
	suite.gm = suite.NewTestGitManager()

	err := suite.gm.Clone(suite.gitURL, suite.origin, suite.cloneDir)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "gitOpen error", err.Error())
}

func (suite *GitManagerPublicTestSuite) TestCloneErrorFetchingConfig() {
	suite.gitClone = func(path string, _ *gogit.CloneOptions) (*gogit.Repository, error) {
		assert.Equal(suite.T(), suite.cloneDir, path)
		repo, err := gogit.Init(gogitmemfs.NewStorage(), nil)
		assert.NoError(suite.T(), err)
		repo.Storer = &failStorer{repo.Storer, errConfig}
		return repo, nil
	}
	suite.gm = suite.NewTestGitManager()

	err := suite.gm.Clone(suite.gitURL, suite.origin, suite.cloneDir)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "Config error", err.Error())
}

// NOTE(nic): Skip the "happy path" test for now -- there are a few options
// for testing it with this fixture, but none of them are particularly good,
// and all involve getting further into the weeds of go-git internals
// than can be considered wise.  The most straighforward approach is to create
// a file:// remote pointing to a repository in a tempdir, but we want to avoid
// touching the OS filesystem in unit tests for as long as we can.
func (suite *GitManagerPublicTestSuite) TestUpdate() {
	suite.T().Skip("need to synthesize repo.Fetch()")
	suite.gitOpen = func(path string) (*gogit.Repository, error) {
		assert.Equal(suite.T(), suite.cloneDir, path)
		repo, _ := gogit.Init(gogitmemfs.NewStorage(), nil)
		_, _ = repo.CreateRemote(&gogitconfig.RemoteConfig{
			Name: suite.origin,
			URLs: []string{suite.gitURL},
		})
		return repo, nil
	}
	suite.gm = suite.NewTestGitManager()

	err := suite.gm.Update(suite.origin, suite.cloneDir)
	assert.NoError(suite.T(), err)
}

func (suite *GitManagerPublicTestSuite) TestUpdateNoOp() {
	suite.gitOpen = func(path string) (*gogit.Repository, error) {
		assert.Equal(suite.T(), suite.cloneDir, path)
		repo, _ := gogit.Init(gogitmemfs.NewStorage(), nil)
		repo.Storer = &failStorer{repo.Storer, gogit.NoErrAlreadyUpToDate}
		return repo, nil
	}
	suite.gm = suite.NewTestGitManager()

	err := suite.gm.Update(suite.origin, suite.cloneDir)
	assert.NoError(suite.T(), err)
}

func (suite *GitManagerPublicTestSuite) TestUpdateErrorOnOpen() {
	suite.gitOpen = func(path string) (*gogit.Repository, error) {
		assert.Equal(suite.T(), suite.cloneDir, path)
		return nil, errors.New("gitOpen error")
	}
	suite.gm = suite.NewTestGitManager()

	err := suite.gm.Update(suite.origin, suite.cloneDir)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "gitOpen error", err.Error())
}

func (suite *GitManagerPublicTestSuite) TestUpdateErrorOnFetch() {
	suite.gitOpen = func(path string) (*gogit.Repository, error) {
		assert.Equal(suite.T(), suite.cloneDir, path)
		repo, _ := gogit.Init(gogitmemfs.NewStorage(), nil)
		repo.Storer = &failStorer{repo.Storer, errConfig}
		return repo, nil
	}
	suite.gm = suite.NewTestGitManager()

	err := suite.gm.Update(suite.origin, suite.cloneDir)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "Config error", err.Error())
}

func (suite *GitManagerPublicTestSuite) TestWorktree() {
	suite.gitOpen = func(path string) (*gogit.Repository, error) {
		assert.Equal(suite.T(), suite.cloneDir, path)
		repo, _ := gogit.Init(gogitmemfs.NewStorage(), nil)
		// Set up a repo with a resolvable treeish in it
		commit := &gogitobject.Commit{}
		obj := repo.Storer.NewEncodedObject()
		_ = commit.Encode(obj)
		hash, _ := repo.Storer.SetEncodedObject(obj)
		_, _ = repo.CreateTag(suite.gitVersion, hash, nil)
		return repo, nil
	}
	suite.gitWorktree = func(gogitstorage.Storer) (git.Worktree, error) {
		return &mockWorktree{}, nil
	}
	suite.gm = suite.NewTestGitManager()

	err := suite.gm.Worktree(suite.cloneDir, suite.gitVersion, suite.dstDir)
	assert.NoError(suite.T(), err)
}

func (suite *GitManagerPublicTestSuite) TestWorktreeErrorWhenAbsErrors() {
	// Make Abs() calls fail
	vfs := failfs.New(suite.appFs)
	_ = vfs.SetFailFunc(func(_ avfs.VFSBase, fn avfs.FnVFS, _ *failfs.FailParam) error {
		if fn == avfs.FnAbs {
			return errors.New("failFS")
		}
		return nil
	})
	suite.appFs = vfs

	gm := suite.NewTestGitManager()

	err := gm.Worktree(suite.cloneDir, suite.gitVersion, suite.dstDir)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "failFS", err.Error())
}

func (suite *GitManagerPublicTestSuite) TestWorktreeErrorOnOpen() {
	suite.gitOpen = func(path string) (*gogit.Repository, error) {
		assert.Equal(suite.T(), suite.cloneDir, path)
		return nil, errors.New("gitOpen error")
	}
	suite.gm = suite.NewTestGitManager()
	err := suite.gm.Worktree(suite.cloneDir, suite.gitVersion, suite.dstDir)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "gitOpen error", err.Error())
}

func (suite *GitManagerPublicTestSuite) TestWorktreeErrorOnResolve() {
	suite.gitOpen = func(path string) (*gogit.Repository, error) {
		assert.Equal(suite.T(), suite.cloneDir, path)
		repo, _ := gogit.Init(gogitmemfs.NewStorage(), nil)
		return repo, nil
	}
	suite.gm = suite.NewTestGitManager()
	err := suite.gm.Worktree(suite.cloneDir, suite.gitVersion, suite.dstDir)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "reference not found", err.Error())
}

func (suite *GitManagerPublicTestSuite) TestWorktreeErrorCreatingWorktree() {
	suite.gitOpen = func(path string) (*gogit.Repository, error) {
		assert.Equal(suite.T(), suite.cloneDir, path)
		repo, _ := gogit.Init(gogitmemfs.NewStorage(), nil)
		commit := &gogitobject.Commit{}
		obj := repo.Storer.NewEncodedObject()
		_ = commit.Encode(obj)
		hash, _ := repo.Storer.SetEncodedObject(obj)
		_, _ = repo.CreateTag(suite.gitVersion, hash, nil)
		return repo, nil
	}
	suite.gitWorktree = func(gogitstorage.Storer) (git.Worktree, error) {
		return nil, errors.New("gitWorktree error")
	}
	suite.gm = suite.NewTestGitManager()
	err := suite.gm.Worktree(suite.cloneDir, suite.gitVersion, suite.dstDir)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "gitWorktree error", err.Error())
}

func (suite *GitManagerPublicTestSuite) TestRemoteExists() {
	suite.gitOpen = func(path string) (*gogit.Repository, error) {
		assert.Equal(suite.T(), suite.cloneDir, path)
		repo, _ := gogit.Init(gogitmemfs.NewStorage(), nil)
		_, _ = repo.CreateRemote(&gogitconfig.RemoteConfig{
			Name: suite.origin,
			URLs: []string{suite.gitURL},
		})
		return repo, nil
	}
	suite.gm = suite.NewTestGitManager()
	exists, err := suite.gm.RemoteExists(suite.cloneDir, suite.origin)
	assert.True(suite.T(), exists)
	assert.NoError(suite.T(), err)
}

func (suite *GitManagerPublicTestSuite) TestRemoteExistsFail() {
	suite.gitOpen = func(path string) (*gogit.Repository, error) {
		assert.Equal(suite.T(), suite.cloneDir, path)
		repo, _ := gogit.Init(gogitmemfs.NewStorage(), nil)
		return repo, nil
	}
	suite.gm = suite.NewTestGitManager()
	exists, err := suite.gm.RemoteExists(suite.cloneDir, suite.origin)
	assert.False(suite.T(), exists)
	assert.NoError(suite.T(), err)
}

func (suite *GitManagerPublicTestSuite) TestRemoteExistsErrorOnOpen() {
	suite.gitOpen = func(path string) (*gogit.Repository, error) {
		assert.Equal(suite.T(), suite.cloneDir, path)
		return nil, errors.New("gitOpen error")
	}
	suite.gm = suite.NewTestGitManager()
	exists, err := suite.gm.RemoteExists(suite.cloneDir, suite.origin)
	assert.False(suite.T(), exists)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "gitOpen error", err.Error())
}

func (suite *GitManagerPublicTestSuite) TestRemoteErrorOnLookup() {
	suite.gitOpen = func(path string) (*gogit.Repository, error) {
		assert.Equal(suite.T(), suite.cloneDir, path)
		repo, _ := gogit.Init(gogitmemfs.NewStorage(), nil)
		repo.Storer = &failStorer{repo.Storer, errConfig}
		return repo, nil
	}
	suite.gm = suite.NewTestGitManager()
	exists, err := suite.gm.RemoteExists(suite.cloneDir, suite.origin)
	assert.False(suite.T(), exists)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "Config error", err.Error())
}

// In order for `go test` to run this suite, we need to create
// a normal test function and pass our suite to suite.Run.
func TestGitManagerPublicTestSuite(t *testing.T) {
	suite.Run(t, new(GitManagerPublicTestSuite))
}
