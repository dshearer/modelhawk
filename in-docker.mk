PROTO_DIR             := proto/v0
PROTO_FILES           := $(wildcard $(PROTO_DIR)/*.proto)
PROTO_MSG_FILES       := $(filter-out %_message.proto,$(PROTO_FILES))
PROTO_SERVICE_FILES   := $(wildcard $(PROTO_DIR)/*_service.proto)

CP_WITH_COMMENT := script/cp-with-comment

.PHONY: all
all: gen-ts gen-go gen/docs/docs.md

.PHONY: gen-ts
gen-ts:
	@mkdir -p gen/ts/src
	$(CP_WITH_COMMENT) package-config/ts/* gen/ts/
	protoc \
		-I "$(PROTO_DIR)" \
		"--ts_out=gen/ts/src" \
		$(PROTO_FILES)
	cd gen/ts && npm install

.PHONY: gen-go
gen-go:
	@mkdir -p gen/go/v0
	$(CP_WITH_COMMENT) package-config/go/* gen/go/
	protoc \
		-I "$(PROTO_DIR)" \
		"--go_out=gen/go/v0" \
		--go_opt=paths=source_relative \
		"--go-grpc_out=gen/go/v0" \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_FILES)

gen/ts/%: package-config/ts/%
	$(CP_WITH_COMMENT) $^ $@

gen/go/%: package-config/go/%
	$(CP_WITH_COMMENT) $^ $@

gen/docs/docs.md: $(PROTO_FILES)
	@mkdir -p gen/docs
	protoc -I "$(PROTO_DIR)" --doc_out=gen/docs --doc_opt=markdown,docs.md $(PROTO_FILES)
