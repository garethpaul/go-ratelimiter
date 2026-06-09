# Header-Only Blank Request Values

status: completed

## Context

Header-only limiter configuration treats an empty configured value list as
"match any non-empty request value." The matcher already skipped missing
headers, but a present blank or whitespace-only request header could still
produce a limiter key.

## Completed Scope

- Updated header-only matching to require at least one non-empty request header
  value before returning a match.
- Added a focused `BuildKeys` regression test for a blank request value in
  header-only mode.
- Extended the static baseline and docs so blank header-only request values
  remain excluded from limiter keys.

## Verification

- `go test ./...`
- `scripts/check-baseline.sh`
- `make check`
- `git diff --check`
