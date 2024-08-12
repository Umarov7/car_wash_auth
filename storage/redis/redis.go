package redis

import (
	"auth-service/config"
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

func ConnectDB(cfg *config.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.REDIS_ADDRESS,
		Password: cfg.REDIS_PASSWORD,
		DB:       cfg.REDIS_DB,
	})

	return rdb
}

func StoreToken(cfg *config.Config, ctx context.Context, userID, token string) error {
	rdb := ConnectDB(cfg)

	err := rdb.Set(ctx, "car-wash:"+userID, token, time.Hour*24).Err()
	if err != nil {
		return errors.Wrap(err, "failed to store token in redis")
	}

	return nil
}

func GetToken(cfg *config.Config, ctx context.Context, userID string) (string, error) {
	rdb := ConnectDB(cfg)

	token, err := rdb.Get(ctx, "car-wash:"+userID).Result()
	if err != nil {
		if err == redis.Nil {
			return "", errors.New("token not found for " + userID)
		}
		return "", errors.Wrap(err, "failed to get token from redis")
	}

	return token, nil
}

func DeleteToken(cfg *config.Config, ctx context.Context, userID string) error {
	rdb := ConnectDB(cfg)

	err := rdb.Del(ctx, "car-wash:"+userID).Err()
	if err != nil {
		return errors.Wrap(err, "failed to delete token from redis")
	}

	return nil
}
