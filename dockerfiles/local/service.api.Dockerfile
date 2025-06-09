# Start from the official Go image as a build stage
FROM golang:latest AS builder

# Set necessary environment variables for Go
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go binary
RUN go build -o api .

# Final stage: a small image with just the binary
FROM alpine:latest

# Set working directory inside container
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/api .

EXPOSE 4000

CMD ["./api"]
