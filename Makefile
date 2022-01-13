
create-proto: ## create api go files using proto file
	protoc --proto_path=internal/controller/grpc/proto --go_out=internal/controller/grpc internal/controller/grpc/proto/user.proto
	protoc --proto_path=internal/controller/grpc/proto --go-grpc_out=internal/controller/grpc internal/controller/grpc/proto/user.proto
.PHONY: create-proto

clean: ## clean api files
	rm internal/controller/grpc/api/*.go
.PHONY: clean

generate-mocks:
	go generate ./...
	go mod tidy -compat=1.17
.PHONY: generate-mocks

docker-up:
	docker compose up
.PHONY: docker-up

