// Copyright 2020. feat.Me Networks. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package services

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"net"
	"os"

	"github.com/featme-inc/agoradb/examples/basic/dsl/protobufs"
	"github.com/featme-inc/agoradb/internal/handler"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"github.com/kr/pretty"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type databaseService struct {
	database   Database
	listener   net.Listener
	grpcServer *grpc.Server
	logger     *logrus.Entry
}

func newDatabaseService(database Database) *databaseService {
	svc := &databaseService{
		database: database,
		grpcServer: grpc.NewServer(),
		logger: logrus.WithField("service", database.Name),
	}
	svc.registerServices()
	return svc
}

func (s *databaseService) Name() string {
	return s.database.Name
}

func (s *databaseService) Address() string {
	return s.listener.Addr().String()
}

func (s *databaseService) PassthroughAddress() string {
	return fmt.Sprintf("passthrough:///unix://%s", s.Address())
}

func (s *databaseService) ListenAndServe() error {
	name := s.Name()
	sock := "/tmp/" + name + "_1.sock"
	_ = os.Remove(sock)
	var err error
	s.listener, err = net.Listen("unix", sock)
	if err != nil {
		return err
	}
	s.logger.Infof("%s database is listening: %s", name, s.Address())
	return s.grpcServer.Serve(s.listener)
}

func (s *databaseService) GracefulStop() {
	if s.grpcServer == nil {
		return
	}
	s.grpcServer.GracefulStop()
}

func (s *databaseService) registerServices() {
	for _, svcD := range s.database.Descriptor.GetServices() {
		svcName := svcD.GetFullyQualifiedName()
		var methods []grpc.MethodDesc
		var streams []grpc.StreamDesc
		for _, mDesc := range svcD.GetMethods() {
			methodName := mDesc.GetName()
			// simple Name + handler
			if mDesc.IsClientStreaming() || mDesc.IsServerStreaming() {
				sd := grpc.StreamDesc{StreamName: methodName}
				if mDesc.IsClientStreaming() && mDesc.IsServerStreaming() {
					sd.Handler = handler.ClientServerStreamHandler(mDesc)
					sd.ClientStreams = true
					sd.ServerStreams = true
				} else if mDesc.IsClientStreaming() {
					sd.Handler = handler.ClientStreamHandler(mDesc)
					sd.ClientStreams = true
				} else if mDesc.IsServerStreaming() {
					sd.Handler = handler.ServerStreamHandler(mDesc)
					sd.ServerStreams = true
				}
				streams = append(streams, sd)
			} else {
				// unary
				methods = append(methods, grpc.MethodDesc{MethodName: methodName, Handler: handler.UnaryMethodHandler(mDesc)})
			}
		}

		s.logger.Debugf("Adding methods and streams to %s: %s", svcName, pretty.Sprint(methods, streams))
		fileName := svcD.GetFile().GetName()
		buf := getServiceMetadata(svcD, fileName)
		s.grpcServer.RegisterService(&grpc.ServiceDesc{
			ServiceName: svcName,
			HandlerType: (*interface{})(nil),
			Methods:     methods,
			Streams:     streams,
			Metadata:    buf.Bytes(),
		}, &userImpl{})
		s.logger.Infof("Registered %s service on database: %s via: %s", svcName, s.Name(), fileName)
	}
}

func getServiceMetadata(svcD *desc.ServiceDescriptor, fileName string) bytes.Buffer {
	bts, _ := proto.Marshal(svcD.GetFile().AsFileDescriptorProto())
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	w.Name = fileName
	w.Write(bts)
	w.Close()
	return buf
}

type userImpl struct{}

func (u *userImpl) Save(ctx context.Context, request *protobufs.SaveUserRequest) (*protobufs.SaveUserResponse, error) {
	return &protobufs.SaveUserResponse{Value: "user service responding!"}, nil
}
