# Build stage
FROM golang:1.24.0-alpine3.21 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main main.go

# Final stage
FROM alpine:3.21
WORKDIR /app

COPY --from=builder /app .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/config.yaml .

EXPOSE 8080

ENTRYPOINT ["/app/main"]
CMD []
