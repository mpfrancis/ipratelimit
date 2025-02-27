package ipratelimit

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

func TestNew(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	limiter := New(ctx)

	if limiter.limit != rate.Every(time.Minute/10) {
		t.Errorf("expected limit to be %v, got %v", rate.Every(time.Minute/10), limiter.limit)
	}
	if limiter.burst != 10 {
		t.Errorf("expected burst to be %d, got %d", 5, limiter.burst)
	}
	if limiter.expiry != 5*time.Minute {
		t.Errorf("expected expiry to be %v, got %v", 5*time.Minute, limiter.expiry)
	}
	if limiter.cleanupInterval != 10*time.Minute {
		t.Errorf("expected cleanup interval to be %v, got %v", 10*time.Minute, limiter.cleanupInterval)
	}
}

func TestAllowIP(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	limiter := New(ctx)
	ip := "192.168.1.1"

	for i := 0; i < 10; i++ {
		if !limiter.AllowIP(ip) {
			t.Errorf("expected AllowIP to return true, got false")
		}
	}

	if limiter.AllowIP(ip) {
		t.Errorf("expected AllowIP to return false, got true")
	}
}

func TestMiddleware(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	limiter := New(ctx)
	handler := limiter.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "192.168.1.1"

	for i := 0; i < 10; i++ {
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("expected status code to be %d, got %d", http.StatusOK, rr.Code)
		}
	}

	{
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusTooManyRequests {
			fmt.Println(rr.Body)
			t.Errorf("expected status code to be %d, got %d", http.StatusTooManyRequests, rr.Code)
		}
	}
}

func TestCleanup(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	limiter := New(ctx, WithCleanupInterval(500*time.Millisecond), WithExpiry(250*time.Millisecond))
	ip := "1.1.1.1"
	limiter.getLimiter(ip)
	if _, ok := limiter.visitors.Get(ip); !ok {
		t.Errorf("expected visitor to exist, but it does not")
	}
	time.Sleep(1 * time.Second)
	if _, ok := limiter.visitors.Get(ip); ok {
		t.Errorf("expected visitor to be cleaned up, but it still exists")
	}
}
