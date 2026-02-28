include versions.mk

IMAGE             := modelhawk-proto
IN_DOCKER         := $(shell test -f /.dockerenv && echo yes)

DOCKER_RUN   := docker run --rm -v $(CURDIR):/workspace $(IMAGE)

.PHONY: all
all: generate

.PHONY: ref-impls
ref-impls: server


# --- Docker ---

.PHONY: docker-build
docker-build:
	docker build \
		--build-arg PROTOC_GEN_GO_VERSION=$(PROTOC_GEN_GO_VERSION) \
		--build-arg PROTOC_GEN_GO_GRPC_VERSION=$(PROTOC_GEN_GO_GRPC_VERSION) \
		--build-arg PROTOBUF_TS_PLUGIN_VERSION=$(PROTOBUF_TS_PLUGIN_VERSION) \
		--build-arg PROTOC_GEN_DOC_VERSION=$(PROTOC_GEN_DOC_VERSION) \
		-t $(IMAGE) .


# --- Code Generation ---

.PHONY: generate
generate: $(docker-build)
	$(DOCKER_RUN) make -f in-docker.mk


# --- server ref impl ---

reference-impls/server/server : generate-go
	@cd reference-impls/server && go build .


# --- Cleanup ---

.PHONY: clean-gen
clean-gen:
	@rm -rf gen

.PHONY: clean
clean:
	@rm -rf reference-impls/server/server package-config/ts/node_modules
