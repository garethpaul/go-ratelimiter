# Empty Configuration Fallback

status: completed

## Context

`Limiter.Methods` documents an empty list as limiting all methods, and
`Limiter.Headers` documents an empty map as skipping header checks. `BuildKeys`
currently distinguishes only nil from non-nil collections, so decoded or
programmatically initialized empty collections enter constrained branches,
derive no keys, and bypass rate limiting.

## Scope

- Treat methods, headers, and Basic Auth user constraints as active only when
  their collections contain entries.
- Preserve configured header keys with empty value lists as the existing
  "match any nonblank request value" behavior.
- Add key-construction and middleware regressions for empty slices/maps, alone
  and in mixed empty configurations.
- Extend the static baseline and repository guidance with mutation-sensitive
  contracts.

## Non-Goals

- Do not change token-bucket capacity, refill, LRU, or atomic batch semantics.
- Do not change proxy IP selection, path construction, or response headers.
- Do not redefine non-empty method, header, or Basic Auth filters.

## Verification Completed

- Focused empty-constraint key and middleware tests passed.
- Uncached full tests, race tests, vet, build, module verification, and
  no-change module tidy passed.
- All four Make gates passed from the repository root, and the absolute
  Makefile path passed from `/tmp`.
- Eight isolated mutations were rejected: nil-based method, header, and Basic
  Auth activation, default-key-test removal, mixed-filter-test removal,
  middleware-test removal, stale plan status, and missing plan evidence.
- `gofmt`, shell syntax, `git diff --check`, exact intended-path review,
  generated-artifact inspection, and changed-line credential scanning passed.
