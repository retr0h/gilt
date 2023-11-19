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
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"

	"github.com/retr0h/go-gilt/internal"
	"github.com/retr0h/go-gilt/internal/config"
)

// New factory to create a new Repository instance.
func New(
	appFs afero.Fs,
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

// glob searches for files matching pattern in the directory dir
// and appends them to matches, returning the updated slice.
// If the directory cannot be opened, glob returns the existing matches.
// New matches are added in lexicographical order.
// Inspired by io/fs/glob since afero doesn't support filepath.Glob.
func glob(
	fs afero.Fs,
	dir string,
	pattern string,
) ([]string, error) {
	m := []string{}
	infos, err := afero.ReadDir(fs, dir)
	if err != nil {
		return nil, err
	}

	for _, info := range infos {
		n := info.Name()
		matched, err := filepath.Match(pattern, n)
		if err != nil {
			return m, err
		}

		if matched {
			m = append(m, filepath.Join(dir, n))
		}
	}

	return m, nil
}

// Clone clone Repository.Git to Repository.getCloneDir, and hard checkout
// to Repository.Version.
func (r *Repository) Clone(
	c config.Repository,
	cloneDir string,
) error {
	r.logger.Info(
		"cloning",
		slog.String("repository", c.Git),
		slog.String("version", c.Version),
		slog.String("dstDir", cloneDir),
	)

	if _, err := r.appFs.Stat(cloneDir); os.IsNotExist(err) {
		if err := r.gitManager.Clone(c.Git, cloneDir); err != nil {
			return err
		}

		if err := r.gitManager.Reset(cloneDir, c.Version); err != nil {
			return err
		}
	} else {
		r.logger.Warn(
			"clone already exists",
			slog.String("dstDir", cloneDir),
		)
	}

	return nil
}

// CheckoutIndex checkout Repository.Git to Repository.DstDir.
func (r *Repository) CheckoutIndex(
	c config.Repository,
	cloneDir string,
) error {
	return r.gitManager.CheckoutIndex(c.DstDir, cloneDir)
}

// CopySources copy Repository.Src to Repository.DstFile or Repository.DstDir.
func (r *Repository) CopySources(
	c config.Repository,
	cloneDir string,
) error {
	for _, source := range c.Sources {
		parts := strings.Split(
			source.Src,
			string(os.PathSeparator),
		) // break up source.Src path
		head := parts[0 : len(parts)-1] // take all path parts but last
		tail := parts[len(parts)-1]     // take the last path part
		cloneDirWithSrcPath := filepath.Join(
			cloneDir,
			strings.Join(head, string(os.PathSeparator)),
		) // join clone dir with head
		globbedSrc, err := glob(
			r.appFs,
			cloneDirWithSrcPath,
			tail,
		) // tail is used by glob for path matching
		if err != nil {
			return err
		}

		for _, src := range globbedSrc {
			// The source is a file
			if info, err := r.appFs.Stat(src); err == nil && info.Mode().IsRegular() {
				// ... and the dst is declared a directory
				if source.DstFile != "" {
					if err := r.copyManager.CopyFile(src, source.DstFile); err != nil {
						return err
					}
				} else if source.DstDir != "" {
					// ... and create te dst directory
					if err := r.appFs.MkdirAll(source.DstDir, 0o755); err != nil {
						return fmt.Errorf("unable to create dest dir: %s", err)
					}
					srcBaseFile := filepath.Base(src)
					newDst := filepath.Join(source.DstDir, srcBaseFile)
					if err := r.copyManager.CopyFile(src, newDst); err != nil {
						return err
					}
				}
				// The source is a directory
			} else if info, err := r.appFs.Stat(src); err == nil && info.Mode().IsDir() {
				// ... and dst dir exists
				if info, err := r.appFs.Stat(source.DstDir); err == nil && info.Mode().IsDir() {
					if err := r.appFs.RemoveAll(source.DstDir); err != nil {
						return err
					}
				}
				if err := r.copyManager.CopyDir(src, source.DstDir); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
