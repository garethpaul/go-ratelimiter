#!/usr/bin/env sh
set -eu

ROOT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
PLAN="$ROOT_DIR/docs/plans/2026-06-08-go-module-baseline.md"
IPV6_PLAN="$ROOT_DIR/docs/plans/2026-06-09-ipv6-remote-addr.md"
REAL_IP_PLAN="$ROOT_DIR/docs/plans/2026-06-09-real-ip-blank-values.md"
PROXY_IP_PLAN="$ROOT_DIR/docs/plans/2026-06-09-proxy-header-ip-validation.md"
MAKE_GATES_PLAN="$ROOT_DIR/docs/plans/2026-06-09-make-gate-aliases.md"
REMOTE_ADDR_PLAN="$ROOT_DIR/docs/plans/2026-06-09-malformed-remote-addr.md"
HEADER_BLANK_VALUE_PLAN="$ROOT_DIR/docs/plans/2026-06-09-header-blank-value-matching.md"
HEADER_BLANK_CONFIG_PLAN="$ROOT_DIR/docs/plans/2026-06-09-header-blank-configured-values.md"
HEADER_ONLY_BLANK_REQUEST_PLAN="$ROOT_DIR/docs/plans/2026-06-09-header-only-blank-request-values.md"
CI_WORKFLOW="$ROOT_DIR/.github/workflows/check.yml"
CI_PLAN="$ROOT_DIR/docs/plans/2026-06-10-ci-baseline.md"
KEY_CAP_PLAN="$ROOT_DIR/docs/plans/2026-06-10-rate-limiter-key-cap.md"
REFILL_PLAN="$ROOT_DIR/docs/plans/2026-06-12-token-bucket-refill-semantics.md"
KEY_ENCODING_PLAN="$ROOT_DIR/docs/plans/2026-06-12-bounded-key-encoding.md"
CI_POLICY_PLAN="$ROOT_DIR/docs/plans/2026-06-12-ci-policy-hardening.md"
HEADER_IDEMPOTENCE_PLAN="$ROOT_DIR/docs/plans/2026-06-12-idempotent-response-headers.md"
HEADER_DEDUPLICATION_PLAN="$ROOT_DIR/docs/plans/2026-06-13-deduplicate-header-values.md"
REJECTED_PREFLIGHT_PLAN="$ROOT_DIR/docs/plans/2026-06-13-rejected-batch-preflight.md"
LOCATION_INDEPENDENT_MAKE_PLAN="$ROOT_DIR/docs/plans/2026-06-13-location-independent-make.md"
EMPTY_CONFIG_PLAN="$ROOT_DIR/docs/plans/2026-06-14-empty-config-fallback.md"
DETERMINISTIC_HEADER_PLAN="$ROOT_DIR/docs/plans/2026-06-15-deterministic-header-key-order.md"
DIRECT_CONSTRUCTION_PLAN="$ROOT_DIR/docs/plans/2026-06-15-direct-limiter-construction.md"
CONCURRENCY_CLEANUP_PLAN="$ROOT_DIR/docs/plans/2026-06-16-concurrency-cleanup-contract.md"
ERROR_RESPONSE_PLAN="$ROOT_DIR/docs/plans/2026-06-16-error-response-extension-contract.md"
STATUS_CODE_SAFETY_PLAN="$ROOT_DIR/docs/plans/2026-06-16-rejection-status-code-safety.md"
FINAL_STATUS_SEMANTICS_PLAN="$ROOT_DIR/docs/plans/2026-06-16-final-rejection-status-semantics.md"

require_file() {
  path=$1
  if [ ! -f "$ROOT_DIR/$path" ]; then
    printf '%s\n' "Required file missing: $path" >&2
    exit 1
  fi
}

for path in \
  ".github/CODEOWNERS" \
  ".github/workflows/check.yml" \
  ".gitignore" \
  "AGENTS.md" \
  "CHANGES.md" \
  "Makefile" \
  "README.md" \
  "SECURITY.md" \
  "VISION.md" \
  "go.mod" \
  "go.sum" \
  "limiter.go" \
  "limiter_test.go" \
  "config/config.go" \
  "config/config_test.go" \
  "errors/errors.go" \
  "libstring/libstring.go" \
  "libstring/libstring_test.go" \
  "docs/plans/2026-06-08-go-module-baseline.md" \
  "docs/plans/2026-06-09-ipv6-remote-addr.md" \
  "docs/plans/2026-06-09-make-gate-aliases.md" \
  "docs/plans/2026-06-09-proxy-header-ip-validation.md" \
  "docs/plans/2026-06-09-real-ip-blank-values.md" \
  "docs/plans/2026-06-09-malformed-remote-addr.md" \
  "docs/plans/2026-06-09-header-blank-value-matching.md" \
  "docs/plans/2026-06-09-header-blank-configured-values.md" \
  "docs/plans/2026-06-09-header-only-blank-request-values.md" \
  "docs/plans/2026-06-10-ci-baseline.md" \
  "docs/plans/2026-06-10-rate-limiter-key-cap.md" \
  "docs/plans/2026-06-12-token-bucket-refill-semantics.md" \
  "docs/plans/2026-06-12-bounded-key-encoding.md" \
  "docs/plans/2026-06-12-ci-policy-hardening.md" \
  "docs/plans/2026-06-12-idempotent-response-headers.md" \
  "docs/plans/2026-06-13-deduplicate-header-values.md" \
  "docs/plans/2026-06-13-rejected-batch-preflight.md" \
  "docs/plans/2026-06-13-location-independent-make.md" \
  "docs/plans/2026-06-14-empty-config-fallback.md" \
  "docs/plans/2026-06-15-deterministic-header-key-order.md" \
  "docs/plans/2026-06-15-direct-limiter-construction.md" \
  "docs/plans/2026-06-16-concurrency-cleanup-contract.md" \
  "docs/plans/2026-06-16-error-response-extension-contract.md" \
  "docs/plans/2026-06-16-rejection-status-code-safety.md" \
  "docs/plans/2026-06-16-final-rejection-status-semantics.md" \
  "docs/plans/2026-06-08-header-value-matching.md"; do
  require_file "$path"
done

if ! grep -Fq 'ROOT := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))' "$ROOT_DIR/Makefile" ||
  ! grep -Fq '"$(ROOT)/scripts/check-baseline.sh"' "$ROOT_DIR/Makefile"; then
  printf '%s\n' "Makefile verification must resolve the checker from the loaded Makefile." >&2
  exit 1
fi

