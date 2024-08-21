package idem

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"time"
)

const IdempotencyKeyPrefix = "%s:idempotency:"

// MessageQueueIdempotencyHandler 消息队列幂等处理器
type MessageQueueIdempotencyHandler struct {
	rdb       *redis.Client
	keyPrefix string
}

func NewMessageQueueIdempotencyHandler(appName string, rdb *redis.Client) MessageQueueIdempotencyHandler {
	return MessageQueueIdempotencyHandler{
		rdb:       rdb,
		keyPrefix: fmt.Sprintf(IdempotencyKeyPrefix, appName),
	}
}

// IsMessageBeingConsumed 消息是否正在被消费
func (h MessageQueueIdempotencyHandler) IsMessageBeingConsumed(mid string) bool {
	key := h.keyPrefix + mid
	ok, err := h.rdb.SetNX(context.Background(), key, "0", 2*time.Second).Result()
	if err != nil {
		return true
	}
	if ok {
		return false
	}
	return true
}

// HasMessageBeenConsumed 消息是否已经被消费
func (h MessageQueueIdempotencyHandler) HasMessageBeenConsumed(mid string) bool {
	key := h.keyPrefix + mid
	if val := h.rdb.Get(context.Background(), key).String(); val == "1" {
		return true
	}
	return false
}

// MarkMessageAsConsumed 标记消息为已消费
func (h MessageQueueIdempotencyHandler) MarkMessageAsConsumed(mid string) {
	key := h.keyPrefix + mid
	err := h.rdb.Set(context.Background(), key, "1", 0).Err()
	if err != nil {
		slog.Error("mark message as consumed failed", "error", err)
		return
	}
}

// DeleteFlag 如果消息处理遇到异常，删除幂等标记
func (h MessageQueueIdempotencyHandler) DeleteFlag(mid string) {
	key := h.keyPrefix + mid
	err := h.rdb.Del(context.Background(), key).Err()
	if err != nil {
		slog.Error("delete flag failed", "error", err)
		return
	}
}
