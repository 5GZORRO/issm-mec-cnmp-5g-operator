# Sample CRs

This folder includes sample CRs to create and test the 5G operator

## Start core

```
./startCore.sh
```

## Register UE

TBD

## Configure SMF with UPFs and topologies

Run the below scripts in this order

### Create UPFs and add them to SMF model

This script creates the UPFs in the same namespace as the core (`5g-test`)

```
./upfs/addUpf.sh 5g-test http://<k8s master node ipaddress>:<smf-api nodeport>
```

### Define UPF topology

```
./upfs/addTopology.sh http://<k8s master node ipaddress>:<smf-ext nodeport>
```
