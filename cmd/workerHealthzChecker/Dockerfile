# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o worker ./cmd/workerHealthzChecker

# Final stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache sqlite

# Copy the binary from builder
COPY --from=builder /app/worker .
COPY ./config ./config

# Create directory for database
RUN mkdir -p /data

# Expose the port
EXPOSE 8082
ENV CONFIG_PATH=/PLACE_HOLDER

# Run the application
CMD ./worker --config "${CONFIG_PATH}"