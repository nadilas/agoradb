// Copyright 2020. feat.Me Networks. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"github.com/featme-inc/agoradb/examples/basic/dsl/protobufs"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/grpc"
)

func ClientStreamHandler(methodDesc *desc.MethodDescriptor) func(srv interface{}, stream grpc.ServerStream) error {
	return func(srv interface{}, stream grpc.ServerStream) error {
		return srv.(DatabaseHandler).SaveClientStream(&userSaveClientStreamServer{stream})
	}
}

type User_SaveClientStreamServer interface {
	SendAndClose(*protobufs.SaveUserResponse) error
	Recv() (*protobufs.SaveUserRequest, error)
	grpc.ServerStream
}

type userSaveClientStreamServer struct {
	grpc.ServerStream
}

func (x *userSaveClientStreamServer) SendAndClose(m *protobufs.SaveUserResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *userSaveClientStreamServer) Recv() (*protobufs.SaveUserRequest, error) {
	m := new(protobufs.SaveUserRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}
