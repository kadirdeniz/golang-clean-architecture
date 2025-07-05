# syntax=docker/dockerfile:1

# --- Build Stage ---
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install git for go mod and build tools
RUN apk add --no-cache git

# Only copy go.mod and go.sum first for better cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go app (static binary, disables CGO)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server ./cmd/server

# --- Final Stage ---
FROM alpine:3.19

WORKDIR /app

# Add a non-root user for security
RUN adduser -D -g '' appuser

# Copy the built binary from builder
COPY --from=builder /app/server .

# Copy configs, migrations, etc.
COPY configs ./configs
COPY migrations ./migrations

# Expose the app port
EXPOSE 8080

USER appuser

# Start the app
ENTRYPOINT ["./server"] 