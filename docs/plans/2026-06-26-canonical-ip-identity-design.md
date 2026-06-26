# Canonical IP Identity Design

status: completed

## Problem

The limiter validates request IP addresses with `net.ParseIP` but returns the
original text. Equivalent addresses such as `2001:0db8:0:0:0:0:0:1` and
`2001:db8::1` therefore produce different request keys and independent token
buckets. A client able to vary a trusted IP source's spelling can bypass the
configured request budget.

## Evidence

- `libstring.ipAddrFromRemoteAddr` returns the parsed host text unchanged.
- `libstring.ipAddrFromHeaderValue` returns the trimmed header text unchanged.
- `BuildKeys` includes the returned IP string in every request-derived key.
- Existing tests validate IPv6 parsing but do not compare equivalent forms.

## Options Considered

1. Canonicalize in `BuildKeys`. This would protect the current limiter path but
   leave the public `RemoteIP` helper inconsistent and duplicate parsing.
2. Preserve textual identity and document trusted proxies. This leaves a
   concrete bucket-bypass path for every supported IP source.
3. Return `net.IP.String()` from both low-level IP parsers after validation.

## Decision

Use option 3. Canonicalize at the parsing boundary so `RemoteAddr`,
`X-Forwarded-For`, and `X-Real-IP` share one identity rule. Keep lookup order,
malformed-value fallback, and trusted-proxy configuration unchanged.

## Validation

- Add a middleware regression proving expanded and compressed IPv6 forms share
  one exhausted bucket.
- Add table-driven `RemoteIP` coverage for all three supported lookup sources.
- Add source, test, guidance, and plan contracts to the static baseline.
- Run Go 1.25.11 race tests, full Make, external Make, hostile mutation, hosted
  checks, and CodeQL.

## Verification Completed

- Before implementation, equivalent IPv6 spellings received separate buckets;
  the second middleware request returned `204` instead of `429`.
- Focused canonicalization and shared-bucket tests pass on Go 1.25.11.
- Repository-root and external-directory `make check` pass with race tests,
  vet, module-integrity checks, and two canonical IP hostile mutations rejected.
- Hosted checks, CodeQL, exact-head review, and merge verification are pending.
