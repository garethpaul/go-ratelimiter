# Location-Independent Go Verification

status: planned

## Context

The maintained Go baseline passes from the checkout but fails when the
absolute Makefile is invoked from another working directory because the shell
checker path is resolved relative to the caller.

## Priority

This is the next isolated reliability gap because local and hosted automation
should be able to load the repository Makefile without first changing
directories. The checker already roots all Go commands internally.

## Scope

1. Derive the repository root from `MAKEFILE_LIST`.
2. Invoke `scripts/check-baseline.sh` through its repository-rooted path.
3. Add rooted-Makefile, completed-plan, external-run, and synchronized-guidance
   contracts.
4. Reject root, checker, plan-status, plan-evidence, and documentation
   mutations.
5. Preserve limiter behavior, tests, module files, and workflow policy.

## Verification Plan

- Run focused and full uncached tests, race tests, vet, module tidiness, all
  four Make gates, formatting, shell syntax, and `git diff --check`.
- Run all four Make gates from /tmp through the absolute Makefile path.
- Reject isolated hostile mutations for root derivation, checker invocation,
  plan status/evidence, and documentation.
- Inspect exact intended paths, secrets, and generated artifacts.

## Risk And Rollback

The change affects only verification path resolution. Rollback restores the
caller-relative recipe; no runtime state or persistent migration exists.

## Work Completed

Pending implementation.

## Verification Completed

Pending implementation and validation. Run `make check` before completion.
