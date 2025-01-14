FROM golang:1.23.4-alpine AS builder

WORKDIR /app

COPY /src/go.mod /src/go.sum ./

RUN go mod download

COPY . .

RUN go build -C src -o /bin/app


FROM alpine:latest

WORKDIR /app

COPY --from=builder /bin/app .

CMD ["./app"]

