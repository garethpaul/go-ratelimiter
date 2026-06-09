# IPv6 RemoteAddr Parsing

status: completed

## Context

The limiter derived default keys from `Request.RemoteAddr` by trimming at the
last colon. That worked for IPv4 host:port values, but bracketed IPv6 host:port
addresses kept brackets and bare IPv6 addresses were truncated.

## Completed Scope

- Added `RemoteIP` coverage for bracketed and bare IPv6 `RemoteAddr` values.
- Switched host:port parsing to `net.SplitHostPort` with a bare-host fallback.
- Extended the static baseline and docs for IPv6 lookup behavior.

## Verification

- `go test ./...`
- `make check`
- `git diff --check`
