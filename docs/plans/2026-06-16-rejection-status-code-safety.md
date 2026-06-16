# Rejection Status Code Safety

Status: In Progress

## Context

The limiter copies caller-configured `StatusCode` values into `HTTPError`, and
middleware passes that value directly to `http.ResponseWriter.WriteHeader`.
Go's HTTP stack rejects status codes outside the three-digit range, so invalid
configuration can panic during a rate-limit rejection instead of returning the
existing safe default response.

## Goals

- Preserve every configured three-digit status code accepted by `net/http`.
- Fall back to the existing default `429 Too Many Requests` for values below
  100 or above 999.
- Apply the same normalized status to middleware-owned responses and direct
  `HTTPError` results.
- Add mutation-sensitive behavior, source, guidance, and completed-plan
  contracts.

## Non-Goals

- Do not restrict valid custom three-digit codes to the registered IANA range.
- Do not add construction errors, change public types, or mutate caller-owned
  configuration.
- Do not change messages, content types, headers, token accounting, or default
  behavior.

## Implementation

1. Normalize rejection status at the shared `limitError` construction point.
2. Add focused direct and middleware tests for low and high invalid values plus
   preservation of a nonstandard three-digit value.
3. Document the fail-safe boundary in maintained guidance and the changelog.
4. Extend the baseline checker with source, tests, guidance, and completed-plan
   contracts.

## Verification

- Focused and uncached package tests
- `go test -race -count=1 ./...`, `go vet ./...`, `go mod tidy -diff`, and
  `go build ./...`
- Repository-root and external-directory `make check`
- Isolated mutations for low/high bounds, shared normalization, direct and
  middleware regressions, valid three-digit preservation, guidance, and plan
  completion
- Exact-path, formatting, artifact, dependency, sensitive-value,
  conflict-marker, mode, and whitespace audits

## Risks

- Normalizing only middleware output would leave direct callers with a
  divergent contract; the shared error path must own the fallback.
- The fallback must not rewrite valid three-digit extension codes.
- PR #13 will be stacked on open PR #12 and requires base-first ordering;
  neither pull request may be merged or closed without explicit authorization.
