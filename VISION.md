## Go Rate Limiter Vision

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
- Maintain clear examples for proxy-aware IP lookup

Next priorities:

- Add tests for key derivation, token-bucket behavior, and proxy headers
- Document concurrency and cleanup behavior
- Add Go module metadata if the package is revived for modern tooling
- Clarify error responses and extension points

Contribution rules:

- One PR = one focused limiter, config, error, or documentation change.
- Keep default behavior dependency-free.
- Add tests for behavior changes.
- Update README examples when public API changes.

## Security

Rate limiting is often part of abuse prevention. Changes should be careful with
trusted proxy headers, basic-auth usernames, and header-derived keys so callers
do not accidentally trust attacker-controlled values.

## What We Will Not Merge (For Now)

- Remote storage dependencies as the default path
- Header/IP trust changes without documentation
- API-breaking changes without migration notes
- Behavior changes without tests