if ! grep -Fq "status: completed" "$LOCATION_INDEPENDENT_MAKE_PLAN" ||
  ! grep -Fq "from /tmp" "$LOCATION_INDEPENDENT_MAKE_PLAN" ||
  ! grep -Fq "root-derivation mutation failed" "$LOCATION_INDEPENDENT_MAKE_PLAN" ||
  ! grep -Fq "checker-command mutation failed" "$LOCATION_INDEPENDENT_MAKE_PLAN" ||
  ! grep -Fq "plan-status mutation failed" "$LOCATION_INDEPENDENT_MAKE_PLAN" ||
  ! grep -Fq "plan-evidence mutation failed" "$LOCATION_INDEPENDENT_MAKE_PLAN" ||
  ! grep -Fq "documentation mutation failed" "$LOCATION_INDEPENDENT_MAKE_PLAN" ||
  ! grep -Fq "absolute Makefile path" "$ROOT_DIR/README.md" ||
  ! grep -Fq "Made Go verification independent" "$ROOT_DIR/CHANGES.md"; then
  printf '%s\n' "Location-independent Make plan and guidance must record completed external verification." >&2
  exit 1
fi

makefile="$ROOT_DIR/Makefile"
if ! grep -Eq '^\.PHONY: .*build.*check.*lint.*test|^\.PHONY: .*build.*lint.*test.*check' "$makefile" ||
  ! grep -Fq "lint test build: check" "$makefile"; then
  printf '%s\n' "Makefile must expose lint, test, build, and check gate targets." >&2
  exit 1
fi

if ! grep -Fxq "go 1.25.11" "$ROOT_DIR/go.mod"; then
  printf '%s\n' "go.mod must retain the patched Go 1.25.11 toolchain baseline." >&2
  exit 1
fi

if command -v go >/dev/null 2>&1; then
  unformatted=$(find "$ROOT_DIR" -name '*.go' -not -path "$ROOT_DIR/.git/*" -print | xargs gofmt -l)
  if [ -n "$unformatted" ]; then
    printf '%s\n' "Go files need gofmt:" >&2
    printf '%s\n' "$unformatted" >&2
    exit 1
  fi
  (cd "$ROOT_DIR" && go vet ./...)
  (cd "$ROOT_DIR" && go test -race ./...)
  (cd "$ROOT_DIR" && go mod tidy -diff)
else
  printf '%s\n' "go is required for go-ratelimiter verification." >&2
  exit 1
fi

if git -C "$ROOT_DIR" grep -nE '^[[:space:]]*"(limiter|limiter/config|limiter/errors|limiter/libstring)"' -- '*.go'; then
  printf '%s\n' "Go imports must use module-qualified local paths." >&2
  exit 1
fi

if ! grep -Fq "TestBuildKeysDefaultUsesRemoteIPAndPath" "$ROOT_DIR/limiter_test.go" ||
  ! grep -Fq "TestBuildKeysCanPreferForwardedFor" "$ROOT_DIR/limiter_test.go" ||
  ! grep -Fq "TestBuildKeysHeaderValuesRequireConfiguredMatch" "$ROOT_DIR/limiter_test.go" ||
  ! grep -Fq "TestBuildKeysHeaderValueMatchIncludesConfiguredValue" "$ROOT_DIR/limiter_test.go" ||
  ! grep -Fq "TestBuildKeysHeaderValueMatchSkipsBlankFirstRequestValue" "$ROOT_DIR/limiter_test.go" ||
  ! grep -Fq "TestBuildKeysHeaderValueMatchSkipsBlankConfiguredValue" "$ROOT_DIR/limiter_test.go" ||
  ! grep -Fq "TestBuildKeysDeduplicatesConfiguredHeaderValues" "$ROOT_DIR/limiter_test.go" ||
  ! grep -Fq "TestBuildKeysHeaderOnlySkipsBlankRequestValue" "$ROOT_DIR/limiter_test.go" ||
  ! grep -Fq "TestBuildKeysMethodHeaderValueMatchIncludesConfiguredValue" "$ROOT_DIR/limiter_test.go" ||
  ! grep -Fq "TestBuildKeysSkipsMalformedRemoteAddr" "$ROOT_DIR/limiter_test.go" ||
  ! grep -Fq "TestBuildKeysFallsBackAfterMalformedRemoteAddr" "$ROOT_DIR/limiter_test.go" ||
  ! grep -Fq "TestLimitFuncHandlerReturnsTooManyRequestsAfterBucketIsEmpty" "$ROOT_DIR/limiter_test.go" ||
  ! grep -Fq "TestLimitHandlerChargesDuplicateConfiguredHeaderValueOnce" "$ROOT_DIR/limiter_test.go" ||
  ! grep -Fq "TestBuildKeysEmptyConstraintsFallBackToRemoteIPAndPath" "$ROOT_DIR/limiter_test.go" ||
  ! grep -Fq "TestBuildKeysIgnoresEmptyConstraintsBesideActiveFilters" "$ROOT_DIR/limiter_test.go" ||
  ! grep -Fq "TestLimitHandlerAppliesDefaultLimitWithMixedEmptyConstraints" "$ROOT_DIR/limiter_test.go" ||
  ! grep -Fq "TestLimitHandlerDoesNotPartiallyChargeRejectedMultiKeyRequest" "$ROOT_DIR/limiter_test.go" ||
  ! grep -Fq "TestRemoteIPTrimsForwardedForList" "$ROOT_DIR/libstring/libstring_test.go" ||
  ! grep -Fq "TestRemoteIPTrimsRealIP" "$ROOT_DIR/libstring/libstring_test.go" ||
  ! grep -Fq "TestRemoteIPFallsBackAfterBlankRealIP" "$ROOT_DIR/libstring/libstring_test.go" ||
  ! grep -Fq "TestRemoteIPFallsBackAfterMalformedRealIP" "$ROOT_DIR/libstring/libstring_test.go" ||
  ! grep -Fq "TestRemoteIPSkipsBlankForwardedForEntries" "$ROOT_DIR/libstring/libstring_test.go" ||
  ! grep -Fq "TestRemoteIPSkipsMalformedForwardedForEntries" "$ROOT_DIR/libstring/libstring_test.go" ||
  ! grep -Fq "TestRemoteIPFallsBackAfterBlankForwardedFor" "$ROOT_DIR/libstring/libstring_test.go" ||
  ! grep -Fq "TestRemoteIPSkipsMalformedRemoteAddr" "$ROOT_DIR/libstring/libstring_test.go" ||
  ! grep -Fq "TestRemoteIPFallsBackAfterMalformedRemoteAddr" "$ROOT_DIR/libstring/libstring_test.go" ||
  ! grep -Fq "TestRemoteIPHandlesIPv6RemoteAddr" "$ROOT_DIR/libstring/libstring_test.go"; then
  printf '%s\n' "Limiter and IP lookup behavior must stay covered by focused tests." >&2
  exit 1
fi

