FROM alpine:3.15.4@sha256:a777c9c66ba177ccfea23f2a216ff6721e78a662cd17019488c417135299cd89
COPY prometheus-hetrixtools-exporter /usr/bin/prometheus-hetrixtools-exporter
ENTRYPOINT ["/usr/bin/prometheus-hetrixtools-exporter"]
