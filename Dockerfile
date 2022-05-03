FROM golang:1.18

ENV ROOT=/go/src/app
WORKDIR ${ROOT}

ENV GO111MODULE=on

RUN apk update && apk add git
COPY go.mod go.sum ./
RUN go mod download
EXPOSE 8080
