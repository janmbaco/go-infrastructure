# Build stage
FROM golang:1.27.4-alpine AS builder
WORKDIR /src

# Copy go.mod and go.sum from repo root to enable module download
# We copy the whole repo to keep relative imports consistent
COPY . /src

# Set CGO disabled for static builds (alpine), adjust as needed
ENV CGO_ENABLED=0

# Build the singlepageapp cmd
RUN cd ./cmd/singlepageapp \
    && go build -o /out/singlepageapp .

# Runtime stage
FROM cgr.dev/chainguard/wolfi-base@sha256:be019c1c74fe69204a42c4e10bc4d886a49bf5403ffa74374ef7baa408fc9a26 AS runtime-packages
# Install signed runtime packages into a fresh root, retaining the APK inventory.
RUN mkdir -p /runtime/etc/apk \
    && cp -R /etc/apk/keys /etc/apk/repositories /runtime/etc/apk/ \
    && apk --root /runtime --initdb --no-cache add wolfi-baselayout busybox ca-certificates-bundle \
    && test ! -e /runtime/usr/lib/libcrypto.so.3 \
    && test ! -e /runtime/usr/lib/libssl.so.3

FROM scratch
COPY --from=runtime-packages /runtime/ /
COPY --from=builder /out/singlepageapp /usr/local/bin/singlepageapp
WORKDIR /app

# default args (can be overridden at docker run or compose)
ENTRYPOINT ["/usr/local/bin/singlepageapp"]
CMD ["-port", ":8080", "-static", "/app/static", "-index", "index.html"]
