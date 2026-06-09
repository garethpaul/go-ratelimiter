#!/usr/bin/env sh
set -eu

ROOT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
PLAN="$ROOT_DIR/docs/plans/2026-06-08-go-module-baseline.md"
IPV6_PLAN="$ROOT_DIR/docs/plans/2026-06-09-ipv6-remote-addr.md"
REAL_IP_PLAN="$ROOT_DIR/docs/plans/2026-06-09-real-ip-blank-values.md"

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
  "docs/plans/2026-06-09-real-ip-blank-values.md" \
  "docs/plans/2026-06-08-header-value-matching.md"; do
  require_file "$path"
done

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
  ! grep -Fq "TestBuildKeysMethodHeaderValueMatchIncludesConfiguredValue" "$ROOT_DIR/limiter_test.go" ||
  ! grep -Fq "TestLimitFuncHandlerReturnsTooManyRequestsAfterBucketIsEmpty" "$ROOT_DIR/limiter_test.go" ||
  ! grep -Fq "TestRemoteIPTrimsForwardedForList" "$ROOT_DIR/libstring/libstring_test.go" ||
  ! grep -Fq "TestRemoteIPTrimsRealIP" "$ROOT_DIR/libstring/libstring_test.go" ||
  ! grep -Fq "TestRemoteIPFallsBackAfterBlankRealIP" "$ROOT_DIR/libstring/libstring_test.go" ||
  ! grep -Fq "TestRemoteIPSkipsBlankForwardedForEntries" "$ROOT_DIR/libstring/libstring_test.go" ||
  ! grep -Fq "TestRemoteIPFallsBackAfterBlankForwardedFor" "$ROOT_DIR/libstring/libstring_test.go" ||
  ! grep -Fq "TestRemoteIPHandlesIPv6RemoteAddr" "$ROOT_DIR/libstring/libstring_test.go"; then
  printf '%s\n' "Limiter and IP lookup behavior must stay covered by focused tests." >&2
  exit 1
fi

if ! grep -Fq "net.SplitHostPort" "$ROOT_DIR/libstring/libstring.go"; then
  printf '%s\n' "RemoteAddr parsing must use net.SplitHostPort for host:port values." >&2
  exit 1
fi

if ! grep -Fq "func ipAddrFromForwardedFor" "$ROOT_DIR/libstring/libstring.go" ||
  ! grep -Fq "strings.TrimSpace(part)" "$ROOT_DIR/libstring/libstring.go"; then
  printf '%s\n' "X-Forwarded-For parsing must skip blank entries before deriving keys." >&2
  exit 1
fi

if ! grep -Fq 'realIP := strings.TrimSpace(r.Header.Get("X-Real-IP"))' "$ROOT_DIR/libstring/libstring.go"; then
  printf '%s\n' "X-Real-IP parsing must trim whitespace and skip blank values before deriving keys." >&2
  exit 1
fi

if ! grep -Fq "go test ./..." "$ROOT_DIR/README.md" ||
  ! grep -Fq "make check" "$ROOT_DIR/README.md" ||
  ! grep -Fq "IPv6 RemoteAddr" "$ROOT_DIR/README.md" ||
  ! grep -Fq "blank X-Forwarded-For" "$ROOT_DIR/README.md" ||
  ! grep -Fq "blank X-Real-IP" "$ROOT_DIR/README.md"; then
  printf '%s\n' "README must document the Go verification baseline." >&2
  exit 1
fi

if ! grep -Fq "scripts/check-baseline.sh" "$ROOT_DIR/VISION.md" ||
  ! grep -Fq "Go module" "$ROOT_DIR/VISION.md" ||
  ! grep -Fq "IPv6 RemoteAddr" "$ROOT_DIR/VISION.md" ||
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

printf '%s\n' "go-ratelimiter module baseline checks passed."
