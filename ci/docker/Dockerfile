# Stage 1: Build the Go application
FROM golang:1.22.7-alpine3.20 AS builder

# Arguments for build
ARG VERSION
ARG GOARCH
ARG GOOS
ARG BUILDFLAGS="-mod=vendor"
ARG LDFLAGS="-X main.Version=${VERSION}"

# Set working directory
WORKDIR /app

# Install dependencies
RUN apk update && apk add --no-cache git make build-base

# Copy Go modules and vendor
COPY backend/go.mod backend/go.sum ./
COPY backend/vendor ./vendor

# Copy application source
COPY backend ./

# Build the application
RUN GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=0 \
    go build -o /out/service -ldflags="${LDFLAGS}" ${BUILDFLAGS} ./internal

# Stage 2: Create a lightweight image
FROM alpine:3.18

# Create runtime directories
RUN mkdir -p /var/service/logs /var/2112/config /var/2112/data

# Copy binary from builder stage
COPY --from=builder /out/service /usr/local/bin/service

# Expose ports
EXPOSE 8081 8080 8079

# Set default command (placeholder or shell)
CMD ["/bin/sh"]
