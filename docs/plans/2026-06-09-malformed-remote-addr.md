# Malformed RemoteAddr Guard

status: completed

## Context

Proxy header values are validated before they become limiter keys, but direct
`RemoteAddr` fallback values could still pass through as arbitrary strings when
they were not parseable IP addresses. Malformed direct addresses should behave
like malformed proxy headers: skip the malformed source and allow later
configured lookup sources to provide the limiter key.

## Completed Scope

- Validated the parsed `RemoteAddr` host with `net.ParseIP`.
- Skipped malformed direct remote addresses before deriving limiter keys.
- Continued to later configured lookup sources after malformed direct remote
  addresses.
- Added focused tests for direct lookup, fallback lookup, and `BuildKeys`
  behavior.
- Extended the static baseline and docs so the guard remains visible.

## Verification

- `go test ./...`
- `scripts/check-baseline.sh`
- `make check`
- `git diff --check`
