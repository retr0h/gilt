// Copyright (c) 2023 Nicolas Simonds

//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//  	http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.
//

package failfs

import (
	"fmt"
	"io/fs"
	"reflect"
	"time"

	"github.com/avfs/avfs"
	"github.com/avfs/avfs/idm/dummyidm"
)

type FailFS struct {
	baseFS            avfs.VFS
	errOpNotPermitted error
	errPermDenied     error
	features          avfs.Features
	failFn            map[string]interface{}
	avfs.Utils[*FailFS]
}

type FailFile struct {
	baseFile avfs.File
	vfs      *FailFS
}

func New(baseFS avfs.VFS, failFn map[string]interface{}) *FailFS {
	return &FailFS{
		baseFS:            baseFS,
		errOpNotPermitted: avfs.ErrOpNotPermitted,
		errPermDenied:     avfs.ErrPermDenied,
		features:          baseFS.Features() &^ avfs.FeatIdentityMgr,
		failFn:            failFn,
	}
}

// callFailFn calls the provided failure function with the given arguments.
// It panics if the function signature does not match the provided arguments.
// The caller is 100% responsible for making sure the function signatures match
func (vfs *FailFS) callFailFn(fn interface{}, args ...interface{}) []reflect.Value {
	failFnValue := reflect.ValueOf(fn)
	if failFnValue.Kind() != reflect.Func {
		panic(fmt.Sprintf("failFn is not a function: %v", fn))
	}

	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		in[i] = reflect.ValueOf(arg)
	}

	results := failFnValue.Call(in)
	return results
}

// FailFS ops

func (vfs *FailFS) Features() avfs.Features {
	return vfs.features
}

func (vfs *FailFS) HasFeature(feature avfs.Features) bool {
	return vfs.features&feature == feature
}

func (vfs *FailFS) Name() string {
	return vfs.baseFS.Name()
}

func (vfs *FailFS) OSType() avfs.OSType {
	return vfs.baseFS.OSType()
}

func (*FailFS) Type() string {
	return "FailFS"
}

func (vfs *FailFS) CreateSystemDirs(basePath string) error {
	return vfs.baseFS.CreateSystemDirs(basePath)
}

func (vfs *FailFS) CreateHomeDir(u avfs.UserReader) (string, error) {
	return vfs.baseFS.CreateHomeDir(u)
}

func (vfs *FailFS) HomeDirUser(u avfs.UserReader) string {
	return vfs.baseFS.HomeDirUser(u)
}

func (vfs *FailFS) SystemDirs(basePath string) []avfs.DirInfo {
	return vfs.baseFS.SystemDirs(basePath)
}

func (vfs *FailFS) Abs(path string) (string, error) {
	if failFn, ok := vfs.failFn["Abs"]; ok {
		results := vfs.callFailFn(failFn, path)
		var abs string
		var err error
		abs, _ = results[0].Interface().(string)
		if !results[1].IsNil() {
			err, _ = results[1].Interface().(error)
		}
		return abs, err
	}
	return vfs.baseFS.Abs(path)
}

func (vfs *FailFS) Base(path string) string {
	return vfs.baseFS.Base(path)
}

func (vfs *FailFS) Chdir(dir string) error {
	return vfs.baseFS.Chdir(dir)
}

func (vfs *FailFS) Chmod(name string, mode fs.FileMode) error {
	if failFn, ok := vfs.failFn["Chmod"]; ok {
		results := vfs.callFailFn(failFn, name, mode)
		var err error
		if !results[0].IsNil() {
			err, _ = results[0].Interface().(error)
		}
		return err
	}
	return vfs.baseFS.Chmod(name, mode)
}

func (vfs *FailFS) Chown(name string, uid, gid int) error {
	return vfs.baseFS.Chown(name, uid, gid)
}

func (vfs *FailFS) Chtimes(name string, atime, mtime time.Time) error {
	return vfs.baseFS.Chtimes(name, atime, mtime)
}

func (vfs *FailFS) Clean(path string) string {
	return vfs.baseFS.Clean(path)
}

func (vfs *FailFS) Create(name string) (avfs.File, error) {
	return vfs.Utils.Create(vfs, name)
}

func (vfs *FailFS) CreateTemp(dir, pattern string) (avfs.File, error) {
	return vfs.baseFS.CreateTemp(dir, pattern)
}

func (vfs *FailFS) Dir(path string) string {
	return vfs.baseFS.Dir(path)
}

func (vfs *FailFS) EvalSymlinks(path string) (string, error) {
	return vfs.baseFS.EvalSymlinks(path)
}

func (vfs *FailFS) FromSlash(path string) string {
	return vfs.baseFS.FromSlash(path)
}

func (vfs *FailFS) Getwd() (dir string, err error) {
	return vfs.baseFS.Getwd()
}

