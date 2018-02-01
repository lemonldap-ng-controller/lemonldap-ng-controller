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

package filesystem

// Inspired by https://talks.golang.org/2012/10things.slide#8

import (
	"os"
)

// Filesystem interface
type Filesystem interface {
	// from "os"
	Mkdir(name string, perm os.FileMode) error
	Open(name string) (File, error)
	Stat(name string) (os.FileInfo, error)

	// from "io/ioutil"
	ReadFile(filename string) ([]byte, error)
	WriteFile(filename string, data []byte, perm os.FileMode) error
}

// File interface
type File interface {
	Close() error
	Readdir(int) ([]os.FileInfo, error)
}
