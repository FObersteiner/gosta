# Build stage
FROM golang:1.24-alpine AS builder

# Install git for fetching dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gost .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy the binary and config from builder stage
COPY --from=builder /app/gost .
COPY --from=builder /app/config.yaml .

# Create non-root user for security
RUN addgroup -g 1001 -S gost && \
    adduser -u 1001 -S gost -G gost

# Change ownership of the app directory
RUN chown -R gost:gost /app

# Switch to non-root user
USER gost

# Expose port
EXPOSE 8080

# Run the application
ENTRYPOINT ["./gost"]
