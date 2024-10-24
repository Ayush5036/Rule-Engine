# Dockerfile for Go backend

# Step 1: Build the Go application
FROM golang:1.20-alpine AS builder

# Set working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire Go project
COPY . .

# Build the Go application
RUN go build -o myapp cmd/main.go

# Step 2: Create a minimal image to run the Go application
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/myapp .

# Expose port 8080 for the backend
EXPOSE 8080

# Command to run the backend
CMD ["./myapp"]
