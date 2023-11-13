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

package repository

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/logrusorgru/aurora/v4"

	"github.com/retr0h/go-gilt/internal/util"
)

var (
	// CopyFile is mocked for tests.
	CopyFile = util.CopyFile
	// CopyDir is mocked for tests.
	CopyDir = util.CopyDir
)

// GetCloneDir returns the path to the Repository's clone directory.
func (r *Repository) GetCloneDir() string {
	return filepath.Join(r.GiltDir, r.getCloneHash())
}

func (r *Repository) getCloneHash() string {
	replacer := strings.NewReplacer(
		"/", "-",
		":", "-",
	)
	replacedGitURL := replacer.Replace(r.Git)

	return fmt.Sprintf("%s-%s", replacedGitURL, r.Version)
}

// CopySources copy Repository.Src to Repository.DstFile or Repository.DstDir.
func (r *Repository) CopySources() error {
	cloneDir := r.GetCloneDir()

	msg := fmt.Sprintf("%-2s [%s]:", "", aurora.Magenta(cloneDir))
	fmt.Println(msg)

	for _, rSource := range r.Sources {
		srcFullPath := filepath.Join(cloneDir, rSource.Src)
		globbedSrc, err := filepath.Glob(srcFullPath)
		if err != nil {
			return err
		}

		for _, src := range globbedSrc {
			// The source is a file.
			if info, err := os.Stat(src); err == nil && info.Mode().IsRegular() {
				// ... and the destination is declared a directory.
				if rSource.DstFile != "" {
					if err := CopyFile(src, rSource.DstFile); err != nil {
						return err
					}
				} else if rSource.DstDir != "" {
					// ... and the destination directory exists.
					if info, err := os.Stat(rSource.DstDir); err == nil && info.Mode().IsDir() {
						srcBaseFile := filepath.Base(src)
						newDst := filepath.Join(rSource.DstDir, srcBaseFile)
						if err := CopyFile(src, newDst); err != nil {
							return err
						}
					} else {
						msg := fmt.Sprintf("DstDir '%s' does not exist", rSource.DstDir)
						return errors.New(msg)
					}
				}
				// The source is a directory.
			} else if info, err := os.Stat(src); err == nil && info.Mode().IsDir() {
				if info, err := os.Stat(rSource.DstDir); err == nil && info.Mode().IsDir() {
					if err := os.RemoveAll(rSource.DstDir); err != nil {
						return err
					}
				}
				if err := CopyDir(src, rSource.DstDir); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
