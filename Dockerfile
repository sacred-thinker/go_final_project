FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

RUN apt-get update && apt-get install -y \
    libsqlite3-dev \
    && rm -rf /var/lib/apt/lists/*

COPY . .

RUN go build -o /app/my_app -a -ldflags '-extldflags "-static"' ./cmd/server/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/my_app .

COPY ./web /app/web

EXPOSE $TODO_PORT

CMD ["/app/my_app"]
