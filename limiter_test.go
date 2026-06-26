package limiter

import (
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/garethpaul/go-ratelimiter/config"
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

func TestSortedHeaderKeysReturnsOrderedCopy(t *testing.T) {
	headers := map[string][]string{
		"X-Zone": {"west"},
		"X-Plan": {"gold"},
		"X-Mode": {"fast"},
	}
	wantHeaders := map[string][]string{
		"X-Zone": {"west"},
		"X-Plan": {"gold"},
		"X-Mode": {"fast"},
	}

	got := sortedHeaderKeys(headers)
	want := []string{"X-Mode", "X-Plan", "X-Zone"}
	if !equalStrings(got, want) {
		t.Fatalf("sortedHeaderKeys() = %#v, want %#v", got, want)
	}

	if !reflect.DeepEqual(headers, wantHeaders) {
		t.Fatal("sortedHeaderKeys mutated caller-owned configuration")
	}
}

func TestBuildKeysOrdersConfiguredHeadersDeterministically(t *testing.T) {
	limiter := NewLimiter(10, time.Minute)
	limiter.Headers = map[string][]string{
		"X-Zone": {"west"},
		"X-Plan": {"gold"},
		"X-Mode": {"fast"},
	}
	request := httptest.NewRequest(http.MethodGet, "/limited", nil)
	request.RemoteAddr = "203.0.113.10:54321"
	request.Header.Set("X-Zone", "west")
	request.Header.Set("X-Plan", "gold")
	request.Header.Set("X-Mode", "fast")

	got := BuildKeys(limiter, request)
	want := [][]string{
		{"203.0.113.10", "/limited", "X-Mode", "fast"},
		{"203.0.113.10", "/limited", "X-Plan", "gold"},
		{"203.0.113.10", "/limited", "X-Zone", "west"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("BuildKeys() = %#v, want %#v", got, want)
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

func TestLimitHandlerSharesBucketAcrossEquivalentIPv6Forms(t *testing.T) {
	limiter := NewLimiter(1, time.Hour)
	handler := LimitFuncHandler(limiter, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	firstRequest := httptest.NewRequest(http.MethodGet, "/limited", nil)
	firstRequest.RemoteAddr = "[2001:0db8:0:0:0:0:0:1]:54321"
	first := httptest.NewRecorder()
	handler.ServeHTTP(first, firstRequest)
	if first.Code != http.StatusNoContent {
		t.Fatalf("first response status = %d, want %d", first.Code, http.StatusNoContent)
	}

	secondRequest := httptest.NewRequest(http.MethodGet, "/limited", nil)
	secondRequest.RemoteAddr = "[2001:db8::1]:54321"
	second := httptest.NewRecorder()
	handler.ServeHTTP(second, secondRequest)
	if second.Code != http.StatusTooManyRequests {
		t.Fatalf("equivalent IPv6 response status = %d, want %d", second.Code, http.StatusTooManyRequests)
	}
}

func TestLimitHandlerSharesBucketAcrossEquivalentScopedIPv6Forms(t *testing.T) {
	limiter := NewLimiter(1, time.Hour)
	handler := LimitFuncHandler(limiter, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	firstRequest := httptest.NewRequest(http.MethodGet, "/limited", nil)
	firstRequest.RemoteAddr = "[fe80:0:0:0:0:0:0:1%eth0]:54321"
	first := httptest.NewRecorder()
	handler.ServeHTTP(first, firstRequest)
	if first.Code != http.StatusNoContent {
		t.Fatalf("first response status = %d, want %d", first.Code, http.StatusNoContent)
	}

	secondRequest := httptest.NewRequest(http.MethodGet, "/limited", nil)
	secondRequest.RemoteAddr = "[fe80::1%eth0]:54321"
	second := httptest.NewRecorder()
	handler.ServeHTTP(second, secondRequest)
	if second.Code != http.StatusTooManyRequests {
		t.Fatalf("equivalent scoped IPv6 response status = %d, want %d", second.Code, http.StatusTooManyRequests)
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

func TestBuildKeysEmptyConstraintsFallBackToRemoteIPAndPath(t *testing.T) {
	tests := []struct {
		name      string
		configure func(*config.Limiter)
	}{
		{
			name: "methods",
			configure: func(limiter *config.Limiter) {
				limiter.Methods = []string{}
			},
		},
		{
			name: "headers",
			configure: func(limiter *config.Limiter) {
				limiter.Headers = map[string][]string{}
			},
		},
		{
			name: "basic auth users",
			configure: func(limiter *config.Limiter) {
				limiter.BasicAuthUsers = []string{}
			},
		},
		{
			name: "mixed empty constraints",
			configure: func(limiter *config.Limiter) {
				limiter.Methods = []string{}
				limiter.Headers = map[string][]string{}
				limiter.BasicAuthUsers = []string{}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limiter := NewLimiter(1, time.Hour)
			tt.configure(limiter)
			request := httptest.NewRequest(http.MethodGet, "/limited", nil)
			request.RemoteAddr = "203.0.113.10:54321"

			got := BuildKeys(limiter, request)
			want := [][]string{{"203.0.113.10", "/limited"}}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("BuildKeys() = %#v, want %#v", got, want)
			}
		})
	}
}

func TestBuildKeysIgnoresEmptyConstraintsBesideActiveFilters(t *testing.T) {
	tests := []struct {
		name      string
		configure func(*config.Limiter)
		request   func() *http.Request
		want      [][]string
	}{
		{
			name: "empty methods preserve header filter",
			configure: func(limiter *config.Limiter) {
				limiter.Methods = []string{}
				limiter.Headers = map[string][]string{"X-Plan": {"gold"}}
			},
			request: func() *http.Request {
				request := httptest.NewRequest(http.MethodGet, "/limited", nil)
				request.Header.Set("X-Plan", "gold")
				return request
			},
			want: [][]string{{"203.0.113.10", "/limited", "X-Plan", "gold"}},
		},
		{
			name: "empty headers and auth preserve method filter",
			configure: func(limiter *config.Limiter) {
				limiter.Methods = []string{http.MethodGet}
				limiter.Headers = map[string][]string{}
				limiter.BasicAuthUsers = []string{}
			},
			request: func() *http.Request {
				return httptest.NewRequest(http.MethodGet, "/limited", nil)
			},
			want: [][]string{{"203.0.113.10", "/limited", http.MethodGet}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limiter := NewLimiter(1, time.Hour)
			tt.configure(limiter)
			request := tt.request()
			request.RemoteAddr = "203.0.113.10:54321"

			if got := BuildKeys(limiter, request); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("BuildKeys() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestLimitHandlerAppliesDefaultLimitWithMixedEmptyConstraints(t *testing.T) {
	limiter := NewLimiter(1, time.Hour)
	limiter.Methods = []string{}
	limiter.Headers = map[string][]string{}
	limiter.BasicAuthUsers = []string{}
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
}

func TestLimitHandlerDoesNotPartiallyChargeRejectedMultiKeyRequest(t *testing.T) {
	limiter := NewLimiter(1, time.Hour)
	limiter.Headers = map[string][]string{"X-Plan": {"gold", "silver"}}
	handler := LimitFuncHandler(limiter, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	request := func(values ...string) *http.Request {
		r := httptest.NewRequest(http.MethodGet, "/limited", nil)
		r.RemoteAddr = "203.0.113.10:54321"
		for _, value := range values {
			r.Header.Add("X-Plan", value)
		}
		return r
	}

	silver := httptest.NewRecorder()
	handler.ServeHTTP(silver, request("silver"))
	if silver.Code != http.StatusNoContent {
		t.Fatalf("silver response status = %d, want %d", silver.Code, http.StatusNoContent)
	}

	combined := httptest.NewRecorder()
	handler.ServeHTTP(combined, request("gold", "silver"))
	if combined.Code != http.StatusTooManyRequests {
		t.Fatalf("combined response status = %d, want %d", combined.Code, http.StatusTooManyRequests)
	}

	gold := httptest.NewRecorder()
	handler.ServeHTTP(gold, request("gold"))
	if gold.Code != http.StatusNoContent {
		t.Fatalf("gold response status = %d, want %d", gold.Code, http.StatusNoContent)
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

func TestLimitHandlerUsesConfiguredRejectionResponse(t *testing.T) {
	limiter := NewLimiter(1, time.Hour)
	limiter.StatusCode = http.StatusServiceUnavailable
	limiter.MessageContentType = "application/problem+json"
	limiter.Message = `{"error":"busy"}`
	nextCalls := 0
	handler := LimitFuncHandler(limiter, func(w http.ResponseWriter, r *http.Request) {
		nextCalls++
		w.WriteHeader(http.StatusNoContent)
	})
	request := httptest.NewRequest(http.MethodGet, "/limited", nil)
	request.RemoteAddr = "203.0.113.10:54321"

	handler.ServeHTTP(httptest.NewRecorder(), request)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	if got, want := recorder.Code, http.StatusServiceUnavailable; got != want {
		t.Fatalf("rejection status = %d, want %d", got, want)
	}
	if got, want := recorder.Header().Get("Content-Type"), "application/problem+json"; got != want {
		t.Fatalf("rejection content type = %q, want %q", got, want)
	}
	if got, want := recorder.Body.String(), `{"error":"busy"}`; got != want {
		t.Fatalf("rejection body = %q, want %q", got, want)
	}
	if got, want := recorder.Header().Get("X-Rate-Limit-Limit"), "1"; got != want {
		t.Fatalf("rate limit header = %q, want %q", got, want)
	}
	if nextCalls != 1 {
		t.Fatalf("wrapped handler calls = %d, want 1", nextCalls)
	}
}

func TestLimitByRequestReturnsConfiguredHTTPError(t *testing.T) {
	limiter := NewLimiter(1, time.Hour)
	limiter.StatusCode = http.StatusServiceUnavailable
	limiter.Message = "temporarily unavailable"
	request := httptest.NewRequest(http.MethodGet, "/limited", nil)
	request.RemoteAddr = "203.0.113.10:54321"

	if httpError := LimitByRequest(limiter, request); httpError != nil {
		t.Fatalf("first request returned an error: %v", httpError)
	}
	httpError := LimitByRequest(limiter, request)
	if httpError == nil {
		t.Fatal("second request did not return HTTPError")
	}
	if got, want := httpError.StatusCode, http.StatusServiceUnavailable; got != want {
		t.Fatalf("HTTPError status = %d, want %d", got, want)
	}
	if got, want := httpError.Message, "temporarily unavailable"; got != want {
		t.Fatalf("HTTPError message = %q, want %q", got, want)
	}
}

func TestLimitByRequestNormalizesRejectionStatusCode(t *testing.T) {
	tests := []struct {
		name       string
		configured int
		want       int
	}{
		{name: "below three digits", configured: 99, want: http.StatusTooManyRequests},
		{name: "continue is informational", configured: 100, want: http.StatusTooManyRequests},
		{name: "success", configured: 200, want: http.StatusTooManyRequests},
		{name: "redirect", configured: 302, want: http.StatusTooManyRequests},
		{name: "upper redirect boundary", configured: 399, want: http.StatusTooManyRequests},
		{name: "lower client error boundary", configured: 400, want: 400},
		{name: "nonstandard client error", configured: 499, want: 499},
		{name: "upper server error boundary", configured: 599, want: 599},
		{name: "above server error boundary", configured: 600, want: http.StatusTooManyRequests},
		{name: "upper three-digit boundary", configured: 999, want: http.StatusTooManyRequests},
		{name: "above three digits", configured: 1000, want: http.StatusTooManyRequests},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			limiter := NewLimiter(1, time.Hour)
			limiter.StatusCode = test.configured
			request := httptest.NewRequest(http.MethodGet, "/limited", nil)
			request.RemoteAddr = "203.0.113.10:54321"

			if httpError := LimitByRequest(limiter, request); httpError != nil {
				t.Fatalf("first request returned an error: %v", httpError)
			}
			httpError := LimitByRequest(limiter, request)
			if httpError == nil {
				t.Fatal("second request did not return HTTPError")
			}
			if got := httpError.StatusCode; got != test.want {
				t.Fatalf("HTTPError status = %d, want %d", got, test.want)
			}
		})
	}
}

func TestLimitHandlerNormalizesRejectionStatusCode(t *testing.T) {
	tests := []struct {
		name       string
		configured int
		want       int
	}{
		{name: "below three digits", configured: 99, want: http.StatusTooManyRequests},
		{name: "continue is informational", configured: 100, want: http.StatusTooManyRequests},
		{name: "success", configured: 200, want: http.StatusTooManyRequests},
		{name: "redirect", configured: 302, want: http.StatusTooManyRequests},
		{name: "upper redirect boundary", configured: 399, want: http.StatusTooManyRequests},
		{name: "lower client error boundary", configured: 400, want: 400},
		{name: "nonstandard client error", configured: 499, want: 499},
		{name: "upper server error boundary", configured: 599, want: 599},
		{name: "above server error boundary", configured: 600, want: http.StatusTooManyRequests},
		{name: "upper three-digit boundary", configured: 999, want: http.StatusTooManyRequests},
		{name: "above three digits", configured: 1000, want: http.StatusTooManyRequests},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			limiter := NewLimiter(1, time.Hour)
			limiter.StatusCode = test.configured
			handler := LimitFuncHandler(limiter, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			})
			request := httptest.NewRequest(http.MethodGet, "/limited", nil)
			request.RemoteAddr = "203.0.113.10:54321"

			handler.ServeHTTP(httptest.NewRecorder(), request)
			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, request)
			if got := recorder.Code; got != test.want {
				t.Fatalf("rejection status = %d, want %d", got, test.want)
			}
		})
	}
}

func TestLimitHandlerNormalizesNonErrorStatusForRealServer(t *testing.T) {
	for _, configured := range []int{http.StatusContinue, http.StatusOK, http.StatusFound, 399, 600, 999} {
		t.Run(strconv.Itoa(configured), func(t *testing.T) {
			limiter := NewLimiter(1, time.Hour)
			limiter.StatusCode = configured
			limiter.IPLookups = []string{"X-Real-IP"}
			handler := LimitFuncHandler(limiter, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			})
			server := httptest.NewServer(handler)
			defer server.Close()

			client := server.Client()
			client.CheckRedirect = func(_ *http.Request, _ []*http.Request) error {
				return http.ErrUseLastResponse
			}
			request := func() *http.Request {
				r, err := http.NewRequest(http.MethodGet, server.URL+"/limited", nil)
				if err != nil {
					t.Fatalf("create request: %v", err)
				}
				r.Header.Set("X-Real-IP", "203.0.113.10")
				return r
			}

			first, err := client.Do(request())
			if err != nil {
				t.Fatalf("first request failed: %v", err)
			}
			first.Body.Close()
			if got, want := first.StatusCode, http.StatusNoContent; got != want {
				t.Fatalf("first response status = %d, want %d", got, want)
			}

			second, err := client.Do(request())
			if err != nil {
				t.Fatalf("second request failed: %v", err)
			}
			defer second.Body.Close()
			body, err := io.ReadAll(second.Body)
			if err != nil {
				t.Fatalf("read rejection body: %v", err)
			}
			if got, want := second.StatusCode, http.StatusTooManyRequests; got != want {
				t.Fatalf("client-visible rejection status = %d, want %d", got, want)
			}
			if got, want := string(body), limiter.Message; got != want {
				t.Fatalf("rejection body = %q, want %q", got, want)
			}
		})
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
