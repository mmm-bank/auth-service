FROM golang:latest AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o auth .

FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /build/auth .

ENTRYPOINT ["/app/auth"]
