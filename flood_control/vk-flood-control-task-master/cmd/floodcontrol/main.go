package main

import (
	"context"
	"flag"
	"github.com/joho/godotenv"
	"log"
	"math/rand"
	"task/internal/domain/floodcontrol"
	redisstorage "task/internal/infra/storage/redis"
	"task/pkg/database/redis"
	"time"
)

var appConfigPaths = flag.String("config.files", "not-found.yaml", "List of application config files separated by comma")
var environmentVariablesPath = flag.String("env.vars.file", "not-found.env", "Path to environment variables file")

func main() {
	flag.Parse()
	if err := loadEnvironmentVariables(); err != nil {
		panic(err)
	}

	appConfig := mustGetAppConfig(*appConfigPaths)

	pgManagedDatabase, err := redis.New(&appConfig.Redis)
	if err != nil {
		panic(err)
	}

	redisStorage := redisstorage.NewRedisStorage(pgManagedDatabase.RedisDb)
	floodControl, err := floodcontrol.New(redisStorage, &appConfig.FloodControlCfg)
	if err != nil {
		panic(err)
	}

	runFloodControlCheck(floodControl)
}

func runFloodControlCheck(floodControl FloodControl) {
	log.Println("running floodControlCheck")

	ctx := context.Background()
	userId := int64(1)
	var fakeRequestIntervalFromUser time.Duration
	fakeRequestCountFromUser := 20

	for i := 0; i < fakeRequestCountFromUser; i++ {
		passedCheckControl, err := floodControl.Check(ctx, userId)
		log.Printf("i=%d user_id=%d passedCheckControl=%v err=%v", i, userId, passedCheckControl, err)

		fakeRequestIntervalFromUser = time.Duration(1+rand.Intn(3)) * time.Second

		time.Sleep(fakeRequestIntervalFromUser)
	}
}

func loadEnvironmentVariables() error {
	return godotenv.Load(*environmentVariablesPath)
}

// FloodControl интерфейс, который нужно реализовать.
// Рекомендуем создать директорию-пакет, в которой будет находиться реализация.
type FloodControl interface {
	// Check возвращает false если достигнут лимит максимально разрешенного
	// кол-ва запросов согласно заданным правилам флуд контроля.
	Check(ctx context.Context, userID int64) (bool, error)
}
