# Proxy Header IP Validation

status: completed

## Context

The limiter can derive request keys from `X-Forwarded-For` and `X-Real-IP`.
Recent guards skip blank values, but malformed non-IP header values could still
be used as rate-limit key material instead of falling back through the
configured lookup order.

## Completed Scope

- Added proxy-header IP normalization backed by `net.ParseIP`.
- Skipped malformed `X-Forwarded-For` entries and malformed `X-Real-IP` values.
- Added focused tests for malformed proxy header fallback behavior.
- Extended the static baseline and docs so malformed proxy IP headers remain
  covered.

## Verification

- `go test ./...`
- `scripts/check-baseline.sh`
- `make check`
- `git diff --check`
