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

// Package exec provides command execution abstractions.
package exec

import (
	"log/slog"
	"os/exec"
	"strings"

	"github.com/avfs/avfs"
)

// New factory to create a new Exec instance.
func New(
	appFs avfs.VFS,
	logger *slog.Logger,
) *Exec {
	return &Exec{
		appFs:  appFs,
		logger: logger,
	}
}

// RunCmd execute the provided command with args.
func (e *Exec) RunCmd(
	name string,
	args []string,
) (string, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	e.logger.Debug(
		"exec",
		slog.String("command", strings.Join(cmd.Args, " ")),
		slog.String("output", string(out)),
		slog.Any("error", err),
	)
	return string(out), err
}

// RunInTempDir creates a temporary directory, and runs the provided function
// with the name of the directory as input.  Then it cleans up the temporary
// directory.
func (e *Exec) RunInTempDir(dir, pattern string, fn func(string) error) error {
	tmpDir, err := e.appFs.MkdirTemp(dir, pattern)
	if err != nil {
		return err
	}
	e.logger.Debug("created tempdir", slog.String("dir", tmpDir))

	// Ignoring errors as there's not much we can do and this is a cleanup function.
	defer func() {
		e.logger.Debug("removing tempdir", slog.String("dir", tmpDir))
		_ = e.appFs.RemoveAll(tmpDir)
	}()

	// Run the provided function.
	return fn(tmpDir)
}
