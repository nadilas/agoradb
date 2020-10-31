package handler

import (
	"context"

	"github.com/featme-inc/agoradb/examples/basic/dsl/protobufs"
	"github.com/featme-inc/agoradb/internal/schema"
	"github.com/sirupsen/logrus"
)

type DatabaseHandler interface {
	Save(ctx context.Context, request *protobufs.SaveUserRequest) (*protobufs.SaveUserResponse, error)
	SaveClientStream(stream User_SaveClientStreamServer) error
	SaveServerStream(request *protobufs.SaveUserRequest, stream User_SaveServerStreamServer) error
	SaveBiStream(stream User_SaveBiStreamServer) error
}

type databaseHandler struct {
	database schema.Database
	logger   *logrus.Entry
}

// Creates a database handler for serving gRPC requests
func New(database schema.Database) *databaseHandler {
	return &databaseHandler{
		database: database,
		logger:   logrus.WithField("service", database.Name),
	}
}

func (u *databaseHandler) Save(ctx context.Context, request *protobufs.SaveUserRequest) (*protobufs.SaveUserResponse, error) {
	u.logger.Debugf("Databasehandler called")
	return &protobufs.SaveUserResponse{Value: "databaseHandler.Save responding!"}, nil
}

func (u *databaseHandler) SaveClientStream(stream User_SaveClientStreamServer) error {
	panic("implement me")
}

func (u *databaseHandler) SaveServerStream(request *protobufs.SaveUserRequest, stream User_SaveServerStreamServer) error {
	panic("implement me")
}

func (u *databaseHandler) SaveBiStream(stream User_SaveBiStreamServer) error {
	panic("implement me")
}
