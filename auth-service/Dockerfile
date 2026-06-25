FROM golang:1.25-bookworm AS builder

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o app main.go

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/app .

CMD ["./app"]