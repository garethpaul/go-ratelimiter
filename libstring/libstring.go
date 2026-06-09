// Package libstring provides various string related functions.
package libstring

import (
	"net"
	"net/http"
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

func ipAddrFromRemoteAddr(s string) string {
	host, _, err := net.SplitHostPort(s)
	if err != nil {
		host = strings.Trim(s, "[]")
	}
	if net.ParseIP(host) == nil {
		return ""
	}
	return host
}

func ipAddrFromHeaderValue(s string) string {
	ip := strings.Trim(strings.TrimSpace(s), "[]")
	if net.ParseIP(ip) == nil {
		return ""
	}
	return ip
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
			return ipAddrFromRemoteAddr(r.RemoteAddr)
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
