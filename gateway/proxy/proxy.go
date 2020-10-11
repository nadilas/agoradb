// Copyright 2020. feat.Me Networks. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proxy

import (
	"fmt"

	"google.golang.org/grpc/encoding"
)

// Codec returns a proxying codec with conforms both to encoding.Codec and grpc.Codec with the default protobuf codec as parent.
//
// See CodecWithParent.
func Codec() codec {
	// retrieve default proto codec
	return CodecWithParent(encoding.GetCodec("proto"))
}

// CodecWithParent returns a proxying encoding.Codec with a user provided codec as parent.
//
// This codec is *crucial* to the functioning of the proxy. It allows the proxy server to be oblivious
// to the schema of the forwarded messages. It basically treats a gRPC message frame as raw bytes.
// However, if the server handler, or the client caller are not proxy-internal functions it will fall back
// to trying to decode the message using a fallback codec.
func CodecWithParent(fallback encoding.Codec) codec {
	return codec{fallback}
}

type codec struct {
	parentCodec encoding.Codec
}

type frame struct {
	payload []byte
}

func (c codec) Marshal(v interface{}) ([]byte, error) {
	out, ok := v.(*frame)
	if !ok {
		return c.parentCodec.Marshal(v)
	}
	return out.payload, nil
}

func (c codec) Unmarshal(data []byte, v interface{}) error {
	dst, ok := v.(*frame)
	if !ok {
		return c.parentCodec.Unmarshal(data, v)
	}
	dst.payload = data
	return nil
}

// Satisfies grpc.Codec
func (c codec) String() string {
	return fmt.Sprintf("proxy>%s", c.parentCodec.Name())
}

// Name is the name registered for the proxy codec.
// Satisfies encoding.Codec
func (c codec) Name() string {
	return fmt.Sprintf("proxy>%s", c.parentCodec.Name())
}