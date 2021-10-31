# Network Function Automation Packages

This directory contains CP4NA packages for installing the 5G network using the 5G Operator and CP4NA. 

Note that the structure of the Kubernetes directories is similar; there is a lot of _common_ boilerplate here that can be
handled by the Kubernetes driver in a future release. For example, the `check-ready` scripts are all very similar due
to the standardised CR status handling implemented in the 5G Operator. Also, the `kubernetes/objects` yaml files defining
the CR to create are similar - there is a lot of common CR boilerplate that can be handled by the Kubernetes driver.