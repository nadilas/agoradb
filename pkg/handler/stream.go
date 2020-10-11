// Copyright 2020. feat.Me Networks. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/grpc"
)

func ClientStreamHandler(methodDesc *desc.MethodDescriptor) func(srv interface{}, stream grpc.ServerStream) error {
	return func(srv interface{}, stream grpc.ServerStream) error {
		return srv.(UserServer).SaveClientStream(&userSaveClientStreamServer{stream})
	}
}

type User_SaveClientStreamServer interface {
	SendAndClose(*SaveUserResponse) error
	Recv() (*SaveUserRequest, error)
	grpc.ServerStream
}

type userSaveClientStreamServer struct {
	grpc.ServerStream
}

func (x *userSaveClientStreamServer) SendAndClose(m *SaveUserResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *userSaveClientStreamServer) Recv() (*SaveUserRequest, error) {
	m := new(SaveUserRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}


func ServerStreamHandler(methodDesc *desc.MethodDescriptor) func(srv interface{}, stream grpc.ServerStream) error {
	return func(srv interface{}, stream grpc.ServerStream) error {
		m := new(SaveUserRequest)
		if err := stream.RecvMsg(m); err != nil {
			return err
		}
		return srv.(UserServer).SaveServerStream(m, &userSaveServerStreamServer{stream})
	}
}

type User_SaveServerStreamServer interface {
	Send(*SaveUserResponse) error
	grpc.ServerStream
}

type userSaveServerStreamServer struct {
	grpc.ServerStream
}

func (x *userSaveServerStreamServer) Send(m *SaveUserResponse) error {
	return x.ServerStream.SendMsg(m)
}

func ClientServerStreamHandler(methodDesc *desc.MethodDescriptor) func(srv interface{}, stream grpc.ServerStream) error {
	return func(srv interface{}, stream grpc.ServerStream) error {
		return srv.(UserServer).SaveBiStream(&userSaveBiStreamServer{stream})
	}
}

type User_SaveBiStreamServer interface {
	Send(*SaveUserResponse) error
	Recv() (*SaveUserRequest, error)
	grpc.ServerStream
}

type userSaveBiStreamServer struct {
	grpc.ServerStream
}

func (x *userSaveBiStreamServer) Send(m *SaveUserResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *userSaveBiStreamServer) Recv() (*SaveUserRequest, error) {
	m := new(SaveUserRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}
