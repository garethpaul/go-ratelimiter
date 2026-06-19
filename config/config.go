package config

import (
	"container/list"
	"crypto/sha256"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const defaultMaxTrackedKeys = 10000

// NewLimiter is a constructor for Limiter.
func NewLimiter(max int64, ttl time.Duration) *Limiter {
	limiter := &Limiter{Max: max, TTL: ttl}
	limiter.MessageContentType = "text/plain; charset=utf-8"
	limiter.Message = "You have reached the maximum request limit for this tool"
	limiter.StatusCode = 429
	limiter.tokenBuckets = make(map[string]*rate.Limiter)
	limiter.tokenBucketOrder = list.New()
	limiter.tokenBucketEntries = make(map[string]*list.Element)
	limiter.maxTrackedKeys = defaultMaxTrackedKeys
	limiter.IPLookups = []string{"RemoteAddr", "X-Forwarded-For", "X-Real-IP"}

	return limiter
}

// Limiter is a config struct to limit a particular request handler.
type Limiter struct {
	// HTTP message when limit is reached.
	Message string

	// Content-Type for Message
	MessageContentType string

	// HTTP status code when limit is reached.
	StatusCode int

	// Maximum burst size and number of tokens refilled per TTL.
	Max int64

	// Duration over which Max tokens refill.
	TTL time.Duration

	// List of places to look up IP address.
	// Default is "RemoteAddr", "X-Forwarded-For", "X-Real-IP".
	// You can rearrange the order as you like.
	IPLookups []string

	// List of HTTP Methods to limit (GET, POST, PUT, etc.).
	// Empty means limit all methods.
	Methods []string

	// List of HTTP headers to limit.
	// Empty means skip headers checking.
	Headers map[string][]string

	// List of basic auth usernames to limit.
	BasicAuthUsers []string

	// Throttler struct
	tokenBuckets map[string]*rate.Limiter

	// LRU bookkeeping bounds request-controlled rate-limiter keys.
	tokenBucketOrder   *list.List
	tokenBucketEntries map[string]*list.Element
	maxTrackedKeys     int

	sync.RWMutex
}

// LimitReached returns a bool indicating if the Bucket identified by key ran out of tokens.
func (l *Limiter) LimitReached(key string) bool {
	return l.LimitReachedForKeys([]string{key})
}

// LimitReachedForKeys atomically checks and consumes one token from each distinct key.
func (l *Limiter) LimitReachedForKeys(keys []string) bool {
	l.Lock()
	defer l.Unlock()
	if len(keys) == 0 {
		return false
	}
	if l.Max <= 0 || l.TTL <= 0 || uint64(l.Max) > uint64(^uint(0)>>1) {
		return true
	}

	storageKeys := make([]string, 0, len(keys))
	seenStorageKeys := make(map[string]struct{}, len(keys))
	for _, key := range keys {
		storageKey := bucketStorageKey(key)
		if _, seen := seenStorageKeys[storageKey]; seen {
			continue
		}
		seenStorageKeys[storageKey] = struct{}{}
		storageKeys = append(storageKeys, storageKey)
	}
	if l.maxTrackedKeys > 0 && len(storageKeys) > l.maxTrackedKeys {
		return true
	}

	buckets := make([]*rate.Limiter, 0, len(storageKeys))
	for _, storageKey := range storageKeys {
		bucket := l.bucketForStorageKey(storageKey)
		buckets = append(buckets, bucket)
	}
	now := time.Now()
	for _, bucket := range buckets {
		if bucket.TokensAt(now) < 1 {
			return true
		}
	}
	for _, bucket := range buckets {
		_ = bucket.AllowN(now, 1)
	}

	return false
}

func (l *Limiter) bucketForStorageKey(storageKey string) *rate.Limiter {
	bucket, found := l.tokenBuckets[storageKey]
	if found {
		l.tokenBucketOrder.MoveToFront(l.tokenBucketEntries[storageKey])
	} else {
		if l.maxTrackedKeys > 0 && len(l.tokenBuckets) >= l.maxTrackedKeys {
			oldest := l.tokenBucketOrder.Back()
			oldestKey := oldest.Value.(string)
			delete(l.tokenBuckets, oldestKey)
			delete(l.tokenBucketEntries, oldestKey)
			l.tokenBucketOrder.Remove(oldest)
		}

		refillRate := rate.Limit(float64(l.Max) / l.TTL.Seconds())
		bucket = rate.NewLimiter(refillRate, int(l.Max))
		l.tokenBuckets[storageKey] = bucket
		l.tokenBucketEntries[storageKey] = l.tokenBucketOrder.PushFront(storageKey)
	}

	return bucket
}

func bucketStorageKey(key string) string {
	digest := sha256.Sum256([]byte(key))
	return string(digest[:])
}
