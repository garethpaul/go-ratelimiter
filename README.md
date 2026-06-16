# go-ratelimiter

<!-- README-OVERVIEW-IMAGE -->
![Project overview](docs/readme-overview.svg)

## Overview

`garethpaul/go-ratelimiter` is a Go project. A golang rate limiter.

This README is based on the checked-in source, manifests, scripts, and repository metadata on the `master` branch. The project language mix found during review was: Go (4).

## Repository Contents

- `CHANGES.md` - concise history of maintenance changes
- `Makefile` - local verification entry point
- `README.md` - project overview and local usage notes
- `config` - source or example code
- `errors` - source or example code
- `go.mod` and `go.sum` - Go module dependency metadata
- `libstring` - source or example code
- `scripts/check-baseline.sh` - Go formatting, tests, import, and documentation guardrails
- `SECURITY.md` - security reporting and disclosure guidance
- `VISION.md` - project direction and maintenance guardrails

Additional scan context:

- Source directories: config, errors, libstring
- Dependency and build manifests: go.mod, go.sum
- Entry points or build surfaces: `make check`, `go test ./...`
- Test-looking files: limiter_test.go, libstring/libstring_test.go

## Getting Started

### Prerequisites

- Git
- Go 1.25.11 or a compatible patched Go toolchain

### Setup

```bash
git clone https://github.com/garethpaul/go-ratelimiter.git
cd go-ratelimiter
go mod download
```

The setup commands above are derived from repository files. Legacy mobile, Python, or JavaScript samples may require older SDKs or package versions than a modern workstation uses by default.

## Running or Using the Project

- Import the package as `github.com/garethpaul/go-ratelimiter`.
- Use `LimitFuncHandler` or `LimitHandler` to wrap an HTTP handler with an in-memory token-bucket limiter.
- `NewLimiter(max, ttl)` permits a burst of `max` requests and refills `max`
  tokens across each `ttl`; non-positive values fail closed.
- Middleware-owned rate-limit and rejection content-type headers replace stale
  values so nested or repeated application remains unambiguous.
- Middleware rejections use the configured `StatusCode`, `MessageContentType`, and `Message`; callers needing extra headers or custom serialization should call `LimitByRequest` or `LimitByKeys` and write the returned `HTTPError` themselves.
- Limiter rejection status codes below 100 or above 999 fall back to 429; valid three-digit configured codes remain unchanged.

## Testing and Verification

Run the baseline:

```bash
make lint
make test
make build
make check
```

The `lint`, `test`, and `build` targets currently delegate to the static
baseline so every local gate entry point runs the same checks. The baseline runs
`go test ./...`, verifies Go formatting, checks module-qualified imports, and
ensures the behavior tests for key derivation, proxy-aware IP lookup, blank
X-Forwarded-For entries, blank X-Real-IP values, malformed proxy IP headers,
malformed RemoteAddr values, IPv6 RemoteAddr parsing, header-value matching,
blank first header value matching, blank configured header values, blank header-only request values,
token-bucket refill semantics, invalid limit rejection, and 429 responses
remain in place. Keep the exact guard phrases
"blank X-Forwarded-For", "blank X-Real-IP", "malformed RemoteAddr", and
"IPv6 RemoteAddr" visible for the static baseline.

Use the absolute Makefile path to run the same gates from another working
directory. Verification resolves the checker relative to the loaded Makefile
rather than the caller's directory.
GitHub Actions installs the exact Go version from `go.mod` and runs formatting,
vet, race-enabled tests, module-integrity checks, and static guardrails. Tracked
buckets use a 10,000-entry LRU cap and fixed-length SHA-256 identifiers, while
length-prefixed component encoding keeps delimiter-containing key parts in
independent buckets.

When the required SDK or runtime is unavailable, use static checks and source review first, then verify on a machine that has the matching platform toolchain.

## Configuration and Secrets

- No required secret or credential file was identified in the repository scan. If you add integrations later, keep secrets out of git.

## Security and Privacy Notes

