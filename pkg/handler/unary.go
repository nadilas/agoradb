// Copyright 2020. feat.Me Networks. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"context"
	"fmt"

	"github.com/featme-inc/agoradb/examples/basic/dsl/protobufs"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/grpc"
)

func UnaryMethodHandler(d *desc.MethodDescriptor) func (srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	fullyQualifiedName := fmt.Sprintf("/%s/%s", d.GetService().GetFullyQualifiedName(), d.GetName())
	return func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
		in := new(protobufs.SaveUserRequest)
		if err := dec(in); err != nil {
			return nil, err
		}
		if interceptor == nil {
			return srv.(protobufs.UserServer).Save(ctx, in)
		}
		info := &grpc.UnaryServerInfo{
			Server:     srv,
			FullMethod: fullyQualifiedName,
		}
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.(protobufs.UserServer).Save(ctx, req.(*protobufs.SaveUserRequest))
		}
		return interceptor(ctx, in, info, handler)
	}
}
