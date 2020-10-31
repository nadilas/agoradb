// Copyright 2020. feat.Me Networks. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"github.com/featme-inc/agoradb/examples/basic/dsl/protobufs"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/grpc"
)

func ServerStreamHandler(methodDesc *desc.MethodDescriptor) func(srv interface{}, stream grpc.ServerStream) error {
	return func(srv interface{}, stream grpc.ServerStream) error {
		m := new(protobufs.SaveUserRequest)
		if err := stream.RecvMsg(m); err != nil {
			return err
		}
		return srv.(DatabaseHandler).SaveServerStream(m, &userSaveServerStreamServer{stream})
	}
}

type User_SaveServerStreamServer interface {
	Send(*protobufs.SaveUserResponse) error
	grpc.ServerStream
}

type userSaveServerStreamServer struct {
	grpc.ServerStream
}

func (x *userSaveServerStreamServer) Send(m *protobufs.SaveUserResponse) error {
	return x.ServerStream.SendMsg(m)
}
