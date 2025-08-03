# Multi-stage build for reciept-invoice-ai-tool
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o reciept-invoice-ai-tool main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Create a non-root user
RUN addgroup -g 1000 appgroup && \
    adduser -u 1000 -G appgroup -s /bin/sh -D appuser

# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/reciept-invoice-ai-tool .

# Create data directory and set ownership
RUN mkdir -p /app/data && chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Set entrypoint
ENTRYPOINT ["./reciept-invoice-ai-tool"]