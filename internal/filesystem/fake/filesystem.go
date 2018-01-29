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
	"sync"
	"time"

	"github.com/lemonldap-ng-controller/lemonldap-ng-controller/internal/filesystem"
)

// FileSystem implements FileSystem interface
type FileSystem struct {
	sync.RWMutex // protects all File's entries

	root *File
}

// NewFileSystem creates a new FileSystem
func NewFileSystem() *FileSystem {
	fs := &FileSystem{}
	fs.root = NewFile(fs, nil, "/", 0755, time.Now(), true)
	fs.Mkdir("/var", 0755)
	fs.Mkdir("/var/lib", 0755)
	fs.Mkdir("/var/lib/lemonldap-ng", 0755)
	fs.Mkdir("/var/lib/lemonldap-ng/conf", 0755)
	content := []byte(`{
		"cfgNum": 1,
		"exportedHeaders": {},
		"locationRules": {}
	}`)
	fs.WriteFile("/var/lib/lemonldap-ng/conf/lmConf-1.js", content, 0644)
	return fs
}

// Mkdir creates a new directory with the specified name and permission bits
func (fs *FileSystem) Mkdir(name string, perm os.FileMode) error {
	fs.Lock()
	defer fs.Unlock()
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
	NewFile(fs, fParent, path.Base(name), perm, time.Now(), true)
	return nil
}

// Open opens the named file for reading
func (fs *FileSystem) Open(name string) (filesystem.File, error) {
	fs.RLock()
	defer fs.RUnlock()
	return fs.root.lookupFile(name, name)
}

// ReadFile reads a file and returns the contents
func (fs *FileSystem) ReadFile(filename string) ([]byte, error) {
	f, err := fs.Open(filename)
	if err != nil {
		return []byte(""), err
	}
	ff := f.(*File)
	ff.RLock()
	defer ff.RUnlock()
	c := make([]byte, len(ff.content))
	copy(c, ff.content)
	return c, nil
}

// WriteFile reads a file and returns the contents
func (fs *FileSystem) WriteFile(filename string, data []byte, perm os.FileMode) error {
	f, err := fs.Open(filename)
	if err != nil {
		fs.Lock()
		defer fs.Unlock()
		fParent, errParent := fs.root.lookupFile(filename, path.Dir(filename))
		if errParent != nil {
			return errParent
		}
		f = NewFile(fs, fParent, path.Base(filename), perm, time.Now(), false)
	}
	ff := f.(*File)
	ff.Lock()
	defer ff.Unlock()
	ff.content = make([]byte, len(data))
	copy(ff.content, data)
	return nil
}

// File implements File interface
type File struct {
	sync.RWMutex // protects this File's content and metadata

	fs      *FileSystem
	parent  *File
	name    string
	mode    os.FileMode
	modTime time.Time
	isDir   bool
	content []byte
	entries map[string]*File
}

// NewFile creates a new File
func NewFile(fs *FileSystem, parent *File, name string, mode os.FileMode, modTime time.Time, isDir bool) *File {
	ff := &File{
		fs:      fs,
		parent:  parent,
		name:    name,
		mode:    mode,
		modTime: modTime,
		isDir:   isDir,
		content: []byte(""),
		entries: make(map[string]*File),
	}
	if parent != nil {
		parent.entries[name] = ff
	}
	return ff
}

// Name returns the base name of the file
func (ff *File) Name() string {
	ff.RLock()
	defer ff.RUnlock()
	return ff.name
}

// Size returns the length in bytes
func (ff *File) Size() int64 {
	ff.RLock()
	defer ff.RUnlock()
	return int64(len(ff.content))
}

// Mode returns the file mode bits
func (ff *File) Mode() os.FileMode {
	ff.RLock()
	defer ff.RUnlock()
	return ff.mode
}

// ModTime returns the  modification time
func (ff *File) ModTime() time.Time {
	ff.RLock()
	defer ff.RUnlock()
	return ff.modTime
}

// IsDir returns true if directory
func (ff File) IsDir() bool {
	ff.RLock()
	defer ff.RUnlock()
	return ff.isDir
}

// Sys returns the underlying File
func (ff *File) Sys() interface{} {
	ff.RLock()
	defer ff.RUnlock()
	return ff.fs
}

func (ff *File) lookupFile(fullpath, relativepath string) (*File, error) {
	parts := strings.SplitN(relativepath, "/", 2)
	if ff.parent == nil && parts[0] == "" { // root
		parts = strings.SplitN(parts[1], "/", 2)
	}
	if len(parts) == 1 && parts[0] == "" {
		return ff, nil
	}
	if nextEntry, ok := ff.entries[parts[0]]; ok {
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
func (ff *File) Close() error {
	return nil
}

// Readdir reads the contents of the directory associated with file and returns a slice of up to n FileInfo values, as would be returned by Lstat, in directory order
func (ff *File) Readdir(n int) ([]os.FileInfo, error) {
	if n > 0 {
		return nil, &os.PathError{
			Op:   "readdir",
			Path: ff.Name(),
			Err:  errors.New("Sliced call to Readdir not supported"),
		}
	}
	ff.fs.RLock()
	defer ff.fs.RUnlock()
	ret := make([]os.FileInfo, len(ff.entries))
	i := 0
	for _, entry := range ff.entries {
		ret[i] = entry
		i++
	}
	return ret, nil
}
