# Malformed RemoteAddr Guard

status: completed

## Context

Proxy header values are validated before they become limiter keys, but direct
`RemoteAddr` fallback values could still pass through as arbitrary strings when
they were not parseable IP addresses. Malformed direct addresses should behave
like malformed proxy headers and produce no limiter key.

## Completed Scope

- Validated the parsed `RemoteAddr` host with `net.ParseIP`.
- Returned an empty IP for malformed direct remote addresses.
- Added focused tests for direct lookup and `BuildKeys` behavior.
- Extended the static baseline and docs so the guard remains visible.

## Verification

- `go test ./...`
- `scripts/check-baseline.sh`
- `make check`
- `git diff --check`
