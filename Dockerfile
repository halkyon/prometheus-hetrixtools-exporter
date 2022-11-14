FROM alpine:3.16.3@sha256:b95359c2505145f16c6aa384f9cc74eeff78eb36d308ca4fd902eeeb0a0b161b
COPY prometheus-hetrixtools-exporter /usr/bin/prometheus-hetrixtools-exporter
ENTRYPOINT ["/usr/bin/prometheus-hetrixtools-exporter"]
