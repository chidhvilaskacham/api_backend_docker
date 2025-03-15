# Use the official Golang image to build the app
FROM golang:1.22 AS builder

# Set working directory inside container
WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Build the Go app
RUN go build -o api .

# Use a minimal Alpine image to run the app
FROM ubuntu:22.04
WORKDIR /app
# Copy the binary from the builder stage
COPY --from=builder /app/api .
EXPOSE 8080
# Expose port 8080 for HTTP traffic
RUN chmod +x /app/api

# Run the app
CMD ["./api"]