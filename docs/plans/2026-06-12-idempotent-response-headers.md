# Idempotent Rate-Limit Response Headers

status: completed

## Context

The HTTP middleware currently appends its rate-limit metadata and rejection
content type with `Header.Add`. If a response writer already contains those
headers, or if limiter middleware is nested, the response carries multiple
values. HTTP clients commonly select the first value, so stale upstream
metadata can obscure the limiter configuration that actually handled the
request.

These fields are authoritative response metadata for this middleware, not
independent list entries. Reapplying the middleware should therefore replace
the previous value and produce the same headers as applying it once.

## Priority

Ambiguous limit and content-type headers make throttling behavior harder for
clients to interpret and can mislabel a rejection body. The fix is small,
backward-compatible for normal single-value responses, and directly improves
middleware composability.

## Prioritized Engineering Backlog

1. Make rate-limit metadata and rejection content type idempotent now.
2. Make multi-key request accounting atomic so a rejected request cannot
   consume tokens from keys checked earlier in the same request.
3. Add optional remaining-token and retry metadata if the public API gains a
   response-observability surface.

## Requirements

- R1. `SetResponseHeaders` must leave exactly one current
  `X-Rate-Limit-Limit` value and one current `X-Rate-Limit-Duration` value.
- R2. A rate-limit rejection must leave exactly one configured `Content-Type`
  value, replacing any stale value already present on the writer.
- R3. Successful requests, rejection status codes, bodies, token accounting,
  and request-key derivation must remain unchanged.
- R4. Regression tests must cover preexisting rate-limit metadata, repeated
  header application, and preexisting rejection content types.
- R5. Public middleware and constructor signatures must remain compatible.

## Implementation Units

### U1. Replace authoritative headers

- **Files:** `limiter.go`
- Use replacement semantics for rate-limit metadata and the rejection content
  type so repeated middleware application is deterministic.

### U2. Add HTTP regressions

- **Files:** `limiter_test.go`
- Seed response writers with stale values, apply the header helper more than
  once, and verify the rejection path replaces stale content metadata.

### U3. Update behavior documentation

- **Files:** `README.md`, `VISION.md`, `CHANGES.md`
- Record the single-value, authoritative response-header contract.

## Scope Boundaries

- Do not change limiter keys, proxy trust, token refill behavior, or LRU state.
- Do not add new response headers or dependencies.
- Do not redesign error body writing or public configuration fields.

## Verification

- `make check`
- `go test -race -count=1 ./...`
- `go vet ./...`
- `go mod tidy -diff`
- `go mod verify`
- `git diff --check`
- Mutations restoring `Header.Add` for owned response fields must fail the
  regression suite.

Completed on 2026-06-12 with the repository baseline, race-enabled tests, vet,
module tidiness and integrity checks, mutation coverage, and diff hygiene
passing.

## Work Completed

- Replaced append semantics with authoritative replacement for both
  rate-limit metadata fields and the configured rejection content type.
- Preserved successful responses, rejection status and body behavior, request
  key derivation, token accounting, and all public APIs.
- Added regressions for stale preexisting values and repeated header
  application so nested middleware remains deterministic.

## Verification Completed

- `go test -count=1 ./...`, `go test -race -count=1 ./...`, `go vet ./...`,
  `go mod tidy -diff`, `go mod verify`, all four Make gates, and
  `git diff --check` passed locally on Go 1.25.11.
- Implementation push run `27393483036` and pull-request run `27393485015`
  passed at commit `9a3ab6509969add0feb4b0db5cf329adec7e4ddf`.
- Post-merge push run `27393504527` and CodeQL run `27402321986` passed at
  default-branch merge commit `a38f42b10e78af3f904c6b14618d1df22b0fe968`.
- Mutations restoring `Header.Add` for rate-limit metadata or rejection
  content type were rejected by the regression suite and baseline.
