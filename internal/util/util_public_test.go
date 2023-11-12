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

package util_test

import (
	"os"
	"os/user"
	"path"
	"path/filepath"
	"testing"

	capturer "github.com/kami-zh/go-capturer"
	"github.com/stretchr/testify/assert"

	"github.com/retr0h/go-gilt/internal/util"
)

func TestExpandUserReturnsError(t *testing.T) {
	// Located in utils_test for mocking.
}

func TestExpandUser(t *testing.T) {
	got, err := util.ExpandUser("~/foo/bar")
	usr, _ := user.Current()
	want := filepath.Join(usr.HomeDir, "foo", "bar")

	assert.Equal(t, want, got)
	assert.NoError(t, err)
}

func TestExpandUserWithFullPath(t *testing.T) {
	got, err := util.ExpandUser("/foo/bar")
	want := "/foo/bar"

	assert.Equal(t, want, got)
	assert.NoError(t, err)
}

func TestRunCommandReturnsError(t *testing.T) {
	err := util.RunCmd(false, "false")

	assert.Error(t, err)
}

func TestRunCommandPrintsStreamingStdout(t *testing.T) {
	got := capturer.CaptureStdout(func() {
		err := util.RunCmd(true, "echo", "-n", "foo")
		assert.NoError(t, err)
	})
	want := "COMMAND: \x1b[30;41mecho -n foo\x1b[0m\nfoo"

	assert.Equal(t, want, got)
}

func TestRunCommandPrintsStreamingStderr(t *testing.T) {
	got := capturer.CaptureStderr(func() {
		err := util.RunCmd(true, "cat", "foo")
		assert.Error(t, err)
	})
	want := "cat: foo: No such file or directory\n"

	assert.Equal(t, want, got)
}

func TestCopyFile(t *testing.T) {
	srcFile := path.Join("..", "..", "test", "resources", "copy", "file")
	dstFile := path.Join("..", "..", "test", "resources", "copy", "copiedFile")
	err := util.CopyFile(srcFile, dstFile)
	defer func() { _ = os.Remove(dstFile) }()

	assert.NoError(t, err)
	assert.FileExistsf(t, dstFile, "File does not exist")
}

func TestCopyFileReturnsError(t *testing.T) {
	srcFile := path.Join("..", "..", "test", "resources", "copy", "file")
	invalidDstFile := "/super/invalid/path/to/write/to"
	err := util.CopyFile(srcFile, invalidDstFile)

	assert.Error(t, err)
}

func TestCopyDir(t *testing.T) {
	srcDir := path.Join("..", "..", "test", "resources", "copy", "dir")
	dstDir := path.Join("..", "..", "test", "resources", "copy", "copiedDir")
	err := util.CopyDir(srcDir, dstDir)
	defer func() { _ = os.RemoveAll(dstDir) }()

	assert.NoError(t, err)
	assert.DirExistsf(t, dstDir, "Dir does not exist")
}

func TestCopyDirReturnsError(t *testing.T) {
	srcDir := path.Join("..", "..", "test", "resources", "copy", "dir")
	invalidDstDir := "/super/invalid/path/to/write/to"
	err := util.CopyDir(srcDir, invalidDstDir)

	assert.Error(t, err)
}
