package service

import (
	"context"
	"encoding/hex"
	"icontext-test-task/internal/entity"
	"icontext-test-task/internal/interfaces"
)

type UserService struct {
	userRepo interfaces.UserGateway
}

func (u *UserService) IncrementValue(ctx context.Context, p entity.Value) (int, error) {
	res, err := u.userRepo.GetRedisKey(ctx, p.Key)
	if err != nil {
		return 0, err
	}

	p.Value += res

	result, err := u.userRepo.SetRedisKey(ctx, p)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (u *UserService) SignBody(p entity.Sign) (string, error) {
	result, err := u.userRepo.SignBody(p)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(result), nil
}

func (u *UserService) CreateUser(ctx context.Context, p entity.User) (int, error) {
	result, err := u.userRepo.CreateUser(ctx, p)
	if err != nil {
		return 0, err
	}

	return result.ID, nil
}

var _ interfaces.UserService = (*UserService)(nil)

func NewUserService(
	userRepo interfaces.UserGateway,
) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}