if ! grep -Fq "limitMethods := len(limiter.Methods) > 0" "$ROOT_DIR/limiter.go" ||
  ! grep -Fq "limitHeaders := len(limiter.Headers) > 0" "$ROOT_DIR/limiter.go" ||
  ! grep -Fq "limitBasicAuth := len(limiter.BasicAuthUsers) > 0" "$ROOT_DIR/limiter.go" ||
  grep -Fq "limiter.Methods != nil" "$ROOT_DIR/limiter.go" ||
  grep -Fq "limiter.Headers != nil" "$ROOT_DIR/limiter.go" ||
  grep -Fq "limiter.BasicAuthUsers != nil" "$ROOT_DIR/limiter.go"; then
  printf '%s\n' "Empty limiter constraint collections must preserve default rate limiting." >&2
  exit 1
fi

if ! grep -Fq "func (l *Limiter) LimitReachedForKeys" "$ROOT_DIR/config/config.go" ||
  ! grep -Fq "bucket.TokensAt(now) < 1" "$ROOT_DIR/config/config.go" ||
  ! grep -Fq "limiter.LimitReachedForKeys(encodedKeys)" "$ROOT_DIR/limiter.go" ||
  ! grep -Fq "TestLimiterDoesNotPartiallyConsumeBatchWhenOneKeyIsLimited" "$ROOT_DIR/config/config_test.go" ||
  ! grep -Fq "TestLimiterAllowsEmptyBatchWithoutTrackingKeys" "$ROOT_DIR/config/config_test.go"; then
  printf '%s\n' "Multi-key requests must preflight every bucket before consuming any tokens." >&2
  exit 1
fi

python3 - "$ROOT_DIR/config/config.go" "$ROOT_DIR/config/config_test.go" <<'PY'
import sys
from pathlib import Path

source = Path(sys.argv[1]).read_text()
tests = Path(sys.argv[2]).read_text()
batch = source[source.index("func (l *Limiter) LimitReachedForKeys"):source.index("func (l *Limiter) bucketForStorageKey")]
required = (
    "existingBuckets := make(map[string]*rate.Limiter, len(storageKeys))",
    "if bucket, found := l.tokenBuckets[storageKey]; found",
    "for _, bucket := range existingBuckets",
    "l.tokenBucketOrder.MoveToFront(l.tokenBucketEntries[storageKey])",
    "if !found {",
    "bucket = l.bucketForStorageKey(storageKey)",
)
if any(item not in batch for item in required):
    raise SystemExit("Rejected batches must preflight existing buckets before creating or touching tracked state.")
if not (batch.index("for _, bucket := range existingBuckets") <
        batch.index("l.tokenBucketOrder.MoveToFront") <
        batch.index("bucket = l.bucketForStorageKey(storageKey)")):
    raise SystemExit("Existing capacity preflight must precede LRU touches and missing-bucket creation.")
required_tests = (
    "TestLimiterRejectedBatchDoesNotAllocateMissingKeys",
    "TestLimiterRejectedBatchDoesNotEvictTrackedKeys",
    'bucketStorageKey("new")',
    'bucketStorageKey("unrelated")',
)
if any(item not in tests for item in required_tests):
    raise SystemExit("Rejected-batch state preservation must retain focused allocation and eviction regressions.")
PY

if ! grep -Fq "defaultMaxTrackedKeys = 10000" "$ROOT_DIR/config/config.go" ||
  ! grep -Fq "list.New()" "$ROOT_DIR/config/config.go" ||
  ! grep -Fq "MoveToFront" "$ROOT_DIR/config/config.go" ||
  ! grep -Fq "tokenBucketOrder.Back()" "$ROOT_DIR/config/config.go" ||
  ! grep -Fq "delete(l.tokenBuckets, oldestKey)" "$ROOT_DIR/config/config.go" ||
  ! grep -Fq "TestLimiterCapsTrackedKeys" "$ROOT_DIR/config/config_test.go" ||
  ! grep -Fq "TestLimiterEvictsLeastRecentlyUsedKey" "$ROOT_DIR/config/config_test.go"; then
  printf '%s\n' "Limiter keys must remain capped with recency-sensitive eviction coverage." >&2
  exit 1
fi

if ! grep -Fq 'l.Max <= 0 || l.TTL <= 0' "$ROOT_DIR/config/config.go" ||
  ! grep -Fq 'rate.Limit(float64(l.Max) / l.TTL.Seconds())' "$ROOT_DIR/config/config.go" ||
  ! grep -Fq "TestLimiterRefillsMaxTokensPerTTL" "$ROOT_DIR/config/config_test.go" ||
  ! grep -Fq "TestLimiterRejectsInvalidConfigurationWithoutTrackingKeys" "$ROOT_DIR/config/config_test.go"; then
  printf '%s\n' "Token buckets must refill Max tokens per TTL and reject invalid limits without tracking keys." >&2
  exit 1
fi

if ! grep -Fq "TestLimiterStoresBoundedKeyIdentifiers" "$ROOT_DIR/config/config_test.go" ||
  ! grep -Fq "sha256.Sum256" "$ROOT_DIR/config/config.go" ||
  ! grep -Fq "TestLimitByKeysKeepsDelimitedComponentsDistinct" "$ROOT_DIR/limiter_test.go" ||
  ! grep -Fq "func encodeKeys" "$ROOT_DIR/limiter.go"; then
  printf '%s\n' "Limiter keys must use bounded storage identifiers and collision-safe component encoding." >&2
  exit 1
fi

if ! grep -Fq "net.SplitHostPort" "$ROOT_DIR/libstring/libstring.go" ||
  ! grep -Fq "net.ParseIP(host)" "$ROOT_DIR/libstring/libstring.go" ||
  ! grep -Fq "ipAddrFromRemoteAddr(r.RemoteAddr); ip != \"\"" "$ROOT_DIR/libstring/libstring.go"; then
  printf '%s\n' "RemoteAddr parsing must use net.SplitHostPort and skip malformed IPs before deriving keys." >&2
  exit 1
fi

if ! grep -Fq "func ipAddrFromForwardedFor" "$ROOT_DIR/libstring/libstring.go" ||
  ! grep -Fq "net.ParseIP" "$ROOT_DIR/libstring/libstring.go" ||
  ! grep -Fq "ipAddrFromHeaderValue(part)" "$ROOT_DIR/libstring/libstring.go"; then
  printf '%s\n' "X-Forwarded-For parsing must skip blank and malformed entries before deriving keys." >&2
  exit 1
fi

if ! grep -Fq 'realIP := ipAddrFromHeaderValue(r.Header.Get("X-Real-IP"))' "$ROOT_DIR/libstring/libstring.go"; then
  printf '%s\n' "X-Real-IP parsing must trim whitespace and skip blank or malformed values before deriving keys." >&2
  exit 1
fi

