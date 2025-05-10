FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git
RUN apk add --no-cache make
RUN apk add --no-cache bash

WORKDIR /app

##COPY go.mod go.sum ./
##RUN go mod download

COPY . .
RUN make mod-tidy
RUN make build-auth

FROM alpine:3.18

COPY --from=builder /app/build/auth /app/
WORKDIR /app

# Запуск приложения
ENTRYPOINT ["/app/auth"]
