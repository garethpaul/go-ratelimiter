// Package libstring provides various string related functions.
package libstring

import (
	"net"
	"net/http"
	"net/netip"
	"strings"
)

// StringInSlice finds needle in a slice of strings.
func StringInSlice(sliceString []string, needle string) bool {
	for _, b := range sliceString {
		if b == needle {
			return true
		}
	}
	return false
}

func canonicalIP(s string) string {
	addr, err := netip.ParseAddr(s)
	if err != nil {
		return ""
	}
	return addr.Unmap().String()
}

func ipAddrFromRemoteAddr(s string) string {
	host, _, err := net.SplitHostPort(s)
	if err != nil {
		host = strings.Trim(s, "[]")
	}
	return canonicalIP(host)
}

func ipAddrFromHeaderValue(s string) string {
	value := strings.Trim(strings.TrimSpace(s), "[]")
	return canonicalIP(value)
}

func ipAddrFromForwardedFor(s string) string {
	parts := strings.Split(s, ",")
	for _, part := range parts {
		ip := ipAddrFromHeaderValue(part)
		if ip != "" {
			return ip
		}
	}
	return ""
}

// RemoteIP finds IP Address given http.Request struct.
func RemoteIP(ipLookups []string, r *http.Request) string {
	realIP := ipAddrFromHeaderValue(r.Header.Get("X-Real-IP"))
	forwardedFor := r.Header.Get("X-Forwarded-For")

	for _, lookup := range ipLookups {
		if lookup == "RemoteAddr" {
			if ip := ipAddrFromRemoteAddr(r.RemoteAddr); ip != "" {
				return ip
			}
		}
		if lookup == "X-Forwarded-For" && forwardedFor != "" {
			if ip := ipAddrFromForwardedFor(forwardedFor); ip != "" {
				return ip
			}
		}
		if lookup == "X-Real-IP" && realIP != "" {
			return realIP
		}
	}

	return ""
}
