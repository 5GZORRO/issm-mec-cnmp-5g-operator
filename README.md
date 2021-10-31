# 5GOperator

This is a Kubernetes Operator for installing components of a [5G Core Network](https://www.free5gc.org/) based on work
[here](https://github.ibm.com/Brian-P-Naughton/free5g). It is intended to be installed with and managed by
[CP4NA](https://www.ibm.com/docs/en/cloud-paks/cp-network-auto/2.1) and a companion Zero Touch Orchestration system (to be
implemented).

It demonstrates an approach to CP4NA integration with Kubernetes Operators. The only part that is necessary for it to work
in the CP4NA/Zero Touch ecosystem is the `generation` and `observedGeneration` [handling](api/v1alpha1/status.go), as outlined
in the Zero Touch Reference Architecture document. In particular, the vendor is free to implement whatever Conditions they see fit;
the Zero Touch Orchestrator will provide mechanisms yet to be worked out on how to apply policy to changes in Conditions.

# Docs

- [Architecture](docs/architecture.md)
- [Installation Guide](docs/installation.md)
- [Developer Guide](docs/development.md)

# Limitations and Further Work

* The 5G operator is not production ready code; there are no tests, the code has not been reviewed, it has not had a lot of testing, 
  it could do with some refactoring (to extract common reconciliation code relating to the 5G components), etc
* The handling of operations e.g. AddNode for SMF may not ensure continued availability of the component (SMF) because the operator may restart the pod. CP4NA should define replicas for the resource to handle this.
* A lot of [boilerplate](cp4na/xnfs/5G/amf/Lifecycle/kubernetes) is currently required in resource packages - need to update kubedriver to handle a lot of this
* It's probable that we can extract most of the Transition code (and perhaps other parts) in to an SDK so that the vendor
  has to implement only the actual transition or operation itself (a lot of the Transition handling code would be transferable
  to other implementations)
