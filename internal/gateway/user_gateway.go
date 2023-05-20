package gateway

import (
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"icontext-test-task/internal/entity"
	"icontext-test-task/internal/interfaces"
	"strconv"
	"time"
)

type UserRepository struct {
	db    *sqlx.DB
	redis *redis.Client
}

func (r *UserRepository) GetRedisKey(ctx context.Context, key string) (int, error) {
	res, err := r.redis.Get(ctx, key).Result()

	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}

	return strconv.Atoi(res)
}

func (r *UserRepository) SetRedisKey(ctx context.Context, p entity.Value) (int, error) {
	_, err := r.redis.Set(ctx, p.Key, p.Value, 72*time.Hour).Result()
	if err != nil {
		return 0, err
	}

	return p.Value, nil
}

func (r *UserRepository) SignBody(p entity.Sign) ([]byte, error) {
	hash := hmac.New(sha512.New, []byte(p.Key))
	if _, err := hash.Write([]byte(p.Text)); err != nil {
		return []byte{}, err
	}
	return hash.Sum(nil), nil
}

func (r *UserRepository) CreateUser(ctx context.Context, p entity.User) (entity.User, error) {
	var result entity.User

	q := `
		INSERT INTO public.user (name,age)
		VALUES ($1, $2)
		RETURNING id, name, age
	`

	err := r.db.GetContext(ctx, &result, q, p.Name, p.Age)
	if err != nil {
		return entity.User{}, err
	}

	return result, nil
}

var _ interfaces.UserGateway = (*UserRepository)(nil)

func NewUserRepository(
	db *sqlx.DB,
	redis *redis.Client,
) *UserRepository {
	return &UserRepository{
		db:    db,
		redis: redis,
	}
}
