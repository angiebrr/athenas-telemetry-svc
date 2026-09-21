# Progress Log

Living state for `athenas-telemetry-svc`. **Claude maintains this file** — see "Keep the Log Current" in SKILL.md. Everything here is mutable and dated; SKILL.md holds only the durable method, so the skill stays reusable.

Published tracker: `.claude/skills/self-study-mentor/tracker.html` (Artifact).

---

## Current Position

- **Milestone:** M2 — In-Memory Concurrency & Buffering (phase 1 of 3 complete)
- **Branch:** `abreuer-M2-add-query`
- **GOOS:** ch. 7 in progress (M2's range is 6–8)
- **Educative:** Course 1, "What distributed systems achieve for us" (early)
- **Last verified green:** 2026-09-21 — unit tests `--race`, acceptance container suite, lint 0 issues both modules

---

### M2 progress (as of 2026-09-21)

M2 was rescoped by Angie on 2026-08-26 into three ordered phases, because the charter's wording bolts a make-it-fast objective onto an early make-it-work milestone. Treat this order as settled; don't push the ring earlier.

| Phase | Scope | State |
|---|---|---|
| 1. Work | In-memory store + `GET /v1/telemetry/:device_id`, so ingested telemetry has an observable destination | **Done** — Spec green at all three levels, lint clean, `--race` clean |
| 2. Right | Dispatch ring / worker pool — the answer to making `202 Accepted` *true*, not a data-race fix (the store's mutex already handles that) | Not started |
| 3. Fast | Benchmark (`--bench --memprofile`), then the charter's 100k RPS / allocation work | Not started |

A seam refactor sits between phases 1 and 2: `HandleQueryTelemetry` currently holds `data.Storer` and threads it into `telemetry.Query`, so transport carries the storage port through itself. The domain should own its store and the handler should hold only the domain. Do it before the ring — the ring lives behind that same seam.

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
- **Config decision before M2 exits.** `main.go:18`'s TODO — CLAUDE.md has always scoped config as M2 work. It lands, or it moves out explicitly. (Angie plans to include a quick config in the query PR — 2026-09-21.)

---

## Open Debt

Mirrors CLAUDE.md's Known Debt; kept here with dates and triggers.

- **Spec matches on error prose, not error class.** `ErrorContains(err, "missing device ID")` couples the container acceptance suite to sentinel wording. Durable fix is error classification carried by `httpserver.Driver`. Open since 2026-08-26.
- **Device IDs only checked for emptiness.** No format/length/charset rules. TODOs in `api_handler.go` and `api_handler_test.go`.
- **Complexity linting off.** Whole codebase measures ≤6 on cyclomatic/cognitive/function-length, so no conventional threshold could fire. Revisit when the ring lands; set the limit from measurement.

---

## Interview Question Bank

**This is the project's actual exit criterion.** The repo is the evidence; answering these fluently is the qualification. Each milestone adds questions. A good answer names the trade, the failure mode it doesn't protect against, and where you'd measure — not a definition.

Answers aren't stored here on purpose. The point is rehearsal out loud, and the reasoning gets captured at decision time in the Decision Record above, where it's honest.

### Cross-cutting — TDD and test design

Not milestone-scoped; these recur at every milestone and the answers get richer as the suite grows. **Terminology worth having straight** (it's the thing most candidates fumble): the *level* axis is GOOS ch. 1's acceptance / integration / unit. The *reusability* axis is separate — a test body written once and run against many implementations is a **contract test**; one written for a single implementation is just a unit test. Fowler's **solitary vs. sociable** splits those further: solitary isolates with doubles, sociable uses real collaborators. And a working stand-in implementation is a **fake**, not a mock (Meszaros: dummy / fake / stub / spy / mock).

- Walk me through the kinds of test in this repo. How many distinct kinds are there, and what does each buy that the others can't?
- `TelemetrySpec` runs at three levels. What does each level catch that the one beneath it can't? If the suite got too slow, which would you delete first, and what risk would you be accepting?
- `specifications.TelemetrySpec` and `data.DataStoreSpec` are both called "spec" and are not the same pattern. What's the difference?
- Why does `specifications` live in an exported, non-test package instead of a `_test.go` file? What does that cost you, and who else does this?
- Your validation table has five cases; the Spec has one invalid-input case. Why don't all five run at the acceptance level?
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
- Why a mutex here rather than a single owning goroutine? When would you switch?

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

---

## Decision Record

- **2026-08-26 — M2 rescoped to destination-first.** Angie pushed back on building the ring first: a dispatch ring with no destination has no externally observable behavior and can't be driven by an acceptance test. Three phases: work (store + GET), right (ring), fast (benchmark + allocation). Don't re-litigate.
- **2026-09-03 — Empty device ID moved out of the Spec.** `GET /v1/telemetry/` 404s; gin won't bind `:device_id` to an empty segment. The HTTP driver structurally cannot express that request, so it's a `telemetry.Query` unit test, not a system promise. (Option B of three; option C — query params instead of path params — stays open for when filtering is needed.)
- **2026-09-21 — Transport dropped out of Spec coverage.** `api_handler_test.go` no longer runs `TelemetrySpec`; it has hand-written status-code cases instead. Defensible (spec = behavior, handler tests = translation), but CLAUDE.md's "three levels" claim was stale and is corrected to two. Revisit if transport behavior starts drifting from the spec.
- **2026-09-21 — Doc maintenance delegated to Claude.** Prose upkeep isn't the learning target; stale docs are worse than none.
