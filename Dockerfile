# Build the manager binary
FROM registry.redhat.io/rhel8/go-toolset as builder

WORKDIR /workspace

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

USER root

RUN chown -R 1001:0 /workspace

USER 1001

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY *.go .

USER root

RUN chown -R 1001:0 /workspace

USER 1001

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o manager .

FROM registry.access.redhat.com/ubi8-micro 
WORKDIR /
COPY --from=builder /workspace/manager .
USER 65532:65532

ENTRYPOINT ["/manager"]
