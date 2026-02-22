FROM golang:1.24-alpine

RUN apk add --no-cache protobuf protobuf-dev nodejs npm make

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.11 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.6.1 && \
    npm install -g @protobuf-ts/plugin@v2.11.1

WORKDIR /workspace
