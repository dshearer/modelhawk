include versions.mk

IMAGE           := modelhawk-proto
PROTO_DIR       := modelhawk/v1
PROTO_FILES     := $(wildcard $(PROTO_DIR)/*.proto)
GO_OUT          := gen/go/v1
TS_OUT          := gen/ts/src

DOCKER_RUN   := docker run --rm -v $(CURDIR):/workspace $(IMAGE)

.PHONY: all generate generate-without-docker generate-go generate-go-local generate-ts generate-ts-local clean docker-build opencode-plugin install-opencode-plugin server ref-impls generate-proto-docs generate-proto-docs-local

all: generate

ref-impls: opencode-plugin server


# --- Docker ---

docker-build:
	docker build \
		--build-arg PROTOC_GEN_GO_VERSION=$(PROTOC_GEN_GO_VERSION) \
		--build-arg PROTOC_GEN_GO_GRPC_VERSION=$(PROTOC_GEN_GO_GRPC_VERSION) \
		--build-arg PROTOBUF_TS_PLUGIN_VERSION=$(PROTOBUF_TS_PLUGIN_VERSION) \
		--build-arg PROTOC_GEN_DOC_VERSION=$(PROTOC_GEN_DOC_VERSION) \
		-t $(IMAGE) .


# --- Code Generation ---

generate: generate-go generate-ts generate-proto-docs

generate-without-docker: generate-go-local generate-ts-local generate-proto-docs-local

generate-go: docker-build $(PROTO_FILES)
	$(DOCKER_RUN) make generate-go-local

generate-go-local:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION)
	@rm -rf "$(GO_OUT)"
	@mkdir -p "$(GO_OUT)"
	protoc \
		-I "$(PROTO_DIR)" \
		"--go_out=$(GO_OUT)" \
		--go_opt=paths=source_relative \
		"--go-grpc_out=$(GO_OUT)" \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_FILES)

generate-ts: docker-build $(PROTO_FILES)
	$(DOCKER_RUN) make generate-ts-local

generate-ts-local:
	npm install -g @protobuf-ts/plugin@$(PROTOBUF_TS_PLUGIN_VERSION)
	@rm -rf "$(TS_OUT)"
	@mkdir -p "$(TS_OUT)"
	protoc \
		-I "$(PROTO_DIR)" \
		"--ts_out=$(TS_OUT)" \
		$(PROTO_FILES)

generate-proto-docs: docker-build
	$(DOCKER_RUN) make generate-proto-docs-local

generate-proto-docs-local:
	@mkdir -p gen/docs
	go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@$(PROTOC_GEN_DOC_VERSION)
	protoc -I "$(PROTO_DIR)" --doc_out=gen/docs --doc_opt=markdown,docs.md $(PROTO_FILES)


# --- opencode-plugin ref impl ---

opencode-plugin : generate-ts
	@cd gen/ts && npm install
	@cd gen/ts && npm run build
	@cd reference-impls/opencode-plugin && bun install
	@cd reference-impls/opencode-plugin && bun run build


# --- server ref impl ---

server : generate-go
	@cd reference-impls/server && go build .

# --- Cleanup ---

clean:
	rm -rf gen/ts/dist gen/ts/node_modules reference-impls/opencode-plugin/dist reference-impls/opencode-plugin/node_modules reference-impls/server/server opencode.log
