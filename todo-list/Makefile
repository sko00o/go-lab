
.PHONY: run-server
run-server:
	@go run cmd/server/main.go \
        -db-host=localhost:3308 \
        -db-schema=lab \
        -db-user=root \
        -db-password=toor \
        -grpc-port=9090

.PHONY: run-client
run-client:
	@go run cmd/client-grpc/main.go \
        -server=localhost:9090

.PHONY: api
api:
	@protoc --proto_path=api/proto/v1 \
		--go_out=plugins=grpc:pkg/api/v1 \
		todo-service.proto
