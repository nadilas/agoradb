// Copyright 2020. feat.Me Networks. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package proxy provides a reverse proxy handler for gRPC.
//
//The implementation allows a `grpc.Server` to pass a received ServerStream to a ClientStream without understanding
// the semantics of the messages exchanged. It basically provides a transparent reverse-proxy.
//
// This package is intentionally generic, exposing a `StreamDirector` function that allows users of this package
// to implement whatever logic of backend-picking, dialing and service verification to perform.
//
// Package is a revised version of https://github.com/mwitkow/grpc-proxy
package proxy
