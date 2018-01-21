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
	"time"

	"github.com/golang/glog"
	"gopkg.in/yaml.v2"

	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	"github.com/lemonldap-ng-controller/lemonldap-ng-controller/pkg/lemonldapng/config"
)

// IngressController watches the kubernetes api for changes to ingresses
type IngressController struct {
	ingressInformer cache.SharedIndexInformer
	kclient         *kubernetes.Clientset
	llngConfig      *config.Config
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *IngressController) Run(stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()

	// Start the informer factories to begin populating the informer caches
	glog.Info("Starting LemonLDAP::NG controller")

	glog.Info("Starting workers")
	go c.ingressInformer.Run(stopCh)

	glog.Info("Started workers")
	<-stopCh
	glog.Info("Shutting down workers")

	return nil
}

// NewIngressController returns a new ingress controller
func NewIngressController(kclient *kubernetes.Clientset, namespace string, llngConfigDir string) *IngressController {
	ingressWatcher := &IngressController{}
	ingressWatcher.llngConfig = config.NewConfig(llngConfigDir)

	// Create informer for watching Ingresses
	ingressInformer := cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				return kclient.ExtensionsV1beta1().Ingresses(namespace).List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return kclient.ExtensionsV1beta1().Ingresses(namespace).Watch(options)
			},
		},
		&extensionsv1beta1.Ingress{},
		3*time.Minute,
		cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
	)

	ingressInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    ingressWatcher.ingressAdded,
		DeleteFunc: ingressWatcher.ingressDeleted,
		UpdateFunc: ingressWatcher.ingressUpdated,
	})

	ingressWatcher.kclient = kclient
	ingressWatcher.ingressInformer = ingressInformer

	return ingressWatcher
}

// parseIngress returns the ingress namespace, the ingress name, and a map of VHosts
func (c *IngressController) parseIngress(obj interface{}) (string, string, map[string]*config.VHost, error) {
	ingressObj := obj.(*extensionsv1beta1.Ingress)
	ingressNamespace := ingressObj.Namespace
	ingressName := ingressObj.Name
	ingressAnnotations := ingressObj.GetAnnotations()
	vhosts := make(map[string]*config.VHost)

	locationRulesAnnotation := "kubernetes-controller.lemonldap-ng.org/location-rules"
	locationRules := make(map[string]string)
	locationRulesYaml, ok := ingressAnnotations[locationRulesAnnotation]
	if ok {
		err := yaml.Unmarshal([]byte(locationRulesYaml), &locationRules)
		if err != nil {
			return ingressNamespace, ingressName, vhosts, fmt.Errorf("Unable to parse locationRules annotation %s of Ingress %s/%s, ignoring Ingress: %s", locationRulesAnnotation, ingressNamespace, ingressName, err)
		}
	} else {
		locationRules = config.DefaultLocationRules
	}

	exportedHeadersAnnotation := "kubernetes-controller.lemonldap-ng.org/exported-headers"
	exportedHeaders := make(map[string]string)
	exportedHeadersYaml, ok := ingressAnnotations[exportedHeadersAnnotation]
	if ok {
		err := yaml.Unmarshal([]byte(exportedHeadersYaml), &exportedHeaders)
		if err != nil {
			return ingressNamespace, ingressName, vhosts, fmt.Errorf("Unable to parse exportedHeaders annotation %s of Ingress %s/%s, ignoring Ingress: %s", exportedHeadersAnnotation, ingressNamespace, ingressName, err)
		}
	} else {
		exportedHeaders = config.DefaultExportedHeaders
	}

	for _, rule := range ingressObj.Spec.Rules {
		serverName := rule.Host
		if serverName == "" || serverName == "*" {
			serverName = "default"
		}
		if rule.HTTP == nil {
			continue
		}
		vhosts[serverName] = config.NewVHost(serverName, locationRules, exportedHeaders)
	}
	return ingressNamespace, ingressName, vhosts, nil
}

func (c *IngressController) ingressAdded(obj interface{}) {
	ingressNamespace, ingressName, vhosts, err := c.parseIngress(obj)
	if err != nil {
		glog.Error(err)
		return
	}
	glog.Infof("An ingress was created: %s/%s", ingressNamespace, ingressName)
	c.llngConfig.AddVhosts(vhosts)
	c.llngConfig.Save() // FIXME async + batch
}

func (c *IngressController) ingressDeleted(obj interface{}) {
	ingressNamespace, ingressName, vhosts, err := c.parseIngress(obj)
	if err != nil {
		glog.Error(err)
		return
	}
	glog.Infof("An ingress was deleted: %s/%s", ingressNamespace, ingressName)
	c.llngConfig.DeleteVhosts(vhosts)
	c.llngConfig.Save() // FIXME async + batch
}

func (c *IngressController) ingressUpdated(old, cur interface{}) {
	_, _, oldVhosts, oldErr := c.parseIngress(cur)
	if oldErr != nil {
		glog.Error(oldErr)
		return
	}
	curIngressNamespace, curIngressName, curVhosts, curErr := c.parseIngress(cur)
	if curErr != nil {
		glog.Error(curErr)
		return
	}
	glog.Infof("An ingress was updated: %s/%s", curIngressNamespace, curIngressName)
	c.llngConfig.DeleteVhosts(oldVhosts)
	c.llngConfig.AddVhosts(curVhosts)
	c.llngConfig.Save() // FIXME async + batch
}
