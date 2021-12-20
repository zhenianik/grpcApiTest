API_NAME 		    := api
MIGRATE_NAME   		:= migrate

.PHONY: config
create: ## create api go files using proto file
	protoc --proto_path=api/proto --go_out=pkg api/proto/user.proto
	protoc --proto_path=api/proto --go-grpc_out=pkg api/proto/user.proto
	##protoc --goclay_out=:pkg api/proto/user.proto

clean: ## clean api files
	rm pkg/api/*.go

goclay: ## generate handler with goclay
	protoc --goclay_out=:pkg api/proto/user.proto
