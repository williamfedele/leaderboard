PROTO_DIR := ./protos
OUT_DIR := ./protogen

PROTO_FILES := $(wildcard $(PROTO_DIR)/*.proto)

SERVER_FILE := ./server/main.go
CLIENT_FILE := ./client/main.go

run: compile_proto run_server

compile_proto: $(PROTO_FILES)
	@mkdir -p $(OUT_DIR)
	@for proto_file in $(PROTO_FILES); do \
		protoc --go_out=$(OUT_DIR) --go-grpc_out=$(OUT_DIR) $$proto_file; \
	done
	@echo "Protobuf files compiled."

run_server:
	@go run $(SERVER_FILE)

run_client:
	@go run $(CLIENT_FILE)

clean:
	@rm -rf $(OUT_DIR)
	@echo "Cleaned up generated files."

.PHONY: run compile_proto run_server run_client clean
