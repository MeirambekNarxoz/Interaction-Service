FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o interaction-service ./cmd/server

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/interaction-service .
COPY .env .env

EXPOSE 8084

CMD ["./interaction-service"]