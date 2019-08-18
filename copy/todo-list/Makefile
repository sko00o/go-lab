
.PHONY: run-server
run-server:
	@echo "run server"
	@go run cmd/server/main.go \
        -db-host=localhost:3308 \
        -db-schema=lab \
        -db-user=root \
        -db-password=toor \
        -grpc-port=9090 \
		-http-port=9091 \
		-log-level=-1 \
		-log-time-format=2006-01-02T15:04:05.999999999Z07:00

.PHONY: run-grpc-client
run-grpc-client:
	@echo "run grpc client"
	@go run cmd/client-grpc/main.go \
        -server=localhost:9090

.PHONY: run-rest-client
run-rest-client:
	@echo "run rest client"
	@go run cmd/client-rest/main.go \
        -server=http://localhost:9091

.PHONY: api
api:
	@echo "generate api"
	@go generate pkg/api/api.go

.PHONY: requirements
requirements:
	@echo "set dev requirements"
	@export GOPROXY=https://goproxy.io
	go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
	go get -u github.com/golang/protobuf/protoc-gen-go
