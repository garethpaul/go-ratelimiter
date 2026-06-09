# Header Blank Configured Values

status: completed

## Context

Header value rate-limit keys include the configured header value when a request
matches. Explicit blank configured values can create keys with empty header
components, which makes configuration intent ambiguous and weakens the key
shape.

## Completed Scope

- Skipped blank configured header values before matching request header values.
- Added a focused `BuildKeys` regression test for blank configured header
  values.
- Extended the static baseline so the code guard, test, docs, and completed plan
  remain in place.
- Documented the behavior in README, VISION, SECURITY, and CHANGES.

## Verification

- `go test ./...`
- `./scripts/check-baseline.sh`
- `make lint`
- `make test`
- `make build`
- `make check`
- `git diff --check`
