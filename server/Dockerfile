# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o smart-fridge ./cmd/smart-fridge

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/smart-fridge .

# Copy configs
COPY --from=builder /app/configs ./configs

# Set default environment variables
ENV HTTP_PORT=8080
ENV NATS_PORT=4222
ENV OPENAI_API_KEY=""

# Expose default ports (can be overridden)
EXPOSE 8080 4222

# Run the application directly
CMD ["./smart-fridge"]
