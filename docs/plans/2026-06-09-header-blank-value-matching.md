---
title: Header Blank Value Matching
date: 2026-06-09
status: completed
execution: code
---

## Context

`matchingHeaderValues` checked `Header.Get` before reading all request header
values. In Go, a request can carry multiple values for one header key. If the
first value was blank and a later value matched the configured limiter value,
the early `Header.Get` check could skip the match entirely.

## Goals

- Match configured limiter header values across all request values.
- Preserve the existing behavior for missing headers and non-matching values.
- Add a regression test for a blank leading request header value.
- Extend the static baseline and docs for the header matching boundary.

## Implementation

- Read request header values with `Header.Values` before deciding whether the
  header is absent.
- Kept configured empty-value matching limited to present request headers.
- Added `TestBuildKeysHeaderValueMatchSkipsBlankFirstRequestValue`.
- Updated README, SECURITY, VISION, CHANGES, and `scripts/check-baseline.sh`.

## Verification

- `gofmt -w limiter.go limiter_test.go`
- `scripts/check-baseline.sh`
- `make lint`
- `make test`
- `make build`
- `make check`
- `git diff --check`
