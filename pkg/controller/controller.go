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
	"time"

	"github.com/golang/glog"

	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

// IngressController watches the kubernetes api for changes to ingresses
type IngressController struct {
	ingressInformer cache.SharedIndexInformer
	kclient         *kubernetes.Clientset
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *IngressController) Run(stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()

	// Start the informer factories to begin populating the informer caches
	glog.Info("Starting Foo controller")

	glog.Info("Starting workers")
	go c.ingressInformer.Run(stopCh)

	glog.Info("Started workers")
	<-stopCh
	glog.Info("Shutting down workers")

	return nil
}

func NewIngressController(kclient *kubernetes.Clientset, namespace string) *IngressController {
	ingressWatcher := &IngressController{}

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
		AddFunc: ingressWatcher.ingressAdded,
	})

	ingressWatcher.kclient = kclient
	ingressWatcher.ingressInformer = ingressInformer

	return ingressWatcher
}

func (c *IngressController) ingressAdded(obj interface{}) {
	ingressObj := obj.(*extensionsv1beta1.Ingress)
	ingressNamespace := ingressObj.Namespace
	ingressName := ingressObj.Name

	glog.Infof("An ingress was created: %s", ingressNamespace, ingressName)
}
