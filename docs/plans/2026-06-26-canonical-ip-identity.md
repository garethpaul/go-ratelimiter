# Canonical IP Identity Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use executing-plans to implement this plan task-by-task.

**Goal:** Ensure equivalent textual IP addresses consume the same rate-limit bucket.

**Architecture:** Canonicalize validated IPs inside the existing `libstring` parsing helpers. Keep limiter key construction and proxy lookup ordering unchanged, while behavior tests prove the canonical value reaches token-bucket accounting.

**Tech Stack:** Go 1.25.11, `net/http`, `net.IP`, `golang.org/x/time/rate`, GNU Make.

---

status: completed

### Task 1: Prove the bucket bypass

**Files:**
- Modify: `limiter_test.go`
- Modify: `libstring/libstring_test.go`

1. Send one request using expanded IPv6 text and another using the equivalent compressed form.
2. Require the second request to receive `429` from the same one-token bucket.
3. Require all supported IP lookup sources to return the same canonical text.
4. Run focused tests and confirm the new behavior fails before implementation.

### Task 2: Canonicalize validated addresses

**Files:**
- Modify: `libstring/libstring.go`

1. Parse each candidate once with `net.ParseIP`.
2. Return the parsed address's canonical `String()` representation.
3. Preserve empty results for malformed inputs and all lookup fallback behavior.
4. Rerun focused tests.

### Task 3: Preserve maintained contracts

**Files:**
- Modify: `scripts/check-baseline.sh`
- Modify: `README.md`
- Modify: `SECURITY.md`
- Modify: `VISION.md`
- Modify: `AGENTS.md`
- Modify: `CHANGES.md`

1. Require canonical source and regression fragments.
2. Document one identity across equivalent IP spellings.
3. Record red/green, full-gate, hosted, and review evidence.

### Task 4: Validate and merge

**Files:**
- Verify only.

1. Run focused tests, hostile mutation, `make check`, external Make, and clean diff checks.
2. Push a focused PR and attempt Codex review.
3. Merge only the exact final head after hosted checks pass.

## Verification Completed

- Red: the expanded/compressed middleware pair returned `204` twice, and all
  three supported lookup sources returned expanded IPv6 text unchanged.
- Green: focused and full tests pass on Go 1.25.11.
- Repository-root and external-directory `make check` pass with race tests,
  vet, module-integrity checks, and two canonical IP hostile mutations rejected.
- Implementation head `a839b8b3304feec5319fb91a02409cfd012b0c7b` passes hosted
  push and pull-request checks (`28245972824`, `28245975267`) and CodeQL
  (`28245973551`).
- `codex review --base origin/master` was attempted and skipped after HTTP 401
  authentication errors.
- Pending: evidence-only final-head hosted checks and merge verification.
