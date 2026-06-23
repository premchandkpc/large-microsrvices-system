package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/premchandkpc/large-microsrvices-system/services/api-gateway/internal/config"
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	visitors map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	return &RateLimiter{
		visitors: make(map[string]*rate.Limiter),
		rate:     r,
		burst:    b,
	}
}

func (rl *RateLimiter) GetLimiter(ip string) *rate.Limiter {
	rl.mu.RLock()
	limiter, exists := rl.visitors[ip]
	rl.mu.RUnlock()

	if exists {
		return limiter
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter = rate.NewLimiter(rl.rate, rl.burst)
	rl.visitors[ip] = limiter
	return limiter
}

func RateLimiterMiddleware(cfg *config.Config) gin.HandlerFunc {
	rl := NewRateLimiter(rate.Limit(cfg.RateLimit), cfg.RateLimitBurst)

	go func() {
		for {
			time.Sleep(10 * time.Minute)
			rl.mu.Lock()
			rl.visitors = make(map[string]*rate.Limiter)
			rl.mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := rl.GetLimiter(ip)
		if !limiter.Allow() {
			c.AbortWithStatusJSON(429, gin.H{
				"error":       "too many requests",
				"retry_after": "1s",
			})
			return
		}
		c.Next()
	}
}
