package config

import (
	"fmt"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

func TestLimiterRefillsMaxTokensPerTTL(t *testing.T) {
	limiter := NewLimiter(60, time.Minute)

	if reached := limiter.LimitReached("key"); reached {
		t.Fatal("first request unexpectedly reached the limit")
	}

	bucket := limiter.tokenBuckets["key"]
	if got, want := bucket.Limit(), rate.Limit(1); got != want {
		t.Fatalf("refill rate = %v tokens/second, want %v", got, want)
	}
	if got, want := bucket.Burst(), 60; got != want {
		t.Fatalf("burst = %d, want %d", got, want)
	}
}

func TestLimiterRejectsInvalidConfigurationWithoutTrackingKeys(t *testing.T) {
	tests := []struct {
		name string
		max  int64
		ttl  time.Duration
	}{
		{name: "zero max", max: 0, ttl: time.Minute},
		{name: "negative max", max: -1, ttl: time.Minute},
		{name: "zero ttl", max: 1, ttl: 0},
		{name: "negative ttl", max: 1, ttl: -time.Minute},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			limiter := NewLimiter(test.max, test.ttl)

			for i := 0; i < 3; i++ {
				if reached := limiter.LimitReached(fmt.Sprintf("key-%d", i)); !reached {
					t.Fatal("invalid limiter configuration allowed a request")
				}
			}

			if got := len(limiter.tokenBuckets); got != 0 {
				t.Fatalf("tracked %d buckets for invalid configuration, want 0", got)
			}
			if got := limiter.tokenBucketOrder.Len(); got != 0 {
				t.Fatalf("LRU order contains %d invalid buckets, want 0", got)
			}
			if got := len(limiter.tokenBucketEntries); got != 0 {
				t.Fatalf("LRU index contains %d invalid buckets, want 0", got)
			}
		})
	}
}

func TestLimiterCapsTrackedKeys(t *testing.T) {
	limiter := NewLimiter(1, time.Hour)
	limiter.maxTrackedKeys = 3

	for i := 0; i < 20; i++ {
		limiter.LimitReached(fmt.Sprintf("key-%d", i))

		if got := len(limiter.tokenBuckets); got > limiter.maxTrackedKeys {
			t.Fatalf("tracked %d keys, want at most %d", got, limiter.maxTrackedKeys)
		}
	}

	if got := limiter.tokenBucketOrder.Len(); got != limiter.maxTrackedKeys {
		t.Fatalf("LRU order contains %d keys, want %d", got, limiter.maxTrackedKeys)
	}
	if got := len(limiter.tokenBucketEntries); got != limiter.maxTrackedKeys {
		t.Fatalf("LRU index contains %d keys, want %d", got, limiter.maxTrackedKeys)
	}
}

func TestLimiterEvictsLeastRecentlyUsedKey(t *testing.T) {
	limiter := NewLimiter(1, time.Hour)
	limiter.maxTrackedKeys = 2

	limiter.LimitReached("oldest")
	limiter.LimitReached("recent")
	limiter.LimitReached("oldest")
	limiter.LimitReached("new")

	if _, found := limiter.tokenBuckets["recent"]; found {
		t.Fatal("least recently used key was retained")
	}
	if _, found := limiter.tokenBuckets["oldest"]; !found {
		t.Fatal("recently accessed key was evicted")
	}
	if reached := limiter.LimitReached("recent"); reached {
		t.Fatal("evicted key did not receive a fresh token bucket")
	}
}
