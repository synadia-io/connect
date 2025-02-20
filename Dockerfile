FROM golang:1.24 AS build

ENV CGO_ENABLED=0
ENV GOOS=linux
RUN useradd -u 10001 connect
RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b /usr/local/bin

WORKDIR /go/src/github.com/synadia-io/connect/
# Update dependencies: On unchanged dependencies, cached layer will be reused
COPY .. /go/src/github.com/synadia-io/connect/
RUN go mod tidy

# Build
RUN task build TAGS="timetzdata"

# Pack
FROM busybox AS package

LABEL maintainer="Synadia <code@synadia.com>"
LABEL org.opencontainers.image.source="https://github.com/synadia-io/connect"

WORKDIR /

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /go/src/github.com/synadia-io/connect/target/connect .

RUN mkdir -p /home/connect/.synadia && chown -R connect /home/connect/

USER connect

ENTRYPOINT ["/connect"]
