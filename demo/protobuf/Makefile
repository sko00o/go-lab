GRPC_GW_PATH:="$(shell go list -f '{{ .Dir }}' github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway)/../third_party/googleapis"

all:
	@protoc -I${GRPC_GW_PATH} \
		--proto_path=. \
		--go_out=plugins=grpc:. \
		--go_out=plugins=grpc:. \
		srv.proto

	@protoc -I${GRPC_GW_PATH} \
		--proto_path=. \
		--go_out=plugins=grpc:. \
		--grpc-gateway_out=logtostderr=true,grpc_api_configuration=srv.yaml:. \
		srv.proto