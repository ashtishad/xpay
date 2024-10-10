package middlewares

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// IPRateLimiter manages rate limiting for multiple IP addresses.
type IPRateLimiter struct {
	ips   map[string]*rate.Limiter
	mu    *sync.RWMutex
	rps   rate.Limit
	burst int
}

// NewIPRateLimiter creates a new IPRateLimiter with specified rate and burst.
// rps: rate limit (requests per second)
// burst: burst size (maximum number of requests allowed to exceed the rate)
func NewIPRateLimiter(rps rate.Limit, burst int) *IPRateLimiter {
	return &IPRateLimiter{
		ips:   make(map[string]*rate.Limiter),
		mu:    &sync.RWMutex{},
		rps:   rps,
		burst: burst,
	}
}

// GetLimiter retrieves or creates a rate limiter for a given IP address.
// It ensures thread-safe access to the map of limiters.
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, ok := i.ips[ip]
	if !ok {
		limiter = rate.NewLimiter(i.rps, i.burst)
		i.ips[ip] = limiter
	}

	return limiter
}

// RateLimiter creates a Gin middleware that applies client IP-based rate limiting using Token bucket algorithm.
// It uses the provided IPRateLimiter to manage rate limits for each IP.
func RateLimiter(limiter *IPRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		l := limiter.GetLimiter(ip)

		if !l.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
			})

			c.Abort()
			return
		}

		c.Next()
	}
}
