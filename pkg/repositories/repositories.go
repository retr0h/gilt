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

package repositories

import (
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/avfs/avfs/vfs/osfs"
	"github.com/danjacques/gofslock/fslock"

	"github.com/retr0h/gilt/v2/internal/exec"
	"github.com/retr0h/gilt/v2/internal/git"
	intPath "github.com/retr0h/gilt/v2/internal/path"
	intRepos "github.com/retr0h/gilt/v2/internal/repositories"
	"github.com/retr0h/gilt/v2/internal/repository"
	"github.com/retr0h/gilt/v2/pkg/config"
)

// New factory to create a new Repository instance.
func New(
	c config.Repositories,
	logger *slog.Logger,
) *Repositories {
	appFs := osfs.NewWithNoIdm()

	copyManager := repository.NewCopy(
		appFs,
		logger,
	)

	execManager := exec.New(
		appFs,
		logger,
	)

	gitManager := git.New(
		appFs,
		execManager,
		logger,
	)

	repoManager := repository.New(
		appFs,
		copyManager,
		gitManager,
		logger,
	)

	reposManager := intRepos.New(
		appFs,
		c,
		repoManager,
		execManager,
		logger,
	)

	return &Repositories{
		appFs:        appFs,
		c:            c,
		reposManager: reposManager,
		logger:       logger,
	}
}

// getGiltDir create the GiltDir if it doesn't exist.
func (r *Repositories) getGiltDir() (string, error) {
	dir, err := intPath.ExpandUser(r.c.GiltDir)
	if err != nil {
		return "", err
	}

	if err := r.appFs.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}

	return dir, nil
}

// withLock is a convenience function to create a lock, execute a function while
// holding that lock, and then release the lock on completion.
func (r *Repositories) withLock(fn func() error) error {
	lockDir, err := r.getGiltDir()
	if err != nil {
		return err
	}

	blocker := func() error {
		time.Sleep(time.Millisecond)
		return nil
	}
	lockFile := r.appFs.Join(lockDir, "gilt.lock")
	r.logger.Info(
		"acquiring lock",
		slog.String("lockfile", lockFile),
	)
	err = fslock.WithBlocking(lockFile, blocker, fn)
	if err != nil {
		if errors.Is(err, fslock.ErrLockHeld) {
			return fmt.Errorf("could not acquire lock on %s: %s", lockFile, err)
		}
	}

	return err
}

// logRepositoriesGroup log Repositories config.
func (r *Repositories) logRepositoriesGroup() []any {
	logGroups := make([]any, 0, len(r.c.Repositories))

	for i, repo := range r.c.Repositories {
		var sourceGroups []any
		for i, s := range repo.Sources {
			group := slog.Group(strconv.Itoa(i),
				slog.String("Src", s.Src),
				slog.String("DstFile", s.DstFile),
				slog.String("DstDir", s.DstDir),
			)
			sourceGroups = append(sourceGroups, group)
		}
		var cmdGroups []any
		for i, s := range repo.Commands {
			group := slog.Group(strconv.Itoa(i),
				slog.String("Cmd", s.Cmd),
			)
			cmdGroups = append(cmdGroups, group)
		}

		group := slog.Group(strconv.Itoa(i),
			slog.String("Git", repo.Git),
			slog.String("Version", repo.Version),
			slog.String("DstDir", repo.DstDir),
			slog.Group("Sources", sourceGroups...),
			slog.Group("Commands", cmdGroups...),
		)
		logGroups = append(logGroups, group)
	}

	return logGroups
}

// Overlay clone and extract the Repository items.
func (r *Repositories) Overlay() error {
	if err := r.withLock(func() error {
		r.logger.Debug(
			"current configuration",
			slog.String("GiltDir", r.c.GiltDir),
			slog.String("GiltFile", r.c.GiltFile),
			slog.Bool("Debug", r.c.Debug),
			slog.Bool("Parallel", r.c.Parallel),
			slog.Group("Repository", r.logRepositoriesGroup()...),
		)

		return r.reposManager.Overlay()
	}); err != nil {
		r.logger.Error(
			"error overlaying repositories",
			slog.String("err", err.Error()),
		)
		return err
	}

	return nil
}
