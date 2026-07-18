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

// Package git package needs reworked into proper Git libraries.  However, this
// package will remain using exec as it was easiest to port from gilt's
// python counterpart.
package git

import (
	"log/slog"

	"github.com/avfs/avfs"
	"github.com/go-git/go-billy/v6"
	"github.com/go-git/go-billy/v6/osfs"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/config"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/protocol/packp"
	"github.com/go-git/go-git/v6/storage"

	xworktree "github.com/go-git/go-git/v6/x/plumbing/worktree"
)

// Worktree is a minimalist interface for the go-git worktree
// operations we actually use, that we may stub them out in tests.
type Worktree interface {
	Add(billy.Filesystem, string, ...xworktree.Option) error
	Remove(string) error
}

// NewWorktree wraps the go-git worktree constructor for use with our interface
func NewWorktree(s storage.Storer) (Worktree, error) {
	return xworktree.New(s)
}

// New factory to create a new Git instance.
func New(
	appFs avfs.VFS,
	logger *slog.Logger,
) *Git {
	return &Git{
		appFs:       appFs,
		logger:      logger,
		gitClone:    git.PlainClone,
		gitOpen:     git.PlainOpen,
		gitWorktree: NewWorktree,
	}
}

// NewWithOverrides for testing with fixtures.
func NewWithOverrides(
	appFs avfs.VFS,
	logger *slog.Logger,
	gitClone func(string, *git.CloneOptions) (*git.Repository, error),
	gitOpen func(string) (*git.Repository, error),
	gitWorktree func(storer storage.Storer) (Worktree, error),
) *Git {
	g := New(appFs, logger)
	g.gitClone = gitClone
	g.gitOpen = gitOpen
	g.gitWorktree = gitWorktree
	return g
}

// Clone the repo.  This is a bare repo, with only metadata to start with.
func (g *Git) Clone(gitURL, origin, cloneDir string) error {
	opts := &git.CloneOptions{
		URL:        gitURL,
		RemoteName: origin,
		Bare:       true,
		NoCheckout: true,
		Filter:     packp.FilterBlobNone(),
	}
	// NOTE(nic): blobless clones don't work quite right yet, so turn them off for
	//  now.  This regression makes Gilt nigh-unusable, but we can knock all the
	//  rough edges off of everything else while we wait for upstream to finish
	//  implementing lazy-fetches
	opts.Filter = packp.Filter("")
	repo, err := g.gitClone(cloneDir, opts)
	g.logger.Debug(
		"git.Clone",
		slog.String("gitURL", gitURL),
		slog.String("origin", origin),
		slog.String("cloneDir", cloneDir),
		slog.Any("repo", repo),
		slog.Any("err", err),
	)
	if err != nil {
		return err
	}
	// NOTE(nic): the above clone should have set these, but doesn't.  Future
	//  versions of go-git will likely fix this, but work around for now
	cfg, err := repo.Storer.Config()
	if err != nil {
		return err
	}
	cfg.Raw.Section("remote").Subsection(origin).
		SetOption("promisor", "true").
		SetOption("partialclonefilter", string(opts.Filter))
	err = repo.Storer.SetConfig(cfg)
	return err
}

// Update the repo.  Fetch the current HEAD and any new tags that may have
// appeared, and update the cache.
func (g *Git) Update(origin, cloneDir string) error {
	repo, err := g.gitOpen(cloneDir)
	if err != nil {
		return err
	}
	err = repo.Fetch(&git.FetchOptions{
		RemoteName: origin,
		RefSpecs: []config.RefSpec{
			`+refs/heads/*:refs/heads/*`,
			`+refs/tags/*:refs/tags/*`,
		},
		Tags:  git.AllTags,
		Force: true,
	})
	g.logger.Debug(
		"git.Update",
		slog.String("cloneDir", cloneDir),
		slog.String("origin", origin),
		slog.Any("error", err),
	)
	if err == git.NoErrAlreadyUpToDate {
		return nil
	}
	return err
}

// Worktree create a working tree from the repo in `cloneDir` at `version` in `dstDir`.
// Under the covers, this will download any/all required objects from origin
// into the cache
func (g *Git) Worktree(
	cloneDir string,
	version string,
	dstDir string,
) error {
	dst, err := g.appFs.Abs(dstDir)
	if err != nil {
		return err
	}

	g.logger.Info(
		"extracting",
		slog.String("from", cloneDir),
		slog.String("version", version),
		slog.String("to", dst),
	)

	repo, err := g.gitOpen(cloneDir)
	if err != nil {
		return err
	}

	hash, err := repo.ResolveRevision(plumbing.Revision(version))
	g.logger.Debug(
		"repo.ResolveRevision",
		slog.String("version", version),
		slog.String("hash", hash.String()),
		slog.Any("err", err),
	)
	if err != nil {
		return err
	}

	// Create a worktree for `version` at `dst`.  We'll call it `gilt` since it
	// needs a name; the lock file should prevent collisions
	const worktreeName = "gilt"
	wt, err := g.gitWorktree(repo.Storer)
	if err != nil {
		return err
	}
	_ = wt.Remove(worktreeName)

	_ = g.appFs.RemoveAll(dst)
	dstFS := osfs.New(dst)
	err = wt.Add(dstFS, worktreeName, xworktree.WithCommit(*hash), xworktree.WithDetachedHead())
	// NOTE(nic): we never need to track these, so remove the worktree and .git breadcrumb
	defer func() {
		_ = wt.Remove(worktreeName)
		_ = dstFS.Remove(git.GitDirName)
	}()
	g.logger.Debug("wt.Add", slog.Any("hash", hash), slog.Any("err", err))
	return err
}

// RemoteExists checks if the remote exists in the given cloneDir.
func (g *Git) RemoteExists(cloneDir string, remote string) (bool, error) {
	repo, err := g.gitOpen(cloneDir)
	if err != nil {
		return false, err
	}
	_, err = repo.Remote(remote)
	if err != nil {
		g.logger.Debug(
			"git.RemoteExists",
			slog.String("cloneDir", cloneDir),
			slog.String("remote", remote),
			slog.Any("error", err),
		)
		if err == git.ErrRemoteNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
