# Installation

## Prerequisites

* kubectl on your desktop
* [lmctl](https://github.com/IBM/lmctl) on your desktop
* Docker on your desktop
* A `target` Kubernetes cluster in to which the 5G network will be installed (can be different from the one running CP4NA)
* A runnning CP4NA installation with Kubernetes driver installed and onboarded:

```
lmctl resourcedriver add --type kubernetes --url https://kubedriver:8294 dev --certificate kubedriver.pem
```

## Build xnf images

Run this command on your target Kubernetes cluster so that the 5G images are built there

```
cd cp4na/xnf-images
make
```

## Install

Install the 5G CP4NA resources in to CP4NA.

```
cd cp4na/xnfs/5G/amf
lmctl project push [environment]
```

Do the same for all the other 5G components.

Install the 5G Core assembly in to CP4NA.

```
cd cp4na/Core
lmctl project push [environment]
```

Deploy the 5G operator to your target Kubernetes cluster (assumes kubectl has been configured to point to your target cluster)

```
make deploy
```


