# Harden Hosted Workflow Policy

status: completed

## Context

The hosted gate was pinned and read-only, but checkout still persisted its token
and substring checks could be satisfied by additional unsafe workflow content.

## Changes

- Update checkout to the immutable v6.0.3 commit and set
  `persist-credentials: false`.
- Keep `check.yml` as the only hosted workflow.
- Compare the entire workflow with one canonical definition.
- Assign repository-wide ownership in `.github/CODEOWNERS`.
- Run the baseline on Go 1.25.11, which includes the standard-library fixes
  identified by `govulncheck` for the prior Go 1.25.3 toolchain.

## Verification

- `make check`
- `git diff --check`
- `go run golang.org/x/vuln/cmd/govulncheck@latest ./...`
- Isolated hostile workflow mutations covering persisted credentials, extra
  steps, permission drift, mutable actions, and hidden workflow files must fail.
