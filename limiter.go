// Package limiter provides rate-limiting logic to HTTP request handler.
package limiter

import (
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/garethpaul/go-ratelimiter/config"
	"github.com/garethpaul/go-ratelimiter/errors"
	"github.com/garethpaul/go-ratelimiter/libstring"
)

// NewLimiter is a convenience function to config.NewLimiter.
func NewLimiter(max int64, ttl time.Duration) *config.Limiter {
	return config.NewLimiter(max, ttl)
}

// LimitByKeys tracks requests using an unambiguous encoding of key components.
// It returns HTTPError when limit is exceeded.
func LimitByKeys(limiter *config.Limiter, keys []string) *errors.HTTPError {
	if limiter.LimitReached(encodeKeys(keys)) {
		return limitError(limiter)
	}

	return nil
}

func limitError(limiter *config.Limiter) *errors.HTTPError {
	return &errors.HTTPError{
		Message:    limiter.Message,
		StatusCode: rejectionStatusCode(limiter.StatusCode),
	}
}

func rejectionStatusCode(statusCode int) int {
	if statusCode < 100 || statusCode > 999 {
		return http.StatusTooManyRequests
	}

	return statusCode
}

func encodeKeys(keys []string) string {
	var encoded strings.Builder
	for _, key := range keys {
		encoded.WriteString(strconv.Itoa(len(key)))
		encoded.WriteByte(':')
		encoded.WriteString(key)
	}
	return encoded.String()
}

// LimitByRequest builds request keys and checks all derived buckets atomically.
func LimitByRequest(limiter *config.Limiter, r *http.Request) *errors.HTTPError {
	sliceKeys := BuildKeys(limiter, r)
	encodedKeys := make([]string, 0, len(sliceKeys))
	for _, keys := range sliceKeys {
		encodedKeys = append(encodedKeys, encodeKeys(keys))
	}
	if limiter.LimitReachedForKeys(encodedKeys) {
		return limitError(limiter)
	}

	return nil
}

// BuildKeys generates a slice of keys to rate-limit by given config and request structs.
func BuildKeys(limiter *config.Limiter, r *http.Request) [][]string {
	remoteIP := libstring.RemoteIP(limiter.IPLookups, r)
	path := r.URL.Path
	sliceKeys := make([][]string, 0)
	limitMethods := len(limiter.Methods) > 0
	limitHeaders := len(limiter.Headers) > 0
	limitBasicAuth := len(limiter.BasicAuthUsers) > 0
	headerKeys := sortedHeaderKeys(limiter.Headers)

	// Don't BuildKeys if remoteIP is blank.
	if remoteIP == "" {
		return sliceKeys
	}

	if limitMethods && limitHeaders && limitBasicAuth {
		// Limit by HTTP methods and HTTP headers+values and Basic Auth credentials.
		if libstring.StringInSlice(limiter.Methods, r.Method) {
			for _, headerKey := range headerKeys {
				headerValues := limiter.Headers[headerKey]
				matchedHeaderValues := matchingHeaderValues(r, headerKey, headerValues)
				if (headerValues == nil || len(headerValues) <= 0) && len(matchedHeaderValues) > 0 {
					// If header values are empty, rate-limit all request with headerKey.
					username, _, ok := r.BasicAuth()
					if ok && libstring.StringInSlice(limiter.BasicAuthUsers, username) {
						sliceKeys = append(sliceKeys, []string{remoteIP, path, r.Method, headerKey, username})
					}

				} else if len(headerValues) > 0 {
					// If header values are not empty, rate-limit all request with headerKey and headerValues.
					for _, headerValue := range matchedHeaderValues {
						username, _, ok := r.BasicAuth()
						if ok && libstring.StringInSlice(limiter.BasicAuthUsers, username) {
							sliceKeys = append(sliceKeys, []string{remoteIP, path, r.Method, headerKey, headerValue, username})
						}
					}
				}
			}
		}

	} else if limitMethods && limitHeaders {
		// Limit by HTTP methods and HTTP headers+values.
		if libstring.StringInSlice(limiter.Methods, r.Method) {
			for _, headerKey := range headerKeys {
				headerValues := limiter.Headers[headerKey]
				matchedHeaderValues := matchingHeaderValues(r, headerKey, headerValues)
				if (headerValues == nil || len(headerValues) <= 0) && len(matchedHeaderValues) > 0 {
					// If header values are empty, rate-limit all request with headerKey.
					sliceKeys = append(sliceKeys, []string{remoteIP, path, r.Method, headerKey})

				} else if len(headerValues) > 0 {
					// If header values are not empty, rate-limit all request with headerKey and headerValues.
					for _, headerValue := range matchedHeaderValues {
						sliceKeys = append(sliceKeys, []string{remoteIP, path, r.Method, headerKey, headerValue})
					}
				}
			}
		}

	} else if limitMethods && limitBasicAuth {
		// Limit by HTTP methods and Basic Auth credentials.
		if libstring.StringInSlice(limiter.Methods, r.Method) {
			username, _, ok := r.BasicAuth()
			if ok && libstring.StringInSlice(limiter.BasicAuthUsers, username) {
				sliceKeys = append(sliceKeys, []string{remoteIP, path, r.Method, username})
			}
		}

	} else if limitMethods {
		// Limit by HTTP methods.
		if libstring.StringInSlice(limiter.Methods, r.Method) {
			sliceKeys = append(sliceKeys, []string{remoteIP, path, r.Method})
		}

	} else if limitHeaders {
		// Limit by HTTP headers+values.
		for _, headerKey := range headerKeys {
			headerValues := limiter.Headers[headerKey]
			matchedHeaderValues := matchingHeaderValues(r, headerKey, headerValues)
			if (headerValues == nil || len(headerValues) <= 0) && len(matchedHeaderValues) > 0 {
				// If header values are empty, rate-limit all request with headerKey.
				sliceKeys = append(sliceKeys, []string{remoteIP, path, headerKey})

			} else if len(headerValues) > 0 {
				// If header values are not empty, rate-limit all request with headerKey and headerValues.
				for _, headerValue := range matchedHeaderValues {
					sliceKeys = append(sliceKeys, []string{remoteIP, path, headerKey, headerValue})
				}
			}
		}

	} else if limitBasicAuth {
		// Limit by Basic Auth credentials.
		username, _, ok := r.BasicAuth()
		if ok && libstring.StringInSlice(limiter.BasicAuthUsers, username) {
			sliceKeys = append(sliceKeys, []string{remoteIP, path, username})
		}
	} else {
		// Default: Limit by remoteIP and path.
		sliceKeys = append(sliceKeys, []string{remoteIP, path})
	}

	return sliceKeys
}

