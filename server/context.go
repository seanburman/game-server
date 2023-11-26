package server

import (
	"context"

	db "github.com/seanburman/game-ws-server/db/clients"
)

type ServerContext struct {
	context.Context
	*db.Clients
}

func NewServerContext() *ServerContext {
	return &ServerContext{
		Context: context.Background(),
	}
}
