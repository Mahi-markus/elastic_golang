# Step 1: Build the Go app
FROM golang:1.20-alpine AS builder

# Set the working directory in the container
WORKDIR /app

# Copy the Go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod tidy

# Copy the rest of the Go application code
COPY . .

# Build the Go app (adjust the binary name as needed)
RUN go build -o app .

# Step 2: Set up the runtime container
FROM alpine:latest

# Install necessary packages
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /root/

# Copy the built binary from the builder container
COPY --from=builder /app/app .

# Expose the port the app will run on
EXPOSE 8080

# Command to run the app
CMD ["./app"]
