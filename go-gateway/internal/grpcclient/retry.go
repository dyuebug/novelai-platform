package grpcclient

import "time"

// RetryPolicy 重试策略
type RetryPolicy struct {
	MaxRetries int
	Backoff    time.Duration
}

func (r RetryPolicy) nextBackoff(attempt int) time.Duration {
	if r.Backoff <= 0 {
		return 0
	}
	// 指数退避
	return time.Duration(attempt) * r.Backoff
}
