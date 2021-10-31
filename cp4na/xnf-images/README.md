# CNF images

To build the CNF docker images you must have a working docker environment. Run the following command from a console in the __free5g/install/xnf-images__ directory. 

```
make
```

## free5gc test

A sample docker-compose and free5gc set of configuration files can be used to test the free5g core and gnb-gw.

```
Vagrant up
vagrant/vagrant
```

open a console on the VM and type the following to start the 5g core

```
cd /vagrant
make
docker-compose up
```

Open firefox in the VM at http://10.100.200.30:5000
Login with __admin/free5gc__
Click on the Subscribers tab and create the default UE subscriber

Run gnb-gw with the following commands in a console on the Vagrant machine.

```
docker exec -ti gnb-gw1 bash
./gnb-gw
```

on another terminal add a route to the IT web server in the gnb-gw contaner. 

```
docker exec -ti gnb-gw1 bash
route add --net 172.168.14.0/24 dev gtp-gnb
```

on the upf1 container configure the iptables NAT rules for external traffic to the web server

```
docker exec -ti upf1 bash
iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE
```

You should now be able to access the nginx web server from the gnb pod, routing through the upf gtp tunnel path. Run the following commands in a console:

```
docker exec -ti gnb-gw1 bash
wget 172.168.14.10
```
