FROM golang:1.24-alpine

ARG PROTOC_GEN_GO_VERSION
ARG PROTOC_GEN_GO_GRPC_VERSION
ARG PROTOBUF_TS_PLUGIN_VERSION
ARG PROTOC_GEN_DOC_VERSION

RUN apk add --no-cache protobuf protobuf-dev nodejs npm make

# Install Go tools
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@${PROTOC_GEN_GO_VERSION} && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@${PROTOC_GEN_GO_GRPC_VERSION} && \
    go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@${PROTOC_GEN_DOC_VERSION}

# Install npm tools
RUN npm install -g @protobuf-ts/plugin@${PROTOBUF_TS_PLUGIN_VERSION}

WORKDIR /workspace
