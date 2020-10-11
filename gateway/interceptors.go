// Copyright 2020. feat.Me Networks. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gateway

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// ContextClientInterceptor passes around headers for tracing and linkerd
func ContextClientInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, resp interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		pairs := make([]string, 0)

		if md, ok := metadata.FromIncomingContext(ctx); ok {
			for key, values := range md {
				if strings.HasPrefix(strings.ToLower(key), "x-") {
					for _, value := range values {
						pairs = append(pairs, key, value)
					}
				}
			}
		}

		ctx = metadata.AppendToOutgoingContext(ctx, pairs...)
		return invoker(ctx, method, req, resp, cc, opts...)
	}
}
