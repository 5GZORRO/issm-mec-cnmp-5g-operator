# issm-mec-cnmp-5g-operator

This project acts as a VNFM for free5gc network functions.

It implements the Kubernetes Operator which is based on [operator-sdk](https://sdk.operatorframework.io) and is responsible to deploy and manage  [Free 5GC](https://www.free5gc.org/) NF components by using declarative [Kubernetes Custom Resource Definitions](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources).

## Prerequisites

* [golang](https://go.dev/)
* [operator-sdk](https://sdk.operatorframework.io)

## Installation

Log into kubernetes master

### golang

Install golang **v1.16**: https://golang.org/doc/install

then issue

```
source ~/.profile
```

validate

```
go version
```

### operator-sdk

Install operator-sdk **v1.8.0** from [install-from-github-release](https://sdk.operatorframework.io/docs/installation/#install-from-github-release)

### 5g-operator

Clone

```bash
cd ~
git clone https://github.com/5GZORRO/issm-mec-cnmp-5g-operator.git
cd issm-mec-cnmp-5g-operator
git checkout v3.0.6-dynamic-load
```

Deploy the operator

```bash
make generate
make manifests
make deploy
```

Wait for controller pod to start

```
kubectl get pod -n 5g
```

**Notes:** 

* before using 'make', load your profile: `source ~/.profile`
* to un-install the operator: `make undeploy`

## Build (**relevant for developers only**)

1. Edit Makefile with `VERSION ?= temp` so that the resulted image tag does not collide with the existing one.

1. Edit Makefile with `IMAGE_TAG_BASE` with the proper image registry. Note: current version uses an internal registry to hold the operator and 5G network function images.

1. Build and push the image.

    ```
    make generate
    ```
    
    ```
    make manifests
    ```
    
    ```
    make docker-build docker-push
    ```

1. Deploy the operator

   ```
   make deploy
   ```

## Maintainers
**Avi Weit** - weit@il.ibm.com

**David Breitgand** - davidbr@il.ibm.com

## Licensing

This 5GZORRO component is published under Apache 2.0 license. Please see the [LICENSE](./LICENSE) file for further details.