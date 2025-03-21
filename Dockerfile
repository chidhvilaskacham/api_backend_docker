# Build stage
FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Ensure the binary is statically compiled
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api .

# Verify the binary was created in the build stage
RUN ls -l /app/api

# Use a minimal runtime image
FROM alpine:latest

WORKDIR /app

# Install necessary OS dependencies
RUN apk update && apk add --no-cache ca-certificates

# Copy the binary
COPY --from=builder /app/api /app/api

COPY .env .env

# Make the binary executable
RUN chmod +x /app/api

# Verify the binary exists in the runtime stage
RUN ls -l /app/api

# Set environment variables
ENV PORT=8080
ENV CORS_ORIGINS=http://localhost:3000

# Expose port 8080
EXPOSE 8080

# Run the application
CMD ["./api"]