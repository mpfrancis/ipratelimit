package ipratelimit

import (
	"testing"
	"time"

	"golang.org/x/time/rate"
)

func TestWithLimit(t *testing.T) {
	limiter := &IPRateLimiter{}
	option := WithLimit(rate.Limit(10))
	option(limiter)

	if limiter.limit != rate.Limit(10) {
		t.Errorf("expected limit to be %v, got %v", rate.Limit(10), limiter.limit)
	}
}

func TestWithBurst(t *testing.T) {
	limiter := &IPRateLimiter{}
	option := WithBurst(20)
	option(limiter)

	if limiter.burst != 20 {
		t.Errorf("expected burst to be %d, got %d", 20, limiter.burst)
	}
}

func TestWithExpiry(t *testing.T) {
	limiter := &IPRateLimiter{}
	option := WithExpiry(30 * time.Second)
	option(limiter)

	if limiter.expiry != 30*time.Second {
		t.Errorf("expected expiry to be %v, got %v", 30*time.Second, limiter.expiry)
	}
}

func TestWithCleanupInterval(t *testing.T) {
	limiter := &IPRateLimiter{}
	option := WithCleanupInterval(1 * time.Minute)
	option(limiter)

	if limiter.cleanupInterval != 1*time.Minute {
		t.Errorf("expected cleanup interval to be %v, got %v", 1*time.Minute, limiter.cleanupInterval)
	}
}
