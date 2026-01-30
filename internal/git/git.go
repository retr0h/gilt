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
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/config"
	"github.com/go-git/go-git/v6/plumbing/protocol/packp"

	"github.com/retr0h/gilt/v2/internal"
)

// New factory to create a new Git instance.
func New(
	appFs avfs.VFS,
	execManager internal.ExecManager,
	logger *slog.Logger,
) *Git {
	return &Git{
		appFs:       appFs,
		execManager: execManager,
		logger:      logger,
	}
}

// Clone the repo.  This is a bare repo, with only metadata to start with.
func (g *Git) Clone(gitURL, origin, cloneDir string) error {
	_, err := git.PlainClone(cloneDir, &git.CloneOptions{
		URL:        gitURL,
		RemoteName: origin,
		Bare:       true,
		NoCheckout: true,
		Filter:     packp.FilterBlobNone(),
	})
	g.logger.Debug(
		"git.Clone",
		slog.String("gitURL", gitURL),
		slog.String("origin", origin),
		slog.String("cloneDir", cloneDir),
		slog.Any("err", err),
	)
	return err
}

// Update the repo.  Fetch the current HEAD and any new tags that may have
// appeared, and update the cache.
func (g *Git) Update(origin, cloneDir string) error {
	repo, err := git.PlainOpen(cloneDir)
	if err != nil {
		return err
	}
	err = repo.Fetch(&git.FetchOptions{
		RemoteName: origin,
		RefSpecs: []config.RefSpec{
			"+refs/heads/*:refs/heads/*",
			"+refs/tags/*:refs/tags/*",
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

	_, err = g.execManager.RunCmdInDir(
		"git",
		[]string{"worktree", "add", "--force", dst, version},
		cloneDir,
	)
	// `git worktree add` creates a breadcrumb file back to the original repo;
	// this is just junk data in our use case, so get rid of it
	if err == nil {
		_ = g.appFs.Remove(g.appFs.Join(dst, ".git"))
		_, _ = g.execManager.RunCmdInDir(
			"git",
			[]string{"worktree", "prune", "--verbose"},
			cloneDir,
		)
	}
	return err
}

// Check if the remote exists in the given cloneDir.
func (g *Git) RemoteExists(cloneDir string, remote string) (bool, error) {
	repo, err := git.PlainOpen(cloneDir)
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
