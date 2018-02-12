/*
Copyright 2018 Mathieu Parent <math.parent@gmail.com>.

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

package setting

import (
	"fmt"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/parnurzeal/gorequest"

	"github.com/lemonldap-ng-controller/lemonldap-ng-controller/test/e2e/framework"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var _ = framework.IngressNginxDescribe("Portal URL", func() {
	f := framework.NewDefaultFramework("portal")

	BeforeEach(func() {
		err := f.NewEchoDeployment()
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
	})

	It("should respect portal parameter", func() {
		host := "goto-portal.example.com"

		setting := "portal"
		oldValue := updateConfigmap(setting, "http://portal.example.com/", f.KubeClientSet)
		defer updateConfigmap(setting, oldValue, f.KubeClientSet)

		bi := buildIngress(host, f.Namespace.Name)

		ing, err := f.EnsureIngress(bi)
		Expect(err).NotTo(HaveOccurred())
		Expect(ing).NotTo(BeNil())

		err = f.WaitForNginxServer(host,
			func(server string) bool {
				return Expect(server).Should(ContainSubstring("server_name goto-portal.example.com")) &&
					Expect(server).ShouldNot(ContainSubstring("return 503"))
			})
		Expect(err).NotTo(HaveOccurred())

		time.Sleep(2 * time.Second) // FIXME wait for LLNG config reload

		resp, _, errs := gorequest.New().
			Get(f.NginxHTTPURL).
			RedirectPolicy(func(req gorequest.Request, via []gorequest.Request) error { return http.ErrUseLastResponse }).
			Set("Host", host).
			End()

		if len(errs) > 0 {
			Expect(errs[0]).NotTo(HaveOccurred())
		}
		Expect(resp.StatusCode).Should(Equal(http.StatusFound))
		Expect(resp.Header.Get("Location")).Should(ContainSubstring("http://portal.example.com/?url="))
	})
})

func updateConfigmap(k, v string, c kubernetes.Interface) string {
	By(fmt.Sprintf("updating configuration configmap setting %v to '%v'", k, v))
	config, err := c.CoreV1().ConfigMaps("ingress-nginx").Get("lemonldap-ng-configuration", metav1.GetOptions{})
	Expect(err).NotTo(HaveOccurred())
	Expect(config).NotTo(BeNil())

	if config.Data == nil {
		config.Data = map[string]string{}
	}
	oldValue := config.Data[k]

	if oldValue == v {
		return oldValue
	}

	config.Data[k] = v
	_, err = c.CoreV1().ConfigMaps("ingress-nginx").Update(config)
	Expect(err).NotTo(HaveOccurred())
	return oldValue
}
