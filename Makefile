IMAGE           := modelhawk-proto
PROTO_DIR       := modelhawk/v1
PROTO_FILES     := $(wildcard $(PROTO_DIR)/*.proto)
GO_OUT          := gen/go/v1
TS_OUT          := gen/ts/src

DOCKER_RUN   := docker run --rm -v $(CURDIR):/workspace $(IMAGE)

.PHONY: all generate generate-without-docker generate-go generate-go-local generate-ts generate-ts-local clean docker-build opencode-plugin install-opencode-plugin server ref-impls

all: generate

ref-impls: opencode-plugin server


# --- Docker ---

docker-build:
	docker build -t $(IMAGE) .


# --- Code Generation ---

generate: generate-go generate-ts

generate-without-docker: generate-go-local generate-ts-local

generate-go: docker-build $(PROTO_FILES)
	$(DOCKER_RUN) make generate-go-local

generate-go-local:
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
	@rm -rf "$(TS_OUT)"
	@mkdir -p "$(TS_OUT)"
	protoc \
		-I "$(PROTO_DIR)" \
		"--ts_out=$(TS_OUT)" \
		$(PROTO_FILES)


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
