# ---------- Builder ----------
FROM golang:1.24-bullseye AS builder

WORKDIR /app

RUN apt-get update && apt-get install -y gcc sqlite3 libsqlite3-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app ./cmd/api


# ---------- Runtime ----------
FROM debian:bullseye-slim

WORKDIR /app

RUN apt-get update && apt-get install -y sqlite3 libsqlite3-0 && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/app .

EXPOSE 8080

CMD ["./app"]
