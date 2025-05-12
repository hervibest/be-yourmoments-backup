# Stage 1: Builder
FROM golang:1.24-alpine AS builder

# Install dependencies for bimg (libvips) + CGO build
RUN apk add --no-cache git build-base vips-dev

ENV CGO_ENABLED=1
ENV GOOS=linux

WORKDIR /app

# Copy module and dependencies
COPY upload-svc/go.mod upload-svc/go.sum ./
COPY pb ./pb
RUN go mod download

# Copy source
COPY upload-svc/. .

# Build binary
RUN go build -o upload-svc ./cmd/web/main.go && chmod +x upload-svc

# Stage 2: Runtime (Alpine, not distroless, to support libvips at runtime)
FROM alpine:3.19

# Install runtime-only dependencies (libvips for bimg)
RUN apk add --no-cache vips

WORKDIR /

# Copy binary from builder
COPY --from=builder /app/upload-svc .

# Ensure it's executable (in case permission lost during COPY)
RUN chmod +x /upload-svc

# Set entrypoint
ENTRYPOINT ["/upload-svc"]
