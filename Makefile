.PHONY: compile
compile: ## Compile the proto file.
	protoc --go_out=simplegossip --go_opt=paths=source_relative -I simplegossip simplegossip/simplegossip.proto  --go-grpc_out=simplegossip --go-grpc_opt=paths=source_relative 
.PHONY: server
server: ## Build and run server.
	go build -race -ldflags "-s -w" -o bin/gossipserver server/main.go
 
.PHONY: client
client: ## Build and run client.
	go build -race -ldflags "-s -w" -o bin/client client/main.go

