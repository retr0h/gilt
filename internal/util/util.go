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

package util

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/logrusorgru/aurora/v4"

	"github.com/retr0h/go-gilt/internal/io"
)

// currentUser is mocked for tests.
var currentUser = user.Current

// ExpandUser expands the path to user's home directory when path begins with '~'.
func ExpandUser(path string) (string, error) {
	if len(path) == 0 || path[0] != '~' {
		return path, nil
	}

	usr, err := currentUser()
	if err != nil {
		return "", err
	}

	return filepath.Join(usr.HomeDir, path[1:]), nil
}

// RunCmd execute the provided command with args.
// Yeah, yeah, yeah, I know I cheated by using Exec in this package.
func RunCmd(debug bool, name string, args ...string) error {
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

// CopyFile copies src file to dst.
func CopyFile(src string, dst string) error {
	baseSrc := filepath.Base(src)
	msg := fmt.Sprintf(
		"%-4s - Copying file '%s' to '%s'",
		"",
		aurora.Cyan(baseSrc),
		aurora.Cyan(dst),
	)
	fmt.Println(msg)

	if err := io.CopyFile(src, dst); err != nil {
		return err
	}

	return nil
}

// CopyDir copies src directory to dst.
func CopyDir(src string, dst string) error {
	baseSrc := filepath.Base(src)
	msg := fmt.Sprintf(
		"%-4s - Copying dir '%s' to '%s'",
		"",
		aurora.Cyan(baseSrc),
		aurora.Cyan(dst),
	)
	fmt.Println(msg)

	if err := io.CopyDir(src, dst); err != nil {
		return err
	}

	return nil
}
