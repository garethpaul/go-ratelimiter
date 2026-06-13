# Make Multi-Key Request Accounting Atomic

status: planned

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

Pending implementation.

## Verification Completed

Pending implementation and verification.
