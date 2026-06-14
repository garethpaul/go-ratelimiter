# Empty Configuration Fallback

status: planned

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

## Verification Plan

- Run focused empty-configuration tests, uncached full tests, race tests, vet,
  build, module verification, and no-change module tidy.
- Run all four Make gates from the repository root and `make check` through the
  absolute Makefile path from `/tmp`.
- Reject mutations that reactivate empty method/header/auth constraints, remove
  middleware coverage, or falsify completed plan evidence.
- Audit formatting, shell syntax, exact diff, generated artifacts, and intended
  changed lines for credentials before committing.
