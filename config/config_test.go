package config

import (
	"crypto/sha256"
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

	bucket := limiter.tokenBuckets[bucketStorageKey("key")]
	if got, want := bucket.Limit(), rate.Limit(1); got != want {
		t.Fatalf("refill rate = %v tokens/second, want %v", got, want)
	}
	if got, want := bucket.Burst(), 60; got != want {
		t.Fatalf("burst = %d, want %d", got, want)
	}
}

func TestLimiterDoesNotPartiallyConsumeBatchWhenOneKeyIsLimited(t *testing.T) {
	limiter := NewLimiter(1, time.Hour)

	if reached := limiter.LimitReached("exhausted"); reached {
		t.Fatal("first exhausted-key request unexpectedly reached the limit")
	}
	if reached := limiter.LimitReachedForKeys([]string{"available", "exhausted"}); !reached {
		t.Fatal("batch containing an exhausted key was allowed")
	}
	if reached := limiter.LimitReached("available"); reached {
		t.Fatal("rejected batch consumed the available key")
	}
}

func TestLimiterRejectedBatchDoesNotAllocateMissingKeys(t *testing.T) {
	limiter := NewLimiter(1, time.Hour)
	limiter.maxTrackedKeys = 3

	if reached := limiter.LimitReached("exhausted"); reached {
		t.Fatal("first exhausted-key request unexpectedly reached the limit")
	}
	if reached := limiter.LimitReachedForKeys([]string{"new", "exhausted"}); !reached {
		t.Fatal("batch containing an exhausted key was allowed")
	}
	if _, found := limiter.tokenBuckets[bucketStorageKey("new")]; found {
		t.Fatal("rejected batch allocated a missing key")
	}
	if got := len(limiter.tokenBuckets); got != 1 {
		t.Fatalf("tracked %d buckets after rejected batch, want 1", got)
	}
}

func TestLimiterRejectedBatchDoesNotEvictTrackedKeys(t *testing.T) {
	limiter := NewLimiter(1, time.Hour)
	limiter.maxTrackedKeys = 2

	limiter.LimitReached("unrelated")
	limiter.LimitReached("exhausted")

	if reached := limiter.LimitReachedForKeys([]string{"new", "exhausted"}); !reached {
		t.Fatal("batch containing an exhausted key was allowed")
	}
	if _, found := limiter.tokenBuckets[bucketStorageKey("unrelated")]; !found {
		t.Fatal("rejected batch evicted an unrelated key")
	}
	if _, found := limiter.tokenBuckets[bucketStorageKey("new")]; found {
		t.Fatal("rejected batch retained a newly allocated key")
	}
}

func TestLimiterAllowsEmptyBatchWithoutTrackingKeys(t *testing.T) {
	limiter := NewLimiter(0, 0)

	if reached := limiter.LimitReachedForKeys(nil); reached {
		t.Fatal("empty batch unexpectedly reached the limit")
	}
	if got := len(limiter.tokenBuckets); got != 0 {
		t.Fatalf("tracked %d buckets for an empty batch, want 0", got)
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

func TestDirectlyConfiguredLimiterInitializesAccountingState(t *testing.T) {
	limiter := &Limiter{Max: 1, TTL: time.Hour}

	if reached := limiter.LimitReached("client"); reached {
		t.Fatal("first directly configured request unexpectedly reached the limit")
	}
	if reached := limiter.LimitReached("client"); !reached {
		t.Fatal("second directly configured request unexpectedly bypassed the limit")
	}

	if got := len(limiter.tokenBuckets); got != 1 {
		t.Fatalf("tracked %d buckets, want 1", got)
	}
	if got := limiter.tokenBucketOrder.Len(); got != 1 {
		t.Fatalf("LRU order contains %d buckets, want 1", got)
	}
	if got := len(limiter.tokenBucketEntries); got != 1 {
		t.Fatalf("LRU index contains %d buckets, want 1", got)
	}
	if got := limiter.maxTrackedKeys; got != defaultMaxTrackedKeys {
		t.Fatalf("maxTrackedKeys = %d, want %d", got, defaultMaxTrackedKeys)
	}
}

func TestDirectlyConfiguredInvalidLimiterDoesNotInitializeAccountingState(t *testing.T) {
	limiter := &Limiter{Max: 0, TTL: time.Hour}

	if reached := limiter.LimitReached("client"); !reached {
		t.Fatal("invalid directly configured limiter allowed a request")
	}
	if limiter.tokenBuckets != nil || limiter.tokenBucketOrder != nil || limiter.tokenBucketEntries != nil {
		t.Fatal("invalid directly configured limiter initialized private accounting state")
	}
	if limiter.maxTrackedKeys != 0 {
		t.Fatalf("maxTrackedKeys = %d, want uninitialized zero", limiter.maxTrackedKeys)
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

	if _, found := limiter.tokenBuckets[bucketStorageKey("oldest")]; found {
		t.Fatal("least recently used key was retained")
	}
	if _, found := limiter.tokenBuckets[bucketStorageKey("recent")]; !found {
		t.Fatal("most recently accepted key was evicted")
	}
	if reached := limiter.LimitReached("oldest"); reached {
		t.Fatal("evicted key did not receive a fresh token bucket")
	}
}

func TestLimiterStoresBoundedKeyIdentifiers(t *testing.T) {
	limiter := NewLimiter(1, time.Hour)
	requestKey := string(make([]byte, 1<<20))

	if reached := limiter.LimitReached(requestKey); reached {
		t.Fatal("first request unexpectedly reached the limit")
	}

	if got := len(limiter.tokenBuckets); got != 1 {
		t.Fatalf("tracked %d keys, want 1", got)
	}
	for storedKey := range limiter.tokenBuckets {
		if got, want := len(storedKey), sha256.Size; got != want {
			t.Fatalf("stored key length = %d, want %d", got, want)
		}
	}
}
