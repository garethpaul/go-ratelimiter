package libstring

import (
	"net/http/httptest"
	"testing"
)

func TestRemoteIPUsesConfiguredLookupOrder(t *testing.T) {
	request := httptest.NewRequest("GET", "/", nil)
	request.RemoteAddr = "10.0.0.1:1234"
	request.Header.Set("X-Real-IP", "203.0.113.9")

	got := RemoteIP([]string{"X-Real-IP", "RemoteAddr"}, request)

	if got != "203.0.113.9" {
		t.Fatalf("RemoteIP = %q, want X-Real-IP value", got)
	}
}

func TestRemoteIPTrimsRealIP(t *testing.T) {
	request := httptest.NewRequest("GET", "/", nil)
	request.Header.Set("X-Real-IP", " 203.0.113.9 ")

	got := RemoteIP([]string{"X-Real-IP"}, request)

	if got != "203.0.113.9" {
		t.Fatalf("RemoteIP = %q, want trimmed X-Real-IP value", got)
	}
}

func TestRemoteIPFallsBackAfterBlankRealIP(t *testing.T) {
	request := httptest.NewRequest("GET", "/", nil)
	request.RemoteAddr = "10.0.0.1:1234"
	request.Header.Set("X-Real-IP", " ")

	got := RemoteIP([]string{"X-Real-IP", "RemoteAddr"}, request)

	if got != "10.0.0.1" {
		t.Fatalf("RemoteIP = %q, want fallback RemoteAddr value", got)
	}
}

func TestRemoteIPFallsBackAfterMalformedRealIP(t *testing.T) {
	request := httptest.NewRequest("GET", "/", nil)
	request.RemoteAddr = "10.0.0.1:1234"
	request.Header.Set("X-Real-IP", "not-an-ip")

	got := RemoteIP([]string{"X-Real-IP", "RemoteAddr"}, request)

	if got != "10.0.0.1" {
		t.Fatalf("RemoteIP = %q, want fallback RemoteAddr value", got)
	}
}

func TestRemoteIPTrimsForwardedForList(t *testing.T) {
	request := httptest.NewRequest("GET", "/", nil)
	request.RemoteAddr = "10.0.0.1:1234"
	request.Header.Set("X-Forwarded-For", "198.51.100.7, 10.0.0.1")

	got := RemoteIP([]string{"X-Forwarded-For"}, request)

	if got != "198.51.100.7" {
		t.Fatalf("RemoteIP = %q, want first forwarded IP", got)
	}
}

func TestRemoteIPSkipsBlankForwardedForEntries(t *testing.T) {
	request := httptest.NewRequest("GET", "/", nil)
	request.RemoteAddr = "10.0.0.1:1234"
	request.Header.Set("X-Forwarded-For", " , 198.51.100.7, 10.0.0.1")

	got := RemoteIP([]string{"X-Forwarded-For"}, request)

	if got != "198.51.100.7" {
		t.Fatalf("RemoteIP = %q, want first non-empty forwarded IP", got)
	}
}

func TestRemoteIPSkipsMalformedForwardedForEntries(t *testing.T) {
	request := httptest.NewRequest("GET", "/", nil)
	request.RemoteAddr = "10.0.0.1:1234"
	request.Header.Set("X-Forwarded-For", "not-an-ip, 198.51.100.7")

	got := RemoteIP([]string{"X-Forwarded-For", "RemoteAddr"}, request)

	if got != "198.51.100.7" {
		t.Fatalf("RemoteIP = %q, want first valid forwarded IP", got)
	}
}

func TestRemoteIPFallsBackAfterBlankForwardedFor(t *testing.T) {
	request := httptest.NewRequest("GET", "/", nil)
	request.RemoteAddr = "10.0.0.1:1234"
	request.Header.Set("X-Forwarded-For", " , ")

	got := RemoteIP([]string{"X-Forwarded-For", "RemoteAddr"}, request)

	if got != "10.0.0.1" {
		t.Fatalf("RemoteIP = %q, want fallback RemoteAddr value", got)
	}
}

func TestRemoteIPHandlesIPv6RemoteAddr(t *testing.T) {
	for _, remoteAddr := range []string{"[2001:db8::1]:1234", "2001:db8::1"} {
		request := httptest.NewRequest("GET", "/", nil)
		request.RemoteAddr = remoteAddr

		got := RemoteIP([]string{"RemoteAddr"}, request)

		if got != "2001:db8::1" {
			t.Fatalf("RemoteIP(%q) = %q, want IPv6 host", remoteAddr, got)
		}
	}
}

func TestRemoteIPSkipsMalformedRemoteAddr(t *testing.T) {
	request := httptest.NewRequest("GET", "/", nil)
	request.RemoteAddr = "not-an-ip"

	got := RemoteIP([]string{"RemoteAddr"}, request)

	if got != "" {
		t.Fatalf("RemoteIP = %q, want empty value for malformed RemoteAddr", got)
	}
}
