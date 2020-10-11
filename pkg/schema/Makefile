
all: schema

schema: schema.proto
	protoc -I. \
	--go_out=plugins=grpc:. \
	schema.proto