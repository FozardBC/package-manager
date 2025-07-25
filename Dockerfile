
FROM golang:1.24.2-alpine AS builder


WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download


RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o pm .