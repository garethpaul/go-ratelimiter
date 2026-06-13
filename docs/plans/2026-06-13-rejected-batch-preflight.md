# Keep Rejected Batch Preflight Side-Effect Free

status: planned

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

Pending implementation.

## Verification Completed

Pending implementation and verification.
