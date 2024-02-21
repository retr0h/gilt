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
	"strings"

	"github.com/avfs/avfs"

	"github.com/retr0h/gilt/v2/internal"
	intPath "github.com/retr0h/gilt/v2/internal/path"
	"github.com/retr0h/gilt/v2/pkg/config"
)

// New factory to create a new Repository instance.
func New(
	appFs avfs.VFS,
	c config.Repositories,
	repoManager internal.RepositoryManager,
	execManager internal.ExecManager,
	logger *slog.Logger,
) *Repositories {
	return &Repositories{
		appFs:       appFs,
		config:      c,
		repoManager: repoManager,
		execManager: execManager,
		logger:      logger,
	}
}

func (r *Repositories) getGiltDir() (string, error) {
	giltDir, err := intPath.ExpandUser(r.config.GiltDir)
	return giltDir, err
}

// getCacheDir create the cacheDir if it doesn't exist.
func (r *Repositories) getCacheDir() (string, error) {
	giltDir, err := r.getGiltDir()
	if err != nil {
		return "", err
	}

	cacheDir := r.appFs.Join(giltDir, "cache")
	if err := r.appFs.MkdirAll(cacheDir, 0o700); err != nil {
		if i, e := r.appFs.Stat(cacheDir); e == nil && i.IsDir() {
			return cacheDir, nil
		}
		return "", err
	}
	return cacheDir, nil
}

// Overlay clone and extract the Repository items.
func (r *Repositories) Overlay() error {
	cacheDir, err := r.getCacheDir()
	if err != nil {
		r.logger.Error(
			"error expanding dir",
			slog.String("giltDir", r.config.GiltDir),
			slog.String("cacheDir", cacheDir),
			slog.String("err", err.Error()),
		)
		return err
	}

	for _, c := range r.config.Repositories {
		targetDir, err := r.repoManager.Clone(c, cacheDir)
		if err != nil {
			return err
		}

		// Easy mode: create a full worktree, directly in DstDir
		if c.DstDir != "" {
			// delete dstDir since `git worktree add` will not replace existing directories
			if info, err := r.appFs.Stat(c.DstDir); err == nil && info.Mode().IsDir() {
				if err := r.appFs.RemoveAll(c.DstDir); err != nil {
					return err
				}
			}
			if err := r.repoManager.Worktree(c, targetDir, c.DstDir); err != nil {
				return err
			}
		}

		// Hard mode: copy subtrees of the worktree from Repository.Src to
		// Repository.DstDir (or Repository.DstFile)
		if len(c.Sources) > 0 {
			giltDir, err := r.getGiltDir()
			if err != nil {
				return err
			}
			err = r.execManager.RunInTempDir(giltDir, "tmp", func(tmpDir string) error {
				tmpClone := r.appFs.Join(tmpDir, r.appFs.Base(targetDir))
				if err := r.repoManager.Worktree(c, targetDir, tmpClone); err != nil {
					return err
				}
				return r.repoManager.CopySources(c, tmpClone)
			})
			if err != nil {
				return err
			}
		}

		// run post commands
		if len(c.Commands) > 0 {
			for _, command := range c.Commands {
				r.logger.Info(
					"executing command",
					slog.String("cmd", command.Cmd),
					slog.String("args", strings.Join(command.Args, " ")),
				)
				if err := r.execManager.RunCmd(command.Cmd, command.Args); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