func sortedHeaderKeys(headers map[string][]string) []string {
	keys := make([]string, 0, len(headers))
	for key := range headers {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func matchingHeaderValues(r *http.Request, headerKey string, headerValues []string) []string {
	requestValues := r.Header.Values(headerKey)
	if len(requestValues) == 0 {
		return nil
	}
	if len(headerValues) == 0 {
		for _, requestValue := range requestValues {
			if strings.TrimSpace(requestValue) != "" {
				return []string{""}
			}
		}
		return nil
	}

	matchedValues := make([]string, 0, len(headerValues))
	seenValues := make(map[string]struct{}, len(headerValues))
	for _, headerValue := range headerValues {
		if strings.TrimSpace(headerValue) == "" {
			continue
		}
		if _, seen := seenValues[headerValue]; seen {
			continue
		}
		if libstring.StringInSlice(requestValues, headerValue) {
			matchedValues = append(matchedValues, headerValue)
			seenValues[headerValue] = struct{}{}
		}
	}
	return matchedValues
}

// SetResponseHeaders configures X-Rate-Limit-Limit and X-Rate-Limit-Duration
func SetResponseHeaders(limiter *config.Limiter, w http.ResponseWriter) {
	w.Header().Set("X-Rate-Limit-Limit", strconv.FormatInt(limiter.Max, 10))
	w.Header().Set("X-Rate-Limit-Duration", limiter.TTL.String())
}

// LimitHandler is a middleware that performs rate-limiting given http.Handler struct.
func LimitHandler(limiter *config.Limiter, next http.Handler) http.Handler {
	middle := func(w http.ResponseWriter, r *http.Request) {
		SetResponseHeaders(limiter, w)

		httpError := LimitByRequest(limiter, r)
		if httpError != nil {
			w.Header().Set("Content-Type", limiter.MessageContentType)
			w.WriteHeader(httpError.StatusCode)
			w.Write([]byte(httpError.Message))
			return
		}

		// There's no rate-limit error, serve the next handler.
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(middle)
}

// LimitFuncHandler is a middleware that performs rate-limiting given request handler function.
func LimitFuncHandler(limiter *config.Limiter, nextFunc func(http.ResponseWriter, *http.Request)) http.Handler {
	return LimitHandler(limiter, http.HandlerFunc(nextFunc))
}