func (vfs *FailFS) Glob(pattern string) (matches []string, err error) {
	return vfs.baseFS.Glob(pattern)
}

func (vfs *FailFS) Idm() avfs.IdentityMgr {
	return dummyidm.NotImplementedIdm
}

func (vfs *FailFS) IsAbs(path string) bool {
	return vfs.baseFS.IsAbs(path)
}

func (vfs *FailFS) IsPathSeparator(c uint8) bool {
	return vfs.baseFS.IsPathSeparator(c)
}

func (vfs *FailFS) Join(elem ...string) string {
	return vfs.baseFS.Join(elem...)
}

func (vfs *FailFS) Lchown(name string, uid, gid int) error {
	return vfs.baseFS.Lchown(name, uid, gid)
}

func (vfs *FailFS) Link(oldname, newname string) error {
	return vfs.baseFS.Link(oldname, newname)
}

func (vfs *FailFS) Lstat(name string) (fs.FileInfo, error) {
	return vfs.baseFS.Lstat(name)
}

func (vfs *FailFS) Match(pattern, name string) (matched bool, err error) {
	return vfs.baseFS.Match(pattern, name)
}

func (vfs *FailFS) Mkdir(name string, perm fs.FileMode) error {
	return vfs.baseFS.Mkdir(name, perm)
}

func (vfs *FailFS) MkdirAll(path string, perm fs.FileMode) error {
	if failFn, ok := vfs.failFn["MkdirAll"]; ok {
		results := vfs.callFailFn(failFn, path, perm)
		var err error
		if !results[0].IsNil() {
			err, _ = results[0].Interface().(error)
		}
		return err
	}
	return vfs.baseFS.MkdirAll(path, perm)
}

func (vfs *FailFS) MkdirTemp(dir, prefix string) (name string, err error) {
	return vfs.baseFS.MkdirTemp(dir, prefix)
}

func (vfs *FailFS) Open(name string) (avfs.File, error) {
	if failFn, ok := vfs.failFn["Open"]; ok {
		results := vfs.callFailFn(failFn, name)
		var file avfs.File
		var err error
		if !results[0].IsNil() {
			file, _ = results[0].Interface().(avfs.File)
		}
		if !results[1].IsNil() {
			err, _ = results[1].Interface().(error)
		}
		return file, err
	}
	return vfs.Utils.Open(vfs, name)
}

func (vfs *FailFS) PathSeparator() uint8 {
	return vfs.baseFS.PathSeparator()
}

func (vfs *FailFS) ReadDir(name string) ([]fs.DirEntry, error) {
	if failFn, ok := vfs.failFn["ReadDir"]; ok {
		results := vfs.callFailFn(failFn, name)
		var entries []fs.DirEntry
		var err error
		if !results[0].IsNil() {
			entries, _ = results[0].Interface().([]fs.DirEntry)
		}
		if !results[1].IsNil() {
			err, _ = results[1].Interface().(error)
		}
		return entries, err
	}
	return vfs.baseFS.ReadDir(name)
}

func (vfs *FailFS) ReadFile(filename string) ([]byte, error) {
	return vfs.baseFS.ReadFile(filename)
}

func (vfs *FailFS) Readlink(name string) (string, error) {
	return vfs.baseFS.Readlink(name)
}

func (vfs *FailFS) Rel(basepath, targpath string) (string, error) {
	return vfs.baseFS.Rel(basepath, targpath)
}

func (vfs *FailFS) Remove(name string) error {
	return vfs.baseFS.Remove(name)
}

func (vfs *FailFS) RemoveAll(path string) error {
	return vfs.baseFS.RemoveAll(path)
}

func (vfs *FailFS) Rename(oldname, newname string) error {
	return vfs.baseFS.Rename(oldname, newname)
}

func (vfs *FailFS) SameFile(fi1, fi2 fs.FileInfo) bool {
	return vfs.baseFS.SameFile(fi1, fi2)
}

func (vfs *FailFS) SetUMask(mask fs.FileMode) {
	vfs.baseFS.SetUMask(mask)
}

func (vfs *FailFS) SetUser(name string) (avfs.UserReader, error) {
	return vfs.baseFS.SetUser(name)
}

func (vfs *FailFS) Split(path string) (dir, file string) {
	return vfs.baseFS.Split(path)
}

func (vfs *FailFS) SplitAbs(path string) (dir, file string) {
	return vfs.baseFS.SplitAbs(path)
}

func (vfs *FailFS) Stat(name string) (fs.FileInfo, error) {
	if failFn, ok := vfs.failFn["Stat"]; ok {
		results := vfs.callFailFn(failFn, name)
		var fileInfo fs.FileInfo
		var err error
		if !results[0].IsNil() {
			fileInfo, _ = results[0].Interface().(fs.FileInfo)
		}
		if !results[1].IsNil() {
			err, _ = results[1].Interface().(error)
		}
		return fileInfo, err
	}
	return vfs.baseFS.Stat(name)
}

