# Changes

## 2026-06-08

- Added Go module metadata and lockfile for modern Go tooling.
- Updated local imports to use the module path.
- Added tests for default key derivation, proxy-aware IP lookup, header-value matching, and rate-limit rejection behavior.
- Added IPv6 `RemoteAddr` parsing coverage and switched host:port parsing to
  `net.SplitHostPort`.
- Skipped blank `X-Forwarded-For` entries before deriving limiter keys.
- Fixed configured header-value limiting so non-matching request values do not reuse configured values as keys.
- Added `make check` and static guardrails for formatting, tests, module imports, and plan completion.
