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

func TestRemoteIPTrimsForwardedForList(t *testing.T) {
	request := httptest.NewRequest("GET", "/", nil)
	request.RemoteAddr = "10.0.0.1:1234"
	request.Header.Set("X-Forwarded-For", "198.51.100.7, 10.0.0.1")

	got := RemoteIP([]string{"X-Forwarded-For"}, request)

	if got != "198.51.100.7" {
		t.Fatalf("RemoteIP = %q, want first forwarded IP", got)
	}
}
