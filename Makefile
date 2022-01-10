
create: ## create api go files using proto file
	protoc --proto_path=api/proto --go_out=internal/controller api/proto/user.proto
	protoc --proto_path=api/proto --go-grpc_out=internal/controller api/proto/user.proto
.PHONY: create

clean: ## clean api files
	rm pkg/api/*.go
.PHONY: clean

# generates mocks
.PHONY: generate-mocks
generate-mocks:
	go generate ./...
	go mod tidy -compat=1.17
