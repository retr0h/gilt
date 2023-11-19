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

package repository_test

import (
	"github.com/spf13/afero"
)

type FileSpec struct {
	appFs    afero.Fs
	srcDir   string
	srcFile  string
	srcFiles []string
}

func createFileSpecs(specs []FileSpec) {
	for _, s := range specs {
		if s.srcDir != "" {
			_ = s.appFs.MkdirAll(s.srcDir, 0o755)
		}

		if s.srcFile != "" {
			_, _ = s.appFs.Create(s.srcFile)
		}

		if len(s.srcFiles) > 0 {
			for _, f := range s.srcFiles {
				_, _ = s.appFs.Create(f)
			}
		}
	}
}
