#!/bin/bash

# Copyright 2018 Mathieu Parent <math.parent@gmail.com>
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

GOPATH=$(go env GOPATH)

if [ "$(pwd)" != "$GOPATH/src/github.com/lemonldap-ng-controller/lemonldap-ng-controller" ]; then
  echo "Working dir not in GOPATH. Should be '$GOPATH/src/github.com/lemonldap-ng-controller/lemonldap-ng-controller'" >&2
  exit 1
fi

for entry in \
  'k8s.io/ingress-nginx/deploy/namespace.yaml deploy' \
  'k8s.io/ingress-nginx/deploy/default-backend.yaml deploy' \
  'k8s.io/ingress-nginx/deploy/configmap.yaml deploy' \
  'k8s.io/ingress-nginx/deploy/tcp-services-configmap.yaml deploy' \
  'k8s.io/ingress-nginx/deploy/udp-services-configmap.yaml deploy' \
  'k8s.io/ingress-nginx/deploy/without-rbac.yaml deploy' \
  'k8s.io/ingress-nginx/deploy/provider/baremetal/service-nodeport.yaml deploy/provider/baremetal' \
  'k8s.io/ingress-nginx/.github/* .github' \
  'k8s.io/ingress-nginx/hack/kube-env.sh hack' \
  'k8s.io/ingress-nginx/hack/verify-all.sh hack' \
  'k8s.io/ingress-nginx/hack/verify-gofmt.sh hack' \
  'k8s.io/ingress-nginx/hack/verify-golint.sh hack' \
  'k8s.io/ingress-nginx/test/e2e/*.go test/e2e' \
  'k8s.io/ingress-nginx/test/e2e/framework/*.go test/e2e/framework' \
  'k8s.io/ingress-nginx/test/e2e/up.sh test/e2e' \
  'k8s.io/ingress-nginx/test/e2e/wait-for-nginx.sh test/e2e' \
  'k8s.io/sample-controller/LICENSE .' \
  'k8s.io/sample-controller/pkg/signals/*.go internal/signals'
do
  src="$GOPATH/src/$(echo $entry | awk '{print $1}')"
  dst_dir="$(echo $entry | awk '{print $2}')"
  mkdir -p $dst_dir
  cp $src $dst_dir/
done

# Keep original image
sed -i 's@^kubectl set image@true no set image@' test/e2e/wait-for-nginx.sh

# Change ingress-nginx e2e import path
sed -i \
  -e /defaultbackend/d \
  -e /ssl/d \
  test/e2e/e2e.go
sed -i \
  's@k8s.io/ingress-nginx/test/e2e@github.com/lemonldap-ng-controller/lemonldap-ng-controller/test/e2e@' \
  test/e2e/*.go \
  test/e2e/framework/*.go
sed -i \
  's@\t"k8s.io/ingress-nginx/internal/file"@\n\tfile "github.com/lemonldap-ng-controller/lemonldap-ng-controller/internal/filesystem/os"@' \
  test/e2e/framework/util.go

# ex: ts=2 sw=2 et filetype=sh
