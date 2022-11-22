FROM alpine:3.17.0@sha256:839e8c0d8c70b6587dd546b3a92357da4db3c36c59285a409826a569c3c58994
COPY prometheus-hetrixtools-exporter /usr/bin/prometheus-hetrixtools-exporter
ENTRYPOINT ["/usr/bin/prometheus-hetrixtools-exporter"]
