---
title: Require Error-Class Rejection Statuses
type: bugfix
status: completed
date: 2026-06-18
execution: code
---

# Require Error-Class Rejection Statuses

## Status

Completed. Implementation, local verification, mutation checks, review, and
exact implementation-head hosted verification are complete.

## Context

The limiter blocks a request but currently preserves every configured status
from `200` through `999`. An isolated real-server reproduction showed that a
configured `200` returns a client-visible successful response containing the
limiter rejection body, while a configured `302` returns a redirect-class
response for the same rejected request.

A rejection must not appear successful, redirect a client, or use an
unclassified extension range. HTTP client and server error classes `400..599`
are the appropriate final response boundary. Custom `4xx` and `5xx` values
remain useful extension points and should continue to pass through unchanged.

## Goals

- Normalize configured rejection statuses outside `400..599` to the existing
  `429 Too Many Requests` fallback.
- Preserve configured standard and nonstandard `4xx` and `5xx` status codes.
- Keep direct `HTTPError` results, recorder behavior, and real-server behavior
  aligned at the shared status-normalization function.
- Add mutation-sensitive repository contracts and synchronized guidance.

## Non-Goals

- Do not change public types, constructors, messages, content types, headers,
  limiter accounting, token-bucket behavior, or successful request handling.
- Do not restrict accepted error statuses to registered IANA values.
- Do not merge or close this branch or its stacked predecessors.

## Implementation Units

### U1: Enforce the HTTP error-class boundary

Files:
- `limiter.go`
- `limiter_test.go`

Change the shared rejection-status normalization boundary to `400..599`.
Cover direct errors and middleware responses at representative lower, upper,
success, redirect, informational, and extension boundaries. Add real-server
coverage proving rejected requests cannot produce client-visible `2xx` or
`3xx` statuses.

### U2: Preserve the contract in project verification

Files:
- `scripts/check-baseline.sh`
- `AGENTS.md`
- `README.md`
- `SECURITY.md`
- `VISION.md`
- `CHANGES.md`

Replace the superseded `200..999` contract with the `400..599` error-class
boundary. Require the source condition, focused cases, real-server regressions,
guidance, and planned/completed plan evidence through structured checks.

## Verification Required

- Focused direct, recorder, and real-server rejection-status tests.
- Uncached and race-enabled package tests.
- `go vet ./...`, `go mod verify`, `go mod tidy -diff`, and `go build ./...`.
- Repository-root and external-directory `make check` plus all Make aliases.
- Isolated mutations for the lower and upper bounds, success/redirect cases,
  real-server coverage, guidance, and plan status/evidence.
- Exact diff, formatting, generated-artifact, file-mode, whitespace,
  conflict-marker, and credential-shaped addition audits.
- Successful canonical push and pull-request checks on the exact implementation
  head and final completed-plan evidence head.

## Risks

- Callers that intentionally configured `2xx`, `3xx`, or `600..999` statuses
  for rejected requests will now receive `429`; this is an intentional contract
  correction because those values do not represent HTTP client/server errors.
- Custom `4xx` and `5xx` codes remain supported and need explicit boundary
  regressions to prevent accidental narrowing.
- The change is stacked on PR #14 and requires base-first integration order.

## Verification Results

- Direct `HTTPError`, recorder, and real-server tests passed across
  informational, success, redirect, client-error, server-error, and extension
  boundaries. Rejected requests configured with `200`, `302`, or other
  non-error statuses now produce client-visible `429` responses.
- `go test -count=1 ./...`, `go test -race -count=1 ./...`, `go vet ./...`,
  `go mod verify`, `go mod tidy -diff`, and `go build ./...` passed.
- `make check`, `make lint`, `make test`, and `make build` passed from the
  repository, and absolute-Makefile `make check` passed from `/tmp`.
- Six isolated mutations were rejected across the lower and upper status
  bounds, real-server success coverage, per-table boundary registration,
  maintained guidance, and plan status.
- Plan-aware review found one mutation-sensitivity weakness in the static test
  contract. The checker now requires every shared boundary case exactly once
  in each direct and recorder table, and the targeted removal mutation fails.
- Exact diff, formatting, generated-artifact, executable-mode, whitespace,
  conflict-marker, and credential-shaped addition audits passed.
- Exact implementation head `1d421d39b7b3515e402f479f99b849a251acfb98`
  passed push run `27748640435` and pull-request run `27748660259`.
- Browser validation was not applicable because this repository is a Go HTTP
  middleware library with no browser route or UI. No live external service is
  required for the real-server regression.
