# Progress Log

Living state for `athenas-telemetry-svc`. **Claude maintains this file** — see "Keep the Log Current" in SKILL.md. Everything here is mutable and dated; SKILL.md holds only the durable method, so the skill stays reusable.

Published tracker: `.claude/skills/self-study-mentor/tracker.html` → https://claude.ai/artifact/4DEea8wrLBwpdkRu9mV1Qk (republish to this URL; don't create a new one).

Study notes: `.claude/skills/self-study-mentor/notes/` (shell `index.html` + one fragment per chapter under `chapters/`, fetched by hash route) → https://claude.ai/artifact/TN7oyZ8z4RVFyy6YSRQake. Republish to this URL with `files` for each chapter. Kept out of `docs/` on purpose: the repo is recruiter-facing, and Claude-written chapters there would blur what she built vs. what was explained. Chapters explain behaviour, diagrams, and trade-offs, and end in a behaviour-list assignment; no implementation code, pseudocode only at a high level (her ask, 2026-10-21). Homework sets live in `notes/homework.js` (routes `#hw01`–`#hw07` and stretch `#hws1`–`#hws3`; chapters `#ch01`–`#ch07`, `#chs1`–`#chs3`; section deep links are `#chNN-sN`). Chapters are numbered in **reading order** (Part I built: 1 testing, 2 config; Part II the ring: 3 queues, 4 races, 5 overload; Part III: 6 observability; Part IV: 7 road ahead, planned 8 M3 and 9 M4; Stretch S1–S3 video). Homework N goes with chapter N.. **Rule: every question in the Interview Question Bank below appears verbatim in exactly one homework set, with refs to the chapter sections that answer it. When a question is added to the bank, add it to a set and make sure a chapter answers it, writing a new section or chapter if none does.**

---

## Current Position

