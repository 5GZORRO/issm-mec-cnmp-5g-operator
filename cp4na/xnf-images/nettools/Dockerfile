FROM ubuntu:18.04

ENV DEBIAN_FRONTEND=noninteractive DEBCONF_NONINTERACTIVE_SEEN=true
ENV TZ=Europe/London
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
RUN apt-get update && apt-get install -y vim strace net-tools curl iperf collectd python python-pip python-setuptools python3 python3-pip iputils-ping
RUN pip install skydive-client 
RUN pip3 install lmctl
