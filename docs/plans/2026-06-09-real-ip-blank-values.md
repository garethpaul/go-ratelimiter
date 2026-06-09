# X-Real-IP Blank Value Guard

status: completed

## Context

`Limiter.IPLookups` can prioritize `X-Real-IP` before `RemoteAddr`. The
`X-Forwarded-For` parser already trims entries and falls through on blank
values, but `X-Real-IP` was returned directly from the header. Padded values
could create distinct limiter keys, and whitespace-only values could prevent
fallback to later lookup sources.

## Completed Scope

- Trimmed `X-Real-IP` before using it as a limiter key source.
- Let lookup continue when `X-Real-IP` contains only whitespace.
- Added focused tests for trimmed and blank `X-Real-IP` behavior.
- Extended the static baseline and docs to preserve the proxy-header guard.

## Verification

- `go test ./...`
- `scripts/check-baseline.sh`
- `make check`
- `git diff --check`
