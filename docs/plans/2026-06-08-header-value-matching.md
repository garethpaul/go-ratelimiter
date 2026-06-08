# Go Rate Limiter Header Value Matching Plan

status: completed

## Context

`BuildKeys` supports rate-limit keys for configured HTTP headers and optional
header values. Before this pass, configured values could be appended to keys
when the request merely had the header, even when the actual request value did
not match the configured value.

## Objectives

- Only build header-value limiter keys for request values that match configured values.
- Preserve the existing behavior where an empty configured value list means any non-empty header value matches.
- Cover header-only, header-value, and method+header-value behavior with tests.
- Keep the existing `make check` baseline as the verification entry point.

## Work Items

1. Added a `matchingHeaderValues` helper for request/configured header comparison.
2. Updated all header-value key-building branches to use matched values only.
3. Added tests for non-matching header values, matching header values, and method+header-value keys.
4. Extended `scripts/check-baseline.sh` to require the new behavior tests and plan.

## Verification

- `make check`
- `go test ./...`
- `git diff --check`
