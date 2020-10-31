// Copyright 2020. feat.Me Networks. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"github.com/featme-inc/agoradb/examples/basic/dsl/protobufs"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/grpc"
)

func ClientServerStreamHandler(methodDesc *desc.MethodDescriptor) func(srv interface{}, stream grpc.ServerStream) error {
	return func(srv interface{}, stream grpc.ServerStream) error {
		return srv.(DatabaseHandler).SaveBiStream(&userSaveBiStreamServer{stream})
	}
}

type User_SaveBiStreamServer interface {
	Send(*protobufs.SaveUserResponse) error
	Recv() (*protobufs.SaveUserRequest, error)
	grpc.ServerStream
}

type userSaveBiStreamServer struct {
	grpc.ServerStream
}

func (x *userSaveBiStreamServer) Send(m *protobufs.SaveUserResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *userSaveBiStreamServer) Recv() (*protobufs.SaveUserRequest, error) {
	m := new(protobufs.SaveUserRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}
