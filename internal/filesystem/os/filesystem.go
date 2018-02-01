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

package os

// Inspired by https://talks.golang.org/2012/10things.slide#8

import (
	"io/ioutil"
	"os"

	"github.com/lemonldap-ng-controller/lemonldap-ng-controller/internal/filesystem"
)

// FileSystem implements FileSystem interface
type FileSystem struct{}

// Mkdir creates a new directory with the specified name and permission bits
func (FileSystem) Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}

// Open opens the named file for reading
func (FileSystem) Open(name string) (filesystem.File, error) {
	return os.Open(name)
}

// Stat returns a FileInfo describing the named file
func (FileSystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

// ReadFile reads a file and returns the contents
func (FileSystem) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

// WriteFile reads a file and returns the contents
func (FileSystem) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return ioutil.WriteFile(filename, data, perm)
}
