# Changes

## 2026-06-26 07:52 PDT

Priority: correctness and rate-limit identity integrity.

Summary:

- Canonicalized validated request IPs so equivalent textual forms cannot obtain
  independent token buckets.

Work completed:

- Returned `net.IP.String()` from both direct-address and proxy-header parsers.
- Added middleware coverage proving expanded and compressed IPv6 forms share
  one exhausted bucket.
- Added table-driven coverage for `RemoteAddr`, `X-Forwarded-For`, and
  `X-Real-IP`, plus two hostile source mutations.

Threads:

- Request identity, IPv6 normalization, trusted proxy inputs, and maintained
  static verification.

Files changed:

- `libstring/libstring.go`, focused Go tests, mutation and baseline scripts,
  project guidance, and canonical-IP plans.

Validation:

- Red phase: equivalent expanded/compressed IPv6 requests returned `204` twice
  instead of exhausting one shared bucket; all three lookup sources returned
  the expanded text unchanged.
- Green focused tests pass on Go 1.25.11.
- Repository-root and external-directory `make check` pass on Go 1.25.11 with
  race tests, vet, module-integrity checks, and two canonical IP hostile mutations.
- Hosted checks, CodeQL, and review are pending.

Bugs and findings:

- Validating without canonicalizing allowed one client address to rotate IPv6
  spellings and receive fresh process-local buckets.
- Equivalent textual IP addresses share one canonical limiter identity across `RemoteAddr`, `X-Forwarded-For`, and `X-Real-IP`.

Blockers:

- None for local implementation.

Next action:

- Push the focused PR, then verify hosted checks and the exact final head.

- Limiter rejection status codes outside 400 through 599 fall back to 429; configured client and server error codes remain unchanged.
- Documented and regression-tested middleware-owned custom rejection responses
  and caller-owned `HTTPError` serialization through the direct limiter APIs.
- Middleware rejections use the configured `StatusCode`, `MessageContentType`, and `Message`; callers needing extra headers or custom serialization should call `LimitByRequest` or `LimitByKeys` and write the returned `HTTPError` themselves.
- Documented that limiter accounting is serialized per limiter, remains
  process-local, has no background cleanup, and uses capacity-driven LRU
  eviction that resets an evicted key to a fresh bucket if it is admitted again.
- `LimitReached` calls on directly configured valid limiters lazily initialize private accounting state with the same 10,000-key cap as `NewLimiter`.
- Configured header names are sorted before limiter keys are derived, while configured value order remains unchanged.
- Prevented empty limiter constraint collections from bypassing the default
  remote-IP/path rate limit.

## 2026-06-13

- Made Go verification independent of the caller's working directory by
  resolving the baseline checker from the loaded Makefile.
- Made rejected multi-key preflight leave tracked bucket state unchanged until
  every existing bucket confirms capacity.
- Made multi-key request accounting atomic so a rejected request cannot consume
  tokens from buckets that still have capacity.
- Deduplicated repeated configured header values before deriving limiter keys so
  one request cannot charge the same token bucket multiple times.

## 2026-06-12

- Made rate-limit metadata and rejection content types replace stale response
  values instead of accumulating ambiguous duplicates.
- Encoded limiter key components unambiguously and stored only fixed-length
  SHA-256 identifiers so request-controlled key bytes remain bounded.
- Disabled persisted checkout credentials, enforced one canonical hosted
  workflow, and added repository-wide ownership.
- Raised the hosted and documented Go toolchain from 1.25.3 to 1.25.11 after
  `govulncheck` found three reachable standard-library vulnerabilities.
- Corrected token-bucket refill rates to restore `Max` requests across each
  `TTL` instead of restoring only one request per complete duration.
- Made non-positive and platform-unrepresentable limiter configurations fail
  closed without allocating request-derived key state.
- Added race-tested refill, burst, invalid-configuration, and allocation
  regressions plus static maintenance contracts.

## 2026-06-10

- Bounded each limiter to 10,000 tracked keys with least-recently-used eviction
  and race-tested capacity coverage.
- Added a pinned, least-privilege Go 1.25.11 workflow with vet, race-enabled
  tests, module-integrity checks, and static guardrails.

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
