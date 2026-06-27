# Changes

## 2026-06-26 18:14 PDT - P2 - Spaced Makefile paths

### Summary

- Preserved the repository root when an absolute Makefile is loaded from
  spaced checkout paths under GNU Make 4.2 and 4.4.
- Added a recursive-safe full-gate regression that invokes the absolute
  Makefile from an external directory under a spaced checkout path.
- Made Go formatting discovery null-delimited so repository paths are not
  split before `gofmt` receives them.

### Files changed

- `Makefile` — resolves one validated loaded Makefile without tokenizing its
  path and rejects `MAKEFILES` or `MAKEFILE_LIST` authority injection.
- `scripts/check-baseline.sh` and `scripts/test-make-spaced-path.py` — preserve
  Go paths and exercise spaced-root plus hostile startup behavior.
- `AGENTS.md`, `README.md`, and
  `docs/plans/2026-06-13-location-independent-make.md` — synchronized guidance
  and completed evidence.

### Validation

- The pre-change gate failed from a spaced checkout, first during root
  resolution and then at whitespace-splitting `gofmt` discovery.
- Hosted Go checks cover formatting, vet, race-enabled tests, build, module
  integrity, the recursive spaced-path gate, and hostile startup rejection.
- Local parser reproduction confirms the spaced absolute path resolves exactly;
  hostile `MAKEFILES` and `MAKEFILE_LIST` inputs fail before recipes run. Full
  local Go execution is unavailable because Go is not installed here.

### Bugs / findings

- Replacing every `MAKEFILE_LIST` space with one sentinel also collapsed
  separators between multiple loaded files, producing a concatenated false
  root. Validated single-file resolution now fails closed instead.

### Blockers

- The exact Go 1.25.11 runtime matrix is hosted-only in this environment.

### Next action

- Review the refreshed exact PR head and merge only after every hosted check is
  green.

## 2026-06-26 - Scoped IPv6 limiter identities

- Fixed valid scoped IPv6 addresses that previously produced no request keys
  and bypassed rate limiting.
- Scoped IPv6 addresses retain their zone and consume limiter buckets across `RemoteAddr`, `X-Forwarded-For`, and `X-Real-IP`.
- Added cross-source parser coverage, middleware bucket-sharing coverage, and
  zone-preservation mutation checks.
- Go 1.25.11 focused tests and both canonical-IP hostile mutations pass; full
  repository and external-directory verification are recorded in the plan.

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
- Implementation head `a839b8b3304feec5319fb91a02409cfd012b0c7b` passed hosted
  push and pull-request checks (`28245972824`, `28245975267`) and CodeQL
  (`28245973551`).
- `codex review --base origin/master` was attempted and skipped after HTTP 401
  authentication errors, per the maintenance loop policy.

Bugs and findings:

- Validating without canonicalizing allowed one client address to rotate IPv6
  spellings and receive fresh process-local buckets.
- Equivalent textual IP addresses share one canonical limiter identity across `RemoteAddr`, `X-Forwarded-For`, and `X-Real-IP`.

Blockers:

- Codex review authentication is unavailable; all executable local and hosted
  gates are green.

Next action:

- Verify the evidence-only final PR head, then merge that exact commit.

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
