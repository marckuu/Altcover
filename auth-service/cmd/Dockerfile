FROM golang:1.25-bookworm AS builder

WORKDIR /app

COPY . .

RUN go mod tidy

RUN ["go", "build", "-o", "app", "./cmd"]
RUN ["go", "build", "-o", "createAdmin", "./cmd/admin"]

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/app .
COPY --from=builder /app/createAdmin .

CMD ["./app"]