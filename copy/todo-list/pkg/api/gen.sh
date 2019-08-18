#!/bin/sh

# no need for `go mod vendor`
GRPC_GW_PATH=`go list -f '{{ .Dir }}' github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway`
GRPC_GW_PATH="${GRPC_GW_PATH}/../third_party/googleapis"
SWAGGER_PATH=`go list -f '{{ .Dir }}' github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger`
SWAGGER_PATH="${SWAGGER_PATH}/.."

protoc -I${GRPC_GW_PATH} -I${SWAGGER_PATH} \
        --proto_path=../../api/proto/v1 \
		--go_out=plugins=grpc:./v1 \
		todo-service.proto

# protoc -I${GRPC_GW_PATH} -I${SWAGGER_PATH} \
#         --proto_path=../../api/proto/v1 \
#         --grpc-gateway_out=logtostderr=true,grpc_api_configuration=../../api/proto/v1/todo-service.yaml:./v1 \
#         todo-service.proto

# protoc -I${GRPC_GW_PATH} -I${SWAGGER_PATH} \
#         --proto_path=../../api/proto/v1 \
#         --swagger_out=logtostderr=true,grpc_api_configuration=../../api/proto/v1/todo-service.yaml:../../api/swagger/v1 \
#         todo-service.proto

protoc -I${GRPC_GW_PATH} -I${SWAGGER_PATH} \
        --proto_path=../../api/proto/v1 \
        --grpc-gateway_out=logtostderr=true:./v1 \
        todo-service.proto

protoc -I${GRPC_GW_PATH} -I${SWAGGER_PATH} \
        --proto_path=../../api/proto/v1 \
        --swagger_out=logtostderr=true:../../api/swagger/v1 \
        todo-service.proto
