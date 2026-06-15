# Make Rate-Limiter Header Key Order Deterministic

Status: Completed

## Context

`BuildKeys` ranges directly over `Limiter.Headers` in each header-bearing
constraint branch. Go map iteration order is intentionally unspecified, so the
same limiter configuration and request can return derived bucket keys in a
different order. That public output also feeds atomic batch accounting and LRU
recency updates, making equivalent header buckets receive nondeterministic
relative ordering.

## Requirements

- Derive configured header names in ascending lexical order before building
  request keys.
- Reuse the same ordered header list in method+header+auth, method+header, and
  header-only branches.
- Preserve configured header-value order, matching semantics, blank handling,
  duplicate-value suppression, key encoding, and atomic limiter accounting.
- Add deterministic unit coverage for the ordering helper and complete
  `BuildKeys` output.
- Add portable mutation-sensitive contracts and maintenance documentation.

## Implementation Units

### U1: Ordered header derivation

- Add a small internal helper in `limiter.go` that copies configured map keys
  and sorts the copy without mutating caller-owned configuration.
- Compute the ordered header list once per `BuildKeys` call and use it in every
  branch that evaluates headers.

### U2: Regression coverage

- Extend `limiter_test.go` with direct helper ordering coverage and an
  end-to-end multi-header `BuildKeys` assertion.
- Preserve all existing request and token-accounting tests.

### U3: Repository evidence

- Extend `scripts/check-baseline.sh` with helper, branch-reuse, regression,
  documentation, and completed-plan contracts.
- Update `AGENTS.md`, `README.md`, `SECURITY.md`, `VISION.md`, and `CHANGES.md`
  with the deterministic header-order boundary.

## Verification

- Run focused limiter tests, `go test -race ./...`, repository and
  external-directory `make check`, and `go vet ./...`.
- Reject hostile mutations that remove sorting, mutate the configured map,
  bypass the helper in any header branch, weaken either regression, remove
  documentation, or reopen the plan.
- Audit the exact diff, Go formatting, generated artifacts, credential
  patterns, conflict markers, and whitespace before commit.

## Risks

- Lexical ordering is an internal determinism rule, not HTTP header precedence;
  configured value order and matching remain unchanged.
- Header aliases that differ only by case remain distinct configuration keys;
  canonicalization or alias deduplication is outside this change.
- Existing stacked pull requests remain open and require explicit owner
  authorization before merge or closure.

## Verification Results

Completed on 2026-06-15:

- Focused deterministic helper and `BuildKeys` tests passed.
- `go test -race ./...` and `go vet ./...` passed across all packages.
- Repository and external-directory `make check` passed all package tests,
  module baseline contracts, formatting, vet, race, and build gates.
- Nine hostile mutations were rejected across sorting removal or reversal,
  caller-map mutation, shared-order bypass, both regressions, documentation,
  and completed-plan status.
- Exact diff, Go formatting, generated-artifact, credential-pattern,
  conflict-marker, module-integrity, and whitespace audits passed before
  commit.
- Configured value order and matching behavior remain unchanged, and header
  aliases that differ only by case remain distinct configuration keys.
