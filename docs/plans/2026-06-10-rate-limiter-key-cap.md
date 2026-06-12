# Rate Limiter Key Cap

status: completed

## Context

The limiter stores one token bucket for every distinct request key and never
removes entries. Request keys include remote addresses, paths, headers, and
basic-auth usernames, so rotating request-controlled values can grow process
memory for the lifetime of the service.

## Priority

The library is intended for HTTP request boundaries. Its default storage must
remain bounded even when callers expose high-cardinality or attacker-controlled
key dimensions.

## Implementation

- Cap each limiter at 10,000 tracked keys by default.
- Maintain least-recently-used ordering with `container/list` for O(1) access,
  refresh, and eviction.
- Evict the least recently used key before inserting beyond the cap.
- Preserve token-bucket behavior for retained keys.
- Add package tests for the hard cap, LRU refresh, and fresh buckets after
  eviction.
- Extend the module baseline and operational documentation.

## Verification

- `go test ./...`
- `go test -race ./...`
- `go vet ./...`
- `make check`
- `make lint`
- `make test`
- `make build`
- `git diff --check`
- Mutations disabling the cap or recency refresh must fail.
- Hosted Go race and module-integrity workflow.
