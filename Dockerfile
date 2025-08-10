FROM alpine:latest as os

RUN apk update
RUN apk add curl git gcc musl-dev bash

FROM golang:1.23.5-alpine3.21 as build
RUN apk add --update bash

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o spe-transaction-gateway-api main.go

COPY . /app

FROM os
WORKDIR /app
RUN apk add --no-cache tzdata
ENV TZ=Asia/Jakarta

COPY --from=build /app/spe-transaction-gateway-api .

EXPOSE 8181

CMD ["./spe-transaction-gateway-api"]