package database

import (
	"context"
	"time"
	"alienx/lib/utility"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
)

var clientLocal *redis.Client

//ConnectToRedis need call before use redis
func ConnectToRedis(options ...int) error {
	indexDB := 0
	if len(options) >= 1 {
		indexDB = options[0]
	}

	redisHost := utility.GetEnv("MAIN_REDIS_HOST", "127.0.0.1")
	redisPort := utility.GetEnv("MAIN_REDIS_PORT", "6379")

	redisAddressLocal := redisHost + ":" + redisPort
	redisPWLocal := utility.GetEnv("MAIN_REDIS_PASS", "")
	clientLocal = redis.NewClient(&redis.Options{
		Addr:         redisAddressLocal,
		Password:     redisPWLocal,
		DB:           indexDB,
		PoolSize:     12,
		MinIdleConns: 2,
		ReadTimeout:  1 * time.Second,
		PoolTimeout:  3 * time.Second,
	})

	_, errLocal := clientLocal.Ping(context.Background()).Result()
	return errLocal
}

var miniRedisServer *miniredis.Miniredis

// ConnectToRedisTest need call before use redis
func ConnectToRedisTest(options ...int) error {

	miniRedisServer, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	indexDB := 0
	if len(options) >= 1 {
		indexDB = options[0]
	}

	clientLocal = redis.NewClient(&redis.Options{
		Addr:         miniRedisServer.Addr(),
		DB:           indexDB,
		PoolSize:     12,
		MinIdleConns: 2,
		PoolTimeout:  3 * time.Second,
	})

	_, errLocal := clientLocal.Ping(context.Background()).Result()
	return errLocal
}

//GetRedisConn ...
func GetRedisConn() *redis.Client {
	return clientLocal
}

//DisconnectFromRedis call when server close
func DisconnectFromRedis() error {
	// MiniRedis is fake redis server for testing
	if miniRedisServer != nil {
		miniRedisServer.Close()
	}
	return clientLocal.Close()
}

