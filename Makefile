.PHONY: compile
compile: ## Compile the proto file.
	protoc -I pkg/proto/simplegossip/ pkg/proto/simplegossip/simplegossip.proto  --go-grpc_out=pkg/proto/simplegossip/
.PHONY: server
server: ## Build and run server.
	go build -race -ldflags "-s -w" -o bin/server server/main.go
	bin/server
 
.PHONY: client
client: ## Build and run client.
	go build -race -ldflags "-s -w" -o bin/client client/main.go
	bin/client

