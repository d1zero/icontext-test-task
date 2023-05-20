package interfaces

import (
	"context"
	"icontext-test-task/internal/entity"
)

type UserGateway interface {
	GetRedisKey(context.Context, string) (int, error)
	SetRedisKey(context.Context, entity.Value) (int, error)
	SignBody(entity.Sign) ([]byte, error)
	CreateUser(context.Context, entity.User) (entity.User, error)
}
