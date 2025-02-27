package ipratelimit

import (
	"time"

	"golang.org/x/time/rate"
)

type Option func(*IPRateLimiter)

// WithLimit sets the rate limit
func WithLimit(limit rate.Limit) Option {
	return func(l *IPRateLimiter) {
		l.limit = limit
	}
}

// WithBurst sets the burst limit
func WithBurst(burst int) Option {
	return func(l *IPRateLimiter) {
		l.burst = burst
	}
}

// WithExpiry sets the expiry duration
func WithExpiry(expiry time.Duration) Option {
	return func(l *IPRateLimiter) {
		l.expiry = expiry
	}
}

// WithCleanupInterval sets the cleanup interval
func WithCleanupInterval(cleanupInterval time.Duration) Option {
	return func(l *IPRateLimiter) {
		l.cleanupInterval = cleanupInterval
	}
}