if grep -Fq 'if r.Header.Get(headerKey) == ""' "$ROOT_DIR/limiter.go" ||
  ! grep -Fq "requestValues := r.Header.Values(headerKey)" "$ROOT_DIR/limiter.go" ||
  ! grep -Fq "if len(requestValues) == 0" "$ROOT_DIR/limiter.go"; then
  printf '%s\n' "Header value matching must inspect all request header values before rejecting a header." >&2
  exit 1
fi

if ! grep -Fq 'strings.TrimSpace(headerValue) == ""' "$ROOT_DIR/limiter.go"; then
  printf '%s\n' "Header value matching must skip blank configured header values before deriving keys." >&2
  exit 1
fi

if ! grep -Fq 'seenValues := make(map[string]struct{}, len(headerValues))' "$ROOT_DIR/limiter.go" ||
  ! grep -Fq 'if _, seen := seenValues[headerValue]; seen {' "$ROOT_DIR/limiter.go" ||
  ! grep -Fq 'seenValues[headerValue] = struct{}{}' "$ROOT_DIR/limiter.go"; then
  printf '%s\n' "Configured header values must be deduplicated before limiter keys are derived." >&2
  exit 1
fi

if ! grep -Fq 'strings.TrimSpace(requestValue) != ""' "$ROOT_DIR/limiter.go"; then
  printf '%s\n' "Header-only matching must skip blank request header values before deriving keys." >&2
  exit 1
fi

workflow_files=$(find "$ROOT_DIR/.github/workflows" -type f -print)
if [ "$workflow_files" != "$CI_WORKFLOW" ]; then
  printf '%s\n' "check.yml must remain the only hosted workflow." >&2
  exit 1
fi

expected_workflow=$(mktemp "${TMPDIR:-/tmp}/go-ratelimiter-check.XXXXXX")
trap 'rm -f "$expected_workflow"' EXIT HUP INT TERM
cat >"$expected_workflow" <<'EOF'
name: Check

on:
  pull_request:
  push:
  workflow_dispatch:

permissions:
  contents: read

concurrency:
  group: check-${{ github.workflow }}-${{ github.ref }}
  cancel-in-progress: true

jobs:
  check:
    runs-on: ubuntu-24.04
    timeout-minutes: 10
    steps:
      - name: Check out repository
        uses: actions/checkout@9f698171ed81b15d1823a05fc7211befd50c8ae0 # v6.0.3
        with:
          persist-credentials: false

      - name: Set up Go
        uses: actions/setup-go@4a3601121dd01d1626a1e23e37211e3254c1c06c # v6.4.0
        with:
          go-version-file: go.mod
          cache: true

      - name: Run baseline
        run: make check
EOF

if ! cmp -s "$expected_workflow" "$CI_WORKFLOW"; then
  printf '%s\n' "GitHub Actions must match the canonical pinned, credential-free Go race-test contract." >&2
  exit 1
fi

if [ "$(cat "$ROOT_DIR/.github/CODEOWNERS")" != "* @garethpaul" ]; then
  printf '%s\n' "CODEOWNERS must assign repository-wide ownership." >&2
  exit 1
fi
if ! grep -Fq "go test ./..." "$ROOT_DIR/README.md" ||
  ! grep -Fq "GitHub Actions" "$ROOT_DIR/README.md" ||
  ! grep -Fq "make lint" "$ROOT_DIR/README.md" ||
  ! grep -Fq "make test" "$ROOT_DIR/README.md" ||
  ! grep -Fq "make build" "$ROOT_DIR/README.md" ||
  ! grep -Fq "make check" "$ROOT_DIR/README.md" ||
  ! grep -Fq "IPv6 RemoteAddr" "$ROOT_DIR/README.md" ||
  ! grep -Fq "malformed RemoteAddr" "$ROOT_DIR/README.md" ||
  ! grep -Fq "malformed proxy IP headers" "$ROOT_DIR/README.md" ||
  ! grep -Fq "blank first header value" "$ROOT_DIR/README.md" ||
  ! grep -Fq "blank configured header values" "$ROOT_DIR/README.md" ||
  ! grep -Fq "duplicate configured header values" "$ROOT_DIR/README.md" ||
  ! grep -Fq "blank header-only request values" "$ROOT_DIR/README.md" ||
  ! grep -Fq "blank X-Forwarded-For" "$ROOT_DIR/README.md" ||
  ! grep -Fq "10,000 tracked keys" "$ROOT_DIR/README.md" ||
  ! grep -Fq "blank X-Real-IP" "$ROOT_DIR/README.md"; then
  printf '%s\n' "README must document the Go verification baseline." >&2
  exit 1
fi

if ! grep -Fq "scripts/check-baseline.sh" "$ROOT_DIR/VISION.md" ||
  ! grep -Fq "GitHub Actions" "$ROOT_DIR/VISION.md" ||
  ! grep -Fq "make lint" "$ROOT_DIR/VISION.md" ||
  ! grep -Fq "make test" "$ROOT_DIR/VISION.md" ||
  ! grep -Fq "make build" "$ROOT_DIR/VISION.md" ||
  ! grep -Fq "Go module" "$ROOT_DIR/VISION.md" ||
  ! grep -Fq "IPv6 RemoteAddr" "$ROOT_DIR/VISION.md" ||
  ! grep -Fq "malformed RemoteAddr" "$ROOT_DIR/VISION.md" ||
  ! grep -Fq "malformed proxy IP headers" "$ROOT_DIR/VISION.md" ||
  ! grep -Fq "blank first header value" "$ROOT_DIR/VISION.md" ||
  ! grep -Fq "blank configured header values" "$ROOT_DIR/VISION.md" ||
  ! grep -Fq "duplicate configured header values" "$ROOT_DIR/VISION.md" ||
  ! grep -Fq "blank header-only request values" "$ROOT_DIR/VISION.md" ||
  ! grep -Fq "blank X-Forwarded-For" "$ROOT_DIR/VISION.md" ||
  ! grep -Fq "10,000 request-derived keys" "$ROOT_DIR/VISION.md" ||
  ! grep -Fq "blank X-Real-IP" "$ROOT_DIR/VISION.md"; then
  printf '%s\n' "VISION must describe the current module baseline." >&2
  exit 1
fi

if ! grep -Fq "status: completed" "$PLAN"; then
  printf '%s\n' "Plan must be marked completed." >&2
  exit 1
fi

if ! grep -Fq "status: completed" "$ROOT_DIR/docs/plans/2026-06-08-header-value-matching.md"; then
  printf '%s\n' "Header value matching plan must be marked completed." >&2
  exit 1
fi

if ! grep -Fq "status: completed" "$IPV6_PLAN"; then
  printf '%s\n' "IPv6 RemoteAddr plan must be marked completed." >&2
  exit 1
fi

