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
// package will remain using exec as it was easiest to port from go-gilt's
// python counterpart.
package git

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/afero"

	"github.com/retr0h/go-gilt/internal"
)

// New factory to create a new Git instance.
func New(
	appFs afero.Fs,
	debug bool,
	execManager internal.ExecManager,
	logger *slog.Logger,
) *Git {
	return &Git{
		appFs:       appFs,
		debug:       debug,
		execManager: execManager,
		logger:      logger,
	}
}

// Clone as exec manager to clone repo.
func (g *Git) Clone(
	gitURL string,
	cloneDir string,
) error {
	// return g.execManager.Clone(gitURL, cloneDir)
	return g.execManager.RunCmd("git", "clone", gitURL, cloneDir)
}

// Reset to the given git version.
func (g *Git) Reset(
	cloneDir string,
	gitVersion string,
) error {
	return g.execManager.RunCmd("git", "-C", cloneDir, "reset", "--hard", gitVersion)
}

// CheckoutIndex checkout Repository.Git to Repository.DstDir.
func (g *Git) CheckoutIndex(
	dstDir string,
	cloneDir string,
) error {
	dst, err := filepath.Abs(dstDir)
	if err != nil {
		return err
	}

	g.logger.Info(
		"extracting",
		slog.String("to", dst),
	)

	cmdArgs := []string{
		"-C",
		cloneDir,
		"checkout-index",
		"--force",
		"--all",
		"--prefix",
		// Trailing separator needed by git checkout-index.
		dst + string(os.PathSeparator),
	}

	return g.execManager.RunCmd("git", cmdArgs...)
}
