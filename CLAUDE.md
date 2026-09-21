# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Read This First

**This repository is a learning project, and its process constraints override normal helpfulness.** It is Angie's self-directed distributed systems portfolio piece — independently designed and built, not coursework — and the point is that *she* writes the production code, guided by Outside-In TDD. Describe it that way in anything reader-facing (README, commit messages, PR text): self-directed or independent study, never a class, course project, or assignment.

**Before doing anything else in this repo, invoke the `self-study-mentor` skill** (`.claude/skills/self-study-mentor/`). It defines:
- The rule against writing production code for her, and the rationalizations to watch for
- The milestone roadmap with its GOOS (*Growing Object-Oriented Software, Guided by Tests*) cover-to-cover reading order
- The make-it-work → make-it-right → make-it-fast sequencing rule (benchmark before optimizing)
- `references.md`, the annotated course reference library
- `progress-log.md` — the living state (current milestone/phase, reading position, gates, open debt, decision record). **Claude keeps it current without being asked**, and republishes `tracker.html` alongside it.
- The rule to teach the *why*: name the industry term, the trade-off, where the pattern recurs at a later milestone, and how it gets asked in an interview

Supporting context lives in `.claude/project-context.md` (the project charter: mission, engineer profile, milestone definitions).

The short version: act as an advisor. Name patterns, sketch interfaces, outline trade-offs, review what she wrote, point at chapters. Do not produce function bodies for `internal/`, `models/`, or `specifications/`.

## Commands

Tooling is managed by `mise` (`mise.toml`). Angie has explicitly standardized on mise over Makefiles — do not add a Makefile or suggest raw `go test` invocations.

```bash
mise run test                                  # every project, everything (needs Docker)
mise run test --short                          # every project, -short (acceptance suite skips)
mise run test telemetry-svc                    # scope to one project -- the fast TDD inner loop
mise run test telemetry-svc internal/api       # scope to ./internal/api/... within the module
mise run test telemetry-svc . TestFoo          # single test by name/regex
mise run test telemetry-svc --race             # data race detector
mise run test telemetry-svc --bench --memprofile   # benchmarks + allocs, writes mem.out

mise run lint                                  # golangci-lint over both modules
mise run lint --fix                            # ...applying auto-fixes
mise run vulncheck                             # govulncheck over both modules
```

**The project argument defaults to `all`, and skipping is opt-in.** A bare `mise run test` fans out over every project and runs the container-backed acceptance suite — slow, but it cannot silently cover nothing. `--short` passes `-short`, which makes the acceptance test call `testCtx.Skip()`. The old default was the reverse (`-short` unless `--all` was given), and a plain `mise run test acceptance-tests` used to exit green having run nothing at all.

If Docker isn't running, the acceptance suite now **fails** rather than skipping. That's intended: a test that quietly passes when its environment is missing is not evidence of anything.

Tests run with `-count=1` by default (`--no-force` to allow the cache).

The project list lives once in `[vars] projects` in `mise.toml` and is templated into every task, so adding a third go project means editing one line.

Editor tooling (`dlv`, `gopls`) is declared in `mise.dev.toml`, not `mise.toml` — run `MISE_ENV=dev mise install` locally to get it. CI sets no `MISE_ENV`, so `jdx/mise-action` installs only the Go toolchain and golangci-lint.

Linting is configured repo-wide in `.golangci.yml` (golangci-lint **v2** schema; the version is pinned in `mise.toml` so the editor and CI agree). revive runs with `enable-all-rules`, minus a short disabled list documented inline. The `exported` rule is configured with `disableChecksOnMethods` specifically so the repo's `Type.Method does X` comment convention keeps working.

## Architecture

### Two Go modules, joined by a `replace` directive

| Path | Module | Role |
|---|---|---|
| `athenas-telemetry-svc/` | `github.com/angiebrr/athenas-telemetry-svc` | The service |
| `athenas-acceptance-tests/` | `athenas-telemetry-acceptance-tests` | Black-box acceptance suite |

There is **no `go.work`** despite the project notes describing "Go workspaces" — the acceptance module resolves the service via `replace github.com/angiebrr/athenas-telemetry-svc => ../athenas-telemetry-svc`, so it always tests the local working tree rather than a published version. Dependencies flow one way only: acceptance-tests imports the service; the service never imports the tests. Open `.vscode/athenas-telemetry-svc.code-workspace` to get both modules as separate roots.

### The Spec / Driver seam (GOOS outside-in)

This is the structural heart of the repo and spans three files across both modules:

1. **`athenas-telemetry-svc/specifications/telemetry_spec.go`** — the Spec ("the Rules"). Declares the `TelemetryIngester` interface and `TelemetryIngesterSpec(t, ingester)`, a table of behavioral cases written against that interface only. It lives in the service module and is deliberately **exported, not `internal/`**, so the acceptance module can import it.
2. **`athenas-acceptance-tests/internal/drivers/httpserver/httpserver_driver.go`** — a Driver ("the Plumbing"). Implements `TelemetryIngester` by making real HTTP calls, translating a non-`202` response into an `IngestError`.
3. **`athenas-acceptance-tests/tests/httpserver/athenas_server_test.go`** — wires a containerized server to the Driver and runs the Spec against it.

The payoff: one Spec, many Drivers. Future milestones add drivers (SQS, multi-pod) that satisfy the same interface, and in-process adapters can run the same Spec as fast unit tests. **New behavior goes into the Spec (red) before any `internal/` code changes.**

### Acceptance test container lifecycle

