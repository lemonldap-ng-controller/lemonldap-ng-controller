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

package converter

import (
	"bytes"
	"flag"
	"strings"
	"testing"
)

func OneTest(t *testing.T, configMapName, input, expectedError, expectedOutput string) {
	flag.Set("alsologtostderr", "true")
	inputReader := strings.NewReader(input)
	var output bytes.Buffer
	err := Run(configMapName, inputReader, &output)
	if expectedError == "" {
		if err != nil {
			t.Errorf("Convert error: %s", err)
		}
	} else {
		if err == nil {
			t.Errorf("No error, expected: %s", expectedError)
		} else if err.Error() != expectedError {
			t.Errorf("Expected error: `%s` but got `%s`", expectedError, err)
		}
	}
	if output.String() != expectedOutput {
		t.Errorf("Expected: `%s` but got `%s`", expectedOutput, output.String())
	}

}

func TestEmptyHash(t *testing.T) {
	OneTest(t, "", `{}`, "", `metadata:
  creationTimestamp: null
  name: lemonldap-ng-configuration
  namespace: ingress-nginx
`)
}

func TestInvalidJson(t *testing.T) {
	OneTest(t, "", `{`, "Unable to parse input: unexpected end of JSON input", ``)
}

func TestDeepJson(t *testing.T) {
	OneTest(t, "", `{"a":"b","c":42,"d":{"e":"f"}}`, "", `data:
  a: b
  c.yaml: |
    42
  d.yaml: |
    e: f
metadata:
  creationTimestamp: null
  name: lemonldap-ng-configuration
  namespace: ingress-nginx
`)
}
