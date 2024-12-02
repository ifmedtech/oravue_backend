# Stage 1: Build the Go application
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Stage 2: Run the application with a minimal image
FROM gcr.io/distroless/base-debian11

WORKDIR /

# Copy the built application
COPY --from=builder /app/main .

# Copy the configuration directory
COPY --from=builder /app/config ./config

USER nonroot:nonroot

ENTRYPOINT ["./main"]
