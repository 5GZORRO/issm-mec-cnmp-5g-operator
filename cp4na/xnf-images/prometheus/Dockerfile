FROM alpine:latest

ARG ARCH=amd64
ARG VERSION=1.7.1

RUN apk --no-cache add --virtual build-dependencies wget ca-certificates \
    && mkdir -p /tmp/install /tmp/dist \
    && wget -O /tmp/install/prometheus.tar.gz https://github.com/prometheus/prometheus/releases/download/v$VERSION/prometheus-$VERSION.linux-$ARCH.tar.gz \
    && apk del build-dependencies \
    && cd /tmp/install \
    && tar --strip-components=1 -xzf prometheus.tar.gz \
    && mkdir -p /etc/prometheus /usr/share/prometheus \
    && mv prometheus promtool   /bin/ \
    && mv prometheus.yml        /etc/prometheus/prometheus.yml \
    && mv consoles console_libraries NOTICE LICENSE /usr/share/prometheus/ \
    && ln -s /usr/share/prometheus/console_libraries /usr/share/prometheus/consoles/ /etc/prometheus/ \
    && rm -rf /tmp/install

RUN apk add python2

VOLUME     [ "/prometheus" ]
WORKDIR    /prometheus
ENTRYPOINT [ "/bin/prometheus" ]




