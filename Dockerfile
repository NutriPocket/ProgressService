FROM golang:1.23.4-alpine AS builder

WORKDIR /app

COPY /src/go.mod /src/go.sum ./

RUN go mod download

COPY . .

RUN go build -C src -o /bin/app


FROM alpine:latest

WORKDIR /code

COPY --from=builder /bin/app /code/bin/app

COPY .env .

WORKDIR /code/bin

CMD ["/code/bin/app"]

