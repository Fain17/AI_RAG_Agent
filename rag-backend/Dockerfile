# ---------- Build Stage ----------
    FROM golang:1.24-alpine AS builder

    # Enable Go modules and setup build env
    ENV CGO_ENABLED=0 \
        GO111MODULE=on \
        GOOS=linux \
        GOARCH=amd64
    
    # Accept secret DB connection string at build time
    ARG DATABASE_URL
    
    # Install git for go mod
    RUN apk add --no-cache git
    
    # Set working directory
    WORKDIR /app
    
    # Copy go mod and sum files
    COPY go.mod go.sum ./
    
    # Download dependencies
    RUN go mod download
    
    # Copy the rest of the app
    COPY . .
    
    # Optional: Print or use the DATABASE_URL securely (e.g., for codegen)
    # Example: RUN echo "Using DB at: $DATABASE_URL"
    
    # Build the Go binary
    RUN go build -o server .
    
    # ---------- Final Image ----------
    FROM alpine:latest
    
    WORKDIR /app
    
    # Copy the binary from builder stage
    COPY --from=builder /app/server .
    
    # Expose the API port (update as needed)
    EXPOSE 8080
    
    # Run the server
    CMD ["./server"]
    