package cache

import "fmt"

const (
	ShortUriCreateBloomFilter = "shortUriCreateBloomFilter"
)

func InitBloomFilter() {
	capacity, errorRate := int64(1000_000), 0.0001
	_, err := Rdb.BFReserve(ctx, ShortUriCreateBloomFilter, errorRate, capacity).Result()
	if err != nil {
		panic(fmt.Errorf("failed to init bloom filter: %v", err))
	}
}
