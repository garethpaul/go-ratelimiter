# Bound Rate-Limiter Key Storage

status: completed

## Context

The limiter capped the number of tracked buckets, but request-controlled paths,
headers, and usernames could still leave arbitrarily large strings in each map
entry. Joining components with `|` also allowed distinct component sequences to
share one bucket.

## Changes

- Encode key components with explicit byte lengths before rate limiting.
- Hash encoded keys with SHA-256 before storing them in bucket and LRU maps.
- Keep the public limiter behavior unchanged except that previously colliding
  component sequences now receive independent buckets.
- Add regression coverage for delimiter collisions and one-megabyte input keys.

## Verification

- `go test -race ./...`
- `make check`
- Hostile key mutations restoring pipe joins, raw map keys, or variable-length
  stored identifiers must fail.
