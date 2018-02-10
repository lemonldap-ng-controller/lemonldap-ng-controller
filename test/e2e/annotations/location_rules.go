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

package annotations

import (
	"encoding/base64"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/parnurzeal/gorequest"

	"github.com/lemonldap-ng-controller/lemonldap-ng-controller/test/e2e/framework"
)

var _ = framework.IngressNginxDescribe("Annotations - location-rules", func() {
	f := framework.NewDefaultFramework("location-rules")

	BeforeEach(func() {
		err := f.NewEchoDeploymentWithReplicas(2)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
	})

	It("should redirect to LemonLDAP::NG portal when default location-rule is accept", func() {
		host := "default-accept.example.com"

		bi := buildIngress(host, f.Namespace.Name)
		bi.Annotations["kubernetes-controller.lemonldap-ng.org/location-rules"] = `{"default": "accept"}`

		ing, err := f.EnsureIngress(bi)
		Expect(err).NotTo(HaveOccurred())
		Expect(ing).NotTo(BeNil())

		err = f.WaitForNginxServer(host,
			func(server string) bool {
				return Expect(server).Should(ContainSubstring("server_name default-accept.example.com")) &&
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
		expectedLocationHeader := "http://auth.example.org/?url=" + base64.StdEncoding.EncodeToString([]byte("http://default-accept.example.com/"))
		Expect(resp.Header.Get("Location")).Should(Equal(expectedLocationHeader))
	})

	It("should return status code 200 when default location-rule is skip", func() {
		host := "default-skip.example.com"

		bi := buildIngress(host, f.Namespace.Name)
		bi.Annotations["kubernetes-controller.lemonldap-ng.org/location-rules"] = `{"default": "skip"}`

		ing, err := f.EnsureIngress(bi)
		Expect(err).NotTo(HaveOccurred())
		Expect(ing).NotTo(BeNil())

		err = f.WaitForNginxServer(host,
			func(server string) bool {
				return Expect(server).Should(ContainSubstring("server_name default-skip.example.com")) &&
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
		Expect(resp.StatusCode).Should(Equal(http.StatusOK))
		Expect(resp.Header.Get("Location")).Should(Equal(""))
	})
})
