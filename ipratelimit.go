package ipratelimit

import (
	"context"
	"net/http"
	"time"

	"github.com/mpfrancis/safemap"
	"golang.org/x/time/rate"
)

type IPRateLimiter struct {
	visitors        safemap.SafeMap[string, *rateLimiter]
	limit           rate.Limit
	burst           int
	expiry          time.Duration
	cleanupInterval time.Duration
}

type rateLimiter struct {
	*rate.Limiter
	lastUsed time.Time
}

func New(ctx context.Context, options ...Option) *IPRateLimiter {
	l := &IPRateLimiter{
		visitors:        safemap.New[string, *rateLimiter](),
		limit:           rate.Every(time.Minute / 10),
		burst:           10,
		expiry:          5 * time.Minute,
		cleanupInterval: 10 * time.Minute,
	}

	for _, option := range options {
		option(l)
	}

	go l.cleanup(ctx)

	return l
}

func (l *IPRateLimiter) cleanup(ctx context.Context) {
	ticker := time.NewTicker(l.cleanupInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			l.visitors.Range(func(key string, value *rateLimiter) bool {
				if time.Since(value.lastUsed) > l.expiry {
					l.visitors.Delete(key)
				}

				return true
			})
		}
	}
}

func (l *IPRateLimiter) AllowIP(ip string) bool {
	return l.getLimiter(ip).Allow()
}

func (l *IPRateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !l.AllowIP(r.RemoteAddr) {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (l *IPRateLimiter) getLimiter(ip string) *rateLimiter {
	if limiter, ok := l.visitors.Get(ip); ok {
		limiter.lastUsed = time.Now()
		return limiter
	}

	// Allow 5 requests per minute per IP
	limiter := rateLimiter{
		Limiter:  rate.NewLimiter(l.limit, l.burst),
		lastUsed: time.Now(),
	}
	l.visitors.Set(ip, &limiter)

	return &limiter
}
