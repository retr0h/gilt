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
	"log/slog"
	"os"
	"path/filepath"

	"github.com/retr0h/go-gilt/internal/repositories"
	"github.com/retr0h/go-gilt/internal/util"
	"github.com/spf13/cobra"
)

var fileName string

// newRepositories constructs a new `repositories.Repositories`.
func newRepositories(debug bool, fileName string) (*repositories.Repositories, error) {
	expandedFileName, err := util.ExpandUser(fileName)
	if err != nil {
		return nil, err
	}

	return &repositories.Repositories{
		Debug:    debug,
		Filename: expandedFileName,
	}, nil
}

// newGiltDir create the GiltDir if it doesn't exist.
func newGiltDir() error {
	expandedGiltDir, err := util.ExpandUser(giltDir)
	if err != nil {
		return err
	}

	cacheGiltDir := filepath.Join(expandedGiltDir, "cache")
	repositories.GiltDir = cacheGiltDir

	if _, err := os.Stat(cacheGiltDir); os.IsNotExist(err) {
		if err := os.Mkdir(cacheGiltDir, 0o755); err != nil {
			return err
		}
	}

	return nil
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

		r, err := newRepositories(debug, fileName)
		if err != nil {
			logger.Error(
				"error creating new Repositories",
				slog.String("err", err.Error()),
			)
			return err
		}

		if err := newGiltDir(); err != nil {
			logger.Error(
				"error expanding dir",
				slog.String("giltDir", giltDir),
				slog.String("err", err.Error()),
			)
			return err
		}

		if err := r.UnmarshalYAMLFile(); err != nil {
			logger.Error(
				"error occurred unmarshalling",
				slog.String("err", err.Error()),
			)
			return err
		}

		if err := r.Overlay(); err != nil {
			logger.Error(
				"error cloning repository",
				slog.String("err", err.Error()),
			)
			return err

		}

		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().
		StringVarP(&fileName, "filename", "f", "gilt.yml", "Path to config file")
	rootCmd.AddCommand(overlayCmd)
}
