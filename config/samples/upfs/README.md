
## Re-configure SMF with UPFs and topologies

Run the below in this order. Modify ips and ports accordingly

### Create UPFs and add them to SMF model

```
./addUpf.sh 5g-test http://<k8s master node ipaddress>:<smf-api nodeport>
```

### Define UPF topology

```
./addTopology.sh http://<k8s master node ipaddress>:<smf-ext nodeport>
```