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
	"fmt"
	"io"
	"log/slog"

	"github.com/avfs/avfs"
)

// NewCopy factory to create a new copy instance.
func NewCopy(
	appFs avfs.VFS,
	logger *slog.Logger,
) *Copy {
	return &Copy{
		appFs:  appFs,
		logger: logger,
	}
}

// CopyFile copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file. The file mode will be copied from the source and
// the copied data is synced/flushed to stable storage.
func (r *Copy) CopyFile(
	src string,
	dst string,
) (err error) {
	baseSrc := r.appFs.Base(src)

	r.logger.Info(
		"copying file",
		slog.String("srcFile", baseSrc),
		slog.String("dstFile", dst),
	)

	// Open the source file for reading, and record its metadata
	in, err := r.appFs.Open(src)
	if err != nil {
		return err
	}
	defer func() { _ = in.Close() }()

	si, err := in.Stat()
	if err != nil {
		return err
	}

	// Open dest file for writing; make it owner-only perms before putting
	// anything in it
	_ = r.appFs.MkdirAll(r.appFs.Dir(dst), si.Mode())
	out, err := r.appFs.Create(dst)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()

	err = out.Chmod(0o600)
	if err != nil {
		return err
	}

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	err = out.Sync()
	if err != nil {
		return err
	}

	// All done; make the permissions match
	err = out.Chmod(si.Mode())
	if err != nil {
		return err
	}

	return nil
}

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
// Symlinks are implicitly dereferenced.
func (r *Copy) CopyDir(
	src string,
	dst string,
) (err error) {
	src = r.appFs.Clean(src)
	dst = r.appFs.Clean(dst)
	baseSrc := r.appFs.Base(src)

	r.logger.Info(
		"copying dir",
		slog.String("srcDir", baseSrc),
		slog.String("dstDir", dst),
	)

	si, err := r.appFs.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	d, err := r.appFs.Open(dst)
	_ = d.Close()
	if err == nil {
		return fmt.Errorf("destination already exists")
	}

	err = r.appFs.MkdirAll(dst, si.Mode())
	if err != nil {
		return err
	}

	entries, err := r.appFs.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := r.appFs.Join(src, entry.Name())
		dstPath := r.appFs.Join(dst, entry.Name())

		// Dereference any symlinks and copy their contents instead
		target, err := r.appFs.Stat(srcPath)
		if err != nil {
			return err
		}

		if target.IsDir() {
			err = r.CopyDir(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			err = r.CopyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
