# Build stage
FROM golang:1.23.2-alpine3.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o xm main.go

# Final stage
FROM alpine:3.20

RUN adduser -D ash
USER ash
WORKDIR /home/ash

COPY --chown=ash:ash --from=builder /app/ .
COPY --chown=ash:ash --from=builder /app/migrations ./migrations
COPY --chown=ash:ash --from=builder /app/config.yaml .

EXPOSE 8080

ENTRYPOINT ["/home/ash/xm"]
CMD []
