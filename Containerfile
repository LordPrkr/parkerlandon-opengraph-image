# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install templ CLI
RUN go install github.com/a-h/templ/cmd/templ@latest

# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate templ templates and build binary
RUN templ generate && go build -o /app/server ./cmd/server

# Runtime stage
FROM alpine:latest

# Install Chromium and dependencies for rod
RUN apk add --no-cache \
    chromium \
    font-noto \
    font-noto-emoji

# Set Chrome path for rod
ENV ROD_BROWSER=/usr/bin/chromium

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server /app/server

# Copy static assets and configs
COPY --from=builder /app/web /app/web
COPY --from=builder /app/configs /app/configs

EXPOSE 3000

CMD ["/app/server"]
