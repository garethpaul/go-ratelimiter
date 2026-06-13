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

func TestBuildKeysFallsBackAfterMalformedRemoteAddr(t *testing.T) {
	limiter := NewLimiter(10, time.Minute)
	limiter.IPLookups = []string{"RemoteAddr", "X-Real-IP"}
	request := httptest.NewRequest(http.MethodGet, "/limited", nil)
	request.RemoteAddr = "not-an-ip"
	request.Header.Set("X-Real-IP", "203.0.113.9")

	keys := BuildKeys(limiter, request)

	if len(keys) != 1 {
		t.Fatalf("expected one key set, got %d", len(keys))
	}
	if got, want := keys[0], []string{"203.0.113.9", "/limited"}; !equalStrings(got, want) {
		t.Fatalf("keys = %#v, want %#v", got, want)
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

func TestBuildKeysHeaderValueMatchSkipsBlankFirstRequestValue(t *testing.T) {
	limiter := NewLimiter(10, time.Minute)
	limiter.Headers = map[string][]string{"X-Plan": {"gold"}}
	request := httptest.NewRequest(http.MethodGet, "/limited", nil)
	request.RemoteAddr = "203.0.113.10:54321"
	request.Header.Add("X-Plan", "")
	request.Header.Add("X-Plan", "gold")

	keys := BuildKeys(limiter, request)

	if len(keys) != 1 {
		t.Fatalf("expected one key set, got %d", len(keys))
	}
	if got, want := keys[0], []string{"203.0.113.10", "/limited", "X-Plan", "gold"}; !equalStrings(got, want) {
		t.Fatalf("keys = %#v, want %#v", got, want)
	}
}

func TestBuildKeysHeaderValueMatchSkipsBlankConfiguredValue(t *testing.T) {
	limiter := NewLimiter(10, time.Minute)
	limiter.Headers = map[string][]string{"X-Plan": {"", "gold"}}
	request := httptest.NewRequest(http.MethodGet, "/limited", nil)
	request.RemoteAddr = "203.0.113.10:54321"
	request.Header.Set("X-Plan", "")

	keys := BuildKeys(limiter, request)

	if len(keys) != 0 {
		t.Fatalf("keys = %#v, want no keys for blank configured header value", keys)
	}
}

func TestBuildKeysDeduplicatesConfiguredHeaderValues(t *testing.T) {
	limiter := NewLimiter(10, time.Minute)
	limiter.Headers = map[string][]string{"X-Plan": {"gold", "gold"}}
	request := httptest.NewRequest(http.MethodGet, "/limited", nil)
	request.RemoteAddr = "203.0.113.10:54321"
	request.Header.Set("X-Plan", "gold")

	keys := BuildKeys(limiter, request)

	if len(keys) != 1 {
		t.Fatalf("expected one key set for duplicate configured values, got %d: %#v", len(keys), keys)
	}
	if got, want := keys[0], []string{"203.0.113.10", "/limited", "X-Plan", "gold"}; !equalStrings(got, want) {
		t.Fatalf("keys = %#v, want %#v", got, want)
	}
}

func TestBuildKeysHeaderOnlySkipsBlankRequestValue(t *testing.T) {
	limiter := NewLimiter(10, time.Minute)
	limiter.Headers = map[string][]string{"X-Plan": nil}
	request := httptest.NewRequest(http.MethodGet, "/limited", nil)
	request.RemoteAddr = "203.0.113.10:54321"
	request.Header.Set("X-Plan", " ")

	keys := BuildKeys(limiter, request)

	if len(keys) != 0 {
		t.Fatalf("keys = %#v, want no keys for blank request header value", keys)
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

func TestLimitHandlerChargesDuplicateConfiguredHeaderValueOnce(t *testing.T) {
	limiter := NewLimiter(1, time.Hour)
	limiter.Headers = map[string][]string{"X-Plan": {"gold", "gold"}}
	handler := LimitFuncHandler(limiter, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	request := httptest.NewRequest(http.MethodGet, "/limited", nil)
	request.RemoteAddr = "203.0.113.10:54321"
	request.Header.Set("X-Plan", "gold")

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
}

func TestSetResponseHeadersReplacesExistingValues(t *testing.T) {
	limiter := NewLimiter(10, time.Minute)
	recorder := httptest.NewRecorder()
	recorder.Header().Add("X-Rate-Limit-Limit", "stale")
	recorder.Header().Add("X-Rate-Limit-Duration", "stale")

	SetResponseHeaders(limiter, recorder)
	SetResponseHeaders(limiter, recorder)

	if got, want := recorder.Header().Values("X-Rate-Limit-Limit"), []string{"10"}; !equalStrings(got, want) {
		t.Fatalf("limit header values = %#v, want %#v", got, want)
	}
	if got, want := recorder.Header().Values("X-Rate-Limit-Duration"), []string{"1m0s"}; !equalStrings(got, want) {
		t.Fatalf("duration header values = %#v, want %#v", got, want)
	}
}

func TestLimitHandlerReplacesExistingRejectionContentType(t *testing.T) {
	limiter := NewLimiter(1, time.Hour)
	limiter.MessageContentType = "application/problem+json"
	handler := LimitFuncHandler(limiter, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	request := httptest.NewRequest(http.MethodGet, "/limited", nil)
	request.RemoteAddr = "203.0.113.10:54321"

	handler.ServeHTTP(httptest.NewRecorder(), request)

	recorder := httptest.NewRecorder()
	recorder.Header().Add("Content-Type", "text/html")
	handler.ServeHTTP(recorder, request)

	if got, want := recorder.Header().Values("Content-Type"), []string{"application/problem+json"}; !equalStrings(got, want) {
		t.Fatalf("content type values = %#v, want %#v", got, want)
	}
}

func TestLimitByKeysKeepsDelimitedComponentsDistinct(t *testing.T) {
	limiter := NewLimiter(1, time.Hour)

	if err := LimitByKeys(limiter, []string{"a|b", "c"}); err != nil {
		t.Fatalf("first key set was limited: %v", err)
	}
	if err := LimitByKeys(limiter, []string{"a", "b|c"}); err != nil {
		t.Fatalf("distinct key set collided: %v", err)
	}
	if err := LimitByKeys(limiter, []string{"a|b", "c"}); err == nil {
		t.Fatal("repeated key set did not reach the limit")
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
