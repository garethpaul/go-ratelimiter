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
- Maintain clear examples for proxy-aware and IPv6 RemoteAddr IP lookup
- Keep the Go module, `scripts/check-baseline.sh`, and behavior tests passing

Next priorities:

- Document concurrency and cleanup behavior
- Clarify error responses and extension points

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

Current baseline: `go.mod` and `go.sum` define the module, `make check` runs
`scripts/check-baseline.sh`, and tests cover default key derivation,
proxy-aware IP lookup, IPv6 RemoteAddr parsing, configured header-value
matching, blank X-Forwarded-For entries, blank X-Real-IP values, and 429
responses when a bucket is empty. Cases with malformed proxy IP headers are
skipped before limiter keys are derived.

## What We Will Not Merge (For Now)

- Remote storage dependencies as the default path
- Header/IP trust changes without documentation
- API-breaking changes without migration notes
- Behavior changes without tests

This list is a roadmap guardrail, not a permanent rule.
Strong user demand and strong technical rationale can change it.
