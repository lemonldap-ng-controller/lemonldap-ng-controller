/*
Copyright 2018 Mathieu Parent <math.parent@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package fake

// Inspired by https://talks.golang.org/2012/10things.slide#8

import (
	"errors"
	"os"
	"path"
	"strings"
	"time"

	"github.com/lemonldap-ng-controller/lemonldap-ng-controller/pkg/filesystem"
)

// FakeFileSystem implements FileSystem interface
type FakeFileSystem struct {
	root *FakeFile
}

// NewFakeFileSystem creates a new FakeFileSystem
func NewFakeFileSystem() *FakeFileSystem {
	fs := FakeFileSystem{}
	fs.root = NewFakeFile(fs, nil, "/", 0755, time.Now(), true)
	fs.Mkdir("/var", 0755)
	fs.Mkdir("/var/lib", 0755)
	fs.Mkdir("/var/lib/lemonldap-ng", 0755)
	fs.Mkdir("/var/lib/lemonldap-ng/conf", 0755)
	content := []byte(`{
		"exportedHeaders": {},
		"locationRules": {}
	}`)
	fs.WriteFile("/var/lib/lemonldap-ng/conf/lmConf-1.js", content, 0644)
	return &fs
}

// Mkdir creates a new directory with the specified name and permission bits
func (fs FakeFileSystem) Mkdir(name string, perm os.FileMode) error {
	_, err := fs.root.lookupFile(name, name)
	if err == nil {
		return &os.PathError{
			Op:   "mkdir",
			Path: name,
			Err:  errors.New("File exists"), // 0x11
		}
	}
	dirName := path.Dir(name)
	fParent, errParent := fs.root.lookupFile(dirName, dirName)
	if errParent != nil {
		return &os.PathError{
			Op:   "mkdir",
			Path: name,
			Err:  errParent.(*os.PathError).Err,
		}
	}
	NewFakeFile(fs, fParent, path.Base(name), perm, time.Now(), true)
	return nil
}

// Open opens the named file for reading
func (fs FakeFileSystem) Open(name string) (filesystem.File, error) {
	return fs.root.lookupFile(name, name)
}

// ReadFile reads a file and returns the contents
func (fs FakeFileSystem) ReadFile(filename string) ([]byte, error) {
	f, err := fs.Open(filename)
	if err != nil {
		return []byte(""), err
	}
	return f.(*FakeFile).content, nil
}

// WriteFile reads a file and returns the contents
func (fs FakeFileSystem) WriteFile(filename string, data []byte, perm os.FileMode) error {
	f, err := fs.Open(filename)
	if err != nil {
		fParent, errParent := fs.root.lookupFile(filename, path.Dir(filename))
		if errParent != nil {
			return errParent
		}
		f = NewFakeFile(fs, fParent, path.Base(filename), perm, time.Now(), false)
	}
	f.(*FakeFile).content = data
	return nil
}

// FakeFile implements File interface
type FakeFile struct {
	fs      FakeFileSystem
	parent  *FakeFile
	name    string
	mode    os.FileMode
	modTime time.Time
	isDir   bool
	content []byte
	entries map[string]*FakeFile
}

// NewFakeFile creates a new FakeFile
func NewFakeFile(fs FakeFileSystem, parent *FakeFile, name string, mode os.FileMode, modTime time.Time, isDir bool) *FakeFile {
	f := &FakeFile{
		fs:      fs,
		parent:  parent,
		name:    name,
		mode:    mode,
		modTime: modTime,
		isDir:   isDir,
		content: []byte(""),
		entries: make(map[string]*FakeFile),
	}
	if parent != nil {
		parent.entries[name] = f
	}
	return f
}

// Name returns the base name of the file
func (f FakeFile) Name() string {
	return f.name
}

// Size returns the length in bytes
func (f FakeFile) Size() int64 {
	return int64(len(f.content))
}

// Mode returns the file mode bits
func (f FakeFile) Mode() os.FileMode {
	return f.mode
}

// ModTime returns the  modification time
func (f FakeFile) ModTime() time.Time {
	return f.modTime
}

// IsDir returns true if directory
func (f FakeFile) IsDir() bool {
	return f.isDir
}

// Sys returns the underlying FakeFile
func (f FakeFile) Sys() interface{} {
	return f.fs
}

func (f FakeFile) lookupFile(fullpath, relativepath string) (*FakeFile, error) {
	parts := strings.SplitN(relativepath, "/", 2)
	if f.parent == nil && parts[0] == "" { // root
		parts = strings.SplitN(parts[1], "/", 2)
	}
	if len(parts) == 1 && parts[0] == "" {
		return &f, nil
	}
	if nextEntry, ok := f.entries[parts[0]]; ok {
		if len(parts) == 1 {
			return nextEntry, nil
		}
		return nextEntry.lookupFile(fullpath, parts[1])
	}
	return nil, &os.PathError{
		Op:   "open",
		Path: fullpath,
		Err:  errors.New("No such file or directory"), // 0x2
	}
}

// Close closes the File
func (f FakeFile) Close() error {
	return nil
}

// Readdir reads the contents of the directory associated with file and returns a slice of up to n FileInfo values, as would be returned by Lstat, in directory order
func (f FakeFile) Readdir(n int) ([]os.FileInfo, error) {
	if n > 0 {
		return nil, &os.PathError{
			Op:   "readdir",
			Path: f.Name(),
			Err:  errors.New("Sliced call to Readdir not supported"),
		}
	}
	ret := make([]os.FileInfo, len(f.entries))
	i := 0
	for _, entry := range f.entries {
		ret[i] = entry
		i += 1
	}
	return ret, nil
}
