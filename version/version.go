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

package version

import "fmt"

var (
	// RELEASE returns the release version
	RELEASE = "UNKNOWN"
	// REPO returns the git repository URL
	REPO = "UNKNOWN"
	// COMMIT returns the short sha from git
	COMMIT = "UNKNOWN"
)

// String returns information about the release.
func String() string {
	return fmt.Sprintf(`-------------------------------------------------------------------------------
LemonLDAP::NG controller
  Release:    %v
  Build:      %v
  Repository: %v
-------------------------------------------------------------------------------
`, RELEASE, COMMIT, REPO)
}

// Short returns one-line information about the release.
func Short() string {
	return fmt.Sprintf(`LemonLDAP::NG controller %v (build: %v, repository: %v)`, RELEASE, COMMIT, REPO)
}
