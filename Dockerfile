# syntax=docker/dockerfile:1
FROM golang:1.23 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the Cobra CLI (binary name "app")
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

# Minimal final image
FROM gcr.io/distroless/base-debian12
WORKDIR /
COPY --from=builder /app/app /app

# Default command is "serve"
ENTRYPOINT ["/app", "serve"]
