FROM golang:1.14.4-stretch AS builder

LABEL maintainer="Free5GC <support@free5gc.org>"

RUN curl -sL https://deb.nodesource.com/setup_lts.x | bash -
RUN apt-get update
RUN apt remove cmdtest yarn nodejs --purge
RUN apt-get -y install gcc cmake autoconf libtool pkg-config libmnl-dev libyaml-dev apt-transport-https
RUN apt-get clean

RUN curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg > pubkey.gpg && \
    apt-key add pubkey.gpg && \
    echo "deb https://dl.yarnpkg.com/debian/ stable main" >> /etc/apt/sources.list.d/yarn.list && \
    apt-get update && \
    apt-get install -y nodejs yarn

# Get Free5GC
RUN cd $GOPATH/src \
    && git clone --recursive -b v3.0.5 -j `nproc` https://github.com/free5gc/free5gc.git 

RUN cd $GOPATH/src/free5gc \
    && make
RUN cd $GOPATH/src/free5gc/webconsole && \
    git checkout v1.0.1 && \
    cd $GOPATH/src/free5gc && \
    make webconsole

FROM alpine

WORKDIR /free5gc
RUN mkdir -p config/ support/TLS/ public

# Copy executables
COPY --from=builder /go/src/free5gc/bin/* ./
COPY --from=builder /go/src/free5gc/NFs/upf/build/bin/* ./
COPY --from=builder /go/src/free5gc/webconsole/bin/* ./
COPY --from=builder /go/src/free5gc/webconsole/public ./public

# Copy linked libs
COPY --from=builder /go/src/free5gc/NFs/upf/build/updk/src/third_party/libgtp5gnl/lib/libgtp5gnl.so.0 ./
COPY --from=builder /go/src/free5gc/NFs/upf/build/utlt_logger/liblogger.so ./

# Copy configuration files (not used for now)
COPY --from=builder /go/src/free5gc/config/* ./config/
COPY --from=builder /go/src/free5gc/NFs/upf/build/config/* ./config/

# Copy default certificates (not used for now)
COPY --from=builder /go/src/free5gc/support/TLS/* ./support/TLS/
