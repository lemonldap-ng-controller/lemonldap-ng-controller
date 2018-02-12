# LemonLDAP::NG controller

[![Build Status](https://travis-ci.org/lemonldap-ng-controller/lemonldap-ng-controller.svg?branch=master)](https://travis-ci.org/lemonldap-ng-controller/lemonldap-ng-controller)
[![Coverage Status](https://coveralls.io/repos/github/lemonldap-ng-controller/lemonldap-ng-controller/badge.svg?branch=master)](https://coveralls.io/github/lemonldap-ng-controller/lemonldap-ng-controller?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/lemonldap-ng-controller/lemonldap-ng-controller)](https://goreportcard.com/report/github.com/lemonldap-ng-controller/lemonldap-ng-controller)

## Description

This repository contains the LemonLDAP::NG controller built around the [Kubernetes Ingress resource](http://kubernetes.io/docs/user-guide/ingress/) that uses [ConfigMap](https://kubernetes.io/docs/tasks/configure-pod-container/configmap/#understanding-configmaps) to store the LemonLDAP configuration.

It is intended to be used with the [NGINX Ingress Controller](https://github.com/kubernetes/ingress-nginx).

## Deployement

See [Deployment](deploy/README.md).

## Ingress Annotations

The following annotations are supported:

| Name                                                                        | type   |
|-----------------------------------------------------------------------------|--------|
|[kubernetes-controller.lemonldap-ng.org/location-rules](#location-rules)     | string |
|[kubernetes-controller.lemonldap-ng.org/exported-headers](#exported-headers) | string |

### location-rules

YAML or JSON are supported:

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

YAML or JSON are supported:

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

## Config Map

A config map can be used to override lmConf-1.js parameters.

Any key suffixed by `.yaml` will be parsed accordingly:

```yaml
kind: ConfigMap
apiVersion: v1
metadata:
  name: lemonldap-ng-configuration
  namespace: ingress-nginx
data:
  domain: example.org
  globalStorage: Apache::Session::Browseable::Postgres # Default Apache::Session::File
  globalStorageOptions.yaml: |
    DataSource: dbi:Pg:dbname=sessions;host=10.2.3.1
    UserName: lemonldapng
    Password: mysuperpassword
    TableName: sessions
    Commit: 1
    Index: _whatToTrace ipAddr
```

This is the most difficult part of LemonLDAP::NG configuration.
Recommended settings include:
- [Single Sign On cookie, domain and portal URL](https://lemonldap-ng.org/documentation/1.9/ssocookie)
- [authentification, user and password backends](https://lemonldap-ng.org/documentation/1.9/start#authentication_users_and_password_databases)
- [session database](https://lemonldap-ng.org/documentation/1.9/start#sessions_database) (if you have more than one replica)

See also the [example ConfigMap](deploy/llng-configmap.yaml) and the [full parameters list from LemonLDAP::NG documentation](https://lemonldap-ng.org/documentation/1.9/parameterlist).

Note: Make sure to have the following to arg in the deployement:
```yaml
- --configmap=ingress-nginx/lemonldap-ng-configuration
```

## Command line flags

```
Usage of /lemonldap-ng-controller:
      --alsologtostderr                               log to standard error as well as files
      --configmap string                              Name of the ConfigMap that contains the custom configuration to use
      --force-namespace-isolation                     Force namespace isolation. This flag is required to avoid the reference of secrets or configmaps located in a different namespace than the specified in the flag --watch-namespace
      --kubeconfig string                             Path to a kubeconfig. Only required if out-of-cluster
      --lemonldap-ng-configuration-directory string   LemonLDAP::NG configuration directory (default "/var/lib/lemonldap-ng/conf")
      --log_backtrace_at traceLocation                when logging hits line file:N, emit a stack trace (default :0)
      --log_dir string                                If non-empty, write log files in this directory
      --logtostderr                                   log to standard error instead of files
      --master string                                 The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster
      --stderrthreshold severity                      logs at or above this threshold go to stderr (default 2)
      --sync-period duration                          Relist and confirm cloud resources this often (default 10m0s)
  -v, --v Level                                       log level for V logs
      --vmodule moduleSpec                            comma-separated list of pattern=N settings for file-filtered logging
      --watch-namespace string                        Namespace to watch for Ingress. Default is to watch all namespaces
```
