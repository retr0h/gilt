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

// Package repository TODO change to lower cases members.
package repository

import (
	"github.com/spf13/afero"
)

// Sources mapping of files and/or directories needing copied.
type Sources struct {
	// Src source file or directory to copy.
	Src string `mapstructure:"src"`
	// DstFile destination of file copy.
	DstFile string `mapstructure:"dstFile"`
	// DstDir destination of directory copy.
	DstDir string `mapstructure:"dstDir"`
}

// Repository contains the repository's details for cloning.
type Repository struct {
	// Git url of Git repository to clone.
	Git string `mapstructure:"git"`
	// Version of Git repository to use.
	Version string `mapstructure:"version"`
	// DstDir destination directory to copy clone to.
	DstDir string `mapstructure:"dstDir"`
	// Sources containing files and/or directories to copy.
	Sources []Sources `mapstructure:"sources"`
	// GiltDir path to Gilt's clone dir option set from CLI.
	GiltDir string
	// AppFs file system mocking.
	// TODO(retr0h): Needs passed in as repo pattern.
	AppFs afero.Fs
}
