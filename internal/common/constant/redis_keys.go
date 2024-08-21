package constant

const (
	// GotoShortLinkKey 短链接跳转前缀 Key
	GotoShortLinkKey = "short-link:goto:%s"

	// GotoIsNullShortLinkKey 短链接空值跳转前缀 Key
	GotoIsNullShortLinkKey = "short-link:is-null:goto:%s"

	// LockGotoShortLinkKey 短链接跳转锁前缀 Key
	LockGotoShortLinkKey = "short-link:lock:goto:%s"

	// LockGidUpdateKey 短链接修改分组 ID 锁前缀 Key
	LockGidUpdateKey = "short-link:lock:update-gid:%s"

	// DelayQueueStatsKey 短链接延迟队列消费统计 Key
	DelayQueueStatsKey = "short-link:delay-queue:stats"

	// ShortLinkStatsUvKey 短链接统计判断是否新用户缓存标识
	ShortLinkStatsUvKey = "short-link:stats:uv:"

	// ShortLinkStatsUipKey 短链接统计判断是否新 IP 缓存标识
	ShortLinkStatsUipKey = "short-link:stats:uip:"

	// ShortLinkStatsStreamTopicKey 短链接监控消息保存队列 Topic 缓存标识
	ShortLinkStatsStreamTopicKey = "short-link:stats-stream"

	// ShortLinkStatsStreamGroupKey 短链接监控消息保存队列 Group 缓存标识
	ShortLinkStatsStreamGroupKey = "short-link:stats-stream:only-group"

	// ShortLinkCreateLockKey 创建短链接锁标识
	ShortLinkCreateLockKey = "short-link:lock:create:%s"
)
