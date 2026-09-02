# ===== builder =====
FROM golang:1.21-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o rooms-api ./cmd/server

# ===== runtime =====
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/rooms-api /app/rooms-api
EXPOSE 8080
ENTRYPOINT ["/app/rooms-api"]
