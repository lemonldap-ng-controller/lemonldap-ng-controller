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
	"github.com/golang/glog"

	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"

	llngconfig "github.com/lemonldap-ng-controller/lemonldap-ng-controller/internal/lemonldapng/config"
)

// LemonLDAPNGController watches the kubernetes api for changes to ingresses
type LemonLDAPNGController struct {
	controllerConfig         *Configuration
	llngConfig               *llngconfig.Config
	ingressCacheStore        cache.Store
	ingressCacheController   cache.Controller
	configMapCacheStore      cache.Store
	configMapCacheController cache.Controller
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *LemonLDAPNGController) Run(stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()

	// Start the informer factories to begin populating the informer caches
	glog.Info("Starting LemonLDAP::NG controller")

	glog.Info("Starting workers")
	go c.ingressCacheController.Run(stopCh)
	go c.configMapCacheController.Run(stopCh)

	glog.Info("Started workers")
	<-stopCh
	glog.Info("Shutting down workers")

	return nil
}

// NewLemonLDAPNGController returns a new ingress controller
func NewLemonLDAPNGController(controllerConfig *Configuration) *LemonLDAPNGController {
	ingressWatcher := &LemonLDAPNGController{}
	ingressWatcher.controllerConfig = controllerConfig
	ingressWatcher.llngConfig = llngconfig.NewConfig(controllerConfig.FS, controllerConfig.LemonLDAPConfigurationDirectory)

	watchNs := corev1.NamespaceAll
	if controllerConfig.ForceNamespaceIsolation && controllerConfig.Namespace != corev1.NamespaceAll {
		watchNs = controllerConfig.Namespace
	}

	// Create informer for watching Ingresses
	ingEventHandler := cache.ResourceEventHandlerFuncs{
		AddFunc:    ingressWatcher.ingressAdded,
		DeleteFunc: ingressWatcher.ingressDeleted,
		UpdateFunc: ingressWatcher.ingressUpdated,
	}
	ingressWatcher.ingressCacheStore, ingressWatcher.ingressCacheController = cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				return controllerConfig.Client.ExtensionsV1beta1().Ingresses(controllerConfig.Namespace).List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return controllerConfig.Client.ExtensionsV1beta1().Ingresses(controllerConfig.Namespace).Watch(options)
			},
		},
		&extensionsv1beta1.Ingress{}, controllerConfig.ResyncPeriod, ingEventHandler)

	// Create informer for watching ConfigMaps
	mapEventHandler := cache.ResourceEventHandlerFuncs{
		AddFunc:    ingressWatcher.configMapAdded,
		DeleteFunc: ingressWatcher.configMapDeleted,
		UpdateFunc: ingressWatcher.configMapUpdated,
	}
	ingressWatcher.configMapCacheStore, ingressWatcher.configMapCacheController = cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				return controllerConfig.Client.CoreV1().ConfigMaps(watchNs).List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return controllerConfig.Client.CoreV1().ConfigMaps(watchNs).Watch(options)
			},
		},
		&corev1.ConfigMap{}, controllerConfig.ResyncPeriod, mapEventHandler)

	return ingressWatcher
}
