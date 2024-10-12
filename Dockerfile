# Build stage
FROM golang:1.23.2-alpine3.20 AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main main.go

# Final stage
FROM alpine:3.20

RUN adduser -D ash
USER ash
WORKDIR /home/ash

COPY --chown=ash:ash --from=builder /build .
COPY --chown=ash:ash --from=builder /build/migrations ./migrations
COPY --chown=ash:ash --from=builder /build/config.yaml .

EXPOSE 8080

ENTRYPOINT ["/home/ash/main"]
CMD []
