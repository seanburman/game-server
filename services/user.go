package services

import (
	"context"

	"github.com/seanburman/game-ws-server/server"
)

type UserService struct {
	*server.ServerContext
}

func NewUserService(ctx *server.ServerContext) *UserService {
	return &UserService{ServerContext: ctx}
}

func (u UserService) CreateUser(ctx context.Context) error {
	return nil
}
