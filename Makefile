.PHONY: compile
compile: ## Compile the proto file.
	protoc --go_out=pkg/proto --go_opt=paths=source_relative -I pkg/proto pkg/proto/simplegossip/simplegossip.proto  --go-grpc_out=pkg/proto --go-grpc_opt=paths=source_relative 
.PHONY: server
server: ## Build and run server.
	go build -race -ldflags "-s -w" -o bin/server server/main.go
	bin/server
 
.PHONY: client
client: ## Build and run client.
	go build -race -ldflags "-s -w" -o bin/client client/main.go
	bin/client

