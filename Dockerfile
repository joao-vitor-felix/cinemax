FROM golang:1.26-alpine AS builder
WORKDIR /app
# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build a static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /app/server cmd/api/main.go

# Use scratch for the smallest possible image
FROM scratch

# Copy CA certificates for HTTPS calls
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the binary
COPY --from=builder /app/server /server

EXPOSE 8080

ENTRYPOINT ["/server"]
