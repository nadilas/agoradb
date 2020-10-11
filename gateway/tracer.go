// Copyright 2020. feat.Me Networks. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gateway

import (
	"os"

	"github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/reporter"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/sirupsen/logrus"

)

func tracer(name string) (opentracing.Tracer, reporter.Reporter) {
	addr := os.Getenv("ZIPKIN_ADDR")
	if addr == "" {
		return opentracing.GlobalTracer(), nil
	}

	// create our local service endpoint
	endpoint, _ := zipkin.NewEndpoint(name, name)

	logrus.Infof("Using Zipkin HTTP tracer: %s", addr)
	zipkinReporter := zipkinhttp.NewReporter(addr)

	// initialize our tracer
	nativeTracer, err := zipkin.NewTracer(zipkinReporter, zipkin.WithLocalEndpoint(endpoint))
	if err != nil {
		logrus.Fatalf("unable to create tracer: %+v\n", err)
	}

	// use zipkin-go-opentracing to wrap our tracer
	tracer := zipkinot.Wrap(nativeTracer)

	// optionally set as Global OpenTracing tracer instance
	opentracing.SetGlobalTracer(tracer)

	return tracer, zipkinReporter
}
