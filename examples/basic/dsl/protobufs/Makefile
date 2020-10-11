
all: test

test: testdb.proto
	protoc -I. \
	--go_out=plugins=grpc:. \
	testdb.proto