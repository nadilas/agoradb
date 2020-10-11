// Copyright 2020. feat.Me Networks. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gateway

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/featme-inc/agoradb/gateway/proxy"
	"github.com/featme-inc/agoradb/internal/services"
	"github.com/featme-inc/agoradb/schema"
	"github.com/featme-inc/agoradb/storage/memory"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/openzipkin/zipkin-go/reporter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	_ "google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/metadata"
	// _ "google.golang.org/grpc/resolver/passthrough"
	"google.golang.org/grpc/status"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
)

type gatewayServer struct {
	config           Config
	listener         net.Listener
	grpcServer       *grpc.Server
	prometheusServer *http.Server
	serviceManager   *services.Manager
	repository       Repository
	zipkinReporter   reporter.Reporter
}

type RegisterImplementation func(s *grpc.Server)

// ServerConfig is a generic server configuration
type ServerConfig struct {
	Port int
	Host string
}

// Address Gets a logical addr for a ServerConfig
func (c *ServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

type Config struct {
	ServerConfig       ServerConfig
	PrometheusConfig   ServerConfig
	UnaryInterceptors  []grpc.UnaryServerInterceptor
	StreamInterceptors []grpc.StreamServerInterceptor

	// Register implementations
	//
	// The grpc server instance is handed over in this method for additional Server/Service registrations
	GRPCImplementation RegisterImplementation
	GRPCOptions        []grpc.ServerOption
	Repository         Repository
	ReflectionEnabled  bool
}

type ManagedLifeCycle interface {
	Stop() error
}

type Repository interface {
	ManagedLifeCycle
	services.ReadDatabaseRepository
	schema.WriteDatabaseRepository
}

var (
	DefaultConfig = Config{
		PrometheusConfig:   ServerConfig{Port: 9000, Host: "0.0.0.0"},
		ServerConfig:       ServerConfig{Port: 5750, Host: "0.0.0.0"},
		GRPCImplementation: func(s *grpc.Server) {},
		UnaryInterceptors: []grpc.UnaryServerInterceptor{
			grpc_prometheus.UnaryServerInterceptor,
			grpc_recovery.UnaryServerInterceptor(),
		},
		StreamInterceptors: []grpc.StreamServerInterceptor{
			grpc_prometheus.StreamServerInterceptor,
			grpc_recovery.StreamServerInterceptor(),
		},
		Repository: memory.New(),
		ReflectionEnabled: true,
	}
)

func New(config Config) *gatewayServer {
	trcr, zipkinReporter := tracer("gateway")
	config.UnaryInterceptors = append(config.UnaryInterceptors, otgrpc.OpenTracingServerInterceptor(trcr))

	return &gatewayServer{
		config:         config,
		serviceManager: services.NewManager(config.Repository),
		repository:     config.Repository,
		zipkinReporter: zipkinReporter,
	}
}

func (g *gatewayServer) Serve() error {
	g.startPrometheusServer()
	g.serviceManager.StartDatabases()
	return g.serveGrpc()
}

// Shuts down the gateway server in a graceful way
func (g *gatewayServer) GracefulShutdown() {
	logrus.Infof("Gracefully shutting down agoraDB gateway")
	if err := g.serviceManager.Stop(); err != nil {
		logrus.Errorf("ERR: failed to stop database services. %v", err)
	}
	g.grpcServer.GracefulStop()
	if g.zipkinReporter != nil {
		if err := g.zipkinReporter.Close(); err != nil {
			logrus.Errorf("Failed shutting down zipkin reporter. Error: %v", err)
		}
	}
	if err := g.repository.Stop(); err != nil {
		logrus.Errorf("ERR: failed to stop repository: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()
	if err := g.prometheusServer.Shutdown(ctx); err != nil {
		logrus.Infof("Timeout during shutdown of metrics server. Error: %v", err)
	}
}

func (g *gatewayServer) serveGrpc() error {
	var err error
	g.listener, err = net.Listen("tcp", g.config.ServerConfig.Address())
	if err != nil {
		log.Fatal(err)
	}

	logrus.Infof("Serving agoraDD gateway on %s", g.config.ServerConfig.Address())
	return g.createGrpcServer().Serve(g.listener)
}

func (g *gatewayServer) createGrpcServer() *grpc.Server {
	g.config.GRPCOptions = append(g.config.GRPCOptions, grpc.UnaryInterceptor(
		grpc_middleware.ChainUnaryServer(g.config.UnaryInterceptors...),
	))
	g.config.GRPCOptions = append(g.config.GRPCOptions, grpc.StreamInterceptor(
		grpc_middleware.ChainStreamServer(g.config.StreamInterceptors...),
	))
	g.config.GRPCOptions = append(g.config.GRPCOptions, grpc.CustomCodec(proxy.Codec()))
	// here the gateway has no knowledge of the schema services
	// because we cannot register them beforehand:
	//  if we had the schema and all endpoints at compile time, we could:
	//   https://github.com/mwitkow/grpc-proxy/blob/0f1106ef9c766333b9acb4b81e705da4bade7215/proxy/examples_test.go#L20
	g.config.GRPCOptions = append(g.config.GRPCOptions, grpc.UnknownServiceHandler(proxy.TransparentHandler(g.streamDirector)))

	g.grpcServer = grpc.NewServer(
		g.config.GRPCOptions...,
	)
	g.config.GRPCImplementation(g.grpcServer)

	grpc_prometheus.EnableHandlingTimeHistogram(func(opts *prometheus.HistogramOpts) {
		opts.Buckets = prometheus.ExponentialBuckets(0.005, 1.4, 20)
	})

	g.registerMainServices()

	return g.grpcServer
}

func (g *gatewayServer) startPrometheusServer() {
	g.prometheusServer = &http.Server{Addr: g.config.PrometheusConfig.Address()}

	http.Handle("/metrics", promhttp.Handler())
	logrus.Infof("Gateway metrics at http://%s/metrics", g.config.PrometheusConfig.Address())

	go func() {
		if err := g.prometheusServer.ListenAndServe(); err != nil {
			logrus.Errorf("Metrics http server: ListenAndServe() error: %v", err)
		}
	}()
}

func (g *gatewayServer) registerMainServices() {
	grpc_prometheus.Register(g.grpcServer)
	if g.config.ReflectionEnabled {
		registerReflection(g)
	}
	schema.RegisterSchemaServer(g.grpcServer, schema.NewServer(g.repository, g.serviceManager))
	RegisterGatewayServer(g.grpcServer, g)
}

func (g *gatewayServer) streamDirector(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error) {
	// Make sure we never forward internal services
	if strings.HasPrefix(fullMethodName, "/io.agoradb") {
		return nil, nil, status.Errorf(codes.Unimplemented, "Unknown method")
	}
	md, ok := metadata.FromIncomingContext(ctx)
	// copy inbound md explicitly
	outCtx, _ := context.WithCancel(ctx)
	outCtx = metadata.NewOutgoingContext(outCtx, md.Copy())
	if ok {
		// Decide which backend to dial
		if backend := g.serviceManager.BackendForDatabaseService(fullMethodName); backend != "" {
			conn, err := grpc.DialContext(ctx, backend, grpc.WithDefaultCallOptions(grpc.ForceCodec(proxy.Codec())), grpc.WithInsecure())
			return outCtx, conn, err
		}
	}
	if strings.HasPrefix(fullMethodName, "/grpc.reflection") {
		return nil, nil, status.Errorf(codes.Unavailable, "No database services are running")
	}
	dbName := strings.Split(strings.ReplaceAll(fullMethodName, "/", ""), ".")[0]
	return nil, nil, status.Errorf(codes.NotFound, fmt.Sprintf("Unknown database: %s", dbName))
}
