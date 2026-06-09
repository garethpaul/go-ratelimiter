# Forwarded-For Blank Entry Guard

status: completed

## Context

`Limiter.IPLookups` lets callers choose trusted proxy headers before
`RemoteAddr`. When `X-Forwarded-For` began with a blank entry, the parser could
return an empty IP string, causing `BuildKeys` to skip limiter key construction
for that request.

## Completed Scope

- Added a focused `X-Forwarded-For` parser that returns the first non-empty
  trimmed entry.
- Let lookup continue to later sources, such as `RemoteAddr`, when the trusted
  forwarded header contains no usable address.
- Covered both the blank-leading entry and all-blank fallback cases in unit
  tests.
- Extended the baseline and docs to preserve the proxy-header guardrail.

## Verification

- `make check`
- `git diff --check`
