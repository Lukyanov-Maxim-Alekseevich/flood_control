package floodcontrol

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"task/internal/infra/storage"
	"time"
)

var (
	ErrCheckCallLimitExceed = errors.New("call limit exceeded")
)

type Config struct {
	TimeIntervalSeconds time.Duration `yaml:"time-interval-seconds"`
	CallCountLimit      int           `yaml:"call-count-limit"`
}

type FloodControl struct {
	TimeIntervalSeconds time.Duration
	CallCountLimit      int

	inMemoryStorage InMemoryStorage
}

func New(inMemoryStorage InMemoryStorage, config *Config) (*FloodControl, error) {
	err := validateConfig(config)
	if err != nil {
		return nil, err
	}

	return &FloodControl{
		TimeIntervalSeconds: config.TimeIntervalSeconds,
		CallCountLimit:      config.CallCountLimit,
		inMemoryStorage:     inMemoryStorage,
	}, nil
}

func (f *FloodControl) Check(ctx context.Context, userID int64) (bool, error) {
	if f.CallCountLimit == 0 {
		return false, ErrCheckCallLimitExceed
	}

	floodControlData, err := f.getCachedFloodControlData(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrRedisKeyNotFound) {
			err = f.cacheFloodControlData(ctx, userID, time.Now(), 1)
			if err != nil {
				return false, fmt.Errorf("failed to check: %w", err)
			}

			return true, nil
		}

		return false, fmt.Errorf("failed to check: %w", err)
	}

	if f.intervalHasExpired(floodControlData.LastCallAt) {
		err = f.cacheFloodControlData(ctx, userID, time.Now(), 1)
		if err != nil {
			return false, fmt.Errorf("failed to check: %w", err)
		}
		return true, nil
	}

	newCallCount := floodControlData.CallCount + 1

	if f.callCountLimitHasExceeded(newCallCount) {
		return false, ErrCheckCallLimitExceed
	}

	err = f.cacheFloodControlData(ctx, userID, floodControlData.LastCallAt, newCallCount)
	if err != nil {
		return false, fmt.Errorf("failed to check: %w", err)
	}

	return true, nil
}

func (f *FloodControl) cacheFloodControlData(ctx context.Context, userId int64, lastCallAt time.Time, callCount int) error {
	payloadBytes, err := json.Marshal(&Payload{
		LastCallAt: lastCallAt,
		CallCount:  callCount,
	})
	if err != nil {
		return fmt.Errorf("failed to cache flood control data: %w", err)
	}
	err = f.inMemoryStorage.Set(ctx, strconv.Itoa(int(userId)), string(payloadBytes), f.TimeIntervalSeconds)
	if err != nil {
		return fmt.Errorf("failed to cache flood control data: %w", err)
	}
	return nil
}

func (f *FloodControl) getCachedFloodControlData(ctx context.Context, userId int64) (*Payload, error) {
	payloadStr, err := f.inMemoryStorage.Get(ctx, strconv.Itoa(int(userId)))
	if err != nil {
		return nil, fmt.Errorf("failed to get cached flood control data: %w", err)
	}
	var payload Payload
	err = json.Unmarshal([]byte(payloadStr), &payload)
	if err != nil {
		return nil, fmt.Errorf("failed to get cached flood control data: %w", err)
	}
	return &payload, nil
}

func (f *FloodControl) intervalHasExpired(lastCallAt time.Time) bool {
	return time.Now().After(lastCallAt.Add(f.TimeIntervalSeconds))
}

func (f *FloodControl) callCountLimitHasExceeded(callCount int) bool {
	return callCount > f.CallCountLimit
}

func validateConfig(cfg *Config) error {
	if cfg.CallCountLimit < 0 {
		return errors.New("'CallCountLimit' parameter must be greater or equal to 0")
	}
	if cfg.TimeIntervalSeconds < 0 {
		return errors.New("'TimeIntervalSeconds' parameter must be greater or equal to 0")
	}
	return nil
}
