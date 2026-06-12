# Token-Bucket Refill Semantics

status: completed

## Context

`Limiter.Max` is documented as the maximum number of requests per `TTL`, but
new buckets currently use `rate.Every(TTL)`. That refills one token per complete
TTL regardless of `Max`, so a limiter configured for 60 requests per minute
only recovers one request per minute after its initial burst.

`rate.Every` also treats non-positive durations as an infinite rate. A limiter
with `TTL <= 0` therefore allows every request, which fails open in an
abuse-prevention component.

## Priority

Token accounting is the package's core contract. Incorrect refill rates cause
valid traffic to be rejected long after a burst, while invalid durations can
silently disable protection.

## Prioritized Engineering Backlog

1. Correct refill math and fail closed for invalid limits now.
2. Make multi-key request accounting atomic so a rejected request cannot
   consume tokens from keys checked earlier in the same request.
3. Add optional metrics for tracked keys, evictions, and rejected requests if
   the package gains an observability surface.

## Requirements

- R1. A valid limiter must refill at `Max / TTL` tokens per unit of time while
  preserving a burst capacity of `Max`.
- R2. `Max <= 0` or `TTL <= 0` must reject requests without creating tracked
  key state.
- R3. Existing LRU capacity and concurrency safety must remain intact.
- R4. The public `NewLimiter(max, ttl)` and `LimitReached(key)` signatures must
  remain compatible.
- R5. Tests and repository contracts must detect regressions in refill math,
  invalid-configuration handling, and state allocation.
- R6. README, security guidance, vision, and change history must describe the
  corrected behavior and fail-closed contract.

## Implementation Units

### U1. Correct bucket construction

- **Files:** `config/config.go`
- Reject invalid `Max`/`TTL` values before state allocation and derive the
  token rate from the documented maximum-per-duration contract.

### U2. Add deterministic regressions

- **Files:** `config/config_test.go`, `scripts/check-baseline.sh`
- Inspect bucket rate and burst directly, verify invalid configurations reject
  repeatedly, and verify they do not allocate LRU entries.

### U3. Update maintenance documentation

- **Files:** `README.md`, `SECURITY.md`, `VISION.md`, `CHANGES.md`
- Record refill semantics, fail-closed invalid configuration, and verification
  commands.

## Scope Boundaries

- Do not change public method signatures.
- Do not add remote storage, persistence, or new dependencies.
- Do not redesign multi-key accounting in this change.
- Do not change request key derivation or HTTP response formatting.

## Verification

- `make check`
- `go test -race ./...`
- `go vet ./...`
- `go mod verify`
- `git diff --check`
- Mutations restoring `rate.Every(TTL)` or allowing invalid configurations to
  allocate buckets must fail the regression suite.

Completed on 2026-06-12 with `make check`, race-enabled tests, vet, module
verification, and diff hygiene checks passing.
