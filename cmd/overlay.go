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

	"github.com/retr0h/go-gilt/internal/util"
	"github.com/spf13/cobra"
)

// getGiltDir create the GiltDir if it doesn't exist.
func getGiltDir() (string, error) {
	expandedGiltDir, err := util.ExpandUser(r.GiltDir)
	if err != nil {
		return "", err
	}

	cacheGiltDir := filepath.Join(expandedGiltDir, "cache")

	if _, err := os.Stat(cacheGiltDir); os.IsNotExist(err) {
		if err := os.Mkdir(cacheGiltDir, 0o755); err != nil {
			return "", err
		}
	}

	return cacheGiltDir, nil
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

		cacheDir, err := getGiltDir()
		if err != nil {
			logger.Error(
				"error expanding dir",
				slog.String("giltDir", r.GiltDir),
				slog.String("err", err.Error()),
			)
			return err
		}

		r.GiltDir = cacheDir
		if err := r.Overlay(); err != nil {
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
