
MAKEFLAGS += --silent

.PHONY: server-v1 client-v1

server-v1:
	go run ws/v1/server/main.go

client-v1:
	go run ws/v1/client/main.go -conn 500

