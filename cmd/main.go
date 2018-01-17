/*
Copyright 2017 The Kubernetes Authors.
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

package cmd

import (
	"flag"
	"time"

	"github.com/golang/glog"

	kubeinformers "k8s.io/client-go/informers"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"github.com/lemonldap-ng-controller/lemonldap-ng-controller/pkg/controller"
	"github.com/lemonldap-ng-controller/lemonldap-ng-controller/pkg/signals"
)

var (
	masterURL      string
	kubeconfig     string
	watchNamespace string
)

func Execute() {
	flag.Parse()

	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		glog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, time.Second*30)

	ingressController := controller.NewIngressController(kubeClient, watchNamespace)

	go kubeInformerFactory.Start(stopCh)

	if err = ingressController.Run(stopCh); err != nil {
		glog.Fatalf("Error running controller: %s", err.Error())
	}
}

func Init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")

    flag.StringVar(&watchNamespace, "watch-namespace", corev1.NamespaceAll, "Namespace to watch for Ingress. Default is to watch all namespaces")
}
