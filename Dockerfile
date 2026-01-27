# Stage 1: Build
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go.mod and download dependencies (if any)
COPY go.mod ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o category-api main.go

# Stage 2: Run
FROM alpine:latest

WORKDIR /root/

# Copy binary from builder stage
COPY --from=builder /app/category-api .

# Expose port 8080
EXPOSE 8080

# Run the application
CMD ["./category-api"]