func (vfs *FailFS) Sub(dir string) (avfs.VFS, error) {
	return vfs.baseFS.Sub(dir)
}

func (vfs *FailFS) Symlink(oldname, newname string) error {
	return vfs.baseFS.Symlink(oldname, newname)
}

func (vfs *FailFS) TempDir() string {
	return vfs.baseFS.TempDir()
}

func (vfs *FailFS) ToSlash(path string) string {
	return vfs.baseFS.ToSlash(path)
}

func (vfs *FailFS) ToSysStat(info fs.FileInfo) avfs.SysStater {
	return vfs.baseFS.ToSysStat(info)
}

func (vfs *FailFS) Truncate(name string, size int64) error {
	return vfs.baseFS.Truncate(name, size)
}

func (vfs *FailFS) UMask() fs.FileMode {
	return vfs.baseFS.UMask()
}

func (vfs *FailFS) User() avfs.UserReader {
	return vfs.baseFS.User()
}

func (vfs *FailFS) WalkDir(root string, fn fs.WalkDirFunc) error {
	return vfs.baseFS.WalkDir(root, fn)
}

func (vfs *FailFS) WriteFile(filename string, data []byte, perm fs.FileMode) error {
	return vfs.baseFS.WriteFile(filename, data, perm)
}

// FailFile ops

func (vfs *FailFS) OpenFile(name string, flag int, perm fs.FileMode) (avfs.File, error) {
	fBase, err := vfs.baseFS.OpenFile(name, flag, perm)
	if fBase == nil || reflect.ValueOf(fBase).IsNil() {
		return (*FailFile)(nil), err
	}
	f := &FailFile{baseFile: fBase, vfs: vfs}
	return f, err
}

func (f *FailFile) Chdir() error {
	return f.baseFile.Chdir()
}

func (f *FailFile) Chmod(mode fs.FileMode) error {
	if failFn, ok := f.vfs.failFn["file.Chmod"]; ok {
		results := f.vfs.callFailFn(failFn, mode)
		var err error
		if !results[0].IsNil() {
			err, _ = results[0].Interface().(error)
		}
		return err
	}
	return f.baseFile.Chmod(mode)
}

func (f *FailFile) Chown(uid, gid int) error {
	return f.baseFile.Chown(uid, gid)
}

func (f *FailFile) Close() error {
	return f.baseFile.Close()
}

func (f *FailFile) Fd() uintptr {
	return f.baseFile.Fd()
}

func (f *FailFile) Name() string {
	return f.baseFile.Name()
}

func (f *FailFile) Read(b []byte) (int, error) {
	if failFn, ok := f.vfs.failFn["file.Read"]; ok {
		results := f.vfs.callFailFn(failFn, b)
		var n int
		var err error
		n, _ = results[0].Interface().(int)
		if !results[1].IsNil() {
			err, _ = results[1].Interface().(error)
		}
		return n, err
	}
	return f.baseFile.Read(b)
}

func (f *FailFile) ReadAt(b []byte, off int64) (n int, err error) {
	return f.baseFile.ReadAt(b, off)
}

func (f *FailFile) ReadDir(n int) ([]fs.DirEntry, error) {
	return f.baseFile.ReadDir(n)
}

func (f *FailFile) Readdirnames(n int) (names []string, err error) {
	return f.baseFile.Readdirnames(n)
}

func (f *FailFile) Seek(offset int64, whence int) (ret int64, err error) {
	return f.baseFile.Seek(offset, whence)
}

func (f *FailFile) Stat() (fs.FileInfo, error) {
	if failFn, ok := f.vfs.failFn["file.Stat"]; ok {
		results := f.vfs.callFailFn(failFn)
		var fileInfo fs.FileInfo
		var err error
		if !results[0].IsNil() {
			fileInfo, _ = results[0].Interface().(fs.FileInfo)
		}
		if !results[1].IsNil() {
			err, _ = results[1].Interface().(error)
		}
		return fileInfo, err
	}
	return f.baseFile.Stat()
}

func (f *FailFile) Sync() error {
	if failFn, ok := f.vfs.failFn["file.Sync"]; ok {
		results := f.vfs.callFailFn(failFn)
		var err error
		if !results[0].IsNil() {
			err, _ = results[0].Interface().(error)
		}
		return err
	}
	return f.baseFile.Sync()
}

func (f *FailFile) Truncate(size int64) error {
	return f.baseFile.Truncate(size)
}

func (f *FailFile) Write(b []byte) (n int, err error) {
	return f.baseFile.Write(b)
}

func (f *FailFile) WriteAt(b []byte, off int64) (n int, err error) {
	return f.baseFile.WriteAt(b, off)
}

func (f *FailFile) WriteString(s string) (n int, err error) {
	return f.Write([]byte(s))
}
