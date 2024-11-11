# Start with a base image with Go and protoc installed
FROM golang:1.23-alpine AS builder

# Set the environment variables for cross-compilation
ENV GOARCH=amd64 
ENV GOOS=linux

# Install protoc and Go plugins for gRPC
RUN apk add --no-cache git make protobuf protobuf-dev
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Set the working directory
WORKDIR /app

# Copy the Go modules and install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .


# Build the application
RUN go build -o main .

# Use a smaller image to run the app
FROM alpine:3.17

# Copy the built application binary from the builder stage
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/proto /app/proto

# Expose the application port
EXPOSE 50052

# Run the application
CMD ["./main"]