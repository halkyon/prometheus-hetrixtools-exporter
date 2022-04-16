FROM alpine
COPY prometheus-hetrixtools-exporter /usr/bin/prometheus-hetrixtools-exporter
ENTRYPOINT ["/usr/bin/prometheus-hetrixtools-exporter"]
