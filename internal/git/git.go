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
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/logrusorgru/aurora/v4"

	"github.com/retr0h/go-gilt/internal/repository"
)

var (
	// RunCommand is mocked for tests.
	RunCommand = runCmd
	// FilePathAbs is mocked for tests.
	FilePathAbs = filepath.Abs
)

// NewGit factory to create a new Git instance.
func NewGit(debug bool) *Git {
	return &Git{
		debug: debug,
	}
}

// Clone clone Repository.Git to Repository.getCloneDir, and hard checkout
// to Repository.Version.
func (g *Git) Clone(repository repository.Repository) error {
	cloneDir := repository.GetCloneDir()

	msg := fmt.Sprintf(
		"[%s@%s]:",
		aurora.Magenta(repository.Git),
		aurora.Magenta(repository.Version),
	)
	fmt.Println(msg)

	msg = fmt.Sprintf("%-2s - Cloning to '%s'", "", aurora.Cyan(cloneDir))
	fmt.Println(msg)

	if _, err := os.Stat(cloneDir); os.IsNotExist(err) {
		if err := g.clone(repository); err != nil {
			return err
		}

		if err := g.reset(repository); err != nil {
			return err
		}
	} else {
		bang := aurora.Bold(aurora.Red("!"))
		msg := fmt.Sprintf("%-2s %s %s", "", bang, aurora.Yellow("Clone already exists"))
		fmt.Println(msg)
	}

	return nil
}

func (g *Git) clone(repository repository.Repository) error {
	cloneDir := repository.GetCloneDir()
	err := RunCommand(g.debug, "git", "clone", repository.Git, cloneDir)

	return err
}

func (g *Git) reset(repository repository.Repository) error {
	cloneDir := repository.GetCloneDir()
	err := RunCommand(g.debug, "git", "-C", cloneDir, "reset", "--hard", repository.Version)

	return err
}

// CheckoutIndex checkout Repository.Git to Repository.DstDir.
func (g *Git) CheckoutIndex(repository repository.Repository) error {
	cloneDir := repository.GetCloneDir()
	dstDir, err := FilePathAbs(repository.DstDir)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("%-2s - Extracting to '%s'", "", aurora.Cyan(dstDir))
	fmt.Println(msg)

	cmdArgs := []string{
		"-C",
		cloneDir,
		"checkout-index",
		"--force",
		"--all",
		"--prefix",
		// Trailing separator needed by git checkout-index.
		dstDir + string(os.PathSeparator),
	}

	return RunCommand(g.debug, "git", cmdArgs...)
}

// runCmd execute the provided command with args.
// Yeah, yeah, yeah, I know I cheated by using Exec in this package.
func runCmd(debug bool, name string, args ...string) error {
	var stderr bytes.Buffer
	cmd := exec.Command(name, args...)

	if debug {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		commands := strings.Join(cmd.Args, " ")
		msg := fmt.Sprintf("COMMAND: %s", aurora.Colorize(commands, aurora.BlackFg|aurora.RedBg))
		fmt.Println(msg)
	} else {
		cmd.Stderr = &stderr
	}

	if err := cmd.Run(); err != nil {
		return errors.New(stderr.String())
	}

	return nil
}
