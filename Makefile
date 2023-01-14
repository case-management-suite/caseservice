generate:
	protoc --go_out=grpc --go-grpc_out=grpc grpc/api.proto