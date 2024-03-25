package main

import (
	"errors"
	"os"
	"task/internal/domain/floodcontrol"
	"task/pkg/config"
	"task/pkg/database/redis"
)

type AppConfig struct {
	FloodControlCfg floodcontrol.Config
	Redis           redis.Config
}

func mustGetAppConfig(sources ...string) AppConfig {
	config, err := tryGetAppConfig(sources...)
	if err != nil {
		panic(err)
	}

	if config == nil {
		panic(errors.New("config cannot be nil"))
	}

	return *config
}

func tryGetAppConfig(sources ...string) (*AppConfig, error) {
	if len(sources) == 0 {
		return nil, errors.New("at least 1 source must be set for app config")
	}

	provider := config.MergeConfigs(sources)

	var floodControlConfig floodcontrol.Config
	err := config.ParseConfig(provider, "app.flood-control", &floodControlConfig)
	if err != nil {
		return nil, err
	}

	var redisConfig redis.Config
	err = config.ParseConfig(provider, "app.redis", &redisConfig)
	if err != nil {
		return nil, err
	}
	redisConfig.User = os.Getenv("REDIS_USER")
	redisConfig.Password = os.Getenv("REDIS_PASSWORD")

	appConfig := AppConfig{
		FloodControlCfg: floodControlConfig,
		Redis:           redisConfig,
	}

	return &appConfig, nil
}
