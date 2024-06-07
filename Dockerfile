FROM golang:latest

RUN go version
ENV GOPATH=/

COPY ./ ./
RUN go mod download

EXPOSE 8080