- Review changes touching authentication or token handling; examples from the scan include config/config.go, limiter.go.
- Proxy header behavior is caller-configured through `Limiter.IPLookups`; do not change lookup order semantics without tests and documentation.
- Blank X-Forwarded-For entries are skipped before limiter keys are derived,
  so malformed leading commas cannot produce an empty IP key.
- Blank or padded X-Real-IP values are trimmed or skipped before limiter keys
  are derived, allowing later configured lookup sources to be used.
- Malformed proxy IP headers are skipped before limiter keys are derived,
  allowing later configured lookup sources to be used.
- `RemoteAddr` parsing supports IPv4 and IPv6 host:port values before deriving
  limiter keys.
- Malformed RemoteAddr values are skipped before limiter keys are derived,
  allowing later configured lookup sources to be used.
- Configured header values only contribute keys when the request header contains one of those configured values.
- Configured header value matching inspects all request header values, so a
  blank first header value cannot hide a later configured match.
- Blank configured header values are skipped before limiter keys are derived.
- Matching duplicate configured header values produce one limiter key and one
  token charge per request.
- Configured header names are sorted before limiter keys are derived, while configured value order remains unchanged.
- Multi-key requests preflight every derived bucket so rejection by one bucket
  does not consume available capacity from another bucket.
- Rejected multi-key requests leave tracked-key and LRU state unchanged; missing
  buckets are created only after existing capacity succeeds.
- Blank header-only request values are skipped before limiter keys are derived.
- Empty method, header, and Basic Auth constraint collections fall back to the
  default remote-IP/path bucket instead of disabling rate limiting.
- Each limiter retains at most 10,000 tracked keys and evicts the least recently
  used bucket before admitting another request-derived key.
- `LimitReached` calls on directly configured valid limiters lazily initialize private accounting state with the same 10,000-key cap as `NewLimiter`.
- Limiter key accounting is serialized per limiter. Buckets are process-local and have no background cleanup; at the 10,000-key default cap, capacity pressure evicts the least-recently-used key, which starts with a fresh bucket if admitted again.

## Maintenance Notes

- See `SECURITY.md` for vulnerability reporting and safe research guidance.
- See `VISION.md` for project direction and contribution guardrails.
- Run `make lint`, `make test`, `make build`, and `make check` before pushing
  limiter behavior, config, or import changes.
- See `docs/plans/2026-06-09-proxy-header-ip-validation.md` for the malformed
  proxy header IP validation guard.
- See `docs/plans/2026-06-09-make-gate-aliases.md` for local verification
  target guardrails.
- See `docs/plans/2026-06-09-malformed-remote-addr.md` for direct RemoteAddr
  validation.
- See `docs/plans/2026-06-14-empty-config-fallback.md` for empty constraint
  collection fallback behavior.
- See `docs/plans/2026-06-15-direct-limiter-construction.md` for safe private
  key-accounting initialization when callers construct a valid limiter directly.
- See `docs/plans/2026-06-09-header-blank-value-matching.md` for header value
  matching with blank leading request values.
- See `docs/plans/2026-06-09-header-blank-configured-values.md` for blank
  configured header value handling.
- See `docs/plans/2026-06-13-deduplicate-header-values.md` for duplicate
  configured header value handling.
- See `docs/plans/2026-06-09-header-only-blank-request-values.md` for blank
  header-only request value handling.
- See `docs/plans/2026-06-10-ci-baseline.md` for the GitHub Actions baseline.
- See `docs/plans/2026-06-10-rate-limiter-key-cap.md` for bounded token-bucket
  storage and least-recently-used eviction.
- See `docs/plans/2026-06-12-bounded-key-encoding.md` for fixed-length storage
  identifiers and collision-safe component encoding.
- See `docs/plans/2026-06-12-ci-policy-hardening.md` for canonical hosted
  workflow enforcement.

## Contributing

Keep changes small and tied to the project that is already present in this repository. For code changes, document the toolchain used, avoid committing generated dependency directories or local configuration, and update this README when setup or verification steps change.
