package floodcontrol

import (
	"context"
	"time"
)

type InMemoryStorage interface {
	Set(ctx context.Context, key, value string, exp time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}
