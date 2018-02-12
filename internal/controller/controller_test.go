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
	"flag"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/golang/glog"

	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	fakeclient "k8s.io/client-go/kubernetes/fake"

	fakefs "github.com/lemonldap-ng-controller/lemonldap-ng-controller/internal/filesystem/fake"
)

func buildFakeClientSet() *fakeclient.Clientset {
	return fakeclient.NewSimpleClientset(
		&extensionsv1beta1.IngressList{Items: []extensionsv1beta1.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ingress1",
					Namespace: corev1.NamespaceDefault,
					Annotations: map[string]string{
						"kubernetes-controller.lemonldap-ng.org/location-rules": `{"^/admin/": "$uid eq \"bart.simpson\"","default": "accept"}`,
					},
				},
				Spec: extensionsv1beta1.IngressSpec{
					Rules: []extensionsv1beta1.IngressRule{
						{
							Host: "test1.example.org",
							IngressRuleValue: extensionsv1beta1.IngressRuleValue{
								HTTP: &extensionsv1beta1.HTTPIngressRuleValue{
									Paths: []extensionsv1beta1.HTTPIngressPath{
										{
											Path: "/foo",
											Backend: extensionsv1beta1.IngressBackend{
												ServiceName: "test1-backend",
												ServicePort: intstr.FromInt(80),
											},
										},
									},
								},
							},
						},
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ingress2",
					Namespace: "test-ns",
				},
				Spec: extensionsv1beta1.IngressSpec{
					Rules: []extensionsv1beta1.IngressRule{
						{
							Host: "test2.example.org",
							IngressRuleValue: extensionsv1beta1.IngressRuleValue{
								HTTP: &extensionsv1beta1.HTTPIngressRuleValue{
									Paths: []extensionsv1beta1.HTTPIngressPath{
										{
											Path: "/foo",
											Backend: extensionsv1beta1.IngressBackend{
												ServiceName: "test2-backend",
												ServicePort: intstr.FromInt(80),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}},
		&corev1.ConfigMapList{Items: []corev1.ConfigMap{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-cm",
					Namespace: "test-ns",
				},
				Data: map[string]string{
					"domain": "example.org",
					"globalStorageOptions.yaml": `DataSource: dbi:Pg:dbname=sessions;host=10.2.3.1
UserName: lemonldapng
Password: mysuperpassword
TableName: sessions
Commit: 1
Index: _whatToTrace ipAddr`,
					"unsupportedSuffix.raw": "should-be -ignored",
				},
			},
		}},
	)
}

func buildControllerConfig(namespace string, forceNamespaceIsolation bool) *Configuration {
	return &Configuration{
		APIServerHost:           "",
		KubeConfigFile:          "",
		Client:                  buildFakeClientSet(),
		ResyncPeriod:            time.Hour,
		ConfigMapName:           "test-ns/test-cm",
		Namespace:               namespace,
		ForceNamespaceIsolation: forceNamespaceIsolation,
		FS: fakefs.NewFilesystem(),
		LemonLDAPConfigurationDirectory: "/var/lib/lemonldap-ng/conf",
		Command: []string{"/bin/true"},
	}
}

func checkLLConfig(t *testing.T, c *LemonLDAPNGController, cfgNum int, checks []*regexp.Regexp) {
	configName := fmt.Sprintf("lmConf-%d.js", cfgNum)
	configPath := "/var/lib/lemonldap-ng/conf/" + configName
	lmConf, errRead := c.controllerConfig.FS.ReadFile(configPath)
	if errRead != nil {
		lastConfigName, _, _ := c.llngConfig.Last()
		t.Errorf("Unable to read %s (last configuration is %s): %s", configPath, lastConfigName, errRead)
		return
	}

	for _, re := range checks {
		if !re.Match(lmConf) {
			t.Errorf("%s to match %s\n%s", configName, re, lmConf)
		}
	}
}

func TestNewLemonLDAPNGController(t *testing.T) {
	flag.Set("alsologtostderr", "true")

	for _, namespace := range []string{corev1.NamespaceAll, "test-ns", "another-ns"} {
		for _, forceNamespaceIsolation := range []bool{false, true} {
			glog.Infof("With namespace=%s, forceNamespaceIsolation=%v", namespace, forceNamespaceIsolation)
			t.Logf("With namespace=%s, forceNamespaceIsolation=%v", namespace, forceNamespaceIsolation)
			stopCh := make(chan struct{})
			controllerConfig := buildControllerConfig(namespace, forceNamespaceIsolation)
			ingressController := NewLemonLDAPNGController(controllerConfig)

			// FIXME We need a better way than sleeping
			time.AfterFunc(3*time.Second, func() {
				controllerConfig.Client.ExtensionsV1beta1().Ingresses(corev1.NamespaceDefault).Delete("test-ingress1", nil)
				ing2 := &extensionsv1beta1.Ingress{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-ingress2",
						Namespace: "test-ns",
						Annotations: map[string]string{
							"kubernetes-controller.lemonldap-ng.org/location-rules": `{"^/admin/": "$uid eq \"bart.simpson\"","default": "accept"}`,
						},
					},
					Spec: extensionsv1beta1.IngressSpec{
						Rules: []extensionsv1beta1.IngressRule{
							{
								Host: "test2.example.org",
								IngressRuleValue: extensionsv1beta1.IngressRuleValue{
									HTTP: &extensionsv1beta1.HTTPIngressRuleValue{
										Paths: []extensionsv1beta1.HTTPIngressPath{
											{
												Path: "/foo",
												Backend: extensionsv1beta1.IngressBackend{
													ServiceName: "test2-backend",
													ServicePort: intstr.FromInt(80),
												},
											},
										},
									},
								},
							},
						},
					},
				}
				controllerConfig.Client.ExtensionsV1beta1().Ingresses("test-ns").Update(ing2)

				controllerConfig.Client.CoreV1().ConfigMaps("test-ns").Delete("test-cm", nil)
			})
			time.AfterFunc(5*time.Second, func() {
				close(stopCh)
			})

			if err := ingressController.Run(stopCh); err != nil {
				t.Fatalf("Error running controller: %s", err.Error())
			}

			configNum := 1
			var /* const */ domainNoneRE = regexp.MustCompile(`"cfgNum": \d+,\s*"exportedHeaders`)
			var /* const */ domainExampleOrgRE = regexp.MustCompile(`"domain": "example.org",`)
			domainRE := domainNoneRE

			var /* const */ globalStorageOptionsNoneRE = regexp.MustCompile(`"cfgNum": \d+,\s*"exportedHeaders`)
			var /* const */ globalStorageOptionsPostgreRE = regexp.MustCompile(`},\s*"globalStorageOptions": {\s*"Commit": 1,\s*"DataSource": "dbi:Pg:dbname=sessions;host=10.2.3.1",\s*"Index": "_whatToTrace ipAddr",\s*"Password": "mysuperpassword",\s*"TableName": "sessions",\s*"UserName": "lemonldapng"\s*},\s*"locationRules": {`)
			globalStorageOptionsRE := globalStorageOptionsNoneRE

			var /* const */ exportedHeadersNoneRE = regexp.MustCompile(`"exportedHeaders": {}`)
			var /* const */ exportedHeadersBothRE = regexp.MustCompile(`"exportedHeaders": {\s*"test1.example.org": {\s*"Auth-User": "\$uid"\s*},\s*"test2.example.org": {\s*"Auth-User": "\$uid"\s*}\s*}`)
			//var /* const */ exportedHeadersTest1RE = regexp.MustCompile(`"exportedHeaders": {\s*"test1.example.org": {\s*"Auth-User": "\$uid"\s*}\s*}`)
			var /* const */ exportedHeadersTest2RE = regexp.MustCompile(`"exportedHeaders": {\s*"test2.example.org": {\s*"Auth-User": "\$uid"\s*}\s*}`)
			exportedHeadersRE := exportedHeadersNoneRE

			var /* const */ locationRulesNoneRE = regexp.MustCompile(`"locationRules": {}`)
			var /* const */ locationRulesBothRE = regexp.MustCompile(`"locationRules": {\s*"test1.example.org": {\s*"\^/admin/": "\$uid eq \\"bart.simpson\\"",\s*"default": "accept"\s*},\s*"test2.example.org": {\s*"default": "accept"\s*}\s*}`)
			//var /* const */ locationRulesTest1RE = regexp.MustCompile(`"locationRules": {\s*"test1.example.org": {\s*"\^/admin/": "\$uid eq \\"bart.simpson\\"",\s*"default": "accept"\s*}\s*}`)
			var /* const */ locationRulesTest2RE = regexp.MustCompile(`"locationRules": {\s*"test2.example.org": {\s*"default": "accept"\s*}\s*}`)
			var /* const */ locationRulesTest2UpdatedRE = regexp.MustCompile(`"locationRules": {\s*"test2.example.org": {\s*"\^/admin/": "\$uid eq \\"bart.simpson\\"",\s*"default": "accept"\s*}\s*}`)
			locationRulesRE := locationRulesNoneRE

			// A ConfigMap was added: test-ns/test-cm
			if namespace == "test-ns" || namespace == corev1.NamespaceAll || !forceNamespaceIsolation {
				configNum++
				domainRE = domainExampleOrgRE
				globalStorageOptionsRE = globalStorageOptionsPostgreRE
			}
			// An ingress was created: test-ns/test-ingress2
			if namespace == corev1.NamespaceAll || namespace == "test-ns" {
				configNum++
				exportedHeadersRE = exportedHeadersTest2RE
				locationRulesRE = locationRulesTest2RE
			}
			// An ingress was created: default/test-ingress1
			if namespace == corev1.NamespaceAll {
				configNum++
				exportedHeadersRE = exportedHeadersBothRE
				locationRulesRE = locationRulesBothRE
			}

			cfgNumRE := regexp.MustCompile(fmt.Sprintf("\"cfgNum\": %d,", configNum))
			checkLLConfig(t, ingressController, configNum, []*regexp.Regexp{
				cfgNumRE,
				domainRE,
				globalStorageOptionsRE,
				exportedHeadersRE,
				locationRulesRE,
			})

			// A ConfigMap was deleted: test-ns/test-cm
			if namespace == "test-ns" || namespace == corev1.NamespaceAll || !forceNamespaceIsolation {
				configNum++
				domainRE = domainNoneRE
				globalStorageOptionsRE = globalStorageOptionsNoneRE
			}
			// An ingress was updated: test-ns/test-ingress2
			if namespace == corev1.NamespaceAll || namespace == "test-ns" {
				configNum++
				exportedHeadersRE = exportedHeadersTest2RE
				locationRulesRE = locationRulesTest2UpdatedRE
			}
			// An ingress was deleted: default/test-ingress1
			if namespace == corev1.NamespaceAll {
				configNum++
			}

			cfgNumRE = regexp.MustCompile(fmt.Sprintf("\"cfgNum\": %d,", configNum))
			checkLLConfig(t, ingressController, configNum, []*regexp.Regexp{
				cfgNumRE,
				domainRE,
				globalStorageOptionsRE,
				exportedHeadersRE,
				locationRulesRE,
			})

			_, lastConfigNum, err := ingressController.llngConfig.Last()
			if err != nil {
				t.Errorf("Unable to get last configuration name: %s", err)
				return
			}
			if lastConfigNum != configNum {
				t.Errorf("configNum mismatch: %d != %d", lastConfigNum, configNum)
			}
		}
	}
}
