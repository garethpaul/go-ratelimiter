# Document Limiter Concurrency and Cleanup

Status: Completed

## Context

The limiter serializes token-bucket accounting with its embedded mutex and
bounds request-derived state with a 10,000-key least-recently-used cache. The
source has no background cleanup loop: inactive buckets remain until capacity
pressure evicts the least-recently-used key, and an evicted key receives a new
full bucket if it appears again. These operational details are visible in the
implementation but are not stated as one stable public contract.

Operators also need an explicit reminder that limiter state is process-local.
Separate processes or replicas do not share counters, so this package alone
does not enforce one aggregate distributed limit.

## Goals

- Document serialized per-limiter accounting and the scope of the mutex.
- Document capacity-driven LRU eviction and the absence of background cleanup.
- State that an evicted key starts with a fresh bucket when admitted again.
- State that counters are process-local and are not coordinated across
  replicas.
- Add static contracts that reject removal of the new guidance or completed
  verification evidence.

## Non-Goals

- Do not change token consumption, refill, key derivation, locking, or eviction
  behavior.
- Do not add timers, goroutines, persistence, remote storage, or a distributed
  limiter interface.
- Do not change the public Go API or dependency graph.

## Implementation

1. Add one canonical concurrency-and-cleanup paragraph to operator-facing and
   contributor-facing guidance.
2. Mark the roadmap item complete while retaining distributed coordination as
   an explicit boundary.
3. Extend `scripts/check-baseline.sh` to require the plan, canonical guidance,
   and completed verification evidence.
4. Update the changelog.

## Validation

- Run `sh -n scripts/check-baseline.sh`.
- Run focused static contracts for the new plan and documentation.
- Run `go test ./...`, `go test -race ./...`, `go vet ./...`, and
  `go mod tidy -diff`.
- Run `make check` from the repository and through the absolute Makefile path
  from `/tmp`.
- Reject mutations that remove each concurrency, cleanup, process-local, or
  completed-plan guarantee.
- Audit the exact diff, formatting, generated artifacts, dependency drift,
  credential patterns, conflict markers, and whitespace before commit.

## Risks

- Documentation must describe the current implementation without promising
  fairness or a time-based retention guarantee that the code does not provide.
- The fresh-bucket behavior after eviction is an availability tradeoff under
  extreme key cardinality and must remain visible to operators.
- The pull request will be stacked on open PR #10 and requires base-first
  ordering; neither pull request may be merged or closed without explicit
  authorization.

## Verification Results

Completed on 2026-06-16:

- `sh -n scripts/check-baseline.sh`, `go test ./...`, `go test -race ./...`,
  `go vet ./...`, `go mod tidy -diff`, and `go build ./...` passed.
- Repository and external-directory `make check` passed the full formatting,
  test, vet, race, module-integrity, build, static-contract, guidance, and plan
  gates.
- Six hostile mutations were rejected across lock removal, cleanup guidance,
  process-local scope, fresh-bucket behavior, completed status, and verification
  evidence.
- Exact diff, formatting, generated-artifact, dependency, credential-pattern,
  conflict-marker, and whitespace audits passed before commit.
