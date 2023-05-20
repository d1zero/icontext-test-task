package interfaces

import (
	"context"
	"icontext-test-task/internal/entity"
)

type UserService interface {
	IncrementValue(context.Context, entity.Value) (int, error)
	SignBody(entity.Sign) (string, error)
	CreateUser(context.Context, entity.User) (int, error)
}
