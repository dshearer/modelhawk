include versions.mk

IMAGE             := modelhawk-proto
IN_DOCKER         := $(shell test -f /.dockerenv && echo yes)
PROTO_DIR         := proto/v0
PROTO_FILES       := $(wildcard $(PROTO_DIR)/*.proto)
PACKAGE_CONFIG_GO := $(wildcard package-config/go/*)
PACKAGE_CONFIG_TS := $(wildcard package-config/ts/*)

DOCKER_RUN   := docker run --rm -v $(CURDIR):/workspace $(IMAGE)
CP_WITH_COMMENT := script/cp-with-comment

.PHONY: all generate generate-without-docker generate-go generate-go-local generate-ts generate-ts-local clean clean-gen docker-build server ref-impls generate-proto-docs generate-proto-docs-local

all: generate

ref-impls: server


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

generate-go-local: $(patsubst package-config/go/%,gen/go/%,$(PACKAGE_CONFIG_GO)) $(PROTO_FILES)
	$(if $(IN_DOCKER),,go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION))
	$(if $(IN_DOCKER),,go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION))
	@rm -rf gen/go
	@mkdir -p gen/go/v0
	$(CP_WITH_COMMENT) package-config/go/* gen/go/
	protoc \
		-I "$(PROTO_DIR)" \
		"--go_out=gen/go/v0" \
		--go_opt=paths=source_relative \
		"--go-grpc_out=gen/go/v0" \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_FILES)

gen/go/%: package-config/go/%
	$(CP_WITH_COMMENT) $^ $@

generate-ts: docker-build
	$(DOCKER_RUN) make generate-ts-local

generate-ts-local: $(patsubst package-config/ts/%,gen/ts/%,$(PACKAGE_CONFIG_GO)) $(PROTO_FILES)
	$(if $(IN_DOCKER),,npm install -g @protobuf-ts/plugin@$(PROTOBUF_TS_PLUGIN_VERSION))
	@rm -rf gen/ts
	@mkdir -p gen/ts/src
	$(CP_WITH_COMMENT) package-config/ts/* gen/ts/
	protoc \
		-I "$(PROTO_DIR)" \
		"--ts_out=gen/ts/src" \
		$(PROTO_FILES)

gen/ts/%: package-config/ts/%
	cd package-config/ts && npm i
	$(CP_WITH_COMMENT) $^ $@

generate-proto-docs: docker-build
	$(DOCKER_RUN) make generate-proto-docs-local

generate-proto-docs-local: gen/docs/docs.md

gen/docs/docs.md: $(PROTO_FILES)
	@mkdir -p gen/docs
	$(if $(IN_DOCKER),,go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@$(PROTOC_GEN_DOC_VERSION))
	protoc -I "$(PROTO_DIR)" --doc_out=gen/docs --doc_opt=markdown,docs.md $(PROTO_FILES)


# --- server ref impl ---

server : generate-go
	@cd reference-impls/server && go build .


# --- Cleanup ---

clean-gen:
	@rm -rf gen

clean:
	@rm -rf reference-impls/server/server package-config/ts/node_modules
