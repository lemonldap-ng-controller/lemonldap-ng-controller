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

package main

import (
	goflag "flag"
	"time"

	"github.com/golang/glog"
	flag "github.com/spf13/pflag"

	corev1 "k8s.io/api/core/v1"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"github.com/lemonldap-ng-controller/lemonldap-ng-controller/internal/controller"
	"github.com/lemonldap-ng-controller/lemonldap-ng-controller/internal/filesystem/os"
	"github.com/lemonldap-ng-controller/lemonldap-ng-controller/internal/signals"
)

var (
	config *controller.Configuration = &controller.Configuration{
		FS: &os.FileSystem{},
	}
)

func main() {
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	flag.Parse()
	if len(flag.Args()) > 0 {
		if flag.Arg(0) != "--" {
			glog.Fatalf("First non-flag parameter should be --, got %s", flag.Arg(0))
		}
		config.Command = flag.Args()[1:]
	} else {
		config.Command = []string{"/usr/sbin/llng-fastcgi-server", "--foreground", "--listen", "127.0.0.1:9000"}
	}

	// Workaround for "ERROR: logging before flag.Parse". See:
	// https://github.com/golang/glog/pull/13
	goflag.CommandLine.Parse([]string{})

	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()

	cfg, err := clientcmd.BuildConfigFromFlags(config.APIServerHost, config.KubeConfigFile)
	if err != nil {
		glog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	config.Client, err = kubernetes.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(config.Client, time.Second*30)

	ingressController := controller.NewLemonLDAPNGController(config)

	go kubeInformerFactory.Start(stopCh)

	if err = ingressController.Run(stopCh); err != nil {
		glog.Fatalf("Error running controller: %s", err.Error())
	}
}

func init() {
	flag.StringVar(&config.APIServerHost, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster")
	flag.StringVar(&config.KubeConfigFile, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster")

	flag.StringVar(&config.ConfigMapName, "configmap", "", "Name of the ConfigMap that contains the custom configuration to use")
	flag.DurationVar(&config.ResyncPeriod, "sync-period", 600*time.Second, "Relist and confirm cloud resources this often")
	flag.StringVar(&config.Namespace, "watch-namespace", corev1.NamespaceAll, "Namespace to watch for Ingress. Default is to watch all namespaces")
	flag.BoolVar(&config.ForceNamespaceIsolation, "force-namespace-isolation", false, "Force namespace isolation. This flag is required to avoid the reference of secrets or configmaps located in a different namespace than the specified in the flag --watch-namespace")
	flag.StringVar(&config.LemonLDAPConfigurationDirectory, "lemonldap-ng-configuration-directory", "/var/lib/lemonldap-ng/conf", "LemonLDAP::NG configuration directory")
}
