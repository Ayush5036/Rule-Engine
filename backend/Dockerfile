# Use the official Golang image as a builder
FROM golang:1.20 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files first to leverage cached layers
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the entire backend source code into the container
COPY . .

# Build the Go app statically linked
RUN CGO_ENABLED=0 GOOS=linux go build -a -o main ./cmd/main.go

# Start a new stage from scratch
FROM alpine:latest

# Install ca-certificates
RUN apk --no-cache add ca-certificates

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
