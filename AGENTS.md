# AGENTS.md

## Repository purpose

`garethpaul/go-ratelimiter` is a dependency-light HTTP middleware package that
derives request keys and applies in-memory token-bucket rate limits.

## Project structure

- `Makefile` - repository verification targets
- `scripts` - baseline checks and helper scripts
- `docs` - plans, notes, and generated README assets
- `go.mod` - Go module definition
- `config` - repository source or sample assets
- `errors` - repository source or sample assets
- `libstring` - repository source or sample assets

## Development commands

- Install dependencies: `go mod download`
- Full baseline: `make check`
- Lint/static analysis: `make lint`
- Tests: `make test`
- Build gate: `make build`
- Go test all packages: `go test ./...`
- Race-enabled tests: `go test -race ./...`
- Go vet all packages: `go vet ./...`
- Go build all packages: `go build ./...`
- If a command above skips because a platform toolchain is missing, verify on a machine with that SDK before claiming platform behavior is tested.

## Coding conventions

- Keep limiter behavior compatible with the documented Go 1.25.11 toolchain
  unless a reviewed compatibility change is intentional.
- Keep imports compatible with module path `github.com/garethpaul/go-ratelimiter`.
- Run gofmt on changed Go files and keep table-driven tests close to the package under change.

## Testing guidance

- Test-related files include `config/config_test.go`,
  `libstring/libstring_test.go`, and `limiter_test.go`.
- Hosted CI runs formatting, vet, race-enabled tests, module-integrity checks,
  and static policy checks through `make check`.
- Start with the narrowest relevant test or Make target, then run `make check` before handing off if the change is not documentation-only.
- Keep README verification notes in sync when commands, fixtures, or supported toolchains change.

## PR / change guidance

- Keep diffs focused on the requested repository and avoid unrelated modernization or formatting churn.
- Preserve public APIs, sample behavior, file formats, and documented environment variables unless the task explicitly changes them.
- Update tests, README notes, or docs/plans when behavior, security posture, or validation commands change.
- Call out skipped platform validation, legacy toolchain assumptions, and any risky files touched in the final summary.

## Safety and gotchas

- No required secret or credential file was identified in the repository scan. If you add integrations later, keep secrets out of git.
- Proxy header behavior is caller-configured through `Limiter.IPLookups`; do not change lookup order semantics without tests and documentation.
- Blank X-Forwarded-For entries are skipped before limiter keys are derived, so malformed leading commas cannot produce an empty IP key.
- Blank or padded X-Real-IP values are trimmed or skipped before limiter keys are derived, allowing later configured lookup sources to be used.
- Malformed proxy IP headers are skipped before limiter keys are derived, allowing later configured lookup sources to be used.
- `RemoteAddr` parsing supports IPv4 and IPv6 host:port values before deriving limiter keys.
- Keep request-derived storage bounded by both the 10,000-key LRU cap and the
  fixed-length SHA-256 storage identifier.
- Keep key component encoding length-prefixed so delimiter-containing values do
  not share buckets accidentally.
- Keep rejected multi-key preflight free of allocation and eviction side effects.
- Treat only non-empty method, header, and Basic Auth collections as active constraints.
- Configured header names are sorted before limiter keys are derived, while configured value order remains unchanged.
- `LimitReached` calls on directly configured valid limiters lazily initialize private accounting state with the same 10,000-key cap as `NewLimiter`.
- Limiter key accounting is serialized per limiter. Buckets are process-local and have no background cleanup; at the 10,000-key default cap, capacity pressure evicts the least-recently-used key, which starts with a fresh bucket if admitted again.
- Keep `check.yml` as the sole pinned, read-only workflow and disable persisted
  checkout credentials.

## Agent workflow

1. Inspect the README, Makefile, manifests, and the files directly related to the request.
2. Make the smallest source or docs change that satisfies the task; avoid generated, vendored, or local-environment files unless required.
3. Run the narrowest useful validation first, then `make check` or the documented package/platform gate when available.
4. If a required SDK, service credential, or external runtime is unavailable, record the skipped command and why.
5. Summarize changed files, commands run, and remaining risks or follow-up validation.
