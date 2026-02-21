FROM golang:1.24-alpine

RUN apk add --no-cache protobuf protobuf-dev nodejs npm

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

RUN npm install -g @protobuf-ts/plugin

WORKDIR /workspace
