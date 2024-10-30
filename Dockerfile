FROM golang:latest

RUN go version
ENV GOPATH=/

COPY go.mod go.sum ./
RUN go mod tidy && go mod download
COPY . .
