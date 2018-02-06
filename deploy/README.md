# Installation Guide

## Contents

- [Requirement: NGINX Ingress Controller](#requirement-nginx-ingress-controller)
- [Generic Install](#generic-install)
- [Verify installation](#verify-installation)

## Requirement: NGINX Ingress Controller

LemonLDAP::NG controller is an additionnal container in the same pod as the NGINX Ingress Controller.

You need to deploy ingress-nginx first. See [their Deployment README](https://github.com/kubernetes/ingress-nginx/blob/master/deploy/README.md).

### Generic install

```console
curl https://raw.githubusercontent.com/lemonldap-ng-controller/lemonldap-ng-controller/master/deploy/llng-configmap.yaml \
    | kubectl apply -f -

curl https://raw.githubusercontent.com/lemonldap-ng-controller/lemonldap-ng-controller/master/deploy/llng-nginx-configmap.yaml \
    | kubectl apply -f -

kubectl patch deployment -n ingress-nginx nginx-ingress-controller \
  --patch="$(curl https://raw.githubusercontent.com/lemonldap-ng-controller/lemonldap-ng-controller/master/deploy/llng-patch-deployement.yaml)"
```

## Verify installation

To check if the ingress controller pods have started, run the following command:

```console
kubectl get pods --all-namespaces -l app=ingress-nginx --watch
```

Once the operator pods are running, you can cancel the above command by typing `Ctrl+C`.
Now, you are ready to create your first ingress.
