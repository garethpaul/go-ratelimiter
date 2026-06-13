# Make Multi-Key Request Accounting Atomic

status: completed

## Context

`LimitByRequest` checks and consumes each derived limiter key sequentially. If
an earlier key has capacity but a later key is exhausted, the request is
rejected after the earlier bucket has already lost a token. Repeated rejected
requests can therefore drain otherwise available buckets and deny unrelated
requests that share only those keys.

## Priority

This is a rate-limit accounting correctness issue on the request path. The
repository's earlier hardening plans identified it as the next limiter-level
improvement, and it can be addressed within the existing in-memory bucket and
locking model without changing dependencies or public configuration.

## Scope

1. Add a limiter operation that resolves a batch of bucket keys under one lock,
   verifies capacity at one timestamp, and consumes tokens only when every
   bucket can admit the request.
2. Route single-key and request-derived limiting through the batch operation
   while preserving existing HTTP errors and key encoding.
3. Add regression coverage proving rejection by one key does not consume a
   token from another key.
4. Extend the maintained baseline with source and test contracts and
   synchronize project documentation.

## Verification Plan

- Run focused and full uncached Go tests, race tests, vet, module tidiness and
  verification, all four Make gates, formatting, shell syntax, diff checks, and
  intended-file secret scans.
- Remove the atomic capacity preflight, remove the regression test, and restore
  sequential request limiting; each hostile mutation must fail.
- Push a stacked pull request and take bounded exact-head workflow, check, and
  CodeQL snapshots without an unbounded polling loop.

## Risk And Rollback

Accepted requests still consume one token from every derived bucket. Only
rejected multi-key requests change: they no longer consume partial capacity.
Rollback restores sequential accounting; there is no persistent data migration
or external storage change.

## Work Completed

- Added `LimitReachedForKeys` to deduplicate bucket identities, resolve them
  under one lock, preflight capacity at one timestamp, and consume only after
  every bucket can admit the request.
- Routed `LimitReached`, `LimitByKeys`, and `LimitByRequest` through the shared
  accounting path while preserving empty derived-key bypass behavior.
- Added config-level and middleware regressions for partial-consumption safety,
  plus an empty-batch compatibility regression.
- Extended the maintained static baseline and synchronized the README, vision,
  and change history.

## Verification Completed

- Focused atomic-accounting and empty-batch tests passed.
- Removing capacity preflight failed both config and middleware regressions.
- Restoring sequential request charging failed the middleware regression after
  the rejected combined request drained the otherwise available bucket.
- Removing the middleware regression tripped the maintained baseline contract.
- Uncached and race-enabled tests, vet, build, module verification, module
  tidiness, formatting, shell syntax, and diff checks passed.
- All four Make gates passed with the maintained baseline.
- Hosted pull-request and CodeQL evidence is recorded separately after push;
  this plan claims only the completed local verification above.
