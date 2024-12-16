FROM golang:alpine

ENV CGO_ENABLED=0

LABEL \
  org.opencontainers.image.title="integrated_exporter" \
  org.opencontainers.image.description="Integrated Exporter" \
  org.opencontainers.image.url="https://github.com/liushunking/integrated_exporter" \
  org.opencontainers.image.documentation="https://github.com/liushunking/integrated_exporter#readme" \
  org.opencontainers.image.source="https://github.com/liushunking/integrated_exporter" \
  org.opencontainers.image.licenses="Apache-2.0 license" \
  maintainer="liushun <liushun0311@zju.edu.cn>"

WORKDIR /app

COPY dist/integrated_exporter_linux_amd64_v1/integrated_exporter /dist/integrated_exporter_linux_amd64/integrated_exporter
COPY dist/integrated_exporter_linux_arm64_v8.0/integrated_exporter /dist/integrated_exporter_linux_arm64/integrated_exporter

RUN if [ `go env GOARCH` = "amd64" ]; then \
      cp /dist/integrated_exporter_linux_amd64/integrated_exporter ./integrated_exporter; \
    elif [ `go env GOARCH` = "arm64" ]; then \
      cp /dist/integrated_exporter_linux_arm64/integrated_exporter ./integrated_exporter; \
    fi

RUN apk update --no-cache \
  && apk add --no-cache tzdata ca-certificates \
  && rm -rf /dist \
  && rm -rf /go/pkg/mod \
  && rm -rf /go/pkg/sumdb

COPY etc/etc.yaml ./etc/etc.yaml
COPY etc/.env.yaml ./etc/.env.yaml

EXPOSE 6070

ENTRYPOINT ["./integrated_exporter",
    "--port=6070",
    "--interval=5s",
    "--route=/metrics",
    "--config=./etc/config.yaml",
    "--env=./etc/env.yaml",
    "--machineConfig.metrics=cpu,memory,disk,process,network",
    "--machineConfig.mounts=/"]