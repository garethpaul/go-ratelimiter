# Changes

## 2026-06-09

- Skipped blank header-only request values before deriving limiter keys.
- Skipped blank configured header values before deriving limiter keys.
- Matched configured header values across all request header values so a blank
  first value cannot hide a later configured match.

## 2026-06-08

- Added `make lint`, `make test`, and `make build` aliases so local verification
  has the expected pre-push gate targets in addition to `make check`.
- Added Go module metadata and lockfile for modern Go tooling.
- Updated local imports to use the module path.
- Added tests for default key derivation, proxy-aware IP lookup, header-value matching, and rate-limit rejection behavior.
- Added IPv6 `RemoteAddr` parsing coverage and switched host:port parsing to
  `net.SplitHostPort`.
- Skipped malformed `RemoteAddr` values before deriving limiter keys so later
  configured lookup sources can be used.
- Skipped blank `X-Forwarded-For` entries before deriving limiter keys.
- Trimmed `X-Real-IP` values and skipped blank values before deriving limiter
  keys.
- Skipped malformed proxy IP headers before deriving limiter keys.
- Fixed configured header-value limiting so non-matching request values do not reuse configured values as keys.
- Added `make check` and static guardrails for formatting, tests, module imports, and plan completion.
