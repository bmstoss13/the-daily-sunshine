package redis

import (
	"context"
	"time"

	"github.com/bmstoss13/the-daily-sunshine/internal/domain"
)

type NoopPostCache struct{}

func NewNoopPostCache() *NoopPostCache {
	return &NoopPostCache{}
}

func (n *NoopPostCache) GetTopPostsOfDay(ctx context.Context, dayKey string) ([]domain.Post, bool, error) {
	return nil, false, nil
}

func (n *NoopPostCache) SetTopPostsOfDay(ctx context.Context, dayKey string, posts []domain.Post, ttl time.Duration) error {
	return nil
}

func (n *NoopPostCache) InvalidateTopPostsOfDay(ctx context.Context, dayKey string) error {
	return nil
}
