# Progress Log

Living state for `athenas-telemetry-svc`. **Claude maintains this file** — see "Keep the Log Current" in SKILL.md. Everything here is mutable and dated; SKILL.md holds only the durable method, so the skill stays reusable.

Published tracker: `.claude/skills/self-study-mentor/tracker.html` → https://claude.ai/artifact/4DEea8wrLBwpdkRu9mV1Qk (republish to this URL; don't create a new one).

Study notes: `.claude/skills/self-study-mentor/notes/` (shell `index.html` + one fragment per chapter under `chapters/`, fetched by hash route) → https://claude.ai/artifact/TN7oyZ8z4RVFyy6YSRQake. Republish to this URL with `files` for each chapter. Kept out of `docs/` on purpose: the repo is recruiter-facing, and Claude-written chapters there would blur what she built vs. what was explained. Chapters explain behaviour, diagrams, and trade-offs, and end in a behaviour-list assignment; no implementation code, pseudocode only at a high level (her ask, 2026-10-21).

---

## Current Position

- **Milestone:** M2 — In-Memory Concurrency & Buffering (phase 1 of 3 complete)
- **Branch:** `abreuer-M2-dispatch-ring` (off `main` after PR #5 merged as `5b40d90`)
- **GOOS:** ch. 7 in progress (M2's range is 6–8)
- **Educative:** Course 1, "What distributed systems achieve for us" (early)
- **Last verified:** 2026-10-09 — `mise run test telemetry-svc --race` green with Angie's uncommitted Spec change (polling via `EventuallyWithT`, rewritten `TelemetryIngester` doc). Lint 0 issues in both modules. Container acceptance suite **not** run against this change.

---

### M2 progress (as of 2026-09-21)

M2 was rescoped by Angie on 2026-08-26 into three ordered phases, because the charter's wording bolts a make-it-fast objective onto an early make-it-work milestone. Treat this order as settled; don't push the ring earlier.

| Phase | Scope | State |
|---|---|---|
| 1. Work | In-memory store + `GET /v1/telemetry/:device_id`, so ingested telemetry has an observable destination | **Done** — Spec green at all three levels, lint clean, `--race` clean |
| 2. Right | Dispatch ring / worker pool — the answer to making `202 Accepted` *true*, not a data-race fix (the store's mutex already handles that) | Not started |
| 3. Fast | Benchmark (`--bench --memprofile`), then the charter's 100k RPS / allocation work | Not started |

The seam refactor between phases 1 and 2 is **done** (merged in PR #5): handlers hold only `*telemetry.Service`, which owns its `data.Storer`. The ring goes behind that seam.

**Phase 2 order (2026-09-28):** decide what a nil `Ingest` promises (read-your-writes vs. accepted-only) → Spec red → ring in `internal/engine` with graceful shutdown (`main` has none today) → error classification in the driver (a full queue is `Ingest`'s first callee-fault). Decision 1 settled 2026-10-09: fire-and-forget (see Decision Record). Decision 2 settled 2026-10-09: shed when full. Next decisions: worker count vs. per-device ordering; queue capacity (deferred to phase 3 measurement).

**Reading position (2026-09-21):** GOOS ch. 7 in progress (M2's range is 6–8). Educative course 1 early — "What distributed systems achieve for us." The build is ahead of the theory after a surgery break; the agreed gate is that course 1 finishes before M3 opens, since SQS and at-least-once delivery are course 2's lab.

---

## Milestone History

| Milestone | State | Notes |
|---|---|---|
| M1 — Walking Skeleton | Complete 2026-08-26 | HTTP 202 + validation, Spec at three levels, four-job CI |
| M2 — In-Memory Concurrency | Active | Phase 1 done; see table above |
| M3–M7 | Not started | Gate: course 1 finishes before M3 opens |

---

## Agreed Gates

- **Course 1 before M3.** SQS and at-least-once delivery are course 2's lab; starting M3 on an unfinished course 1 is where "theory leads" stops being advisory.
- **Benchmark before optimizing.** M2 phase 3 only opens after `--bench --memprofile` has run. "I know it'll allocate" is not a measurement.
- ~~**Config decision before M2 exits.**~~ **Met 2026-09-26** — `internal/config` landed on `abreuer-M2-add-query` (port only). `GIN_MODE` folded in on 2026-09-26 via `ENV` (see Decision Record).

---

## Open Debt

Mirrors CLAUDE.md's Known Debt; kept here with dates and triggers.

- **Spec matches on error prose, not error class.** `ErrorContains(err, "missing device ID")` couples the container acceptance suite to sentinel wording. Durable fix is error classification carried by `httpserver.Driver`. Open since 2026-08-26. The domain unit tests moved to `ErrorIs` on the sentinels on 2026-09-26 (`863e66b`), so this now affects only the Spec and container suite.
- **Device IDs only checked for emptiness.** No format/length/charset rules. TODOs in `api_handler.go` and `api_handler_test.go`.
- ~~**`Ingest`'s success promise.**~~ **Resolved 2026-10-21 (`4a00682`).** The port doc now describes queue-and-retry semantics and the ingest case polls via `EventuallyWithT`. Both review gaps (the positive "will become queryable" promise, and that a nil return can still be lost) closed by her in `13f970a`.
- **Complexity linting off.** Whole codebase measures ≤6 on cyclomatic/cognitive/function-length, so no conventional threshold could fire. Revisit when the ring lands; set the limit from measurement.

---

## Planned: Observability Stack (Prometheus + Loki + Grafana)

Proposed 2026-10-21, Angie's idea (Grafana interview 2026-10-22). **Placement is hers to confirm**; the suggested split below keeps it from jumping ahead of the ring. These are TODOs, not code: she writes the compose file, config, and instrumentation herself, and any TODO comments inside `internal/` are hers to add.

**Suggested placement.** Instrumentation is *measurement*, so it belongs with M2 phase 3: the benchmark gate and "acceptable loss" both need numbers. Queue metrics need the queue, so they land after the ring. The compose stack is dev tooling and can land any time after the ring, without blocking anything. M7 still owns DORA metrics and dashboards-as-code in CI.

- [ ] **HTTP RED metrics via middleware.** Candidates from the gin docs (<https://gin-gonic.com/en/docs/server-config/metrics/>): `github.com/gin-contrib/openmetrics`, or hand-rolled `prometheus/client_golang` counters and histograms. Label by `c.FullPath()` (the route template), never the raw URL. `/v1/telemetry/:device_id` with real IDs would mean one series per device. Decide what the label is for unmatched routes (`FullPath()` is empty there).
- [ ] **Queue metrics from the engine.** Depth (gauge), accepted / shed counts (counters), drain duration at shutdown. These turn the at-most-once trade into a measured loss rate (Study Notes ch. 1 §7–§8).
- [ ] **`/metrics` on its own port.** The gin docs suggest 9090, which is also the Prometheus server's default port; in a compose file they'll collide. Pick another port (2112 is a common convention) and keep it off the public listener.
- [ ] **Structured logs.** Handlers use package-level `slog.Warn` with the default handler, and gin's own access log is unstructured. Decide between `slog`'s text handler (logfmt-shaped, which Loki's `| logfmt` parser reads) and the JSON handler, and whether gin's logger is replaced by an `slog` request logger. Decide in `main`, the composition root.
- [ ] **`docker-compose` stack:** service + Prometheus (scrapes `/metrics`) + Loki + a log shipper + Grafana with provisioned data sources. **Use Grafana Alloy as the shipper, not Promtail**: Promtail is deprecated and past end-of-life. Keep `device_id` out of Loki *labels* (stream cardinality); leave it in the log line and filter at query time.
- [ ] **One provisioned dashboard**, committed as JSON: request rate / error rate / latency (RED) plus queue depth and shed rate (saturation, from USE).
- [ ] **Test level.** Metrics are operational surface, not product behaviour, so they stay out of `TelemetrySpec`. An adapter-level test that a request increments the expected series is enough.

## Interview Question Bank

**This is the project's actual exit criterion.** The repo is the evidence; answering these fluently is the qualification. Each milestone adds questions. A good answer names the trade, the failure mode it doesn't protect against, and where you'd measure — not a definition.

Answers aren't stored here on purpose. The point is rehearsal out loud, and the reasoning gets captured at decision time in the Decision Record above, where it's honest.

### Cross-cutting — TDD and test design

Not milestone-scoped; these recur at every milestone and the answers get richer as the suite grows. **Terminology worth having straight** (it's the thing most candidates fumble): the *level* axis is GOOS ch. 1's acceptance / integration / unit. The *reusability* axis is separate — a test body written once and run against many implementations is a **contract test**; one written for a single implementation is just a unit test. Fowler's **solitary vs. sociable** splits those further: solitary isolates with doubles, sociable uses real collaborators. And a working stand-in implementation is a **fake**, not a mock (Meszaros: dummy / fake / stub / spy / mock).

- Walk me through the kinds of test in this repo. How many distinct kinds are there, and what does each buy that the others can't?
- `TelemetrySpec` runs at three levels. What does each level catch that the one beneath it can't? If the suite got too slow, which would you delete first, and what risk would you be accepting?
- `specifications.TelemetrySpec` and `data.DataStoreSpec` are both called "spec" and are not the same pattern. What's the difference?
- Why does `specifications` live in an exported, non-test package instead of a `_test.go` file? What does that cost you, and who else does this?
- Your validation table has seven cases; the Spec has one invalid-input case. Why don't all seven run at the acceptance level?
- When `telemetry.Query` is tested against a real in-memory store, is that a mock, a stub, or a fake? Why does the distinction matter here?
- Why outside-in rather than inside-out? What goes wrong if you build the store first?
- What does "red for the right reason" mean? Give a case where a test of yours failed for the *wrong* reason and how you noticed.
- Your acceptance cases need no reset between runs. How did you get that, and why is it better than a reset hook on the interface?
- A test was awkward to write and you changed the design instead of the test. Walk me through it.
- What's an ice-cream cone, and what in this repo is most likely to produce one?

### M1 — Walking Skeleton

- Why build a walking skeleton before any feature?
- What does your container acceptance test prove that the `httptest` one doesn't? What does it cost per run?
- Why does `httpserver.Driver` hardcode `/v1/telemetry` instead of importing `api.TelemetryPath`?
- Why is `Validate` not a method on `models.Telemetry`?
- There's no `go.work` despite the project notes saying "Go workspaces." What does the `replace` directive do instead, and why does the dependency only flow one way?

### M2 — In-Memory Concurrency

- What happens when your consumer is slower than your producer?
- How did you size that buffer? *(Little's Law: L = λW)*
- Your endpoint returns `202 Accepted` — what exactly have you promised the caller? What can still be lost?
- Buffered or unbuffered channel, and why? What does each do to tail latency under a burst?
- Who closes the channel, and why can't it be the senders?
- How do you know your store is thread-safe? Name two different failures and which tool catches each.
- What does `-race` actually prove? What does it miss?
- `GetByDeviceID` hands back a slice — what's the risk, and why is a mutex not enough?
- What's the production failure mode of an unbounded queue?
- You chose fire-and-forget over group commit. What did you trade away, and what would make you switch? How would you know how much telemetry you're actually losing?
- Which errors should a telemetry sender retry, and which should it drop? Why does exponential backoff need jitter?
- What is the end-to-end argument, and how does it apply to your `202`?
- When your queue is full you return `503`, not `429`. Why? When would `429` be the right answer?
- How do you test that the service sheds load? How did you make the queue full on purpose?
- Block, shed, or drop: what does each do to latency and memory under sustained overload?
- Why buffer in memory at all when you could write straight to SQS? What does each give up: latency, cost per request, durability of the `202`? Which production clients buffer in-process anyway?
- Why a mutex here rather than a single owning goroutine? When would you switch?

- Where does your service get its configuration, and in what precedence order? What happens if the config file is missing?
- Your first config test passed while production ignored every env var. Why, and what did you change so the test couldn't miss it again?
- What's the risk of logging the parsed config at startup?
- Why is the environment baked into your image as `prod`, and what does 12-factor say about that? Why is it `APP_ENV` and not `ENV`?
- Where does "prod means release mode" live, and why there rather than in the HTTP layer?
- A timestamp of `0` is now rejected. Why is Go's zero value a validation problem at a JSON boundary? What other options were there (pointer fields, `omitempty`, a presence check), and what do you lose by treating zero as "missing"?
- Your library only offers a global setting (gin's mode). Where did you put the call that changes it, and what did that choice hide or expose? What does gin's release mode actually change?

### Observability — seeded 2026-10-21 (Grafana interview)

- Why must `device_id` never be a Prometheus label? What exactly grows, and what does it cost? Where *can* it live instead?
- Loki indexes labels but not log content. What does that buy compared with a full-text index, and how does it change which labels you choose?
- RED vs. USE: which applies to your HTTP handlers, and which to your queue? What's the saturation signal for the ring?
- Histogram or summary for request latency? Why can't you average p99s across pods?
- Your queue sheds with `503`. Which metric and alert tells you that's happening, and at what rate should someone be woken up?
- Why serve `/metrics` on a separate port?

### M3–M7 — seeded as each milestone opens

- **M3:** at-least-once vs. exactly-once; why consumers must be idempotent; what a visibility timeout buys.
- **M4:** why not two-phase commit; what the outbox costs you in exchange for atomicity.
- **M5:** partition ordering vs. throughput; what consumer lag measures; why rebalancing is the hard part.
- **M6:** what linearizable means for rate-limit state; why coordination is expensive.
- **M7:** the four keys; what counts as a successful deploy; which two need incident data.

---

## Reading Assignments

Current, tied to the decision in front of her. Short — one at a time.

| Assigned | Source | For | State |
|---|---|---|---|
| 2026-09-21 | [LMAX Disruptor paper](https://lmax-exchange.github.io/disruptor/disruptor.html) | M2 phase 2. A lock-free ring buffer built around cache-line padding, false sharing, and the single-writer principle — her particle-system optimization background applied directly to a queue. Gives her a production system to name. | Assigned |
| 2026-09-21 | *Learn Go with Tests* — Concurrency, Select, Sync, Context | M2 phase 2, Go idiom | Assigned |
| 2026-09-21 | Google SRE ch. 21, "Handling Overload" | M2 phase 2, load-shedding vocabulary | Assigned |
| 2026-10-21 | Study Notes ch. 1, "Bounded Queues and Worker Pools" (+ its assignment) | M2 phase 2. Synthesis of the three readings above with this repo's decisions; the assignment is the phase 2 build as a behaviour checklist. | Assigned |

---

## Decision Record

- **2026-08-26 — M2 rescoped to destination-first.** Angie pushed back on building the ring first: a dispatch ring with no destination has no externally observable behavior and can't be driven by an acceptance test. Three phases: work (store + GET), right (ring), fast (benchmark + allocation). Don't re-litigate.
- **2026-09-03 — Empty device ID moved out of the Spec.** `GET /v1/telemetry/` 404s; gin won't bind `:device_id` to an empty segment. The HTTP driver structurally cannot express that request, so it's a `telemetry.Query` unit test, not a system promise. (Option B of three; option C — query params instead of path params — stays open for when filtering is needed.)
- **2026-09-21 — Transport dropped out of Spec coverage.** `api_handler_test.go` no longer runs `TelemetrySpec`; it has hand-written status-code cases instead. Defensible (spec = behavior, handler tests = translation), but CLAUDE.md's "three levels" claim was stale and is corrected to two. Revisit if transport behavior starts drifting from the spec.
- **2026-09-21 — Doc maintenance delegated to Claude.** Prose upkeep isn't the learning target; stale docs are worse than none.
- **2026-09-26 — Mechanical test changes opened to Claude, on explicit request only.** Table-driven conversions and testable examples for code she already wrote. The learning is in deciding what to assert, not in typing out the table. The limits keep it mechanical: `_test.go` only, never the Spec, no added or dropped cases, suite green before and after. Any gap Claude spots gets named, and she writes the fix. Each use gets noted here so the record of what she wrote stays accurate. Public note is in `docs/self-study.md`. Same day, a working checklist was added to the skill, because a test refactor has no test of its own. It only runs on green. She defines the table's shape or the example's usage. Subtest names must match one-to-one. One case is broken on purpose to prove the table still asserts. Each change lands as its own `refactor:` commit.
- **2026-09-26 — Config: defaults < `.env` < env vars, and the file is optional.** 12-factor factor III: the environment is authoritative, and a missing `.env` is a warning, not a fatal, because the container never has one (verified by the acceptance suite). The defaults setup function was made non-injectable after an injected test double hid a bug where production defaults went to viper's global instance, so every env var was silently ignored. Tests now exercise the real defaults. Lesson worth keeping: an injectable seam can let a test replace the code it was meant to check.
- **2026-09-26 — Timestamp `0` is invalid (`863e66b`).** `Validate` now rejects `Timestamp <= 0`. Go decodes a missing JSON `timestamp` to `0`, so accepting `0` meant "no timestamp" was silently accepted as the 1970 epoch. The cost: a device that really does report epoch 0 (for example after a clock reset) is rejected, which is the right failure direction. The alternative, a `*int64` field that can tell "absent" from "zero", was not taken. It would have changed the exported `models` type. Note that this shipped in a `refactor:` commit together with new test cases, although it is a behavior change.
- **2026-09-26 — Mechanical change (Claude): `TestAppConfigValidate` port cases rebased on `validCfg`.** Angie set the shape (the `func(cfg) … (validCfg)` override on "system port"); Claude applied it to the other three port cases so they stop failing on an empty `Env`. Same 6 subtests before and after, no cases added or dropped. The sabotage check (setting the max-boundary case to `+1`) failed only that subtest. Done mid-cycle at her request, before the `ENV` work was committed.
- **2026-09-26 — Mechanical change (Claude): `TestTelemetry_Validate` overrides moved to `telemetry_helpers_test.go`.** This follows the shape Angie set in `config_helpers_test.go`: unexported `telemetryWith<Field>` copy helpers kept in a `_test.go` file. `shared.ValidTelemetry` stays exported because the Spec, which is not a test file, calls it. Same 5 subtests before and after. The sabotage check (timestamp `-1` → `1`) failed only that subtest. Lint shows 0 issues. Gaps named and left for her: a zero timestamp, an empty non-nil `Metrics` slice, a bad metric that isn't first in the list, and `ErrorIs` on the sentinels instead of matching the error text.
- **2026-09-26 — `APP_ENV` (`dev`/`prod`) drives gin mode; the image sets `APP_ENV=prod` (`fac44eb`; `ValidTelemetryN` fix in `e113d69`).** This resolves the `GIN_MODE` debt. `main` maps the environment to a gin mode and passes that to `api.NewServer`, so `api` doesn't import `config`. The variable was first named `ENV` and was renamed because POSIX `sh` uses `ENV`. Baking `prod` into the image makes the shipped artifact safe by default, and the acceptance suite now tests the configuration that ships. The cost is that the image claims to be prod wherever it runs unless `ENV` is overridden at deploy time (12-factor V: build, release, run). Also fixed the same day: `ValidTelemetryN` no longer returns 2n items. `gin.SetMode` is still process-global, which is accepted for now. A `TestConfigInit` case now sets `APP_ENV=prod`, which isn't the default. That closes the gap where nothing proved the variable is read, and it would catch the struct tag and `EnvName` falling out of sync.
- **2026-09-28 — `api.ServerMode` + functional options on `NewServer`; `gin.SetMode` stays inside `api` (`c5a2139`).** Gin has no per-engine mode: v1.12.0 keeps it in a package-level atomic, and it only gates gin's own debug output. Options weighed: move `SetMode` to `main` (the global is visible where the process is assembled, but `main` has to import gin), or wrap it in `api` (gin stays in one package, and `NewServer` quietly changes global state). She took the wrapper, with the same validated functional-options shape as `config.InitEnv`. The cost is that `NewServer` has a process-wide side effect, which should be written down where callers will see it. Same day, the `StorerContract` leak case (`ee9bc4c`) now changes a `Metrics` element as well as `DeviceID`, so going back to a shallow `slices.Clone` would fail it. The skipped empty-ID handler case stays on purpose, as a placeholder for format validation.
- **2026-10-09 — `202` means fire-and-forget: accepted into the in-memory buffer, not yet written.** Three options were weighed: fire-and-forget, group commit (the handler waits on a reply channel until its batch is flushed), and a synchronous write. Angie chose fire-and-forget. Her reasons: relying on devices to retry correctly is brittle; without reply channels the handler stays simple; a worker queue with a dead-letter path is a well-known pattern; and the data is informational, so losing some is acceptable. Costs she accepted: at-most-once delivery, since a crash or an undrained deploy loses what's buffered and the device never retries it; reliability after the `202` is the server's job; and the Spec becomes eventually consistent and must poll. Conditions that would reverse this: telemetry that includes counters or alarm events, where a lost reading is gone for good rather than replaced by the next sample, or a measured loss rate that turns out not to be small. Also noted: graceful shutdown belongs in the ring's own commit, not "eventually", because a buffer with no drain loses data on every deploy. Retry, backoff, and the dead-letter path wait for M3, since the M2 destination (the map) can't fail.
- **2026-10-09 — Shed load when the queue is full.** Options weighed: block (no errors, but unbounded latency and handler goroutines piling up under sustained overload), shed (immediate `503`), or drop (`202`, then silently discard). Angie chose shed. This doesn't contradict the fire-and-forget reasoning: a shed happens *before* acceptance, so the device gets an honest failure it can act on or ignore, and the loss is visible at the boundary rather than hidden. Points raised for her to settle while building: `503` (server overloaded) vs. `429` (this client is over its limit, which becomes M6's Redis rate limit); a `Retry-After` value, picked for now and measured later; queue-full as the first retryable, callee-fault error class, which makes the error-classification debt concrete; and how to fill the queue on purpose in a test (it's an engine-level test, not a Spec case; options include a small capacity and a destination that blocks, which would be the repo's first test double).
