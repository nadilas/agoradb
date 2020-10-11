// Copyright 2020. feat.Me Networks. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gateway

import (
	"context"
)

func (g *gatewayServer) Info(ctx context.Context, request *InfoRequest) (*InfoResponse, error) {
	var dbs []*DatabaseService
	for _, ds := range g.serviceManager.AllDatabases() {
		var svcs []*ServiceInfo
		for svc, info := range ds.Info {
			var methods []*MethodInfo
			for _, mi := range info.Methods {
				methods = append(methods, &MethodInfo{
					Name:           mi.Name,
					IsClientStream: mi.IsClientStream,
					IsServerStream: mi.IsServerStream,
				})
			}
			svcs = append(svcs, &ServiceInfo{Name: svc, Metadata: info.Metadata.(string), Methods: methods})
		}
		dbs = append(dbs, &DatabaseService{
			Name:     ds.Name,
			Services: svcs,
		})
	}
	return &InfoResponse{Databases: dbs}, nil
}
