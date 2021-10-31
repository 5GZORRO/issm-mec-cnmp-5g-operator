FROM free5gc/base:latest AS builder
FROM bitnami/minideb:stretch

LABEL description="Free5GC open source 5G Core Network" \
    version="Stage 3"

ENV F5GC_MODULE webconsole
ARG DEBUG_TOOLS

RUN apt-get update && apt-get install -y python3 procps

# Set working dir
WORKDIR /free5gc
RUN mkdir -p config/ webconsole/public

# Copy executable, frontend static files and default configuration
COPY --from=builder /free5gc/${F5GC_MODULE} ./webconsole/webconsole
COPY --from=builder /free5gc/public ./webconsole/public

# Move to the executable location
WORKDIR /free5gc/webconsole

# Config files volume
VOLUME [ "/free5gc/config" ]

# WebUI uses the port 5000
EXPOSE 5000
