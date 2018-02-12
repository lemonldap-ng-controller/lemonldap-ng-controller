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

	"github.com/golang/glog"
	"gopkg.in/yaml.v2"

	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"

	llngconfig "github.com/lemonldap-ng-controller/lemonldap-ng-controller/internal/lemonldapng/config"
)

// parseIngress returns the ingress namespace, the ingress name, and a map of VHosts
func (c *LemonLDAPNGController) parseIngress(obj interface{}) (string, string, map[string]*llngconfig.VHost, *llngconfig.Application, error) {
	ingressObj := obj.(*extensionsv1beta1.Ingress)
	ingressNamespace := ingressObj.Namespace
	ingressName := ingressObj.Name
	ingressAnnotations := ingressObj.GetAnnotations()
	vhosts := make(map[string]*llngconfig.VHost)

	locationRulesAnnotation := "kubernetes-controller.lemonldap-ng.org/location-rules"
	locationRules := make(map[string]string)
	locationRulesYaml, ok := ingressAnnotations[locationRulesAnnotation]
	if ok {
		err := yaml.Unmarshal([]byte(locationRulesYaml), &locationRules)
		if err != nil {
			return ingressNamespace, ingressName, vhosts, nil, fmt.Errorf("Unable to parse locationRules annotation %s of Ingress %s/%s, ignoring Ingress: %s", locationRulesAnnotation, ingressNamespace, ingressName, err)
		}
	} else {
		locationRules = llngconfig.DefaultLocationRules
	}

	exportedHeadersAnnotation := "kubernetes-controller.lemonldap-ng.org/exported-headers"
	exportedHeaders := make(map[string]string)
	exportedHeadersYaml, ok := ingressAnnotations[exportedHeadersAnnotation]
	if ok {
		err := yaml.Unmarshal([]byte(exportedHeadersYaml), &exportedHeaders)
		if err != nil {
			return ingressNamespace, ingressName, vhosts, nil, fmt.Errorf("Unable to parse exportedHeaders annotation %s of Ingress %s/%s, ignoring Ingress: %s", exportedHeadersAnnotation, ingressNamespace, ingressName, err)
		}
	} else {
		exportedHeaders = llngconfig.DefaultExportedHeaders
	}

	var firstVhost *llngconfig.VHost
	for _, rule := range ingressObj.Spec.Rules {
		serverName := rule.Host
		if serverName == "" || serverName == "*" {
			serverName = "default"
		}
		if rule.HTTP == nil {
			continue
		}
		vhosts[serverName] = llngconfig.NewVHost(serverName, locationRules, exportedHeaders)
		if firstVhost == nil {
			firstVhost = vhosts[serverName]
		}
	}

	application := llngconfig.NewApplication(firstVhost, ingressAnnotations, "kubernetes-controller.lemonldap-ng.org")
	return ingressNamespace, ingressName, vhosts, application, nil
}

func (c *LemonLDAPNGController) ingressAdded(obj interface{}) {
	ingressNamespace, ingressName, vhosts, application, err := c.parseIngress(obj)
	if err != nil {
		glog.Error(err)
		return
	}
	glog.Infof("An ingress was created: %s/%s", ingressNamespace, ingressName)
	c.llngConfig.AddVhosts(vhosts)
	c.llngConfig.AddApplication(application)
	err = c.llngConfig.Save() // FIXME async + batch
	if err != nil {
		glog.Error(err)
		return
	}
}

func (c *LemonLDAPNGController) ingressDeleted(obj interface{}) {
	ingressNamespace, ingressName, vhosts, application, err := c.parseIngress(obj)
	if err != nil {
		glog.Error(err)
		return
	}
	glog.Infof("An ingress was deleted: %s/%s", ingressNamespace, ingressName)
	c.llngConfig.DeleteVhosts(vhosts)
	c.llngConfig.DeleteApplication(application)
	err = c.llngConfig.Save() // FIXME async + batch
	if err != nil {
		glog.Error(err)
		return
	}
}

func (c *LemonLDAPNGController) ingressUpdated(old, cur interface{}) {
	_, _, oldVhosts, oldApplication, err := c.parseIngress(old)
	if err != nil {
		glog.Error(err)
		return
	}
	curIngressNamespace, curIngressName, curVhosts, curApplication, err := c.parseIngress(cur)
	if err != nil {
		glog.Error(err)
		return
	}
	if !reflect.DeepEqual(oldVhosts, curVhosts) {
		glog.Infof("An ingress was updated (vhosts): %s/%s", curIngressNamespace, curIngressName)
		c.llngConfig.DeleteVhosts(oldVhosts)
		c.llngConfig.AddVhosts(curVhosts)
	}
	if !reflect.DeepEqual(oldApplication, curApplication) {
		glog.Infof("An ingress was updated (application): %s/%s", curIngressNamespace, curIngressName)
		c.llngConfig.DeleteApplication(oldApplication)
		c.llngConfig.AddApplication(curApplication)
	}
	err = c.llngConfig.Save() // FIXME async + batch
	if err != nil {
		glog.Error(err)
		return
	}
}
