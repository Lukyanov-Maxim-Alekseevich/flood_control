package storage

import "errors"

var (
	ErrRedisKeyNotFound = errors.New("key not found")
)
