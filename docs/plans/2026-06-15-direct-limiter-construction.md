# Support Direct Limiter Construction Safely

Status: Planned

## Context

`config.Limiter` is an exported configuration type with exported `Max` and
`TTL` fields, but its private token-bucket maps and LRU list are initialized
only by `NewLimiter`. A caller that directly constructs a valid limiter, such
as `&config.Limiter{Max: 1, TTL: time.Hour}`, reaches
`LimitReachedForKeys` and panics while assigning the first bucket to a nil map.

## Requirements

- Lazily initialize private token-bucket and LRU accounting state while the
  limiter mutex is held.
- Apply the same 10,000-key default cap used by `NewLimiter` when private state
  has not been initialized.
- Preserve constructor defaults, token refill behavior, atomic multi-key
  accounting, rejection preflight, LRU eviction, and invalid-config fail-closed
  behavior.
- Keep empty batches and invalid configurations free of bucket allocation.
- Add focused regression coverage plus portable mutation-sensitive contracts
  and synchronized maintenance guidance.

## Implementation Units

### U1: Private state initialization

- Add a small locked helper in `config/config.go` that initializes only missing
  private maps, the LRU list, and the default key cap.
- Invoke it after empty-batch and invalid-config preflight, before any private
  accounting state is read or written.

### U2: Regression coverage

- Add a direct-construction test that accepts the first request, rejects the
  second request, and verifies the bounded private accounting structures.
- Preserve all existing constructor-based tests and race coverage.

### U3: Repository evidence

- Extend `scripts/check-baseline.sh` with source, regression, guidance, and
  completed-plan contracts.
- Update `AGENTS.md`, `README.md`, `SECURITY.md`, `VISION.md`, and `CHANGES.md`
  with the direct-construction safety boundary.

## Verification

- Reproduce the pre-fix panic in an isolated copy.
- Run the focused config test, `go test -race ./...`, `go vet ./...`, and
  repository plus external-directory `make check`.
- Reject hostile mutations that remove initialization, move it before invalid
  preflight, omit any private structure, remove the regression, weaken the
  default cap, remove guidance, or reopen the plan.
- Audit the exact diff, Go formatting, generated artifacts, dependency drift,
  credential patterns, conflict markers, and whitespace before commit.

## Risks

- This change initializes private accounting state; it does not synthesize
  public middleware response fields such as `StatusCode` or `Message`.
- Callers should continue to prefer `NewLimiter` when they want all documented
  response defaults.
- PR #10 will remain stacked on open PR #9 and requires base-first ordering;
  neither pull request may be merged or closed without explicit authorization.

