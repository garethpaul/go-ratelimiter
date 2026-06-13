# Keep Rejected Batch Preflight Side-Effect Free

status: completed

## Context

`LimitReachedForKeys` now preflights token capacity before consuming, but it
resolves every bucket through `bucketForStorageKey` first. Missing keys are
created and existing keys are LRU-touched before another exhausted bucket can
reject the batch. A rejected request can therefore allocate or evict tracked
state even though it consumes no tokens.

## Priority

This is the highest-value remaining isolated limiter-state defect because
request-controlled rejected batches can churn the bounded key cache and remove
unrelated limiter history. The fix fits the existing lock, map, LRU, and token
bucket model without changing public configuration or accepted-request limits.

## Scope

1. Preflight all existing distinct buckets at one timestamp before creating or
   touching missing buckets.
2. Return immediately on exhausted existing capacity without allocating,
   evicting, reordering, or consuming bucket state.
3. Create missing buckets and consume all distinct keys only after preflight
   succeeds.
4. Add focused state-preservation and hostile-mutation contracts plus project
   guidance.

## Verification Plan

- Run focused and full uncached Go tests, race tests, vet, module tidiness, all
  four Make gates, formatting, shell syntax, `git diff --check`, and intended
  artifact and secret scans.
- Restore eager bucket creation, remove the no-allocation regression, and
  remove the no-eviction regression; each hostile mutation must fail.
- Push a stacked pull request and take one bounded exact-head workflow and
  code-scanning snapshot without polling.

## Risk And Rollback

Accepted requests still create missing buckets and consume one token from each
distinct key. Only rejected batches change: they leave tracked-key and LRU state
untouched. Rollback restores eager resolution; no persistent migration exists.

## Work Completed

- Split batch resolution into existing-bucket capacity preflight followed by
  accepted-request LRU touches and missing-bucket creation.
- Protected requested existing buckets from accepted-batch eviction by moving
  them to the LRU front before new keys are admitted.
- Added focused no-allocation and no-eviction regressions plus static,
  documentation, and completed-plan contracts.

## Verification Completed

- The focused rejected-batch tests passed.
- Uncached and race-enabled tests, vet, module tidiness, formatting, shell
  syntax, and `git diff --check` passed.
- All four Make gates passed.
- The eager creation mutation failed after restoring bucket creation before
  existing capacity preflight.
- The no-allocation test mutation failed after removing its regression.
- The no-eviction test mutation failed after removing its regression.
- Intended-file artifact and secret scans passed.
- The hosted pull-request and code-scanning snapshot is a post-push evidence
  step; its bounded exact-head result is recorded after the implementation
  commit.
