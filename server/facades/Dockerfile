# Build stage
FROM golang:1.25-alpine AS builder
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
FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=builder /out/singlepageapp /usr/local/bin/singlepageapp
WORKDIR /app

# default args (can be overridden at docker run or compose)
ENTRYPOINT ["/usr/local/bin/singlepageapp"]
CMD ["-port", ":8080", "-static", "/app/static", "-index", "index.html"]
