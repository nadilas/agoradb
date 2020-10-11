// Copyright 2020. feat.Me Networks. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gateway

import (
	"context"
)

func (g *gatewayServer) Ping(ctx context.Context, request *PingRequest) (*PingResponse, error) {
	return &PingResponse{Value: "Gateway service OK"}, nil
}
