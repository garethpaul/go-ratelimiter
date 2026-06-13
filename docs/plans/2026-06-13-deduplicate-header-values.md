# Deduplicate Configured Header Values

status: in_progress

## Context

Configured header values are iterated directly when request values are matched.
If configuration repeats the same value, `BuildKeys` emits the same limiter key
more than once. `LimitByRequest` then charges one token bucket repeatedly for a
single request and can reject the first request when `Max` is one.

## Priority

This is a request-accounting correctness issue at the public configuration
boundary. It is isolated, dependency-free, and can be fixed without changing
the behavior of distinct configured values or distinct header names.

## Scope

1. Deduplicate matching nonblank configured values while preserving their first
   configured order.
2. Add key-construction coverage proving duplicate configured values produce
   one key.
3. Add middleware coverage proving one request consumes one token and the next
   request is rejected when `Max` is one.
4. Extend the baseline with source and test contracts and synchronize project
   documentation.

## Verification Plan

- Run focused and full uncached Go tests, race tests, vet, module tidiness and
  verification, all four Make gates, formatting, shell syntax, diff checks, and
  intended-file secret scans.
- Remove deduplication, remove the key test, and remove the middleware test;
  each hostile mutation must fail.
- Push a stacked pull request and take one bounded exact-head workflow, check,
  and CodeQL snapshot without polling.

## Risk And Rollback

Only duplicate configured values change behavior. Distinct values retain their
existing order and bucket identity. Rollback restores repeated charging; there
is no data migration, external storage, or public type change.
