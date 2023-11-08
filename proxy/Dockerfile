# syntax=docker/dockerfile:1

FROM golang:1.21.3

WORKDIR /app
COPY go.mod ./
RUN go mod download

COPY proxy.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /server

EXPOSE 8080

ENTRYPOINT ["/proxy"]
