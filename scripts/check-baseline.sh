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

require_file() {
  path=$1
  if [ ! -f "$ROOT_DIR/$path" ]; then
    printf '%s\n' "Required file missing: $path" >&2
    exit 1
  fi
}

for path in \
  ".gitignore" \
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
  "docs/plans/2026-06-08-header-value-matching.md"; do
  require_file "$path"
done

makefile="$ROOT_DIR/Makefile"
if ! grep -Eq '^\.PHONY: .*build.*check.*lint.*test|^\.PHONY: .*build.*lint.*test.*check' "$makefile" ||
  ! grep -Fq "lint test build: check" "$makefile"; then
  printf '%s\n' "Makefile must expose lint, test, build, and check gate targets." >&2
  exit 1
fi

if command -v go >/dev/null 2>&1; then
  unformatted=$(find "$ROOT_DIR" -name '*.go' -not -path "$ROOT_DIR/.git/*" -print | xargs gofmt -l)
  if [ -n "$unformatted" ]; then
    printf '%s\n' "Go files need gofmt:" >&2
    printf '%s\n' "$unformatted" >&2
    exit 1
  fi
  (cd "$ROOT_DIR" && go test ./...)
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
  ! grep -Fq "TestBuildKeysMethodHeaderValueMatchIncludesConfiguredValue" "$ROOT_DIR/limiter_test.go" ||
  ! grep -Fq "TestBuildKeysSkipsMalformedRemoteAddr" "$ROOT_DIR/limiter_test.go" ||
  ! grep -Fq "TestBuildKeysFallsBackAfterMalformedRemoteAddr" "$ROOT_DIR/limiter_test.go" ||
  ! grep -Fq "TestLimitFuncHandlerReturnsTooManyRequestsAfterBucketIsEmpty" "$ROOT_DIR/limiter_test.go" ||
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

if ! grep -Fq "go test ./..." "$ROOT_DIR/README.md" ||
  ! grep -Fq "make lint" "$ROOT_DIR/README.md" ||
  ! grep -Fq "make test" "$ROOT_DIR/README.md" ||
  ! grep -Fq "make build" "$ROOT_DIR/README.md" ||
  ! grep -Fq "make check" "$ROOT_DIR/README.md" ||
  ! grep -Fq "IPv6 RemoteAddr" "$ROOT_DIR/README.md" ||
  ! grep -Fq "malformed RemoteAddr" "$ROOT_DIR/README.md" ||
  ! grep -Fq "malformed proxy IP headers" "$ROOT_DIR/README.md" ||
  ! grep -Fq "blank first header value" "$ROOT_DIR/README.md" ||
  ! grep -Fq "blank configured header values" "$ROOT_DIR/README.md" ||
  ! grep -Fq "blank X-Forwarded-For" "$ROOT_DIR/README.md" ||
  ! grep -Fq "blank X-Real-IP" "$ROOT_DIR/README.md"; then
  printf '%s\n' "README must document the Go verification baseline." >&2
  exit 1
fi

if ! grep -Fq "scripts/check-baseline.sh" "$ROOT_DIR/VISION.md" ||
  ! grep -Fq "make lint" "$ROOT_DIR/VISION.md" ||
  ! grep -Fq "make test" "$ROOT_DIR/VISION.md" ||
  ! grep -Fq "make build" "$ROOT_DIR/VISION.md" ||
  ! grep -Fq "Go module" "$ROOT_DIR/VISION.md" ||
  ! grep -Fq "IPv6 RemoteAddr" "$ROOT_DIR/VISION.md" ||
  ! grep -Fq "malformed RemoteAddr" "$ROOT_DIR/VISION.md" ||
  ! grep -Fq "malformed proxy IP headers" "$ROOT_DIR/VISION.md" ||
  ! grep -Fq "blank first header value" "$ROOT_DIR/VISION.md" ||
  ! grep -Fq "blank configured header values" "$ROOT_DIR/VISION.md" ||
  ! grep -Fq "blank X-Forwarded-For" "$ROOT_DIR/VISION.md" ||
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

printf '%s\n' "go-ratelimiter module baseline checks passed."
