FROM free5gc/base:latest AS builder
FROM ubuntu:20.04

LABEL description="Free5GC open source 5G Core Network" \
    version="Stage 3"

ENV F5GC_MODULE free5gc-upfd
ENV DEBIAN_FRONTEND noninteractive
ARG DEBUG_TOOLS

# Install debug tools ~ 100MB (if DEBUG_TOOLS is set to true)
RUN if [ "$DEBUG_TOOLS" = "true" ] ; then apt-get update && apt-get install -y tcpdump vim strace iputils-ping curl netcat ; fi

# Install UPF dependencies
RUN apt-get clean \
    && apt-get update \
    && apt-get install -y libmnl0 libyaml-0-2 iproute2 iptables net-tools wget python \
    && apt-get clean

# install prometheus node collector
RUN cd /opt && \
    wget https://github.com/prometheus/node_exporter/releases/download/v1.0.1/node_exporter-1.0.1.linux-amd64.tar.gz && \
    tar xvfz node_exporter-*.*-amd64.tar.gz

# Set working dir
WORKDIR /free5gc
RUN mkdir -p ${F5GC_MODULE}/config/

# Copy executable and default configuration
COPY --from=builder /free5gc/${F5GC_MODULE} ./${F5GC_MODULE}/

# Copy linked libs
COPY --from=builder /free5gc/libgtp5gnl.so.0 /usr/local/lib
COPY --from=builder /free5gc/liblogger.so /usr/local/lib

# Move to the executable location
WORKDIR /free5gc/${F5GC_MODULE}

# Update links
RUN ldconfig
