# Clarify Error Responses and Extension Points

Status: Completed

## Context

The middleware has two existing rejection-response extension paths. Callers can
set `Message`, `MessageContentType`, and `StatusCode` before wrapping a handler,
or they can call `LimitByRequest` or `LimitByKeys` directly and serialize the
returned `HTTPError` themselves. The defaults and this ownership boundary are
not documented together, and current tests only assert part of a customized
middleware response.

## Goals

- Document the default 429 plain-text rejection response.
- Document all three middleware-owned customization fields and prove they are
  emitted together without invoking the wrapped handler.
- Document that direct limiter functions return message and status only and do
  not write headers or bodies.
- Explain that callers needing structured bodies, extra headers, or logging can
  own serialization by using the direct functions.
- Add static contracts for tests, guidance, roadmap completion, and completed
  verification evidence.

## Non-Goals

- Do not add a callback, renderer, serializer, logger, or new public type.
- Do not change default status, message, content type, rate-limit headers, or
  token accounting.
- Do not define fallback behavior for invalid caller-supplied HTTP status codes.

## Implementation

1. Add focused middleware coverage for a custom status, JSON content type,
   custom body, and wrapped-handler short-circuit.
2. Add focused direct-function coverage for configured `HTTPError` values.
3. Add one canonical extension-point statement to repository guidance.
4. Extend `scripts/check-baseline.sh` to preserve the tests, source ownership
   boundary, guidance, and completed plan.
5. Update the changelog and remove the completed roadmap item.

## Validation

- Run the focused error-response tests, uncached package tests,
  `go test -race ./...`, `go vet ./...`, and `go mod tidy -diff`.
- Run repository and external-directory `make check`.
- Reject mutations that remove custom status, body, content type,
  short-circuiting, direct `HTTPError`, guidance, completed status, or
  verification evidence.
- Audit the exact diff, formatting, artifacts, dependency drift, credentials,
  conflict markers, modes, and whitespace before commit.

## Risks

- Documentation must not imply that `HTTPError` carries content type or that
  direct functions write to an `http.ResponseWriter`.
- Custom response fields are caller-owned configuration and should be set before
  the limiter is shared across concurrent handlers.
- PR #12 will be stacked on open PR #11 and requires base-first ordering;
  neither pull request may be merged or closed without explicit authorization.

## Verification Results

Completed on 2026-06-16:

- Focused error-response tests, uncached package tests, `go test -race ./...`,
  `go vet ./...`, `go mod tidy -diff`, and `go build ./...` passed.
- Repository and external-directory `make check` passed the complete formatting,
  test, vet, race, module-integrity, build, source-contract, guidance, and plan
  gates.
- Eight hostile mutations were rejected across custom status, content type,
  body, wrapped-handler short-circuiting, direct `HTTPError` values, guidance,
  completed status, and verification evidence.
- Exact diff, Go formatting, artifact, dependency, credential-pattern,
  conflict-marker, mode, and whitespace audits passed before commit.
