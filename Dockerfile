FROM golang:1.22-bookworm AS builder

WORKDIR /app

# Copy go.mod and source code
COPY . .

# Initialize go.sum if it's missing (helps in case local machine didn't have go installed)
RUN go mod tidy

# Enable CGO for sqlite3 and build
ENV CGO_ENABLED=1
RUN go build -o ticket-system ./cmd/api

# Final image
FROM debian:bookworm-slim
WORKDIR /app

# Install ca-certificates in case we need to make HTTPS calls
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/ticket-system .

EXPOSE 8080
CMD ["./ticket-system"]
