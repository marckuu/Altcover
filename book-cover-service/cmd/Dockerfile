FROM golang:latest AS builder

WORKDIR /app

COPY . .

RUN go mod tidy

RUN ["go", "build", "-o", "app", "./cmd"]

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/app .

CMD ["./app"]