# Build

FROM golang:1.18-alpine AS build

ENV GO111MODULE=auto

WORKDIR /order_service

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build


# Deploy

FROM alpine:3.8

WORKDIR /root/

COPY --from=build /order_service/order_service .

EXPOSE 9090

ENTRYPOINT ["./order_service"] 