package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"go-gateway/internal/config"
	"go-gateway/pkg/response"
	"golang.org/x/time/rate"
)

type limiterStore struct {
	mu       sync.RWMutex
	limiters map[string]*rate.Limiter
}

var store = &limiterStore{limiters: make(map[string]*rate.Limiter)}

func RateLimit(cfg config.RateLimitConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()
		limiter := getLimiter(key, cfg)
		if !limiter.Allow() {
			response.Error(c, http.StatusTooManyRequests, "rate limit exceeded")
			return
		}
		c.Next()
	}
}

func getLimiter(key string, cfg config.RateLimitConfig) *rate.Limiter {
	store.mu.RLock()
	if l, ok := store.limiters[key]; ok {
		store.mu.RUnlock()
		return l
	}
	store.mu.RUnlock()

	store.mu.Lock()
	defer store.mu.Unlock()

	// 双重检查
	if l, ok := store.limiters[key]; ok {
		return l
	}

	limiter := rate.NewLimiter(rate.Limit(cfg.RPS), cfg.Burst)
	store.limiters[key] = limiter
	return limiter
}
