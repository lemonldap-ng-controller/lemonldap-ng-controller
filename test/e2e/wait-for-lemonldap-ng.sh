#!/bin/bash

# Copyright 2018 Mathieu Parent <math.parent@gmail.com>.
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

export JSONPATH='{range .items[*]}{range @.status.containerStatuses[*]}{@.name}:ready={@.ready}{"\n"}{end}{end}'
export EXPECTEDCONTAINERS='lemonldap-ng-controller:ready=true
nginx-ingress-controller:ready=true'

echo "deploying lemonldap-ng-controller ConfigMaps..."
cat deploy/llng-configmap.yaml | kubectl apply -f -
cat deploy/llng-nginx-configmap.yaml | kubectl apply -f -

echo "Adding lemonldap-ng-controller container to nginx-ingress-controller Deployment..."
kubectl patch deployment \
    --namespace ingress-nginx \
    nginx-ingress-controller\
    --patch "$(cat deploy/llng-patch-deployement.yaml)"

echo "updating image..."
kubectl set image \
    deployments \
    --namespace ingress-nginx \
	--selector app=ingress-nginx \
    lemonldap-ng-controller=lemonldapng/lemonldap-ng-controller:test

sleep 5

echo "waiting LemonLDAP::NG pod..."

function waitForPod() {
    until [ "$(kubectl get pods -n ingress-nginx -l app=ingress-nginx -o jsonpath="$JSONPATH" 2>&1 | grep "ready=true")" = "$EXPECTEDCONTAINERS" ] ;
    do
        sleep 1;
    done
}

export -f waitForPod

timeout 20s bash -c waitForPod

if [ "$(kubectl get pods -n ingress-nginx -l app=ingress-nginx -o jsonpath="$JSONPATH" 2>&1 | grep "ready=true")" = "$EXPECTEDCONTAINERS" ];
then
    echo "Kubernetes deployments started:"
    kubectl get pods -n ingress-nginx
else
    echo "Kubernetes deployments with issues:"
    kubectl get pods -n ingress-nginx

    echo "Reason:"
    kubectl describe pods -n ingress-nginx
    kubectl logs -n ingress-nginx -l app=ingress-nginx
    exit 1
fi
