// Copyright 2020. feat.Me Networks. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"context"

	"github.com/featme-inc/agoradb/internal/schema/schemapb"
	"github.com/jhump/protoreflect/desc/builder"
)

func (s *Server) CreateDatabase(ctx context.Context, request *schemapb.Database) (*schemapb.CreateDatabaseResponse, error) {
	//mdEmpty, _ := desc.LoadMessageDescriptorForMessage((*empty.Empty)(nil))
	//mdAny,_ := desc.LoadMessageDescriptorForMessage((*any.Any)(nil))
	//mbAny, _ := builder.FromMessage(mdAny)

	databaseObject := builder.NewMessage("database").AddField(builder.NewField("createdBy", builder.FieldTypeString()))
	infoRM := builder.NewMessage("InfoRequest")
	infoRS := builder.NewMessage("InfoResponse").AddField(builder.NewField("database", builder.FieldTypeString()))
	database := request.Name
	fd, _ := builder.NewFile(database + ".proto").
		SetProto3(true).
		SetPackageName(database).
		AddMessage(databaseObject).
		AddMessage(infoRM).
		AddMessage(infoRS).
		AddService(builder.NewService("Database").
			AddMethod(builder.NewMethod("Info", builder.RpcTypeMessage(infoRM, false), builder.RpcTypeMessage(infoRS, false)))).
		Build()

	if err := s.repository.CreateDatabase(database, fd); err != nil {
		return &schemapb.CreateDatabaseResponse{ErrorMessage: err.Error()}, nil
	}

	go s.manager.StartDatabase(Database{
		Name:       database,
		Descriptor: fd,
	})

	return &schemapb.CreateDatabaseResponse{Success: true}, nil
}
