package limiter

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestBuildKeysDefaultUsesRemoteIPAndPath(t *testing.T) {
	limiter := NewLimiter(10, time.Minute)
	request := httptest.NewRequest(http.MethodGet, "/api/items?ignored=true", nil)
	request.RemoteAddr = "203.0.113.10:54321"

	keys := BuildKeys(limiter, request)

	if len(keys) != 1 {
		t.Fatalf("expected one key set, got %d", len(keys))
	}
	if got, want := keys[0], []string{"203.0.113.10", "/api/items"}; !equalStrings(got, want) {
		t.Fatalf("keys = %#v, want %#v", got, want)
	}
}

func TestBuildKeysCanPreferForwardedFor(t *testing.T) {
	limiter := NewLimiter(10, time.Minute)
	limiter.IPLookups = []string{"X-Forwarded-For", "RemoteAddr"}
	request := httptest.NewRequest(http.MethodGet, "/limited", nil)
	request.RemoteAddr = "10.0.0.5:54321"
	request.Header.Set("X-Forwarded-For", "198.51.100.7, 10.0.0.5")

	keys := BuildKeys(limiter, request)

	if len(keys) != 1 {
		t.Fatalf("expected one key set, got %d", len(keys))
	}
	if got, want := keys[0], []string{"198.51.100.7", "/limited"}; !equalStrings(got, want) {
		t.Fatalf("keys = %#v, want %#v", got, want)
	}
}

func TestBuildKeysSkipsMalformedRemoteAddr(t *testing.T) {
	limiter := NewLimiter(10, time.Minute)
	request := httptest.NewRequest(http.MethodGet, "/limited", nil)
	request.RemoteAddr = "not-an-ip"

	keys := BuildKeys(limiter, request)

	if len(keys) != 0 {
		t.Fatalf("keys = %#v, want no keys for malformed RemoteAddr", keys)
	}
}

func TestBuildKeysHeaderValuesRequireConfiguredMatch(t *testing.T) {
	limiter := NewLimiter(10, time.Minute)
	limiter.Headers = map[string][]string{"X-Plan": {"gold"}}
	request := httptest.NewRequest(http.MethodGet, "/limited", nil)
	request.RemoteAddr = "203.0.113.10:54321"
	request.Header.Set("X-Plan", "bronze")

	keys := BuildKeys(limiter, request)

	if len(keys) != 0 {
		t.Fatalf("keys = %#v, want no keys for non-matching header value", keys)
	}
}

func TestBuildKeysHeaderValueMatchIncludesConfiguredValue(t *testing.T) {
	limiter := NewLimiter(10, time.Minute)
	limiter.Headers = map[string][]string{"X-Plan": {"gold"}}
	request := httptest.NewRequest(http.MethodGet, "/limited", nil)
	request.RemoteAddr = "203.0.113.10:54321"
	request.Header.Set("X-Plan", "gold")

	keys := BuildKeys(limiter, request)

	if len(keys) != 1 {
		t.Fatalf("expected one key set, got %d", len(keys))
	}
	if got, want := keys[0], []string{"203.0.113.10", "/limited", "X-Plan", "gold"}; !equalStrings(got, want) {
		t.Fatalf("keys = %#v, want %#v", got, want)
	}
}

func TestBuildKeysMethodHeaderValueMatchIncludesConfiguredValue(t *testing.T) {
	limiter := NewLimiter(10, time.Minute)
	limiter.Methods = []string{http.MethodPost}
	limiter.Headers = map[string][]string{"X-Plan": {"gold"}}
	request := httptest.NewRequest(http.MethodPost, "/limited", nil)
	request.RemoteAddr = "203.0.113.10:54321"
	request.Header.Set("X-Plan", "gold")

	keys := BuildKeys(limiter, request)

	if len(keys) != 1 {
		t.Fatalf("expected one key set, got %d", len(keys))
	}
	if got, want := keys[0], []string{"203.0.113.10", "/limited", http.MethodPost, "X-Plan", "gold"}; !equalStrings(got, want) {
		t.Fatalf("keys = %#v, want %#v", got, want)
	}
}

func TestLimitFuncHandlerReturnsTooManyRequestsAfterBucketIsEmpty(t *testing.T) {
	limiter := NewLimiter(1, time.Hour)
	handler := LimitFuncHandler(limiter, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	request := httptest.NewRequest(http.MethodGet, "/limited", nil)
	request.RemoteAddr = "203.0.113.10:54321"

	first := httptest.NewRecorder()
	handler.ServeHTTP(first, request)
	if first.Code != http.StatusNoContent {
		t.Fatalf("first response status = %d, want %d", first.Code, http.StatusNoContent)
	}

	second := httptest.NewRecorder()
	handler.ServeHTTP(second, request)
	if second.Code != http.StatusTooManyRequests {
		t.Fatalf("second response status = %d, want %d", second.Code, http.StatusTooManyRequests)
	}
	if second.Header().Get("X-Rate-Limit-Limit") != "1" {
		t.Fatalf("missing rate limit header: %#v", second.Header())
	}
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
