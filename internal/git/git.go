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
	"os"
	"path/filepath"

	"github.com/spf13/afero"

	"github.com/retr0h/gilt/internal"
)

// New factory to create a new Git instance.
func New(
	appFs afero.Fs,
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
func (g *Git) Clone(
	gitURL string,
	cloneDir string,
) error {
	return g.execManager.RunCmd(
		"git",
		[]string{"clone", "--bare", "--filter=blob:none", gitURL, cloneDir},
	)
}

// Worktree create a working tree from the repo in `cloneDir` at `version` in `dstDir`.
// Under the covers, this will download any/all required objects from origin
// into the cache
func (g *Git) Worktree(
	cloneDir string,
	version string,
	dstDir string,
) error {
	dst, err := filepath.Abs(dstDir)
	if err != nil {
		return err
	}

	g.logger.Info(
		"extracting",
		slog.String("from", cloneDir),
		slog.String("version", version),
		slog.String("to", dst),
	)

	err = g.execManager.RunCmdInDir(
		"git",
		[]string{"worktree", "add", "--force", dst, version},
		cloneDir,
	)
	// `git worktree add` creates a breadcrumb file back to the original repo;
	// this is just junk data in our use case, so get rid of it
	if err == nil {
		_ = os.Remove(filepath.Join(dst, ".git"))
	}
	return err
}
