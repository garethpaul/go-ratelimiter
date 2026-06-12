# Go Rate Limiter CI Baseline

status: completed

## Context

The repository already has a `make check` baseline that runs `gofmt` checks,
`go test ./...`, import guardrails, and documentation-plan checks. The missing
guard was hosted CI that runs the same gate for pushes and pull requests.

## Changes

- Added a least-privilege GitHub Actions workflow that reads the exact Go
  toolchain from `go.mod`.
- Pinned checkout and setup-go by commit, cancelled superseded runs, and
  bounded execution with a timeout.
- Extended `make check` to run formatting, `go vet`, `go test -race ./...`, and
  `go mod tidy -diff`.
- Extended the static baseline and docs so the hosted gate remains part of the
  maintained project contract.

## Verification

- `make check`
- `git diff --check`
- Hosted Go 1.25.11 GitHub Actions run
