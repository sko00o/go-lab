GRPC_GW_PATH:="$(shell go list -f '{{ .Dir }}' github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway)/../third_party/googleapis"

all:
	@protoc -I${GRPC_GW_PATH} \
		--proto_path=. \
		--go_out=plugins=grpc:. \
		service.proto

	@protoc -I${GRPC_GW_PATH} \
		--proto_path=. \
		--grpc-gateway_out=logtostderr=true:. \
		service.proto

	@protoc -I${GRPC_GW_PATH} \
		--proto_path=. \
		--swagger_out=logtostderr=true:. \
		service.proto