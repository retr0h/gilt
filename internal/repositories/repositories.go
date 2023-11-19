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
	"fmt"
	"log/slog"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"

	"github.com/retr0h/go-gilt/internal"
	"github.com/retr0h/go-gilt/internal/config"
)

var currentUser = user.Current

// New factory to create a new Repository instance.
func New(
	appFs afero.Fs,
	c config.Repositories,
	repoManager internal.RepositoryManager,
	logger *slog.Logger,
) *Repositories {
	return &Repositories{
		appFs:       appFs,
		config:      c,
		repoManager: repoManager,
		logger:      logger,
	}
}

func expandUser(
	path string,
) (string, error) {
	if len(path) == 0 || path[0] != '~' {
		return path, nil
	}

	usr, err := currentUser()
	if err != nil {
		return "", err
	}

	return filepath.Join(usr.HomeDir, path[1:]), nil
}

// getCloneDir returns the path to the Repository's clone directory.
func (r *Repositories) getCloneDir(
	giltDir string,
	c config.Repository,
) string {
	return filepath.Join(giltDir, r.getCloneHash(c))
}

func (r *Repositories) getCloneHash(
	c config.Repository,
) string {
	replacer := strings.NewReplacer(
		"/", "-",
		":", "-",
	)
	replacedGitURL := replacer.Replace(c.Git)

	return fmt.Sprintf("%s-%s", replacedGitURL, c.Version)
}

// getGiltDir create the GiltDir if it doesn't exist.
func (r *Repositories) getGiltDir() (string, error) {
	expandedGiltDir, err := expandUser(r.config.GiltDir)
	if err != nil {
		return "", err
	}

	cacheGiltDir := filepath.Join(expandedGiltDir, "cache")
	if _, err := r.appFs.Stat(cacheGiltDir); os.IsNotExist(err) {
		if err := r.appFs.Mkdir(cacheGiltDir, 0o755); err != nil {
			return "", err
		}
	}

	return cacheGiltDir, nil
}

// Overlay clone and extract the Repository items.
func (r *Repositories) Overlay() error {
	cacheDir, err := r.getGiltDir()
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
		cloneDir := r.getCloneDir(cacheDir, c)
		err = r.repoManager.Clone(c, cloneDir)
		if err != nil {
			return err
		}

		// checkout into c.DstDir
		if c.DstDir != "" {
			// delete dstDir since Checkout-Index does not clean old files that may
			// no longer exist in config
			if info, err := r.appFs.Stat(c.DstDir); err == nil && info.Mode().IsDir() {
				if err := os.RemoveAll(c.DstDir); err != nil {
					return err
				}
			}
			if err := r.repoManager.CheckoutIndex(c, cloneDir); err != nil {
				return err
			}
		}

		// copy sources from Repository.Src to Repository.DstDir or Repository.DstFile
		if len(c.Sources) > 0 {
			if err := r.repoManager.CopySources(c, cloneDir); err != nil {
				return err
			}
		}
	}

	return nil
}
