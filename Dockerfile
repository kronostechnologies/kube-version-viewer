FROM --platform=$BUILDPLATFORM golang:1.22 AS builder
RUN echo "nobody:x:65534:65534:nobody:/:" > /tmp/passwd

WORKDIR /go/src/github.com/kronostechnologies/kube-version-viewer/
ARG TARGETOS TARGETARCH
ADD . .
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 go build -ldflags="-w -s" -o version-viewer .

FROM scratch
COPY --from=builder /go/src/github.com/kronostechnologies/kube-version-viewer/version-viewer /bin/
COPY --from=builder /tmp/passwd /etc/passwd

USER 65534:65534

ENTRYPOINT ["/bin/version-viewer"]
CMD ["--serve"]
