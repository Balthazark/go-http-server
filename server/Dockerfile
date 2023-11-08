# syntax=docker/dockerfile:1

FROM golang:1.21.3

WORKDIR /app
COPY go.mod ./
RUN go mod download

COPY server.go ./
COPY content ./content

RUN CGO_ENABLED=0 GOOS=linux go build -o /server

EXPOSE 8080

ENTRYPOINT ["/server"]
