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

package repository

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/avfs/avfs"

	"github.com/retr0h/gilt/v2/internal"
	"github.com/retr0h/gilt/v2/pkg/config"
)

const ORIGIN = "gilt"

// We'll use this to normalize Git URLs as "safe" filenames
var replacer = strings.NewReplacer("/", "-", ":", "-")

// New factory to create a new Repository instance.
func New(
	appFs avfs.VFS,
	copyManager CopyManager,
	gitManager internal.GitManager,
	logger *slog.Logger,
) *Repository {
	return &Repository{
		appFs:       appFs,
		copyManager: copyManager,
		gitManager:  gitManager,
		logger:      logger,
	}
}

// Clone Repository.Git under Repository.getCloneDir
func (r *Repository) Clone(
	c config.Repository,
	cloneDir string,
) (string, error) {
	targetDir := r.appFs.Join(cloneDir, replacer.Replace(c.Git))
	remote, err := r.gitManager.Remote(targetDir)
	if err == nil && !strings.Contains(remote, ORIGIN) {
		r.logger.Info(
			"remote does not exist in clone, invalidating cache",
			slog.Any("remote", ORIGIN),
			slog.String("dstDir", targetDir),
		)
		_ = r.appFs.RemoveAll(targetDir)
		err = errors.New("cache does not exist")
	}
	if err != nil {
		r.logger.Info("cloning", slog.String("repository", c.Git), slog.String("dstDir", targetDir))
		if err := r.gitManager.Clone(c.Git, ORIGIN, targetDir); err != nil {
			return targetDir, err
		}
	} else {
		r.logger.Info("clone already exists", slog.String("dstDir", targetDir))
		if err := r.gitManager.Update(ORIGIN, targetDir); err != nil {
			return targetDir, err
		}
	}
	return targetDir, nil
}

// Worktree create a git workingtree at the given version in Repository.DstDir.
func (r *Repository) Worktree(
	c config.Repository,
	cloneDir string,
	targetDir string,
) error {
	return r.gitManager.Worktree(cloneDir, c.Version, targetDir)
}

// CopySources copy Repository.Src to Repository.DstFile or Repository.DstDir.
func (r *Repository) CopySources(
	c config.Repository,
	cloneDir string,
) error {
	r.logger.Debug("copy", slog.String("origin", cloneDir))
	for _, source := range c.Sources {
		cloneDirWithSrcPath := r.appFs.Join(cloneDir, source.Src) // join clone dir with head
		globbedSrc, err := r.appFs.Glob(cloneDirWithSrcPath)
		if err != nil {
			return err
		}

		for _, src := range globbedSrc {
			// The source is a directory
			if info, err := r.appFs.Stat(src); err == nil && info.IsDir() {
				// ... and dst dir exists
				if info, err := r.appFs.Stat(source.DstDir); err == nil && info.IsDir() {
					if err := r.appFs.RemoveAll(source.DstDir); err != nil {
						return err
					}
				}
				if err := r.copyManager.CopyDir(src, source.DstDir); err != nil {
					return err
				}
			} else {
				if source.DstFile != "" {
					if err := r.copyManager.CopyFile(src, source.DstFile); err != nil {
						return err
					}
				} else if source.DstDir != "" {
					if err := r.appFs.MkdirAll(source.DstDir, 0o755); err != nil {
						return fmt.Errorf("unable to create dest dir: %s", err)
					}
					srcBaseFile := r.appFs.Base(src)
					newDst := r.appFs.Join(source.DstDir, srcBaseFile)
					if err := r.copyManager.CopyFile(src, newDst); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}
