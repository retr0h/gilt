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
	"runtime"
	"strings"
	"sync"

	"github.com/avfs/avfs"

	"github.com/retr0h/gilt/v2/internal"
	intPath "github.com/retr0h/gilt/v2/internal/path"
	"github.com/retr0h/gilt/v2/pkg/config"
)

// This should be a nice upper bound for parallel fetches
const maxSlots = 8

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
		cloneCache:  make(map[string]string),
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
	if err := r.populateCloneCache(r.config.Parallel); err != nil {
		return err
	}

	for _, c := range r.config.Repositories {
		targetDir := r.cloneCache[c.Git]

		// Easy mode: create a full worktree, directly in DstDir
		if err := r.overlayTree(c, targetDir); err != nil {
			return err
		}

		// Hard mode: copy subtrees of the worktree from Repository.Src to
		// Repository.DstDir (or Repository.DstFile)
		if err := r.overlaySubtrees(c, targetDir); err != nil {
			return err
		}

		// run post commands
		if err := r.runCommands(c); err != nil {
			return err
		}
	}

	return nil
}

// populateCloneCache ensure that all named repos exist and are up-to-date
func (r *Repositories) populateCloneCache(parallel bool) error {
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

	// Run all the clones concurrently (1 coroutine per CPU), up to 8 workers
	slots := 1
	if parallel {
		slots = min(maxSlots, runtime.GOMAXPROCS(0))
	}
	var wg sync.WaitGroup
	var mu sync.Mutex                                       // Mutex to protect cloneCache
	errChan := make(chan error, len(r.config.Repositories)) // Channel to collect errors
	semaphore := make(chan struct{}, slots)                 // Semaphore to limit concurrency

	for _, repo := range r.config.Repositories {
		wg.Add(1)
		go func(c config.Repository) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			if err := r.runPopulate(c, cacheDir, &mu); err != nil {
				errChan <- err
			}
		}(repo)
	}

	// Roll up any errors fetching the above
	wg.Wait()
	close(errChan)
	return r.anyErrors(errChan)
}

func (r *Repositories) runPopulate(c config.Repository, cacheDir string, mu *sync.Mutex) error {
	mu.Lock()
	if _, exists := r.cloneCache[c.Git]; exists {
		mu.Unlock()
		return nil
	}
	// Set a "stub" value to claim territory
	// This worker is now responsible for populating the "full" value
	r.cloneCache[c.Git] = ""
	mu.Unlock()

	// Initialize and/or update the clone (long-running operation outside the lock)
	targetDir, err := r.repoManager.Clone(c, cacheDir)
	if err != nil {
		return err
	}

	// Rewrite with the "full" value
	mu.Lock()
	r.cloneCache[c.Git] = targetDir
	mu.Unlock()

	return nil
}

func (r *Repositories) anyErrors(errChan <-chan error) error {
	for err := range errChan {
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repositories) overlayTree(c config.Repository, targetDir string) error {
	if c.DstDir == "" {
		return nil
	}
	// delete DstDir since `git worktree add` will not replace existing directories
	if info, err := r.appFs.Stat(c.DstDir); err == nil && info.IsDir() {
		if err := r.appFs.RemoveAll(c.DstDir); err != nil {
			return err
		}
	}
	if err := r.repoManager.Worktree(c, targetDir, c.DstDir); err != nil {
		return err
	}
	return nil
}

func (r *Repositories) overlaySubtrees(c config.Repository, targetDir string) error {
	if len(c.Sources) == 0 {
		return nil
	}
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
	return nil
}

func (r *Repositories) runCommands(c config.Repository) error {
	if len(c.Commands) == 0 {
		return nil
	}
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
	return nil
}
