# Go Rate Limiter Module Baseline Plan

status: completed

## Context

`go-ratelimiter` is a small reusable Go HTTP middleware package. The source uses GOPATH-era local imports and has no module metadata, lockfile, test command, or behavior tests, so modern Go cannot compile the package directly.

## Objectives

- Add Go module metadata and dependency locking for modern tooling.
- Update local imports to use the module path.
- Add focused tests for default key derivation, proxy-aware IP lookup, and request limiting.
- Add a reproducible `make check` baseline for formatting, tests, and import guardrails.
- Update README, VISION, SECURITY, and CHANGES with the new baseline.

## Verification

- `make check`
- `go test ./...`
- `git diff --check`
