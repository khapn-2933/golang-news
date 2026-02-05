# Build stage
# Sử dụng golang:latest để đảm bảo có Go version mới nhất (>= 1.24)
FROM golang:latest AS builder

# Set working directory
WORKDIR /app

# Install git (cần cho một số dependencies)
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o news main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/news .

# Expose port
EXPOSE 8080

# Run application
CMD ["./news"]

