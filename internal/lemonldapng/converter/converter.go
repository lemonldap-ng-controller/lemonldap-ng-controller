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
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/ghodss/yaml"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Run do the conversion
func Run(configMapName string, r io.Reader, w io.Writer) error {
	configMapNameAndNamespace := strings.SplitN(configMapName, "/", 2)
	if len(configMapNameAndNamespace) == 1 || configMapNameAndNamespace[1] == "" {
		configMapNameAndNamespace = []string{"ingress-nginx", "lemonldap-ng-configuration"}
	}
	inputBuffer, err := ioutil.ReadAll(r)
	lmConf := make(map[string]interface{})
	err = json.Unmarshal(inputBuffer, &lmConf)
	if err != nil {
		return fmt.Errorf("Unable to parse input: %s", err)
	}
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapNameAndNamespace[1],
			Namespace: configMapNameAndNamespace[0],
		},
		Data:       map[string]string{},
		BinaryData: map[string][]byte{},
	}
	for k, v := range lmConf {
		switch v := v.(type) {
		case string:
			cm.Data[k] = v
		default:
			out, err := yaml.Marshal(v)
			if err != nil {
				return fmt.Errorf("Unable to encode key %s: %s", k, err)
			}
			cm.Data[k+".yaml"] = string(out[:])
		}
	}
	out, err := yaml.Marshal(cm)
	if err != nil {
		return fmt.Errorf("Unable to encode ConfigMap: %s", err)
	}
	fmt.Fprintf(w, "%s", out[:])
	return nil
}
