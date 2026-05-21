package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bmstoss13/the-daily-sunshine/internal/domain"
	"github.com/redis/go-redis/v9"
)

type RedisPostCache struct {
	client *redis.Client
}

func NewRedisPostCache(client *redis.Client) *RedisPostCache {
	return &RedisPostCache{
		client: client,
	}
}

func (r *RedisPostCache) GetTopPostsOfDay(ctx context.Context, dayKey string) ([]domain.Post, bool, error) {
	var cachedPosts []domain.Post

	redisPosts, redisErr := r.client.Get(ctx, dayKey).Result()
	if redisErr == redis.Nil {
		return nil, false, nil
	} else if redisErr != nil {
		return nil, false, fmt.Errorf("[post_redis.go] GetTopPostsOfDay: failed to get results from redis: %w", redisErr)
	}

	err := json.Unmarshal([]byte(redisPosts), &cachedPosts)
	if err != nil {
		return nil, false, fmt.Errorf("[post_redis.go] GetTopPostsOfDay: failed to unmarshal results from redis: %w", err)
	}
	return cachedPosts, true, nil
}

func (r *RedisPostCache) SetTopPostsOfDay(ctx context.Context, dayKey string, posts []domain.Post, ttl time.Duration) error {
	marshaledPosts, err := json.Marshal(posts)
	if err != nil {
		return fmt.Errorf("[post_redis.go] SetTopPostsOfDay: failed to marshal posts: %w", err)
	}

	redisErr := r.client.Set(ctx, dayKey, marshaledPosts, ttl).Err()
	if redisErr != nil {
		return fmt.Errorf("[post_redis.go] SetTopPostsOfDay: failed to set cache in redis: %w", err)
	}
	return nil
}

func (r *RedisPostCache) InvalidateTopPostsOfDay(ctx context.Context, dayKey string) error {
	err := r.client.Del(ctx, dayKey).Err()
	if err != nil {
		return fmt.Errorf("[post_redis.go]: InvalidateTopPostsOfDay: failed to invalidate cache for key %s: %w", dayKey, err)
	}
	return nil
}
