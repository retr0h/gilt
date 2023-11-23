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

package git

import (
	"log/slog"
	"os/exec"
	"strings"
)

// NewExecManagerCmd factory to create a new exec manager instance.
func NewExecManagerCmd(
	debug bool,
	logger *slog.Logger,
) *ExecManagerCmd {
	return &ExecManagerCmd{
		debug:  debug,
		logger: logger,
	}
}

// RunCmd execute the provided command with args.
// Yeah, yeah, yeah, I know I cheated by using Exec in this package.
func (e *ExecManagerCmd) RunCmd(
	name string,
	args ...string,
) error {
	cmd := exec.Command(name, args...)

	if e.debug {
		commands := strings.Join(cmd.Args, " ")
		e.logger.Debug(
			"exec",
			slog.String("command", commands),
		)
	}

	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}
