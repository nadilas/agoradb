
all: gateway

gateway: gateway.proto
	protoc -I. \
	--go_out=plugins=grpc:. \
	gateway.proto