- **Active (2026-11-06):** **M2 phase 2**, resumed. The Video track closed early: the SimpliSafe hiring manager is likely to pass, so the video material is now a **stretch goal** (see "Stretch: Video Streaming" below), not a deadline. Re-entry checklist under that section still applies.
- **Milestone:** M2 — In-Memory Concurrency & Buffering (phase 1 of 3 complete; phase 2 design half done, no ring code yet)
- **Branch:** `abreuer-M2-dispatch-ring` (off `main` after PR #5 merged as `5b40d90`). On 2026-11-06: sort-on-read landed (`c64930c` contract case, `8be6d78` store), engine scaffold `4ff937e`, naive goroutine-per-push `PushTelemetry` `34079f2`, FIXMEs `edb7df0`; notes and homework `eaa9422`, `835219d`. Next engine test: full queue rejects. `abreuer-M2-add-buffer` exists but has no commits past `main`.
- **Study plan:** which chapters and homework go with each phase is the "Study plan" table on the tracker; keep it and Reading Assignments below in sync.
- **2026-11-06, later (uncommitted):** gated `BlockingDataStore` double (close-once gate) plus the "Push doesn't wait for the store" test with a guard assertion that the store is empty. `PushTelemetry` is `go store.Insert(...)`: green with `--race`, but it's the predicted goroutine-per-push fake. The full-queue test is next and will rule it out. Lint has 7 issues (errcheck on the dropped `Insert` error, unkeyed `NotFoundError`, missing doc comments).
- **Phase 2 state (2026-11-06):** `internal/engine` scaffolded (`4ff937e`): `Queue`, `NewQueue(store)`, stub `PushTelemetry`, and a first red test "queued readings eventually reach the data store". That test is satisfied by a synchronous `store.Insert`, so on its own it doesn't force the ring. Advised next tests, in order: (1) Push returns while the store is blocked (kills synchronous insert), (2) full queue rejects with a sentinel (kills goroutine-per-push). Both need the repo's first test double, a gated `Storer`, and test 2 should synchronize on "worker entered Insert" signals, not sleeps (GOOS ch. 27 pattern).
- **Phase 2 state (audited 2026-11-05):** Decisions 1 and 2 are settled (fire-and-forget; shed with `503`). The Spec polls with `EventuallyWithT` (`4a00682`, `13f970a`), but it passes against the synchronous store, so there is still **no failing test driving the ring**. `internal/engine` does not exist. Next decision when M2 resumes: worker count vs. per-device ordering.
- **GOOS:** ch. 7 in progress (M2's range is 6–8). Resumed 2026-11-06.
- **Educative:** Course 1, "What distributed systems achieve for us" (early). Resumed 2026-11-06.
- **Last verified:** 2026-11-05 — `mise run test telemetry-svc --race` green on `3f0f3b4` (api, config, data, shared, telemetry). Container acceptance suite **not** run since 2026-10-09; run it first on M2 re-entry.

---

## Stretch: Video Streaming (was the Video track, 2026-11-05 → closed 2026-11-06)

**Status: stretch goal, no deadline.** Kept because it's interesting and reuses the project's ideas. Study notes ch. S1–S3 and HW S1–S3 cover it. Nothing here blocks or schedules M2. The original plan is kept below for reference.

**Why:** a referral for a Staff Software Engineer role on SimpliSafe's Video team, with interviews likely in 2–3 weeks. Their stack is Go, C#, TypeScript, AWS, Kubernetes, SQS, SNS, ffmpeg, and WebRTC. Per a friend on the team: Mist (MistServer) for livestream, Kinesis Video Streams for recording, WebRTC to viewers. The posting asks for real-time streaming, WebRTC/GStreamer/ffmpeg/libav, and streaming standards, codecs, and transport protocols, none of which she has worked with yet. **Budget: ~10 hrs/week × 3 weeks.** Sources are in `references.md` § Video Track.

**Not a milestone.** It doesn't change the M1–M7 roadmap or any gate, and it adds nothing to `internal/`. One point of contact with Athenas: the getStats producer *design* below. Its implementation waits for M2.

### Schedule

| Week | Theme | Work (≈ hours) | Checkpoint |
|---|---|---|---|
| **1** · Nov 5–11 | Pixels to packets | HowVideo.works (2) · Digital Video Introduction through H.264 + "online streaming", with the ffmpeg exercises (3) · WebRTC for the Curious ch. 1–3 (3) · Run Pion `reflect` and `save-to-disk` *unchanged* with `chrome://webrtc-internals` open (2) | Explain, out loud, how one webcam frame reaches a viewer via WebRTC and via HLS, naming every protocol on the way |
| **2** · Nov 12–18 | Real-time transport + the spike | WfC ch. 5, 6, 9; ch. 4 skimmed (4) · **Spike** (5): her Go/Pion recorder → ffmpeg → HLS · RFC 8216 §3–4 skim (1) | 30 s of webcam recorded by her code plays back as HLS in a browser |
| **3** · Nov 19–25 | Managed services + design | KVS dev guide: how it works, data model, retention, HLS playback, kvssink, KVS WebRTC + ingestion/storage (3) · MistServer intro, concepts, protocols (1.5) · getStats → Athenas design on paper (2.5) · Question bank rehearsal + one system-design mock (3) | Answers the Video question bank to the standard (trade-offs both ways, failure mode, where to measure) |

**If the interview lands in week 2:** move week 3's KVS + Mist reading and the rehearsal ahead of the spike. For an interview, reading pays off more than finishing the spike, and a half-finished spike is still worth talking about.

### The spike (hers to write)

A **spike** in XP's sense: throwaway code that answers a question, exempt from TDD *because* it gets thrown away. Keep it **out of this repo**, in a separate scratch repo, so the recruiter-facing Athenas history stays Spec-driven. Steps, each a checkpoint:

1. Browser page: `getUserMedia` → `RTCPeerConnection`, offer posted to her Go server. Signaling choice is hers: a hand-rolled JSON POST, or WHIP (RFC 9725, Pion's `whip-whep` example).
2. Go server (Pion): answer, receive the video track, write the RTP payload to a file with a Pion media writer.
3. ffmpeg: package the recording as HLS (`.m3u8` + segments) and play it back (Safari natively, or hls.js elsewhere).
4. Stretch: log `GetStats()` from the Pion side once a second while recording. This feeds the design below.

**The decision the spike forces (worth more than the code):** the browser will likely negotiate VP8. HLS players expect H.264/HEVC, AV1 at best. So either **transcode** (`-c:v libx264`: CPU per stream, added latency, generation loss) or **negotiate H.264 in SDP** and **transmux** (`-c copy`: nearly free, but the codec choice is pushed to the edge, and every device has to support it). That transcode-vs-transmux trade is exactly what a fleet of cameras plus a media server is built around. Record which one she chose, and why, in the Decision Record.

### Designing a getStats producer for Athenas, Spec-first (design only; build after the ring)

Goal: WebRTC session quality (bitrate, packet loss, jitter, frames dropped, RTT) as a **real producer** for Athenas, instead of synthetic `ValidTelemetry`. These are her decisions. The options are laid out here so they don't get made implicitly:

1. **What is the `device_id`?** Choices: the camera, the viewer session, or the peer connection. Stats are per connection *and* per direction: a sender sees `outbound-rtp` plus `remote-inbound-rtp` (loss and RTT reported back over RTCP), a receiver sees `inbound-rtp`. The identity choice decides what `GET /v1/telemetry/:device_id` can answer.
2. **Counters or rates?** `getStats()` returns mostly **cumulative counters** (`bytesReceived`, `packetsLost`, `framesDropped`) plus a few gauges (`jitter`, `currentRoundTripTime`). Shipping raw counters means a lost sample costs only resolution, and the rate gets computed at query time (Prometheus's `rate()` model). Shipping producer-computed deltas makes each sample carry data that's gone if the sample is lost. **This collides directly with the 2026-10-09 fire-and-forget decision**, whose stated reversal condition was "telemetry that includes counters." Counters actually make at-most-once *safer* than deltas do. She should think that through before choosing.
3. **Dimensions: flatten or extend the model?** `models.MetricReading` is `{Name, Value}` with no labels. Either flatten into names (`video.inbound.packets_lost`), which changes nothing, or add labels (kind, direction, SSRC), which changes the exported model and therefore starts with the Spec. Tie-in to the Observability questions: SSRC as a label is a cardinality bomb.
4. **Edge cases the current validation doesn't know about.** `packetsLost` is **signed** and can go negative (RTP duplicates, RFC 3550). `getStats` timestamps are float **milliseconds**, but `Timestamp` is an `int64` with no stated unit. Counters reset when a connection restarts. Each one is a candidate Spec or unit case, and deciding which level each belongs at is the GOOS exercise.
5. **Where does the producer run?** In the browser (sender-side view, needs CORS, untrusted client) or in the Pion server (receiver-side view, Go, trusted)? It could be both, with different `device_id`s.
6. **Rate, batching, and the M2 tie-in.** At 1 Hz per stream, 100k RPS is 100k concurrent streams. That makes a realistic load generator for **phase 3's benchmark**, rather than a synthetic one. Batch several samples per POST, or one per request?

**Spec-first order when it's built:** a new *producer* is not new *behavior*. If flattened names suffice, no Spec case changes. The producer gets its own tests in its own repo, and Athenas needs nothing. Only behavior Athenas must newly promise (labels round-tripping through query, a unit-bearing timestamp, accepting negative values on purpose) goes into `TelemetrySpec` red, before `internal/`. Saying where that line falls is itself a good interview answer.

### Return to M2

**Exit:** 2026-11-25, or the day after the final interview round if that's later. Any extension beyond one more week is a decision recorded here, not drift. Re-entry, in order:

1. `mise run test` (the full run, including the container suite; it hasn't run since 2026-10-09). Then `mise run lint`.
2. Re-read the two 2026-10-09 Decision Record entries and Study Notes ch. 3 (Bounded Queues). ~~Then decide worker count vs. per-device ordering.~~ Decided 2026-11-05: shared queue, sort on read.
3. Resume GOOS ch. 7 and Educative course 1 where they stopped. The "course 1 before M3" gate is unchanged.
4. The getStats producer is the planned workload for phase 3's benchmark. The observability stack stays where it was planned (phase 3).

**Paused for the track:** GOOS ch. 7–8, Educative course 1, the LMAX / Learn Go with Tests / SRE ch. 21 reading assignments, and the observability TODOs. None of them is cancelled.

---

### M2 progress (resumed 2026-11-06)

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
- [ ] **Queue metrics from the engine.** Depth (gauge), accepted / shed counts (counters), drain duration at shutdown. These turn the at-most-once trade into a measured loss rate (Study Notes ch. 3 §7–§8).
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
- `specifications.TelemetrySpec` and `data.DataStoreSpec` are both called "spec" and are not the same pattern. What's the difference? *(The storage body is now `StorerContract`, in `data_contract_test.go`.)*
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
- How many workers would you need once the store is SQS? Why is that a question about your target throughput and per-call latency rather than about SQS's capacity? What does batching do to the number?
- SRE ch. 22 says keep queues ≤ 50% of the thread pool. Why doesn't that rule transfer to a fire-and-forget queue?
- Shared queue vs. a queue per worker: which one balances load better, and why does routing to the least-loaded worker give up the only thing per-worker queues were for?
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
- Is a `503` with `Retry-After` backpressure or load shedding? Why does it matter that your devices are an *open* system (they keep sending whether or not you answer)? Where in this project do you get real backpressure, and why there?
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

### Stretch: video streaming — seeded 2026-11-05 (HW S1–S3)

Staff-level questions are rarely definitions. Expect "design it" and "what breaks." Each of these has a follow-up that asks what you'd measure.

*Delivery and latency*
- WebRTC vs. HLS for live view: what latency does each give, why, and what does each cost to serve to 10,000 viewers? Why can a CDN cache one and not the other?
- Live vs. recorded: why might the same camera use WebRTC to viewers and HLS (from KVS) for playback? Where does the switch happen in the UI, and what does the user see at the boundary?
- Where does glass-to-glass latency come from? Name each stage, from capture and encode through the network and jitter buffer to decode and render, and which one you'd attack first.
- What is LL-HLS, and what does it trade to get from ~10 s to ~2–3 s? Why isn't it WebRTC?
- UDP vs. TCP for real-time media: why does WebRTC choose loss over delay? When does it fall back to TCP or TURN-over-TLS, and what happens to quality?

*Codecs and packaging*
- What is a GOP, and how does keyframe interval affect join time, seek, bandwidth, and recovery after loss?
- Codec vs. container: H.264 in MPEG-TS, fMP4, MKV, IVF. Which container does each system in the stack use, and why does that matter for recording?
- Transcode vs. transmux: what does each cost per stream, and how would you design a fleet so the server rarely transcodes?
- A viewer's video freezes but audio keeps playing. Walk through the likely causes (lost keyframe, PLI/FIR, decoder stall) and what stats would tell them apart.
- H.264 vs. VP8/VP9 vs. AV1 vs. HEVC for a battery-powered camera: hardware encode, licensing, browser support.

*WebRTC mechanics*
- Walk through a WebRTC session setup: signaling, SDP offer/answer, ICE candidates, STUN, TURN, DTLS, SRTP. Which parts does the standard leave to you?
- What fraction of sessions need TURN, and what does TURN cost you (bandwidth, latency, servers)? How would you find out?
- What do NACK, FEC, and PLI each do, and when is each the wrong tool?
- How does WebRTC congestion control react to a bad uplink from a home camera? What would you see in `getStats()`?
- P2P vs. SFU vs. MCU: when does a home-camera product want each? What changes when two family members watch at once?

*Media servers and managed services*
- What does a media server like MistServer actually do? Why put one between a camera and viewers instead of going P2P?
- What does Kinesis Video Streams handle for you (ingest, durable time-indexed storage, retention, HLS/DASH/clip playback, WebRTC signaling + TURN), and what's left for you to build?
- How would you pick a KVS retention period? What does retention cost, and how does event-triggered recording change the bill?
- Build vs. buy: Mist (self-hosted) for live and KVS (managed) for recording. What does each choice trade, and what would make you consolidate?

*Reliability and operations (where Athenas transfers)*
- Which stream-quality metrics would you put on a dashboard, and which would wake someone up? Why does `device_id` (or SSRC) as a label break Prometheus?
- `getStats()` counters are cumulative. Why ship counters rather than rates, and what does that do to your tolerance for losing telemetry?
- A camera's upload is slower than its encoder. Where's the queue, what's the backpressure signal, and what should drop first? *(the M2 question, wearing a video costume)*
- Design live view and event recording for 1M home cameras: capacity math for ingest bandwidth, storage per day, and concurrent viewers.

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
| 2026-09-21 | *Learn Go with Tests* — Concurrency, Select, Sync, Context | M2 phase 2, Go idiom | **Done** (whole book, reported 2026-11-05) |
| 2026-09-21 | ~~Google SRE ch. 21, "Handling Overload"~~ → swapped 2026-11-05 for [SRE ch. 22, "Addressing Cascading Failures"](https://sre.google/sre-book/addressing-cascading-failures/), focus on "Queue Management" + "Load Shedding and Graceful Degradation" | M2 phase 2. Ch. 21 stalled her on Google-specific client-throttling math; ch. 22 is plainer and maps directly to the build | **Current** |
| 2026-10-21 | Study Notes ch. 3 (was ch. 1), "Bounded Queues and Worker Pools" (+ its assignment) | M2 phase 2. Synthesis of the three readings above with this repo's decisions; the assignment is the phase 2 build as a behaviour checklist. | **Current** (reread on re-entry) |
| 2026-11-06 | Study Notes ch. 1–7, S1–S3 + Homework | Rehearsal for the whole question bank, sequenced by the tracker's Study plan. Now: ch. 1 §3 and ch. 4, HW 4. Then ch. 5 and HW 1 during phase 2; HW 3 and HW 5 after the ring lands; ch. 6 and HW 6 in phase 3; HW 7 as each milestone opens; S1–S3 whenever. | **Current** (ch. 4, HW 4) |
| 2026-11-05 | [HowVideo.works](https://howvideo.works/) | Stretch (video). Whole-pipeline overview | Stretch |
| 2026-11-05 | [Digital Video Introduction](https://github.com/leandromoreira/digital_video_introduction), through H.264 + online streaming | Stretch (video). Codec vocabulary; do the ffmpeg exercises | Stretch |
| 2026-11-05 | [WebRTC for the Curious](https://webrtcforthecurious.com/) ch. 1–3, then 5, 6, 9 | Stretch (video) | Stretch |
| 2026-11-05 | KVS dev guide + MistServer docs (see `references.md` § Video Track) | Stretch (video) | Stretch |

The LMAX and SRE ch. 22 rows above resumed 2026-11-06 with M2.

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
- **2026-11-05 — Video track opened as a side track; M2 paused, not replaced.** Trigger: a referral for a Staff role on SimpliSafe's Video team, with interviews expected in 2–3 weeks, and no streaming background yet. Budget ~10 hrs/week for 3 weeks (2026-11-05 → 2026-11-25). Options weighed: a new milestone (rejected: it would bend the M1–M7 roadmap around one application), interleaving video with M2 (rejected: 10 hrs/week split two ways finishes neither before the interview), or a time-boxed side track with a fixed return date (taken). Why it isn't a detour from the project's goal: the project's real exit criterion is the question bank, and the video questions are distributed-systems questions in another domain (queues and backpressure in the uplink, at-most-once telemetry, cardinality, build vs. buy). Rules: the spike lives in a **separate repo**, because it is throwaway and exempt from TDD, and this repo's history stays Spec-driven. The getStats → Athenas producer is **designed now and built after the ring**, where it becomes phase 3's realistic workload. Nothing goes into `internal/` during the track. Return: re-entry checklist in § Video Track. An extension beyond one week past the last interview round gets recorded here.
- **2026-11-05 — N workers on one shared queue; no ordering through the pipeline; queries sort by timestamp on read.** The §5 decision. Options weighed: one worker, a shared queue, one queue per worker routed by device, or no order kept with sort on read. Angie chose the shared queue with sort on read. Her reasons: telemetry carries its own timestamp, so order is recoverable at the destination, and per-worker queues shed as soon as one queue is full while others sit empty (a hot partition). Two misconceptions were cleared up on the way. A shared Go channel needs no mutex, because the channel is the synchronization; and per-worker queues don't remove the store's mutex unless each worker also owns a shard of the store (the single-writer principle). Routing to the least-loaded worker would give up the ordering that per-worker queues exist for, so a shared queue is the better load balancer. Video doesn't change this: media frames need strict order, but RTP sequence numbers and the receiver's jitter buffer restore it on the media path, which never touches Athenas. `getStats()` samples are counters and gauges with timestamps. The general rule is to stamp order at the source and restore it where it's needed (end-to-end). Costs accepted: every read pays an O(n log n) sort; the order trusts device clocks, which already failed once (the `Timestamp` `0` decision); and ties between equal timestamps need a rule. Reversal condition: telemetry that carries events or state changes, where applying them in the wrong order gives the wrong final state. Carries over unchanged to M3, since SQS standard queues don't guarantee order either. Still open: where the sort lives, in `telemetry.Query` or in the `Storer` contract.
- **2026-11-05 — The engine is an internal of `telemetry.Service`; writes go through it, reads don't.** The queue is separate from `data.Storer`: processing vs. storage, so `StorerContract` stays valid for every adapter. `Service` builds the engine itself (option A) rather than `main` injecting it (option B). Angie's reason: the engine is an implementation detail of `Service`, and tracking it in `main` adds complexity. GOOS ch. 7 calls this internals vs. peers. Costs accepted: `Service` gains lifecycle methods. Mirror `http.Server.Shutdown(ctx)`, so `main` sets the drain deadline. Every `NewService`, including the domain-level Spec run, now starts goroutines and must be shut down in tests. Capacity and N reach `Service`'s constructor, so functional options as on `NewServer`. Shape worth naming: `Ingest` goes through the queue while `Query` reads the `Storer` directly, a command/query split. The gap between them is the eventual consistency the Spec already polls for.
- **2026-11-05 — The full-queue test uses a blocking `Storer` stub (the repo's first test double).** Options: a hook in the engine that fakes "full" (rejected: it's the same trap as the 2026-09-26 config seam, a test replacing the code it means to check), capacity alone (rejected: a ~100 ns map write drains a small queue faster than a test can fill it, so the test is flaky), not starting the workers (simpler, but tests a state production never reaches), or a store whose `Insert` blocks until the test releases it (taken; it exercises the real running engine). With N workers and capacity C, offers 1..N park in `Insert`, the next C fill the queue, and offer N+C+1 must fail as full. Taxonomy correction made at the time: Angie called it part fake, part stub. It's a **stub**, because it controls an indirect input (when `Insert` returns); a fake is a *working* shortcut like `InMemoryStore`. The drain test ("offer K, shut down, count K stored") does care what reached the store, so it uses the real `InMemoryStore` (state verification), not a spy. Update CLAUDE.md's "no test doubles" line when the stub lands.
- **2026-11-05 — Queue capacity and worker count are constructor parameters; sizing is deferred to phase 3.** Angie first proposed SRE ch. 22's "queue ≤ 50% of the thread pool" rule. It doesn't transfer, and she agreed to defer: ch. 22 describes thread-per-request servers where a queued request is a client waiting on the line, and where there's another replica to fail over to. Athenas returns `202` at the put, so a queued reading costs staleness, not client latency. The queue's job is absorbing bursts, device fleets are bursty, and until M6 there's no replica. Ch. 22's own exception for bursty load (size by burst size and frequency) points the same way. Sizing method for phase 3: Little's Law, workers ≈ throughput × per-write latency ÷ batch size, from measured numbers. Note that N > 1 buys almost nothing against `InMemoryStore`, whose writes serialize on one mutex. N earns its keep when the store is a network call (M3/M4).
- **2026-11-06 — Homework sets for the question bank, with chapters to back every answer.** Her ask: turn the bank into homework (matching, fill-in-the-blank, multiple choice, short answer) and make sure every homework and interview question is answerable from the study notes. Ch. 1 covered about a third of the bank, so nine chapters were added: 2 testing outside-in (cross-cutting + M1), 3 races and ownership, 4 overload/retries/the `202` promise, 5 configuration, 6 observability, 7–9 video (codecs, WebRTC, delivery at scale), 10 a preview of M3–M7. The planned M3/M4 deep-dive chapters moved to 11–12. All 86 bank questions are placed (the M3–M7 one-line seeds were expanded into 13 questions). The open parts deliberately have **no answer key**, consistent with "the bank holds prompts, never answers": each points to the section that answers it, and interview items carry a self-grading rubric (trade-off both ways, failure mode not covered, where to measure, production system). Objective parts self-check. Answers first saved in the viewer's browser only; superseded the same day by database sync (next entry).
- **2026-11-06 — Homework answers sync, and Claude grades open answers (critique first).** The notes page declares `db` and `sample`. Answers live in the artifact's database at `answers/<hwNN>` (`fields`, `critiques`, `modelViewed`), readable and writable only by the owner. She chose this over the private `data/users/` area so that a Claude session can read her answers with `ArtifactData` (`list` on collection `answers`) to quiz her or find weak spots. localStorage stays as an offline mirror. Grading: "Grade my answer" sends the question, her answer, and the linked chapter sections to Claude, which grades against the interview standard and returns a verdict, the four-part standard, strengths, gaps (described as concepts plus the section to read, never as copyable prose), and the interviewer's likely follow-up. The critique is stored and marked stale if she edits the answer afterward. Her choice: **critique only, with a model answer behind an "I'm stuck" button** and a confirmation step. A model answer is generated on demand, never stored, and `modelViewed` records that she looked, so the record stays honest about which answers were her own. Each grade costs her own Claude usage.
- **2026-11-06 — Video track closed early; study notes renumbered into reading order.** The SimpliSafe hiring manager told Jay he'll probably pass, so the interview that justified the time-boxed track is unlikely. The video material stays as a **stretch goal** with no deadline (it interests her), and M2 phase 2 is active again. At the same time the study notes were renumbered: they had been numbered in the order they were written (ch. 1 for the build, then whatever the video track needed), which read as random. New order follows the project: Part I what's built (1 testing ← old 2, 2 config ← old 5), Part II the ring (3 queues ← old 1, 4 races ← old 3, 5 overload ← old 4), Part III (6 observability, unchanged), Part IV (7 road ahead ← old 10; planned 8 M3, 9 M4), Stretch (S1 codecs ← old 7, S2 WebRTC ← old 8, S3 video at scale ← old 9). Homework numbers match their chapters. Saved answers moved with them: the one database document (old `hw01` → `hw03`) via `ArtifactData`, and browser-saved answers and checklist ticks by a one-time migration in the page. Decision Record entries written before today use the old numbers.
- **2026-11-05 — Log audit.** The log was current through 2026-10-21, but the Current Position header still read 2026-10-09 / 2026-09-21. Corrected it after checking the branches and running the suite (see Current Position). No code changes since `3f0f3b4`.
