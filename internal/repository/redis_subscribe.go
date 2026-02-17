package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/opusdvs/DonWeather-ms-subscribe/internal/domain"
	"github.com/redis/go-redis/v9"
)

type RedisSubscribeRepository struct {
	redisClient *redis.Client
	ttl         time.Duration
}

func NewRedisSubscribeRepository(redisClient *redis.Client, ttl time.Duration) *RedisSubscribeRepository {
	return &RedisSubscribeRepository{redisClient: redisClient, ttl: ttl}
}

func (r *RedisSubscribeRepository) Save(ctx context.Context, sub domain.PendingSubscribe) error {
	data, err := json.Marshal(sub)
	if err != nil {
		return err
	}
	return r.redisClient.Set(ctx, sub.Token, data, r.ttl).Err()
}

func (r *RedisSubscribeRepository) Get(ctx context.Context, token string) (domain.PendingSubscribe, error) {
	var sub domain.PendingSubscribe
	raw, err := r.redisClient.Get(ctx, token).Bytes()
	if err != nil {
		return domain.PendingSubscribe{}, err
	}
	if err := json.Unmarshal(raw, &sub); err != nil {
		return domain.PendingSubscribe{}, err
	}
	return sub, nil
}

func (r *RedisSubscribeRepository) Delete(ctx context.Context, token string) error {
	return r.redisClient.Del(ctx, token).Err()
}
