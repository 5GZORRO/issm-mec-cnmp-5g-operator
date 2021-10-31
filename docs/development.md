# Developing

## Configuring Development Environment

Install the Operator SDK:

```
brew update
brew install operator-sdk
brew upgrade operator-sdk
```

Add export GO111MODULE=on to ~/.bash_profile

## Creating the 5G Operator project

This is the first step in creating the operator.

```
mkdir 5goperator
cd 5goperator/
operator-sdk init --domain "ibm.com" --license apache2 --owner "IBM" --project-name 5goperator --repo github.ibm.com/tnc/5GOperator --skip-go-version-check
operator-sdk create api --group 5g --version v1alpha1 --kind Amf --resource --controller
operator-sdk create api --group 5g --version v1alpha1 --kind Ausf --resource --controller
operator-sdk create api --group 5g --version v1alpha1 --kind Mongo --resource --controller
operator-sdk create api --group 5g --version v1alpha1 --kind N3iwf --resource --controller
operator-sdk create api --group 5g --version v1alpha1 --kind Nrf --resource --controller
operator-sdk create api --group 5g --version v1alpha1 --kind Nssf --resource --controller
operator-sdk create api --group 5g --version v1alpha1 --kind Pcf --resource --controller
operator-sdk create api --group 5g --version v1alpha1 --kind Smf --resource --controller
operator-sdk create api --group 5g --version v1alpha1 --kind Transition --resource --controller
operator-sdk create api --group 5g --version v1alpha1 --kind Udm --resource --controller
operator-sdk create api --group 5g --version v1alpha1 --kind Udr --resource --controller
operator-sdk create api --group 5g --version v1alpha1 --kind Upf --resource --controller
operator-sdk create api --group 5g --version v1alpha1 --kind Webconsole --resource --controller

make generate
make manifests

go get k8s.io/api/core/v1@v0.19.2
go mod download github.com/openshift/api
go get sigs.k8s.io/yaml@v1.2.0
go get k8s.io/apimachinery/pkg/util/httpstream/spdy@v0.19.2

# Build the Docker image
make docker-build docker-push

# Deploy the operator to a Kubernetes cluster
make deploy
kubectl apply -f config/samples/5g_v1alpha1_amf.yaml
kubectl get amf
kubectl get pods --all-namespaces
kubectl logs -f fivegoperator-controller-manager-5b5c77b4fc-n5clw -n 5goperator-system -c manager

# Undeploy the operator from a Kubernetes cluster
make undeploy
```

## OLM (untested)

```
make bundle bundle-build bundle-push

operator-sdk olm install
operator-sdk olm status
operator-sdk run bundle 9.20.198.112:5000/5goperator-bundle:v0.0.1
```

