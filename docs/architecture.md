# Architecture

Each 5G Core component is represented by a CRD and related controller/reconciler, responsible for creating and removing
component CRs. CP4NA transitions (other than Create) and operations are represented as a [generic CRD](api/v1alpha1/transition_types.go);
the transition controller delegates the implementation of the transition or operation to the specific component reconciler (see,
[for example, the Start and Stop implementations here](pkg/reconcilers/amf/amf.go)). For the 5G components in this example,
the implementation of the transitions and operations is in the form of shell scripts that are run by the Operator using the
[Kubernetes exec api](https://www.openshift.com/blog/executing-commands-in-pods-using-k8s-api). This is an implementation detail;
it could be accomplished using any one of a number of other approaches - for example by implementing them as Ansible scripts
and running them via [Ansible Tower](https://www.ansible.com/products/tower). In addition, we have gone for a generic transition CRD
model here; alternate approaches include:

* representing each 5G component transition and operation as independent CRDs (with their own controllers).
* supporting an `intendedState` on each 5G core component CRD; the controller would then need to determine how to achieve that state.
  The disadvantage here is the difficulty mapping from CP4NA transitions to intended states; in fact, this may not be possible
  for CP4NA operations.

## CP4NA Transition and Operation Implementation

Transitions and operations are implemented as Go methods, which are called when a Transition CR is instantiated. See [Smf](pkg/reconcilers/amf/smf.go)
for example, which implements the following transition and operations: 

* Create
* Start
* Stop
* Attachnode
* Addueroute

See lines 34 - 37:

```
  smf.EntityHandlerImpl.AddTransitionFunction("start", smf.Start)
  smf.EntityHandlerImpl.AddTransitionFunction("stop", smf.Stop)
  smf.EntityHandlerImpl.AddTransitionFunction("attachnode", smf.AddNode)
  smf.EntityHandlerImpl.AddTransitionFunction("addueroute", smf.AddUERoute)
```

## CR Status Handling

Each controller implements standardised `observedGeneration` handling per the Zero Touch Reference Architecture. In particular, the
`observedGeneration` is updated whenever a reconciliation is started to indicate that the controller has noticed any changes
in the spec. Subsequent changes to the status Conditions (either triggered by changes in the spec, or by 'unmanaged' changes)
will be monitored by the Zero Touch Watcher, and events generated that will trigger CP4NA policies only if the `observedGeneration`
is greater than or equal to the `generation`.

The choice of which Conditions to support is up to the vendor - CP4NA does not impose any restrictions here. CP4NA will
support the mapping of condition changes to CP4NA policy using a to be determined mechanism.

In this specific example, the supported status `conditions` are `Ready` and `Reconciling`.

The component is ready:

```
  ...
  metadata:
    generation: 2
  ...
  status:
    ...
    conditions:
    - lastTransitionTime: "2021-06-10T15:14:53Z"
      message: ""
      reason: ""
      status: "False"
      type: Reconciling
    - lastTransitionTime: "2021-06-10T15:15:53Z"
      message: ""
      reason: ready
      status: "True"
      type: Ready
    observedGeneration: 2
    outputs:
      ip_address: 10.0.0.204
      podname: core2-mongodb-afa8ed1a-73c3-4e21-b9a2-549d029cffce
```

An error has occurred:

```
  ...
  metadata:
    generation: 2
  ...
  status:
    ...
    conditions:
    - lastTransitionTime: "2021-06-10T15:14:53Z"
      message: ""
      reason: ""
      status: "False"
      type: Reconciling
    - lastTransitionTime: "2021-06-10T15:14:53Z"
      message: ""
      reason: "Exception has occurred"
      status: "False"
      type: Ready
    observedGeneration: 2
    outputs:
      ip_address: 10.0.0.204
      podname: core2-mongodb-afa8ed1a-73c3-4e21-b9a2-549d029cffce
```

The component is not ready if either:

* it is still being reconciled with the CR spec
* an error has occurred either because the reconciler cannot reconcile with the CR spec or an unmanaged action has occurred
  causing the component to fail. In this case, the message is populated with an appropriate error message. The current implementation
  regards any error as fatal and will stop reconciliation. A more intelligent implementation could better handle transient errors by
  retrying the reconciliation.
  
Note also the outputs which are mapped in to CP4NA (read only) properties of the resource instance that _owns_ i.e. instantiated the CR.