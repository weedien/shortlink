package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

const (
	ShortUriCreateBloomFilter = "shortUriCreateBloomFilter"
)

func setUpBloomFilter(rdb *redis.Client) {
	capacity, errorRate := int64(1000_000), 0.0001
	_, err := rdb.BFReserve(context.Background(), ShortUriCreateBloomFilter, errorRate, capacity).Result()
	if err != nil {
		panic(fmt.Errorf("failed to setup bloom filter: %v", err))
	}
}
