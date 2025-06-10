# Start from the official Go image as a build stage
FROM golang:latest AS builder

# Set necessary environment variables for Go
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set the working directory inside the container to the module root
WORKDIR /app

# Copy go.mod and go.sum
COPY ./services/social-media-aggregator-api/go.mod ./services/social-media-aggregator-api/go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire module code
COPY ./services/social-media-aggregator-api ./

# Build
WORKDIR /app/app

RUN go build -o api .

# Final stage: a small image with just the binary
FROM alpine:latest

#install curl
RUN apk add --no-cache curl

WORKDIR /root/

COPY --from=builder /app/app/api .

EXPOSE 4000

ENTRYPOINT ["./api"]
CMD []