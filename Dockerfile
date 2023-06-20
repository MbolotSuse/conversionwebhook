# Build the manager binary
FROM registry.suse.com/bci/golang:1.19 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY config.go config.go
COPY controllers.go controllers.go
COPY api/ api/

# Build
# the GOARCH has not a default value to allow the binary be built according to the host where the command
# was called. For example, if we call make docker-build in a local env which has the Apple Silicon M1 SO
# the docker BUILDPLATFORM arg will be linux/arm64 when for Apple x86 it will be linux/amd64. Therefore,
# by leaving it empty we can ensure that the container and binary shipped on it will have the same platform.
RUN CGO_ENABLED=0 go build -a -o crd-conversion-webhook main.go config.go controllers.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM registry.suse.com/bci/bci-micro:latest

WORKDIR /
COPY --from=builder /workspace/crd-conversion-webhook .
RUN chmod +x crd-conversion-webhook
USER 65532:65532

CMD crd-conversion-webhook
