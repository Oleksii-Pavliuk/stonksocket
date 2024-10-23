FROM golang:1.23.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o socket .


FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/socket .

EXPOSE 8080

CMD ["./socket"]