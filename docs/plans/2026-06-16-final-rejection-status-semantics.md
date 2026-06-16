# Final Rejection Status Semantics

Status: Completed

## Context

The limiter currently preserves every configured status in the `100..999`
range. That prevents `net/http` from panicking, but informational `1xx`
responses are not final responses. A real `httptest.Server` reproduction with
`StatusCode = 100` delivered the limiter rejection body with client-visible
status `200`, so the middleware can accidentally report a successful response
when the request was rejected.

## Goals

- Normalize configured informational statuses `100..199` to the existing
  `429 Too Many Requests` fallback.
- Preserve configured final status codes `200..999`, including nonstandard
  extension codes already supported by the package.
- Keep direct `HTTPError` results and middleware-owned responses aligned at the
  shared error-construction boundary.
- Add a real-server regression so recorder-only behavior cannot conceal the
  informational-response downgrade.

## Non-Goals

- Do not restrict configured final codes to registered IANA values or only the
  `4xx` and `5xx` classes.
- Do not change public types, construction APIs, messages, content types,
  headers, accounting, or default status behavior.
- Do not alter the prior out-of-range fallback for values below `100` or above
  `999`.

## Implementation

1. Raise the shared rejection-status lower bound from `100` to `200`.
2. Update direct and recorder-based boundary cases for informational values.
3. Add an `httptest.Server` regression proving a configured informational
   status produces a client-visible `429` rejection rather than `200`.
4. Synchronize maintained guidance and extend the portable baseline with
   source, behavior, guidance, and completed-plan contracts.

## Verification

- Focused direct, recorder, and real-server tests.
- Uncached and race-enabled package tests, `go vet`, module verification and
  tidiness, and package build.
- Repository-root and external-directory `make check`.
- Isolated mutations for the lower bound, direct behavior, recorder behavior,
  real-server behavior, guidance, and completed-plan evidence.
- Exact-path, formatting, generated-artifact, sensitive-value,
  conflict-marker, file-mode, and whitespace audits.

## Risks

- Direct callers that intentionally configured informational `HTTPError`
  statuses will now receive `429`; this is necessary to maintain the package's
  documented rejection contract and match client-visible middleware behavior.
- Final nonstandard three-digit codes must remain unchanged.
- This change is stacked on PR #13 and requires base-first ordering; neither
  pull request may be merged or closed without explicit authorization.

## Verification Results

- Focused direct, recorder, and real-server regressions passed; the real-server
  regression confirms configured informational status `100` produces a
  client-visible `429` rejection with the configured body.
- Uncached package tests, race-enabled tests, `go vet`, `go mod verify`,
  `go mod tidy -diff`, and `go build ./...` passed.
- Both repository-root and external-directory `make check` passed.
- The seven isolated final-status mutations were rejected across the lower and
  upper bounds, fallback, informational expectations, real-server regression,
  maintained guidance, and reopened plan status.
- No public API or caller configuration was mutated; final configured codes
  from `200` through `999` remain unchanged.
