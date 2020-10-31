// Copyright 2020. feat.Me Networks. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package services

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"net"
	"os"

	"github.com/featme-inc/agoradb/internal/handler"
	"github.com/featme-inc/agoradb/internal/schema"
	"github.com/golang/protobuf/proto"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/jhump/protoreflect/desc"
	"github.com/kr/pretty"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type databaseService struct {
	database   schema.Database
	listener   net.Listener
	grpcServer *grpc.Server
	logger     *logrus.Entry
}

func newDatabaseService(database schema.Database) *databaseService {
	logger := logrus.WithField("service", database.Name)
	svc := &databaseService{
		database:   database,
		grpcServer: createGrpcServer(logger),
		logger:     logger,
	}
	svc.registerServices()
	return svc
}

func createGrpcServer(logger *logrus.Entry) *grpc.Server {
	// TODO decide we want to log grpc internal calls as well?
	// grpc_logrus.ReplaceGrpcLogger(logger)
	trcr := opentracing.GlobalTracer()
	return grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_prometheus.UnaryServerInterceptor,
			grpc_recovery.UnaryServerInterceptor(),
			otgrpc.OpenTracingServerInterceptor(trcr),
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_logrus.UnaryServerInterceptor(logger),
			// TODO decide if we want to log payload as well
			// grpc_logrus.PayloadUnaryServerInterceptor(logger, grpc_logrus.WithDecider(func (methodFullName string, err error) bool {}),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_prometheus.StreamServerInterceptor,
			grpc_recovery.StreamServerInterceptor(),
			grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_logrus.StreamServerInterceptor(logger),
			// TODO decide if we want to log payload as well
			// grpc_logrus.PayloadStreamServerInterceptor(logger, func(ctx context.Context, fullMethodName string, servingObject interface{}) bool {}),
		),
	)
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
			HandlerType: (*handler.DatabaseHandler)(nil),
			Methods:     methods,
			Streams:     streams,
			Metadata:    buf.Bytes(),
		}, handler.New(s.database))
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
