# LemonLDAP::NG controller

[![Build Status](https://travis-ci.org/lemonldap-ng-controller/lemonldap-ng-controller.svg?branch=master)](https://travis-ci.org/lemonldap-ng-controller/lemonldap-ng-controller)
[![Coverage Status](https://coveralls.io/repos/github/lemonldap-ng-controller/lemonldap-ng-controller/badge.svg?branch=master)](https://coveralls.io/github/lemonldap-ng-controller/lemonldap-ng-controller?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/lemonldap-ng-controller/lemonldap-ng-controller)](https://goreportcard.com/report/github.com/lemonldap-ng-controller/lemonldap-ng-controller)

## Description

This repository contains the LemonLDAP::NG controller built around the [Kubernetes Ingress resource](http://kubernetes.io/docs/user-guide/ingress/) that uses [ConfigMap](https://kubernetes.io/docs/tasks/configure-pod-container/configmap/#understanding-configmaps) to store the LemonLDAP configuration.

It is intended to be used with the [NGINX Ingress Controller](https://github.com/kubernetes/ingress-nginx) with the [nginx.ingress.kubernetes.io/auth-url](https://github.com/kubernetes/ingress-nginx/blob/master/docs/user-guide/annotations.md#external-authentication) annotation.

## Ingress Annotations

The following annotations are supported:

| Name                                                                        | type |
|-----------------------------------------------------------------------------|------|
|[kubernetes-controller.lemonldap-ng.org/location-rules](#location-rules)     | hash |
|[kubernetes-controller.lemonldap-ng.org/exported-headers](#exported-headers) | hash |

### location-rules

```yaml
kubernetes-controller.lemonldap-ng.org/location-rules: |
  {
    "^/admin/": "$uid eq \"bart.simpson\"",
    "default": "accept"
  }
```

If not specified in the Ingress, the default location-rules are:

```yaml
kubernetes-controller.lemonldap-ng.org/location-rules: |
  {
    "default": "accept"
  }
```

Which ensures that the user is authentified.

See also [LemonLDAP::NG documentation](https://www.lemonldap-ng.org/documentation/1.9/writingrulesand_headers#rules).

### exported-headers

```yaml
kubernetes-controller.lemonldap-ng.org/exported-headers: |
  {
    "Display-Name": "$givenName.\" \".$surName"
  }
```

If not specified in the Ingress, the default exported-headers are:

```yaml
kubernetes-controller.lemonldap-ng.org/exported-headers: |
  {
    "Auth-User ": "$uid"
  }
```

See also [LemonLDAP::NG documentation](https://www.lemonldap-ng.org/documentation/1.9/writingrulesand_headers#headers).
