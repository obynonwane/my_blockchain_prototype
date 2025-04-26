# Use Go 1.24 on Alpine Linux
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install make and git
RUN apk add --no-cache git make

# Copy go.mod, go.sum, and vendor folder first (better layer caching)
COPY go.mod go.sum vendor/ ./

# Copy the rest of the source code
COPY . .

# Build the binary using the vendor folder
RUN go build -mod=vendor -o nodeApp ./cmd

# Default command to run the app
CMD ["/app/nodeApp"]