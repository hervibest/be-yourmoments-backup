# Stage 1: Builder
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git build-base

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

# Go module
COPY user-svc/go.mod user-svc/go.sum ./
COPY pb ./pb
RUN go mod download

# Copy full source
COPY user-svc/. .

# Build binary
RUN go build -o user-svc ./cmd/web/main.go

# Stage 2: Runtime (distroless)
FROM gcr.io/distroless/static

WORKDIR /

# Copy binary and migration files only
COPY --from=builder /app/user-svc .
COPY --from=builder /app/db/migrations ./db/migrations

ENTRYPOINT ["/user-svc"]
