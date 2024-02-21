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

func (e *Exec) RunCmdImpl(
	name string,
	args []string,
	cwd string,
) error {
	cmd := exec.Command(name, args...)
	if cwd != "" {
		cmd.Dir = cwd
	}

	commands := strings.Join(cmd.Args, " ")
	e.logger.Debug("exec", slog.String("command", commands), slog.String("cwd", cwd))

	out, err := cmd.CombinedOutput()
	e.logger.Debug("result", slog.String("output", string(out)))
	if err != nil {
		return err
	}

	return nil
}

// RunCmd execute the provided command with args.
// Yeah, yeah, yeah, I know I cheated by using Exec in this package.
func (e *Exec) RunCmd(
	name string,
	args []string,
) error {
	return e.RunCmdImpl(name, args, "")
}

func (e *Exec) RunCmdInDir(
	name string,
	args []string,
	cwd string,
) error {
	return e.RunCmdImpl(name, args, cwd)
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
