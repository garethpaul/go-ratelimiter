# Go Rate Limiter CI Baseline

status: completed

## Context

The repository already has a `make check` baseline that runs `gofmt` checks,
`go test ./...`, import guardrails, and documentation-plan checks. The missing
guard was hosted CI that runs the same gate for pushes and pull requests.

## Changes

- Added `.github/workflows/check.yml` for GitHub Actions.
- Configured `actions/setup-go` to read the Go toolchain from `go.mod`.
- Ran `make check` in the hosted workflow.
- Extended the static baseline and docs so the hosted gate remains part of the
  maintained project contract.

## Verification

- `make check`
- `git diff --check`
