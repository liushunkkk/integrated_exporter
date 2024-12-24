# Use official Alpine image for a lightweight base
FROM alpine:latest

# Add image metadata
LABEL org.opencontainers.image.title="integrated_exporter" \
      org.opencontainers.image.description="Integrated Exporter" \
      org.opencontainers.image.url="https://github.com/liushunkkk/integrated_exporter" \
      org.opencontainers.image.documentation="https://github.com/liushunkkk/integrated_exporter#readme" \
      org.opencontainers.image.source="https://github.com/liushunkkk/integrated_exporter" \
      org.opencontainers.image.licenses="Apache-2.0 license" \
      maintainer="liushun <liushun0311@zju.edu.cn>"

# Set the working directory
WORKDIR /app

# Copy binaries for different architectures
COPY dist/integrated_exporter_linux_amd64_v1/integrated_exporter /dist/integrated_exporter_linux_amd64/integrated_exporter
COPY dist/integrated_exporter_linux_arm64_v8.0/integrated_exporter /dist/integrated_exporter_linux_arm64/integrated_exporter

# Select the appropriate binary based on the architecture
RUN ARCH=$(uname -m) && \
    if [ "$ARCH" = "x86_64" ] || [ "$ARCH" = "amd64" ]; then \
      cp /dist/integrated_exporter_linux_amd64/integrated_exporter ./integrated_exporter; \
    elif [ "$ARCH" = "arm64" ] || [ "$ARCH" = "aarch64" ]; then \
      cp /dist/integrated_exporter_linux_arm64/integrated_exporter ./integrated_exporter; \
    fi

# Install dependencies
RUN apk update --no-cache \
    && apk add --no-cache tzdata ca-certificates

# Copy configuration files
COPY etc/etc.yaml ./etc/etc.yaml
COPY etc/.env.yaml ./etc/.env.yaml

# Expose the application port
EXPOSE 6070

# Set the default command
ENTRYPOINT ["./integrated_exporter", \
            "server", \
            "--port=6070", \
            "--config=./etc/etc.yaml", \
            "--env=./etc/.env.yaml"]