`athenas-acceptance-tests/internal/shared/docker.go` builds `deploy/Dockerfile` via Testcontainers with the **repo root** as build context (the Dockerfile copies `./athenas-telemetry-svc`). Docker **build** logs are bridged to `t.Log` via `BuildLogWriter`, and the running container's stdout/stderr via a `LogConsumer`; cleanup terminates the container. The consumer is mutex-guarded and its `stop` is registered *before* the container's cleanup so it runs *after* it (`t.Cleanup` is last-added-first-called) — that ordering is load-bearing, since terminating the container drains the shutdown lines worth keeping, and `t.Log` panics once a test completes.

`repoRoot()` resolves the root by `runtime.Caller(0)` and walking up three directories — **moving `docker.go` breaks the Docker build context silently**. This is documented in the code as a known, accepted fragility.

### Service internals

`cmd/athenas/main.go` builds a `data.TelemetryDataStorer` and passes it to `api.NewServer()` (a `gin.Engine` wrapper). `api.InitHandlers` registers two routes:

| Route | Handler | Success |
|---|---|---|
| `POST api.TelemetryPath` (`/v1/telemetry`) | `HandleIngestTelemetry` | `202 Accepted` |
| `GET api.QueryTelemetryPath` (`/v1/telemetry/:device_id`) | `HandleQueryTelemetry` | `200` + JSON array |

Three layers, dependencies flowing inward:

- **`internal/api`** — transport only. Binds/serializes, and translates domain errors to status codes: `telemetry.ValidateTelemetryError` → `400`, `data.DataNotFoundError` → `404`, anything else → `500` with the internal text scrubbed and the real error sent to `slog`.
- **`internal/telemetry`** — the domain rules. `Ingest` (→ `Validate` → store) and `Query` (device-ID check → store). Validation is deliberately **not** on `models.Telemetry` — `models` is exported, so a driver could otherwise call `Validate` client-side and pass the Spec without the server doing anything.
- **`internal/data`** — the storage port (`TelemetryDataStorer`) plus `InMemoryDataStore`, a mutex-guarded `map[string][]models.Telemetry`. `GetByDeviceID` returns a **deep** copy (the outer slice *and* each record's `Metrics`), because a shallow `slices.Clone` still leaks the metrics backing array to callers. `data_test.go` holds a contract test run against the port, so M4's Postgres store can be checked against the same suite.

`GET /v1/telemetry/` (empty device ID) is **not routable** — gin's radix tree won't bind `:device_id` to an empty segment, and there's no GET handler at `/v1/telemetry` to redirect to, so it 404s. Empty-device-ID validation is therefore a `telemetry.Query` unit test, not a Spec case: the HTTP driver structurally cannot express that request.

The same `TelemetrySpec` runs at three levels: against `internal/telemetry` directly (microseconds), against the gin engine via `httptest` (`internal/api/handler_test.go`, transport translation only), and against a container over real HTTP. `httpserver.Driver` **deliberately hardcodes** `/v1/telemetry` rather than importing `api.TelemetryPath` — it is a black-box client, and sharing the constant would let a route rename ship green.

## Known Debt (deliberate, tracked in TODOs)

Don't "fix" these unprompted — several are milestone work she plans to do herself. Claude keeps this list current as items are resolved; see the doc-maintenance note in the `self-study-mentor` skill.

- **The Spec matches on error prose, not error class.** `TelemetrySpec` asserts with `ErrorContains(err, "missing device ID")` / `"data not found"`, so rewording a sentinel in `internal/telemetry` breaks the container acceptance suite. The durable fix is error classification carried by `httpserver.Driver`, letting the Spec distinguish caller-fault from callee-fault across any transport without depending on the message text. Still open.
- **No config layer.** The port is Gin's `0.0.0.0:8080` default, and `GIN_MODE=release` is set as a bare `ENV` in the Dockerfile. Both are placeholders for real configuration, nominally M2 work — decide in M2 whether it lands here or moves out explicitly.
- **Device IDs are only checked for emptiness.** `telemetry.Query` rejects `""` and nothing else; no format, length, or charset rules. Tracked by TODOs in `handler.go` and `handler_test.go`.
- **Complexity linting is off.** `cyclomatic`, `cognitive-complexity`, and `function-length` are disabled in `.golangci.yml` because the whole codebase currently measures ≤6 on all three, so any conventional threshold could not fire. Revisit when the dispatch ring lands and set the limit from measurement, not folklore.

## CI

`.github/workflows/validation.yaml` runs four jobs on every push and PR to `main`: `tests`, `acceptance-tests` (Docker), `lint`, and `vulncheck`. Each starts with the local `./.github/actions/setup-go-env` composite action, which runs `jdx/mise-action` and restores one shared Go module/build cache keyed on `hashFiles('**/go.sum')` — mise-action caches the toolchain but not `GOMODCACHE`/`GOCACHE`.

CI sets `MISE_ENV: ci`, which overrides the committed `.miserc.toml` (`env = ["dev"]`) so `mise.dev.toml`'s editor tooling (`dlv`, `gopls`) is skipped. Jobs invoke `mise run <task>` so CI and local run identical commands — the one exception is lint, which uses `golangci/golangci-lint-action` with `install-mode: none` for inline PR annotations while still executing the mise-pinned binary.

**There is no deploy.** The image is built and exercised, never published. Continuous delivery is Milestone 7.

## Conventions

- Section banners: `// ===...===` separates a file's major sections, `// ---...---` separates declarations within one. Match the surrounding density.
- Methods use an `r`-prefixed receiver named for the type (`rDriver *Driver`, `rWriter testWriter`).
- Test context params are named `testCtx` / `subTestCtx`, not `t`.
- Doc comments on methods are written as `Type.Method does X` rather than starting with the bare method name.
- Errors are wrapped with `fmt.Errorf(... %w)` and a context string describing the failed operation.
- Commit messages are sentence-case, present participle ("Adding mise test task"), occasionally with a `refactor:` prefix.