if ! grep -Fq "status: completed" "$ROOT_DIR/docs/plans/2026-06-09-forwarded-for-blank-entries.md"; then
  printf '%s\n' "Forwarded-for blank entry plan must be marked completed." >&2
  exit 1
fi

if ! grep -Fq "status: completed" "$REAL_IP_PLAN"; then
  printf '%s\n' "X-Real-IP blank value plan must be marked completed." >&2
  exit 1
fi

if ! grep -Fq "status: completed" "$PROXY_IP_PLAN"; then
  printf '%s\n' "Proxy header IP validation plan must be marked completed." >&2
  exit 1
fi

if ! grep -Fq "status: completed" "$MAKE_GATES_PLAN"; then
  printf '%s\n' "Make gate alias plan must be marked completed." >&2
  exit 1
fi

if ! grep -Fq "status: completed" "$REMOTE_ADDR_PLAN"; then
  printf '%s\n' "Malformed RemoteAddr plan must be marked completed." >&2
  exit 1
fi

if ! grep -Fq "status: completed" "$HEADER_BLANK_VALUE_PLAN"; then
  printf '%s\n' "Header blank value matching plan must be marked completed." >&2
  exit 1
fi

if ! grep -Fq "status: completed" "$HEADER_BLANK_CONFIG_PLAN"; then
  printf '%s\n' "Header blank configured value plan must be marked completed." >&2
  exit 1
fi

if ! grep -Fq "status: completed" "$HEADER_ONLY_BLANK_REQUEST_PLAN"; then
  printf '%s\n' "Header-only blank request value plan must be marked completed." >&2
  exit 1
fi

if ! grep -Fq "status: completed" "$CI_PLAN" ||
  ! grep -Fq "GitHub Actions" "$CI_PLAN" ||
  ! grep -Fq "make check" "$CI_PLAN"; then
  printf '%s\n' "CI baseline plan must record completed status and make check verification." >&2
  exit 1
fi

if ! grep -Fq "status: completed" "$KEY_CAP_PLAN" ||
  ! grep -Fq "Mutations disabling the cap or recency refresh must fail" "$KEY_CAP_PLAN"; then
  printf '%s\n' "Rate-limiter key-cap plan must record completed mutation verification." >&2
  exit 1
fi

if ! grep -Fq "status: completed" "$REFILL_PLAN" ||
  ! grep -Fq "Mutations restoring \`rate.Every(TTL)\`" "$REFILL_PLAN"; then
  printf '%s\n' "Token-bucket refill plan must record completed mutation verification." >&2
  exit 1
fi

if ! grep -Fq "status: completed" "$KEY_ENCODING_PLAN" ||
  ! grep -Fq "Hostile key mutations" "$KEY_ENCODING_PLAN"; then
  printf '%s\n' "Bounded key encoding plan must record completed mutation verification." >&2
  exit 1
fi

if ! grep -Fq "status: completed" "$CI_POLICY_PLAN" ||
  ! grep -Fq "persist-credentials: false" "$CI_POLICY_PLAN" ||
  ! grep -Fq "hostile workflow mutations" "$CI_POLICY_PLAN"; then
  printf '%s\n' "CI policy hardening plan must record completed mutation verification." >&2
  exit 1
fi

