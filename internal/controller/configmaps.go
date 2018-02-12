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

package controller

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/golang/glog"
	"gopkg.in/yaml.v2"

	corev1 "k8s.io/api/core/v1"
)

func (c *LemonLDAPNGController) parseConfigMap(obj interface{}) (namespace string, name string, match bool, overrides map[string]interface{}, err error) {
	configMapObj := obj.(*corev1.ConfigMap)
	configMapKey := fmt.Sprintf("%s/%s", configMapObj.Namespace, configMapObj.Name)
	if configMapKey != c.controllerConfig.ConfigMapName {
		return configMapObj.Namespace, configMapObj.Name, false, nil, nil
	}
	overrides = make(map[string]interface{})
	for k, v := range configMapObj.Data {
		if strings.HasSuffix(k, ".yaml") {
			vUnmarshaled := make(map[string]interface{})
			err = yaml.Unmarshal([]byte(v), &vUnmarshaled)
			if err != nil {
				glog.Errorf("Unable to decode key %s in ConfigMap %s: %s", k, configMapKey, err)
			}
			overrides[strings.TrimSuffix(k, ".yaml")] = vUnmarshaled
		} else if strings.Contains(k, ".") {
			glog.Errorf("Unsupported suffix for key %s in ConfigMap %s: %s", k, configMapKey, "Use .yaml or none")
		} else {
			overrides[k] = v
		}
	}
	return configMapObj.Namespace, configMapObj.Name, true, overrides, nil
}

func (c *LemonLDAPNGController) configMapAdded(obj interface{}) {
	namespace, name, match, overrides, err := c.parseConfigMap(obj)
	if !match {
		return
	}
	if err != nil {
		glog.Error(err)
		return
	}
	glog.Infof("A ConfigMap was added: %s/%s", namespace, name)
	c.llngConfig.SetOverrides(overrides)
	err = c.llngConfig.Save() // FIXME async + batch
	if err != nil {
		glog.Error(err)
		return
	}
}

func (c *LemonLDAPNGController) configMapDeleted(obj interface{}) {
	namespace, name, match, _, err := c.parseConfigMap(obj)
	if !match {
		return
	}
	if err != nil {
		// glog.Error(err)
		return
	}
	glog.Infof("A ConfigMap was deleted: %s/%s", namespace, name)
	c.llngConfig.SetOverrides(make(map[string]interface{}))
	err = c.llngConfig.Save() // FIXME async + batch
	if err != nil {
		glog.Error(err)
		return
	}
}

func (c *LemonLDAPNGController) configMapUpdated(old, cur interface{}) {
	_, _, oldMatch, oldOverrides, _ := c.parseConfigMap(old)
	curNamespace, curName, curMatch, curOverrides, curErr := c.parseConfigMap(cur)
	if !curMatch && !oldMatch {
		return
	}
	if curErr != nil {
		glog.Error(curErr)
		return
	}
	if reflect.DeepEqual(oldOverrides, curOverrides) {
		return
	}
	glog.Infof("A ConfigMap was updated: %s/%s", curNamespace, curName)
	c.llngConfig.SetOverrides(curOverrides)
	err := c.llngConfig.Save() // FIXME async + batch
	if err != nil {
		glog.Error(err)
		return
	}
}
