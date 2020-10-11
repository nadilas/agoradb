// Copyright 2020. feat.Me Networks. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proxy

import (
	"context"
	"io"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type handler struct {
	director StreamDirector
}

var (
	clientStreamingDescForProxying = &grpc.StreamDesc{
		ServerStreams: true,
		ClientStreams: true,
	}
)

func RegisterService(server *grpc.Server, director StreamDirector, serviceName string, methodNames ...string) {
	streamer := &handler{director}
	fakeDesc := &grpc.ServiceDesc{
		ServiceName: serviceName,
		HandlerType: (*interface{})(nil),
	}
	// register methods
	for _, method := range methodNames {
		streamDesc := grpc.StreamDesc{
			StreamName:    method,
			Handler:       streamer.handler,
			ServerStreams: true,
			ClientStreams: true,
		}
		fakeDesc.Streams = append(fakeDesc.Streams, streamDesc)
	}

	// attach service
	server.RegisterService(fakeDesc, streamer)
}

// TransparentHandler returns a handler that attempts to proxy all requests that are not registered in the server.
// The indented use here is as a transparent proxy, where the server doesn't know about the services implemented by the
// backends. It should be used as a `grpc.UnknownServiceHandler`.
//
// This can *only* be used if the `server` also uses proxy.CustomCodec(proxy.Codec()) ServerOption.
func TransparentHandler(director StreamDirector) grpc.StreamHandler {
	streamer := &handler{director}
	return streamer.handler
}

// handler is where the proxying happens
//
// It is invoked like any gRPC server stream and uses the gRPC server framing to get and receive bytes from the wire,
// forwarding it to a ClientStream established against the relevant ClientConn
func (h handler) handler(srv interface{}, serverStream grpc.ServerStream) error {
	fullMethodName, ok := grpc.MethodFromServerStream(serverStream)
	if !ok {
		return status.Errorf(codes.Internal, "lowLevelServerStream not exists in context")
	}
	outgoingCtx, backendConn, err := h.director(serverStream.Context(), fullMethodName)
	if err != nil {
		return err
	}

	// plug in client IP if available
	if p, ok := peer.FromContext(serverStream.Context()); ok {
		outgoingCtx = metadata.AppendToOutgoingContext(outgoingCtx, "X-Forwarded-For", p.Addr.String())
	}
	clientCtx, clientCancel := context.WithCancel(outgoingCtx)
	clientStream, err := grpc.NewClientStream(clientCtx, clientStreamingDescForProxying, backendConn, fullMethodName)
	if err != nil {
		return err
	}

	// Explicitly *do not close* s2cErrchan and c2sErrChan, otherwise the select below will not terminate
	//   Channels do not have to be closed, it is just a control flow mechanism, see
	//   https://groups.google.com/forum/#!msg/golang-nuts/pZwdYRGxCIk/qpbHxRRPJdUJ
	s2cErrChan := h.forwardServerToClient(serverStream, clientStream)
	c2sErrChan := h.forwardClientToServer(clientStream, serverStream)
	// We don't know which side is going to stop sending first, so we need a select the one which comes earlier
	for i := 0; i < 2; i++ {
		select {
		case s2cErr := <- s2cErrChan:
			// handle s2c err
			if s2cErr == io.EOF {
				// this is the happy case where the sender has encountered io.EOF, and won't be sending anymore.
				//   the clientStream>serverStream may continue pumping though.
				_ = clientStream.CloseSend()
			} else {
				// however, we may have gotten a receive error (stream disconnected, a read error etc) in which case
				//   we need to cancel the clientStream to the backend, let all of its goroutines be freed up by the
				//   CancelFunc and exit with an error to the stack
				clientCancel()
				return status.Errorf(codes.Internal, "failed proxying s2c: %v", s2cErr)
			}
		case c2sErr := <- c2sErrChan:
			// handle c2s err
			// This happends when the clientStream has nothing else to offer (io.EOF) or return a gRPC error. In those two
			// cases we may have received Trailers as part of the call. In case of other errors (stream closed) the trailers
			// will be nil.
			serverStream.SetTrailer(clientStream.Trailer())
			// c2sErr will contain RPC error from client code. If not io.EOF return the RPC error as server stream error.
			if c2sErr != io.EOF {
				return c2sErr
			}
			return nil
		}
	}
	return status.Errorf(codes.Internal, "gRPC proxying: unreachable code")
}

func (h *handler) forwardClientToServer(src grpc.ClientStream, dst grpc.ServerStream) chan error {
	ret := make(chan error, 1)
	go func() {
		f := &frame{}
		for i := 0; ; i++ {
			if err := src.RecvMsg(f); err != nil {
				ret <-err // this can be io.EOF which is happy case
				break
			}
			// This is a bit of a hack, but client to server headers are only readable after first client msg is
			// received but must be written to server stream before the first msg is flushed.
			// This is the only place to do it nicely.
			if i == 0 {
				md, err := src.Header()
				if err != nil {
					ret <- err
					break
				}
				if err := dst.SendHeader(md); err != nil {
					ret <- err
					break
				}
			}
			if err := dst.SendMsg(f); err != nil {
				ret <- err
				break
			}
		}
	}()
	return ret
}

func (h *handler) forwardServerToClient(src grpc.ServerStream, dst grpc.ClientStream) chan error {
	ret := make(chan error, 1)
	go func() {
		f := &frame{}
		for {
			if err := src.RecvMsg(f); err != nil {
				ret <- err // this can be io.EOF which is happy case
				break
			}
			if err := dst.SendMsg(f); err != nil {
				ret <- err
				break
			}
		}
	}()
	return ret
}