completed_statuses=$(grep -c '^status: completed$' "$HEADER_IDEMPOTENCE_PLAN" || true)
all_statuses=$(grep -c '^status:' "$HEADER_IDEMPOTENCE_PLAN" || true)
header_verification=$(awk '
  /^## Verification Completed$/ { in_verification = 1; next }
  in_verification && /^## / { exit }
  in_verification { print }
' "$HEADER_IDEMPOTENCE_PLAN")

if [ "$completed_statuses" -ne 1 ] || [ "$all_statuses" -ne 1 ]; then
  printf '%s\n' "Idempotent response-header plan must record exactly one completed status." >&2
  exit 1
fi

for evidence in \
  'go test -race -count=1 ./...' \
  'push run `27393483036`' \
  'pull-request run `27393485015`' \
  'push run `27393504527`' \
  'CodeQL run `27402321986`' \
  'Mutations restoring `Header.Add`'; do
  if ! printf '%s\n' "$header_verification" | grep -Fq "$evidence"; then
    printf '%s\n' "Idempotent response-header plan must record actual completed verification." >&2
    exit 1
  fi
done

if printf '%s\n' "$header_verification" | grep -Eiq '(^|[^[:alnum:]_])(pending|todo|tbd|not run)([^[:alnum:]_]|$)'; then
  printf '%s\n' "Idempotent response-header verification must not contain placeholders." >&2
  exit 1
fi

dedupe_completed_statuses=$(grep -c '^status: completed$' "$HEADER_DEDUPLICATION_PLAN" || true)
dedupe_all_statuses=$(grep -c '^status:' "$HEADER_DEDUPLICATION_PLAN" || true)
dedupe_verification=$(awk '
  /^## Verification Completed$/ { in_verification = 1; next }
  in_verification && /^## / { exit }
  in_verification { print }
' "$HEADER_DEDUPLICATION_PLAN")

if [ "$dedupe_completed_statuses" -ne 1 ] || [ "$dedupe_all_statuses" -ne 1 ]; then
  printf '%s\n' "Header-value deduplication plan must record exactly one completed status." >&2
  exit 1
fi

for evidence in \
  'focused duplicate-value tests passed' \
  'all four Make gates passed' \
  'deduplication removal mutation failed' \
  'key-test removal mutation failed' \
  'middleware-test removal mutation failed' \
  'hosted pull-request and CodeQL snapshot'; do
  if ! printf '%s\n' "$dedupe_verification" | grep -Fq "$evidence"; then
    printf '%s\n' "Header-value deduplication plan must record actual completed verification." >&2
    exit 1
  fi
done

if printf '%s\n' "$dedupe_verification" | grep -Eiq '(^|[^[:alnum:]_])(pending|todo|tbd|not run)([^[:alnum:]_]|$)'; then
  printf '%s\n' "Header-value deduplication verification must not contain placeholders." >&2
  exit 1
fi

preflight_completed_statuses=$(grep -c '^status: completed$' "$REJECTED_PREFLIGHT_PLAN" || true)
preflight_all_statuses=$(grep -c '^status:' "$REJECTED_PREFLIGHT_PLAN" || true)
preflight_verification=$(awk '
  /^## Verification Completed$/ { in_verification = 1; next }
  in_verification && /^## / { exit }
  in_verification { print }
' "$REJECTED_PREFLIGHT_PLAN")

if [ "$preflight_completed_statuses" -ne 1 ] || [ "$preflight_all_statuses" -ne 1 ]; then
  printf '%s\n' "Rejected-batch preflight plan must record exactly one completed status." >&2
  exit 1
fi

for evidence in \
  'focused rejected-batch tests passed' \
  'All four Make gates passed' \
  'eager creation mutation failed' \
  'no-allocation test mutation failed' \
  'no-eviction test mutation failed' \
  'hosted pull-request and code-scanning snapshot'; do
  if ! printf '%s\n' "$preflight_verification" | grep -Fq "$evidence"; then
    printf '%s\n' "Rejected-batch preflight plan must record actual completed verification." >&2
    exit 1
  fi
done

if printf '%s\n' "$preflight_verification" | grep -Eiq '(^|[^[:alnum:]_])(pending|todo|tbd|not run)([^[:alnum:]_]|$)'; then
  printf '%s\n' "Rejected-batch preflight verification must not contain placeholders." >&2
  exit 1
fi

if ! grep -Fq "Rejected multi-key requests leave tracked-key and LRU state unchanged" "$ROOT_DIR/README.md" ||
  ! grep -Fq "Rejected multi-key preflight should not allocate, evict, or reorder tracked buckets" "$ROOT_DIR/SECURITY.md" ||
  ! grep -Fq "Keep rejected multi-key preflight side-effect free" "$ROOT_DIR/VISION.md" ||
  ! grep -Fq "Made rejected multi-key preflight leave tracked bucket state unchanged" "$ROOT_DIR/CHANGES.md" ||
  ! grep -Fq "Keep rejected multi-key preflight free of allocation and eviction side effects" "$ROOT_DIR/AGENTS.md"; then
  printf '%s\n' "Project guidance must document side-effect-free rejected-batch preflight." >&2
  exit 1
fi

empty_config_completed_statuses=$(grep -c '^status: completed$' "$EMPTY_CONFIG_PLAN" || true)
empty_config_all_statuses=$(grep -c '^status:' "$EMPTY_CONFIG_PLAN" || true)
empty_config_verification=$(awk '
  /^## Verification Completed$/ { in_verification = 1; next }
  in_verification && /^## / { exit }
  in_verification { print }
' "$EMPTY_CONFIG_PLAN")

if [ "$empty_config_completed_statuses" -ne 1 ] || [ "$empty_config_all_statuses" -ne 1 ]; then
  printf '%s\n' "Empty-configuration fallback plan must record exactly one completed status." >&2
  exit 1
fi

for evidence in \
  'Focused empty-constraint key and middleware tests passed' \
  'Uncached full tests, race tests, vet, build, module verification' \
  'All four Make gates passed' \
  'Makefile path passed from `/tmp`' \
  'Eight isolated mutations were rejected' \
  'generated-artifact inspection' \
  'changed-line credential scanning passed'; do
  if ! printf '%s\n' "$empty_config_verification" | grep -Fq "$evidence"; then
    printf '%s\n' "Empty-configuration fallback plan must record actual completed verification." >&2
    exit 1
  fi
done

if printf '%s\n' "$empty_config_verification" | grep -Eiq '(^|[^[:alnum:]_])(pending|todo|tbd|not run)([^[:alnum:]_]|$)'; then
  printf '%s\n' "Empty-configuration fallback verification must not contain placeholders." >&2
  exit 1
fi

if ! grep -Fq "Empty method, header, and Basic Auth constraint collections fall back" "$ROOT_DIR/README.md" ||
  ! grep -Fq "Empty constraint collections must not bypass" "$ROOT_DIR/SECURITY.md" ||
  ! grep -Fq "Preserve default limiting for empty constraint collections" "$ROOT_DIR/VISION.md" ||
  ! grep -Fq "Prevented empty limiter constraint collections from bypassing" "$ROOT_DIR/CHANGES.md" ||
  ! grep -Fq "Treat only non-empty method, header, and Basic Auth collections as active constraints" "$ROOT_DIR/AGENTS.md"; then
  printf '%s\n' "Project guidance must document empty-constraint fallback behavior." >&2
  exit 1
fi

for deterministic_header_contract in \
  '"sort"' \
  'headerKeys := sortedHeaderKeys(limiter.Headers)' \
  'func sortedHeaderKeys(headers map[string][]string) []string' \
  'keys := make([]string, 0, len(headers))' \
  'sort.Strings(keys)' \
  'return keys'; do
  if ! grep -Fq "$deterministic_header_contract" "$ROOT_DIR/limiter.go"; then
    printf '%s\n' "Deterministic header ordering must keep contract: $deterministic_header_contract" >&2
    exit 1
  fi
done
if [ "$(grep -Fc 'for _, headerKey := range headerKeys {' "$ROOT_DIR/limiter.go")" -ne 3 ]; then
  printf '%s\n' "Every header-aware BuildKeys branch must reuse the ordered header list." >&2
  exit 1
fi
if grep -Fq 'for headerKey, headerValues := range limiter.Headers {' "$ROOT_DIR/limiter.go"; then
  printf '%s\n' "BuildKeys must not range directly over configured header maps." >&2
  exit 1
fi
for deterministic_header_test in \
  'TestSortedHeaderKeysReturnsOrderedCopy' \
  'TestBuildKeysOrdersConfiguredHeadersDeterministically' \
  'sortedHeaderKeys mutated caller-owned configuration'; do
  if ! grep -Fq "$deterministic_header_test" "$ROOT_DIR/limiter_test.go"; then
    printf '%s\n' "Limiter tests must keep deterministic header regression: $deterministic_header_test" >&2
    exit 1
  fi
done
for deterministic_header_doc in AGENTS.md README.md SECURITY.md VISION.md CHANGES.md; do
  if ! grep -Fq "Configured header names are sorted before limiter keys are derived, while configured value order remains unchanged." "$ROOT_DIR/$deterministic_header_doc"; then
    printf '%s\n' "$deterministic_header_doc must document deterministic configured-header ordering." >&2
    exit 1
  fi
done
for deterministic_header_plan_contract in \
  'Status: Completed' \
  'go test -race ./...' \
  'make check' \
  'hostile mutations' \
  'Header aliases that differ only by case remain distinct'; do
  if ! grep -Fq "$deterministic_header_plan_contract" "$DETERMINISTIC_HEADER_PLAN"; then
    printf '%s\n' "Deterministic header-order plan must record completed evidence: $deterministic_header_plan_contract" >&2
    exit 1
  fi
done

for direct_construction_contract in \
  'l.initializeAccountingState()' \
  'func (l *Limiter) initializeAccountingState()' \
  'l.tokenBuckets = make(map[string]*rate.Limiter)' \
  'l.tokenBucketOrder = list.New()' \
  'l.tokenBucketEntries = make(map[string]*list.Element)' \
  'l.maxTrackedKeys = defaultMaxTrackedKeys'; do
  if ! grep -Fq "$direct_construction_contract" "$ROOT_DIR/config/config.go"; then
    printf '%s\n' "Direct limiter construction must keep contract: $direct_construction_contract" >&2
    exit 1
  fi
done
invalid_config_line=$(grep -n 'if l.Max <= 0 || l.TTL <= 0' "$ROOT_DIR/config/config.go" | cut -d: -f1)
initialize_state_line=$(grep -n 'l.initializeAccountingState()' "$ROOT_DIR/config/config.go" | cut -d: -f1)
if [ -z "$invalid_config_line" ] || [ -z "$initialize_state_line" ] || [ "$initialize_state_line" -le "$invalid_config_line" ]; then
  printf '%s\n' "Private limiter accounting state must initialize after invalid-config preflight." >&2
  exit 1
fi
for direct_construction_test in \
  'TestDirectlyConfiguredLimiterInitializesAccountingState' \
  'TestDirectlyConfiguredInvalidLimiterDoesNotInitializeAccountingState' \
  'second directly configured request unexpectedly bypassed the limit' \
  'invalid directly configured limiter initialized private accounting state'; do
  if ! grep -Fq "$direct_construction_test" "$ROOT_DIR/config/config_test.go"; then
    printf '%s\n' "Config tests must keep direct-construction regression: $direct_construction_test" >&2
    exit 1
  fi
done
for direct_construction_doc in AGENTS.md README.md SECURITY.md VISION.md CHANGES.md; do
  if ! grep -Fq '`LimitReached` calls on directly configured valid limiters lazily initialize private accounting state with the same 10,000-key cap as `NewLimiter`.' "$ROOT_DIR/$direct_construction_doc"; then
    printf '%s\n' "$direct_construction_doc must document safe direct key accounting." >&2
    exit 1
  fi
done
direct_construction_completed_statuses=$(grep -c '^Status: Completed$' "$DIRECT_CONSTRUCTION_PLAN" || true)
direct_construction_all_statuses=$(grep -c '^Status:' "$DIRECT_CONSTRUCTION_PLAN" || true)
direct_construction_verification=$(awk '
  /^## Verification Results$/ { in_verification = 1; next }
  in_verification && /^## / { exit }
  in_verification { print }
' "$DIRECT_CONSTRUCTION_PLAN")
if [ "$direct_construction_completed_statuses" -ne 1 ] || [ "$direct_construction_all_statuses" -ne 1 ]; then
  printf '%s\n' "Direct-construction plan must record exactly one completed status." >&2
  exit 1
fi
for direct_construction_evidence in \
  'pre-fix panic' \
  'go test -race ./...' \
  'external-directory `make check`' \
  'hostile mutations'; do
  if ! printf '%s\n' "$direct_construction_verification" | grep -Fq "$direct_construction_evidence"; then
    printf '%s\n' "Direct-construction plan must record completed verification: $direct_construction_evidence" >&2
    exit 1
  fi
done
if printf '%s\n' "$direct_construction_verification" | grep -Eiq '(^|[^[:alnum:]_])(pending|todo|tbd|not run)([^[:alnum:]_]|$)'; then
  printf '%s\n' "Direct-construction verification must not contain placeholders." >&2
  exit 1
fi

for concurrency_cleanup_source_contract in \
  'l.Lock()' \
  'defer l.Unlock()' \
  'oldest := l.tokenBucketOrder.Back()' \
  'delete(l.tokenBuckets, oldestKey)' \
  'l.tokenBucketOrder.Remove(oldest)' \
  'bucket = rate.NewLimiter(refillRate, int(l.Max))'; do
  if ! grep -Fq "$concurrency_cleanup_source_contract" "$ROOT_DIR/config/config.go"; then
    printf '%s\n' "Limiter concurrency and cleanup source contract missing: $concurrency_cleanup_source_contract" >&2
    exit 1
  fi
done
concurrency_cleanup_guidance='Limiter key accounting is serialized per limiter. Buckets are process-local and have no background cleanup; at the 10,000-key default cap, capacity pressure evicts the least-recently-used key, which starts with a fresh bucket if admitted again.'
for concurrency_cleanup_doc in AGENTS.md README.md SECURITY.md VISION.md; do
  if ! grep -Fq "$concurrency_cleanup_guidance" "$ROOT_DIR/$concurrency_cleanup_doc"; then
    printf '%s\n' "$concurrency_cleanup_doc must document limiter concurrency and cleanup behavior." >&2
    exit 1
  fi
done
concurrency_cleanup_completed_statuses=$(grep -c '^Status: Completed$' "$CONCURRENCY_CLEANUP_PLAN" || true)
concurrency_cleanup_all_statuses=$(grep -c '^Status:' "$CONCURRENCY_CLEANUP_PLAN" || true)
concurrency_cleanup_verification=$(awk '
  /^## Verification Results$/ { in_verification = 1; next }
  in_verification && /^## / { exit }
  in_verification { print }
' "$CONCURRENCY_CLEANUP_PLAN")
if [ "$concurrency_cleanup_completed_statuses" -ne 1 ] || [ "$concurrency_cleanup_all_statuses" -ne 1 ]; then
  printf '%s\n' "Concurrency-cleanup plan must record exactly one completed status." >&2
  exit 1
fi
for concurrency_cleanup_evidence in \
  'go test -race ./...' \
  'external-directory `make check`' \
  'hostile mutations'; do
  if ! printf '%s\n' "$concurrency_cleanup_verification" | grep -Fq "$concurrency_cleanup_evidence"; then
    printf '%s\n' "Concurrency-cleanup plan must record completed verification: $concurrency_cleanup_evidence" >&2
    exit 1
  fi
done
if printf '%s\n' "$concurrency_cleanup_verification" | grep -Eiq '(^|[^[:alnum:]_])(pending|todo|tbd|not run)([^[:alnum:]_]|$)'; then
  printf '%s\n' "Concurrency-cleanup verification must not contain placeholders." >&2
  exit 1
fi

for error_response_source_contract in \
  'Message:    limiter.Message' \
  'StatusCode: rejectionStatusCode(limiter.StatusCode)' \
  'w.Header().Set("Content-Type", limiter.MessageContentType)' \
  'w.WriteHeader(httpError.StatusCode)' \
  'w.Write([]byte(httpError.Message))'; do
  if ! grep -Fq "$error_response_source_contract" "$ROOT_DIR/limiter.go"; then
    printf '%s\n' "Limiter error-response source contract missing: $error_response_source_contract" >&2
    exit 1
  fi
done
for error_response_test in \
  'TestLimitHandlerUsesConfiguredRejectionResponse' \
  'rejection status' \
  'rejection content type' \
  'rejection body' \
  'wrapped handler calls' \
  'TestLimitByRequestReturnsConfiguredHTTPError' \
  'HTTPError status' \
  'HTTPError message'; do
  if ! grep -Fq "$error_response_test" "$ROOT_DIR/limiter_test.go"; then
    printf '%s\n' "Limiter tests must keep error-response extension contract: $error_response_test" >&2
    exit 1
  fi
done
error_response_guidance='Middleware rejections use the configured `StatusCode`, `MessageContentType`, and `Message`; callers needing extra headers or custom serialization should call `LimitByRequest` or `LimitByKeys` and write the returned `HTTPError` themselves.'
for error_response_doc in AGENTS.md README.md SECURITY.md VISION.md CHANGES.md; do
  if ! grep -Fq "$error_response_guidance" "$ROOT_DIR/$error_response_doc"; then
    printf '%s\n' "$error_response_doc must document limiter error-response extension points." >&2
    exit 1
  fi
done
error_response_completed_statuses=$(grep -c '^Status: Completed$' "$ERROR_RESPONSE_PLAN" || true)
error_response_all_statuses=$(grep -c '^Status:' "$ERROR_RESPONSE_PLAN" || true)
error_response_verification=$(awk '
  /^## Verification Results$/ { in_verification = 1; next }
  in_verification && /^## / { exit }
  in_verification { print }
' "$ERROR_RESPONSE_PLAN")
if [ "$error_response_completed_statuses" -ne 1 ] || [ "$error_response_all_statuses" -ne 1 ]; then
  printf '%s\n' "Error-response plan must record exactly one completed status." >&2
  exit 1
fi
for error_response_evidence in \
  'go test -race ./...' \
  'external-directory `make check`' \
  'hostile mutations'; do
  if ! printf '%s\n' "$error_response_verification" | grep -Fq "$error_response_evidence"; then
    printf '%s\n' "Error-response plan must record completed verification: $error_response_evidence" >&2
    exit 1
  fi
done
if printf '%s\n' "$error_response_verification" | grep -Eiq '(^|[^[:alnum:]_])(pending|todo|tbd|not run)([^[:alnum:]_]|$)'; then
  printf '%s\n' "Error-response verification must not contain placeholders." >&2
  exit 1
fi
for status_code_source_contract in \
  'func rejectionStatusCode(statusCode int) int {' \
  'statusCode < 200 || statusCode > 999' \
  'return http.StatusTooManyRequests' \
  'StatusCode: rejectionStatusCode(limiter.StatusCode)'; do
  if ! grep -Fq "$status_code_source_contract" "$ROOT_DIR/limiter.go"; then
    printf '%s\n' "Limiter rejection status-code source contract missing: $status_code_source_contract" >&2
    exit 1
  fi
done
for status_code_test_contract in \
  'TestLimitByRequestNormalizesRejectionStatusCode' \
  'TestLimitHandlerNormalizesRejectionStatusCode' \
  'TestLimitHandlerNormalizesInformationalStatusForRealServer' \
  'configured: 99' \
  'configured: 100, want: http.StatusTooManyRequests' \
  'configured: 199, want: http.StatusTooManyRequests' \
  'configured: 200, want: 200' \
  'configured: 999, want: 999' \
  'configured: 1000' \
  'configured: 799, want: 799'; do
  if ! grep -Fq "$status_code_test_contract" "$ROOT_DIR/limiter_test.go"; then
    printf '%s\n' "Limiter tests must keep rejection status-code safety contract: $status_code_test_contract" >&2
    exit 1
  fi
done
status_code_safety_guidance='Limiter rejection status codes below 200 or above 999 fall back to 429; final configured codes from 200 through 999 remain unchanged.'
for status_code_safety_doc in AGENTS.md README.md SECURITY.md VISION.md CHANGES.md; do
  if ! grep -Fq "$status_code_safety_guidance" "$ROOT_DIR/$status_code_safety_doc"; then
    printf '%s\n' "$status_code_safety_doc must document rejection status-code safety." >&2
    exit 1
  fi
done
status_code_safety_completed_statuses=$(grep -c '^Status: Completed$' "$STATUS_CODE_SAFETY_PLAN" || true)
status_code_safety_all_statuses=$(grep -c '^Status:' "$STATUS_CODE_SAFETY_PLAN" || true)
status_code_safety_verification=$(awk '
  /^## Verification Results$/ { in_verification = 1; next }
  in_verification && /^## / { exit }
  in_verification { print }
' "$STATUS_CODE_SAFETY_PLAN")
if [ "$status_code_safety_completed_statuses" -ne 1 ] || [ "$status_code_safety_all_statuses" -ne 1 ]; then
  printf '%s\n' "Rejection status-code safety plan must record exactly one completed status." >&2
  exit 1
fi
for status_code_safety_evidence in \
  'repository and external-directory `make check`' \
  'isolated status-code mutations were rejected' \
  'No public type or configuration mutation'; do
  if ! printf '%s\n' "$status_code_safety_verification" | grep -Fq "$status_code_safety_evidence"; then
    printf '%s\n' "Rejection status-code safety plan must record completed verification: $status_code_safety_evidence" >&2
    exit 1
  fi
done
if printf '%s\n' "$status_code_safety_verification" | grep -Eiq '(^|[^[:alnum:]_])(pending|todo|tbd|not run)([^[:alnum:]_]|$)'; then
  printf '%s\n' "Rejection status-code safety verification must not contain placeholders." >&2
  exit 1
fi
final_status_completed_statuses=$(grep -c '^Status: Completed$' "$FINAL_STATUS_SEMANTICS_PLAN" || true)
final_status_all_statuses=$(grep -c '^Status:' "$FINAL_STATUS_SEMANTICS_PLAN" || true)
final_status_verification=$(awk '
  /^## Verification Results$/ { in_verification = 1; next }
  in_verification && /^## / { exit }
  in_verification { print }
' "$FINAL_STATUS_SEMANTICS_PLAN")
if [ "$final_status_completed_statuses" -ne 1 ] || [ "$final_status_all_statuses" -ne 1 ]; then
  printf '%s\n' "Final rejection-status semantics plan must record exactly one completed status." >&2
  exit 1
fi
for final_status_evidence in \
  'real-server regression' \
  'repository-root and external-directory `make check`' \
  'isolated final-status mutations were rejected' \
  'No public API or caller configuration was mutated'; do
  if ! printf '%s\n' "$final_status_verification" | grep -Fq "$final_status_evidence"; then
    printf '%s\n' "Final rejection-status plan must record completed verification: $final_status_evidence" >&2
    exit 1
  fi
done
if printf '%s\n' "$final_status_verification" | grep -Eiq '(^|[^[:alnum:]_])(pending|todo|tbd|not run)([^[:alnum:]_]|$)'; then
  printf '%s\n' "Final rejection-status verification must not contain placeholders." >&2
  exit 1
fi
printf '%s\n' "go-ratelimiter module baseline checks passed."
