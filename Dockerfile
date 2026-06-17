# Step 1: Build the Go application
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies if any
RUN apk add --no-cache git

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Compile the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/app/main.go

# Step 2: Run the binary in a clean minimal container
FROM alpine:3.19

WORKDIR /app

# Copy the compiled binary from builder stage
COPY --from=builder /app/main .


EXPOSE 8080

CMD ["./main"]
