FROM golang:1.24-alpine

RUN apk add --no-cache protobuf protobuf-dev nodejs npm make

WORKDIR /workspace
