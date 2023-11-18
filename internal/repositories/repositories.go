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

package repositories

import (
	"os"

	"github.com/retr0h/go-gilt/internal/git"
	"github.com/spf13/afero"
)

// Overlay clone and extract the Repository items.
func (r *Repositories) Overlay() error {
	g := git.NewGit(r.Debug)

	for _, repository := range r.Repositories {
		repository.GiltDir = r.GiltDir
		repository.AppFs = afero.NewOsFs()
		err := g.Clone(repository)
		if err != nil {
			return err
		}

		// Checkout into repository.DstDir.
		if repository.DstDir != "" {
			// Delete dstDir since Checkout-Index does not clean old files that may
			// no longer exist in repository.
			if info, err := os.Stat(repository.DstDir); err == nil && info.Mode().IsDir() {
				if err := os.RemoveAll(repository.DstDir); err != nil {
					return err
				}
			}
			if err := g.CheckoutIndex(repository); err != nil {
				return err
			}
		}

		// Copy sources from Repository.Src to Repository.DstDir or Repository.DstFile.
		if len(repository.Sources) > 0 {
			if err := repository.CopySources(); err != nil {
				return err
			}
		}
	}

	return nil
}
