# Scoped IPv6 Identity Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use executing-plans to implement this plan task-by-task.

**Goal:** Ensure valid scoped IPv6 client addresses are rate-limited instead of bypassing key construction.

**Architecture:** Replace the existing zone-blind `net.ParseIP` calls with one `net/netip` canonicalization helper shared by `RemoteAddr`, `X-Forwarded-For`, and `X-Real-IP`. Preserve lookup order and malformed-input fallback while retaining zone identifiers and unmapping IPv4-mapped addresses.

**Tech Stack:** Go 1.25.11, `net/http`, `net/netip`, table-driven and middleware tests.

---

status: completed

### Task 1: Prove the scoped-address bypass

**Files:**
- Modify: `libstring/libstring_test.go`
- Modify: `limiter_test.go`

1. Require all three IP lookup sources to canonicalize a scoped IPv6 address.
2. Require expanded and compressed scoped forms to share one exhausted bucket.
3. Run focused tests and confirm current parsing returns no identity.

### Task 2: Add zone-aware canonicalization

**Files:**
- Modify: `libstring/libstring.go`

1. Parse candidates with `netip.ParseAddr`.
2. Return `addr.Unmap().String()` to preserve existing mapped-IPv4 behavior.
3. Keep malformed inputs empty and preserve fallback order.
4. Rerun focused tests.

### Task 3: Lock contracts and verify

**Files:**
- Modify: `AGENTS.md`
- Modify: `README.md`
- Modify: `SECURITY.md`
- Modify: `VISION.md`
- Modify: `CHANGES.md`
- Modify: `scripts/check-baseline.sh`
- Modify: `docs/plans/2026-06-26-scoped-ipv6-identity.md`

1. Document scoped-address rate-limit identity behavior.
2. Add source, test, documentation, and completed-plan contracts.
3. Run focused tests, hostile parser mutation, `make check`, external Make, formatting, and diff checks.
4. Merge only an exact green reviewed head.

## Verification Completed

- Red: all three lookup sources returned no scoped identity (or skipped to a later forwarded address), and equivalent scoped requests returned `204` twice.
- Green: focused parser and middleware tests pass on Go 1.25.11.
- The canonical-text and stripped-zone hostile mutations are rejected.
- Repository-root and external-directory `make check` pass with gofmt, vet, race-enabled tests, module-integrity checks, and static contracts.
- `git diff --check` and shell/Python syntax checks pass.
- Hosted checks and exact-head review remain required before merge.
