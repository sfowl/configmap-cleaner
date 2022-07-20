# Build the manager binary
FROM registry.redhat.io/rhel8/go-toolset as builder

WORKDIR /workspace

# Copy the Go Modules manifests
# 1001 is the default user in the go builder image
COPY --chown=1001:0 go.mod go.mod
COPY --chown=1001:0 go.sum go.sum

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY --chown=1001:0 *.go .

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o manager .

FROM registry.access.redhat.com/ubi8-micro 
WORKDIR /
COPY --from=builder /workspace/manager .
USER 65532:65532

ENTRYPOINT ["/manager"]
