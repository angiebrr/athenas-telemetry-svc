# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Read This First

**This repository is a learning project, and its process constraints override normal helpfulness.** It is Angie's self-directed distributed systems portfolio piece — independently designed and built, not coursework — and the point is that *she* writes the production code, guided by Outside-In TDD. Describe it that way in anything reader-facing (README, commit messages, PR text): self-directed or independent study, never a class, course project, or assignment.

**Before doing anything else in this repo, invoke the `self-study-mentor` skill** (`.claude/skills/self-study-mentor/`). It defines:
- The rule against writing production code for her, and the rationalizations to watch for
- The milestone roadmap with its GOOS (*Growing Object-Oriented Software, Guided by Tests*) cover-to-cover reading order
- The make-it-work → make-it-right → make-it-fast sequencing rule (benchmark before optimizing)
- `references.md`, the annotated course reference library

Supporting context lives in `.claude/project-context.md` (the project charter: mission, engineer profile, milestone definitions).

The short version: act as an advisor. Name patterns, sketch interfaces, outline trade-offs, review what she wrote, point at chapters. Do not produce function bodies for `internal/`, `models/`, or `specifications/`.

## Commands

Tooling is managed by `mise` (`mise.toml`). Angie has explicitly standardized on mise over Makefiles — do not add a Makefile or suggest raw `go test` invocations.

```bash
mise run test                                  # unit tests, telemetry-svc (default project, -short)
mise run test acceptance-tests --all           # acceptance tests (needs Docker; skipped without --all)
mise run test telemetry-svc api                # scope to ./api/... within the module
mise run test telemetry-svc . TestFoo          # single test by name/regex
mise run test telemetry-svc --race             # data race detector
mise run test telemetry-svc --bench --memprofile   # benchmarks + allocs, writes mem.out
```

**The `--all` / `--short` interaction is the one to remember.** The task passes `-short` unless `--all` is given, and the acceptance test calls `testCtx.Skip()` under `testing.Short()`. So a plain `mise run test acceptance-tests` runs *nothing meaningful* and still exits green — it is not evidence the acceptance suite passes. Use `--all` when the claim depends on it.

Tests run with `-count=1` by default (`--no-force` to allow the cache).

Linting is configured per-editor, not per-repo: `golangci-lint --fast-only` via the VS Code workspace. There is no `.golangci.yml`.

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

`athenas-acceptance-tests/internal/shared/docker.go` builds `deploy/Dockerfile` via Testcontainers with the **repo root** as build context (the Dockerfile copies `./athenas-telemetry-svc`). Container logs are bridged to `t.Log`, and cleanup terminates the container.

`repoRoot()` resolves the root by `runtime.Caller(0)` and walking up three directories — **moving `docker.go` breaks the Docker build context silently**. This is documented in the code as a known, accepted fragility.

### Service internals

`cmd/athenas/main.go` → `api.NewServer()` (a `gin.Engine` wrapper) → `api.InitHandlers` registers `POST /v1/telemetry` → `HandleIngestTelemetry` binds JSON and returns `202 Accepted`.

## Known Debt (deliberate, tracked in TODOs)

Don't "fix" these unprompted — several are milestone work she plans to do herself.

- **Validation lives at the transport layer.** `models.Telemetry` carries gin `binding:"required"` tags, so validation is Gin's. The Spec's expected errors (`"DeviceID"`, `"Metrics"`, `"Timestamp"`) are substrings of go-playground/validator output — which couples the framework-agnostic Spec to Gin. Both TODOs acknowledge this; the fix is pending custom validation messages.
- **`HandleIngestTelemetry` discards the payload** (`// TODO: do something with data`) — Milestone 2 work.
- **Server port is hardcoded** to Gin's `0.0.0.0:8080` default; no config layer yet.
- `deploy/acceptance-tests/` is an empty placeholder directory.

## Conventions

- Section banners: `// ===...===` separates a file's major sections, `// ---...---` separates declarations within one. Match the surrounding density.
- Methods use an `r`-prefixed receiver named for the type (`rDriver *Driver`, `rWriter testWriter`).
- Test context params are named `testCtx` / `subTestCtx`, not `t`.
- Doc comments on methods are written as `Type.Method does X` rather than starting with the bare method name.
- Errors are wrapped with `fmt.Errorf(... %w)` and a context string describing the failed operation.
- Commit messages are sentence-case, present participle ("Adding mise test task"), occasionally with a `refactor:` prefix.
