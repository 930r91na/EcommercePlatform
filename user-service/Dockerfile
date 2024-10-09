# Stage 1: Build the Go binary
FROM golang:1.21 AS builder

WORKDIR /app

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . ./

# Disable CGO and build the Go application binary for Linux AMD64
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o user-service cmd/main.go

# Stage 2: Create a minimal image with the compiled binary
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/user-service .

# Ensure the binary has executable permissions
RUN chmod +x user-service

# Copy the environment variables, removed relying in env variables
# COPY .env .env

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./user-service"]
