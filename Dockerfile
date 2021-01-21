FROM golang:1.15 AS builder
RUN apt update ; apt install upx-ucl -y ; apt clean
WORKDIR /go/src/github.com/kronostechnologies/kube-version-viewer/
COPY * ./
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o version-viewer . && upx --best version-viewer
RUN echo "nobody:x:65534:65534:nobody:/nonexistent:" > /tmp/passwd

FROM scratch
LABEL org.opencontainers.image.source=https://github.com/kronostechnologies/kube-version-viewer

COPY --from=builder /go/src/github.com/kronostechnologies/kube-version-viewer/version-viewer /bin/

USER nobody

ENTRYPOINT ["/bin/version-viewer"]
CMD ["--serve"]