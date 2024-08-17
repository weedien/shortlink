package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"shortlink/config"
)

var Rdb *redis.Client

var ctx = context.Background()

func ConnectToRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr.String(),
		Password: config.RedisPassword.String(),
		DB:       config.RedisDB.Int(),
	})

	// 检测连接是否成功
	_, err := Rdb.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Errorf("failed to connect to redis: %v", err))
	}

	// 初始化布隆过滤器
	InitBloomFilter()
}
