## Go Rate Limiter Vision

This document explains the current state and direction of the project.
Project overview and developer docs: [`README.md`](README.md)

Go Rate Limiter is a generic Go middleware for rate-limiting HTTP requests with
a token-bucket algorithm.

The repository is useful as a small reusable package that can limit by remote
IP, path, method, headers, and basic-auth usernames. Usage details live in
[`README.md`](README.md).

The goal is to keep the limiter simple, composable, and predictable under load.

The current focus is:

Priority:

- Preserve per-handler rate limiting through `LimitFuncHandler`
- Keep key selection behavior explicit and documented
- Avoid external storage requirements for the default limiter
- Refill `Max` tokens per `TTL` and fail closed for invalid limit configuration
- Keep middleware-owned response metadata single-valued and authoritative
- Maintain clear examples for proxy-aware and IPv6 RemoteAddr IP lookup
- Equivalent textual IP addresses share one canonical limiter identity across `RemoteAddr`, `X-Forwarded-For`, and `X-Real-IP`.
- Preserve configured header matching when a request has a blank first header value
- Preserve one token charge when duplicate configured header values match
- Preserve available bucket capacity when another key rejects a multi-key request
- Keep rejected multi-key preflight side-effect free for tracked bucket state
- Preserve header-only matching only for non-empty request header values
- Preserve default limiting for empty constraint collections
- `LimitReached` calls on directly configured valid limiters lazily initialize private accounting state with the same 10,000-key cap as `NewLimiter`.
- Limiter key accounting is serialized per limiter. Buckets are process-local and have no background cleanup; at the 10,000-key default cap, capacity pressure evicts the least-recently-used key, which starts with a fresh bucket if admitted again.
- Middleware rejections use the configured `StatusCode`, `MessageContentType`, and `Message`; callers needing extra headers or custom serialization should call `LimitByRequest` or `LimitByKeys` and write the returned `HTTPError` themselves.
- Limiter rejection status codes outside 400 through 599 fall back to 429; configured client and server error codes remain unchanged.
- Keep the Go module, `scripts/check-baseline.sh`, `make lint`, `make test`,
  `make build`, `make check`, and behavior tests passing

Contribution rules:

- One PR = one focused limiter, config, error, or documentation change.
- Keep default behavior dependency-free.
- Add tests for behavior changes.
- Update README examples when public API changes.

## Security

Canonical security policy and reporting:

- [`SECURITY.md`](SECURITY.md)

Rate limiting is often part of abuse prevention. Changes should be careful with
trusted proxy headers, basic-auth usernames, and header-derived keys so callers
do not accidentally trust attacker-controlled values.

Current baseline: `go.mod` and `go.sum` define the module, `make lint`,
`make test`, `make build`, and `make check` run `scripts/check-baseline.sh`,
and tests cover default key derivation, proxy-aware IP lookup, IPv6 RemoteAddr
parsing, configured header-value matching, blank X-Forwarded-For entries, blank
X-Real-IP values, and 429 responses when a bucket is empty. Cases with
malformed RemoteAddr values and malformed proxy IP headers are skipped before
limiter keys are derived, allowing later configured lookup sources to be used.
Configured header matching checks all request values so a blank first header
value cannot hide a later configured match.
The blank configured header values guard skips empty configured values before
limiter keys are derived.
The duplicate configured header values guard derives one key and charges one
token for each distinct matched value.
Configured header names are sorted before limiter keys are derived, while configured value order remains unchanged.
Limiter key accounting is serialized per limiter. Buckets are process-local and have no background cleanup; at the 10,000-key default cap, capacity pressure evicts the least-recently-used key, which starts with a fresh bucket if admitted again.
The blank header-only request values guard skips empty request header values
before limiter keys are derived.
Each limiter retains at most 10,000 request-derived keys and evicts the least
recently used bucket before admitting another. Stored identifiers are
fixed-length hashes of length-prefixed components, bounding retained key bytes
and keeping delimiter-containing values distinct.
Valid token buckets refill `Max` requests across each `TTL`; non-positive or
platform-unrepresentable limits reject requests without tracking keys.
Rate-limit metadata and rejection content types replace stale response values,
so repeated middleware application remains deterministic for HTTP clients.
Keep the exact guard phrases
"blank X-Forwarded-For", "blank X-Real-IP", "malformed RemoteAddr", and
"IPv6 RemoteAddr" visible for the static baseline, along with
"malformed proxy IP headers" and "blank header-only request values". GitHub Actions
runs formatting, vet, race-enabled tests, module-integrity checks, and static
guardrails using the Go version in `go.mod`.
The hosted workflow is pinned, read-only, credential-free after checkout, and
enforced as the repository's sole workflow by the local baseline.

## What We Will Not Merge (For Now)

- Remote storage dependencies as the default path
- Header/IP trust changes without documentation
- API-breaking changes without migration notes
- Behavior changes without tests

This list is a roadmap guardrail, not a permanent rule.
Strong user demand and strong technical rationale can change it.
