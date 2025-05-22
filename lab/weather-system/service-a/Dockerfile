FROM golang:1.24 AS build

WORKDIR /app

# Copy module files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build static optimized binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o service-a ./cmd/server/main.go

# Use minimal image for final build
FROM scratch

# Copy CA certificates for HTTPS calls
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy binary
COPY --from=build /app/service-a /app/service-a

# Set entrypoint
ENTRYPOINT ["/app/service-a"]