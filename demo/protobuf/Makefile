GRPC_GW_PATH := $(shell go list -f '{{ .Dir }}' github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway)
GOOGLEAPIS_PATH := "$(GRPC_GW_PATH)/../third_party/googleapis"

all:
	@protoc -I${GOOGLEAPIS_PATH} -I=. \
		--go_out=plugins=grpc,paths=source_relative:. \
		srv.proto

	@protoc -I${GOOGLEAPIS_PATH} -I=. \
		--go_out=plugins=grpc,paths=source_relative:. \
		--grpc-gateway_out=paths=source_relative,logtostderr=true,grpc_api_configuration=srv.yaml:. \
		srv.proto

requirements:
	go install github.com/golang/protobuf/protoc-gen-go
	go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
