package constant

import "time"

const DefaultTimeOut = 3 * time.Second

const DefaultExpiration = 30 * time.Minute

const NeverExpire = 0

const (
	SafeGetDistributionLock = "safe_get_distributed_lock_get:"

	// GotoShortLinkKey 短链接跳转前缀 Key
	GotoShortLinkKey = "short-link:goto:"

	// GotoIsNullShortLinkKey 短链接空值跳转前缀 Key
	//
	// 在某次查询中，当发现短链接不存在或者已经失效，就会将这个 Key 写入缓存，从而避免查询数据库
	GotoIsNullShortLinkKey = "short-link:is-null:goto:"

	// LockGotoShortLinkKey 短链接跳转锁前缀 Key
	LockGotoShortLinkKey = "short-link:lock:goto:"

	// LockGidUpdateKey 短链接修改分组 ID 锁前缀 Key
	LockGidUpdateKey = "short-link:lock:update-gid"

	// DelayQueueStatKey 短链接延迟队列消费统计 Key
	DelayQueueStatKey = "short-link:delay-queue:stats"

	// ShortLinkStatUvKey 短链接统计判断是否新用户缓存标识
	ShortLinkStatUvKey = "short-link:stats:uv:"

	// ShortLinkStatUipKey 短链接统计判断是否新 IP 缓存标识
	ShortLinkStatUipKey = "short-link:stats:uip:"

	// ShortLinkStatStreamTopicKey 短链接监控消息保存队列 Topic 缓存标识
	ShortLinkStatStreamTopicKey = "short-link:stats-stream"

	// ShortLinkStatStreamGroupKey 短链接监控消息保存队列 Group 缓存标识
	ShortLinkStatStreamGroupKey = "short-link:stats-stream:only-group"

	// ShortLinkCreateLockKey 创建短链接锁标识
	ShortLinkCreateLockKey = "short-link:lock:create:%s"
)
