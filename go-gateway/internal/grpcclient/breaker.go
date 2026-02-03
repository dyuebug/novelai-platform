package grpcclient

import (
	"time"

	"github.com/sony/gobreaker"
)

// NewBreaker 创建熔断器
func NewBreaker(name string) *gobreaker.CircuitBreaker {
	settings := gobreaker.Settings{
		Name:        name,
		MaxRequests: 3,                // 半开状态允许的请求数
		Interval:    10 * time.Second, // 统计周期
		Timeout:     30 * time.Second, // 熔断后恢复时间
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			// 连续失败 5 次或失败率超过 60% 时熔断
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.ConsecutiveFailures >= 5 || (counts.Requests >= 10 && failureRatio >= 0.6)
		},
	}
	return gobreaker.NewCircuitBreaker(settings)
}
