// Copyright 2020. feat.Me Networks. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This package is organizes the schema related use cases
//
// In order to get started with a new schema
package schema

import (
	"github.com/featme-inc/agoradb/internal/services"
	"github.com/jhump/protoreflect/desc"
)

type Server struct {
	repository WriteDatabaseRepository
	manager    DatabaseManager
}

type WriteDatabaseRepository interface {
	CreateDatabase(database string, fd *desc.FileDescriptor) error
}

type DatabaseManager interface {
	StartDatabase(database services.Database)
}

func NewServer(repository WriteDatabaseRepository, manager DatabaseManager) *Server {
	return &Server{repository: repository, manager: manager}
}
