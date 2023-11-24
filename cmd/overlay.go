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

package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"strconv"

	"github.com/danjacques/gofslock/fslock"
	"github.com/spf13/cobra"

	giltpath "github.com/retr0h/go-gilt/internal/path"
)

// getGiltDir create the GiltDir if it doesn't exist.
func getGiltDir() (string, error) {
	dir, err := giltpath.ExpandUser(appConfig.GiltDir)
	if err != nil {
		return "", err
	}

	if err := appFs.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}

	return dir, nil
}

// withLock is a convenience function to create a lock, execute a function while
// holding that lock, and then release the lock on completion.
func withLock(fn func() error) error {
	lockDir, err := getGiltDir()
	if err != nil {
		return err
	}

	lockFile := filepath.Join(lockDir, "gilt.lock")
	logger.Info(
		"acquiring lock",
		slog.String("lockfile", lockFile),
	)

	err = fslock.With(lockFile, fn)
	if err != nil {
		if errors.Is(err, fslock.ErrLockHeld) {
			return fmt.Errorf("could not acquire lock on %s: %s", lockFile, err)
		}
	}

	return err
}

func logRepositoriesGroup() []any {
	logGroups := make([]any, 0, len(appConfig.Repositories))

	for i, repo := range appConfig.Repositories {
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

// overlayCmd represents the overlay command
var overlayCmd = &cobra.Command{
	Use:   "overlay",
	Short: "Install gilt dependencies",
	RunE: func(cmd *cobra.Command, args []string) error {
		// By the time we reach this point, we know that the arguments were
		// properly parsed, and we don't want to show the usage if an error
		// occurs
		cmd.SilenceUsage = true
		// We are logging errors, no need for cobra to re-log the error
		cmd.SilenceErrors = true

		logger.Debug(
			"current configuration",
			slog.String("GiltDir", appConfig.GiltDir),
			slog.String("GiltFile", appConfig.GiltFile),
			slog.Bool("Debug", appConfig.Debug),
			slog.Group("Repository", logRepositoriesGroup()...),
		)

		if err := withLock(func() error {
			return repos.Overlay()
		}); err != nil {
			logger.Error(
				"error overlaying repositories",
				slog.String("err", err.Error()),
			)
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(overlayCmd)
}
