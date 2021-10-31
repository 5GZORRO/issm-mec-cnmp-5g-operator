FROM golang:1.15.6-buster 

RUN apt-get update && apt-get install -y tcpdump vim strace net-tools curl netcat-openbsd iptables python python-pip python-setuptools python3 python3-pip iperf

RUN cd /go/src/ && \
    git clone https://github.com/brianpnaughton/gnbsim.git &&\
    cd gnbsim &&\
    make

# install prometheus node collector
RUN cd /opt && \
    wget https://github.com/prometheus/node_exporter/releases/download/v1.0.1/node_exporter-1.0.1.linux-amd64.tar.gz && \
    tar xvfz node_exporter-*.*-amd64.tar.gz

WORKDIR /go/src/gnbsim/gnb-gw