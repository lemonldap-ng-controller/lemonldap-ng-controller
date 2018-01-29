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

	"github.com/golang/glog"
	"gopkg.in/yaml.v2"

	corev1 "k8s.io/api/core/v1"
)

func (c *LemonLDAPNGController) configMapSmurfed(obj interface{}, verb string) {
	configMapObj := obj.(*corev1.ConfigMap)
	configMapKey := fmt.Sprintf("%s/%s", configMapObj.Namespace, configMapObj.Name)
	if configMapKey == c.controllerConfig.ConfigMapName {
		glog.Infof("A ConfigMap was %s: %s", verb, configMapKey)
		if verb == "deleted" {
			c.llngConfig.SetOverrides(make(map[string]interface{}))
			err := c.llngConfig.Save() // FIXME async + batch
			if err != nil {
				glog.Error(err)
				return
			}
			return
		}
		lmConfYaml, ok := configMapObj.Data["lmConf.js"]
		lmConf := make(map[string]interface{})
		if !ok {
			glog.Error("Missing key in ConfigMap: lmConf.js")
			return
		}
		err := yaml.Unmarshal([]byte(lmConfYaml), &lmConf)
		if err != nil {
			glog.Errorf("Unable to parse lmConf.js of ConfigMap %s, ignoring Ingress: %s", configMapKey, err)
			return
		}
		c.llngConfig.SetOverrides(lmConf)
		err = c.llngConfig.Save() // FIXME async + batch
		if err != nil {
			glog.Error(err)
			return
		}
	}
}

func (c *LemonLDAPNGController) configMapAdded(obj interface{}) {
	c.configMapSmurfed(obj, "added")
}

func (c *LemonLDAPNGController) configMapDeleted(obj interface{}) {
	c.configMapSmurfed(obj, "deleted")
}

func (c *LemonLDAPNGController) configMapUpdated(old, cur interface{}) {
	c.configMapSmurfed(cur, "updated")
}
