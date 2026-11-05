# Self-Study Reference Library

Annotated references for Angie's self-directed distributed systems curriculum, mapped to the milestone where each one earns its keep. Pull the relevant ones into a milestone conversation rather than dumping the list.

Source: `.claude/project-context.md` §"References & Links". URLs below are the **verified working** ones — see [Corrections](#corrections-to-the-source-list) at the bottom for the entries whose links in the source doc are wrong or dead.

---

## Primary Curriculum (self-selected)

**Educative — Become a Distributed Systems Professional** (54 hrs, 404 lessons, 3 courses)
<https://www.educative.io/path/become-a-distributed-systems-professional>

| Course | Covers | Milestones |
|---|---|---|
| 1. Introduction to Distributed Systems for Dummies | Basics; designing/developing distributed systems | 1–2 |
| 2. Distributed Systems for Practitioners | Core algorithms and protocols | 3–5 |
| 3. Distributed Systems: Building Software for the Real World | Architecting and building for production | 6–7 |

Theory course leads, project milestone follows. When she's mid-milestone, ask which lessons she's covered — the milestone is the lab for the course, so a milestone running ahead of its course is a signal to slow the build, not to skip the reading.

**GOOS — Growing Object-Oriented Software, Guided by Tests** (Freeman & Pryce)
Book site: <https://growing-object-oriented-software.com/> · PDF: [Books-3 mirror](https://github.com/GunterMueller/Books-3/blob/master/Growing%20Object%20Oriented%20Software%20Guided%20by%20Tests.pdf)
Chapter-to-milestone reading order lives in SKILL.md. Read cover to cover; the examples are Java/jMock/Hamcrest, so the translation exercise into Go is part of the learning, not an obstacle to route around.

---

## Method & Craft

**Make It Work, Make It Right, Make It Fast** — C2 Wiki
<https://wiki.c2.com/?MakeItWorkMakeItRightMakeItFast>
Kent Beck's sequencing rule (he cites it as "make it *run*, make it right, make it fast," learned from his grandfather). Three ordered phases: get correct behavior by any means, then make it sustainable/well-designed, then and only then optimize. Doing step three early has a name — premature optimization.

**This is the single most load-bearing reference for this project, because it cuts directly against Angie's strengths.** Her background is C++ cache locality, particle-system optimization, and bandwidth reduction at Sonos — the instinct to make it fast is well-trained and will fire early. Milestone 2 as written ("optimize memory allocation to handle 100,000 RPS... bypassing dynamic maps and heap boxing") is a *make it fast* objective sitting inside an early *make it work* milestone. That tension is deliberate and worth naming out loud when M2 starts:

- Get the dispatch ring correct and well-shaped first, with a spec that proves it.
- Then benchmark (`mise run test telemetry-svc --bench --memprofile`) to find where allocation actually hurts.
- Then optimize, with the benchmark as the evidence that it worked.

"I already know it'll allocate" is not the same as having measured it. Ask for the benchmark before the optimization, every time. This maps cleanly onto Red-Green-**Refactor**: the optimization is refactor work, and refactor comes third.

**Learn Go with Tests** (Quii) — <https://quii.gitbook.io/learn-go-with-tests>
The Go-idiomatic companion to GOOS; use it when the gap is "how is this expressed in Go" rather than "what's the design." Chapters worth naming by milestone:

| Chapters | Milestone |
|---|---|
| Dependency Injection, Mocking | 1, 3 |
| Concurrency, Select, Sync, Context | 2, 5 |
| Introduction to acceptance tests, Scaling acceptance tests | 1 (this is the spec/driver pattern already in the repo) |
| Working without mocks, Anti-patterns | 3–4 (read before reaching for GoMock everywhere) |
| HTTP server; JSON, routing and embedding | 1 |
| Revisiting time, with testing/synctest | 4–5 (deterministic tests for pollers and retries) |

"Scaling acceptance tests" is the direct source of the Spec + Driver structure in `specifications/` and `internal/drivers/` — worth a re-read at any milestone that adds a new driver.

**The Practical Test Pyramid** (Ham Vocke / martinfowler.com)
<https://martinfowler.com/articles/practical-test-pyramid.html>
Lots of fast unit tests, some integration tests, very few end-to-end. Two warnings that apply directly here:
- **Ice-cream cone risk is real for this repo.** The acceptance suite boots a Docker image; it's already `testing.Short()`-gated for that reason. As LocalStack, Postgres, and Kafka containers land in M3–M5, the E2E tier gets slow enough to stop being run. Each new milestone should add *more* unit-level coverage than E2E coverage.
- **Don't test the same condition at two levels.** If a spec case is fully covered by a fast unit-level adapter, it doesn't also need a container round-trip — keep it only if the higher-level test genuinely raises confidence.

**Understanding Screenplay** (Cucumber blog, 4-part series)
Part 1: <https://cucumber.io/blog/bdd/understanding-screenplay-part-1/>
Test automation as a stage-play: **Actors** with **Abilities** (the dependencies — an HTTP client, a DB handle) perform **Interactions** (small discrete units: `PostJson`, `InsertInto`), composed into **Tasks** (domain-meaningful chunks), and assert via **Questions** (queries against system state). Explicitly "the command pattern applied to organising test automation code."

Relevance: `httpserver.NewDriver(endpoint, &http.Client{...})` is already an Ability-carrying actor in miniature. When a second driver arrives (M3's SQS driver, M6's multi-pod driver) and the drivers start duplicating setup, Screenplay is the refactor target — composable interactions instead of a widening Driver interface. Don't adopt the vocabulary prematurely in M1; it's a solution to duplication that doesn't exist yet.

---

## Milestone 3–5: Queues, Outbox, Streaming

**Testcontainers for Go — Quickstart** — <https://golang.testcontainers.org/quickstart/>
Already in use via `shared.StartDockerServer`. Extend the same harness for LocalStack (M3), Postgres (M4), Kafka (M5) rather than standing up a parallel docker-compose path.

**GoMock (uber-go/mock)** — <https://github.com/uber-go/mock>
The maintained fork; `golang/mock` is archived, so `go.uber.org/mock` is the correct import path. Pair with GOOS Ch. 8 **"Only Mock Types That You Own"**: generate mocks for *her* `TelemetryPublisher`-style interfaces, not for the AWS SDK's types. The AWS SDK gets exercised for real against LocalStack in the acceptance tier. This is the cleanest illustration in the whole project of a GOOS rule and a tool choice agreeing — worth making explicit at M3.

**Transactional Outbox pattern** (M4)
The source doc's video link is a dead placeholder (see Corrections). Teach it from first principles instead — it's a short argument: you cannot atomically write to Postgres and publish to a broker without distributed transactions, so you write the event and an outbox row in one local transaction and let a poller publish afterward. That buys atomicity at the cost of at-least-once delivery and publish latency, which is why consumers must be idempotent. If a link is wanted, Fowler/microservices.io's outbox writeup is the standard reference.

**How to Handle Backpressure in Kafka Consumers** (OneUptime) — <https://oneuptime.com/blog/post/2026-01-24-handle-backpressure-kafka-consumers/view>
Concrete strategies for M5, in rough order of how much control they give: partition **pause/resume** on a backlog high/low watermark (the article uses pause at 1000, resume at 500 — hysteresis matters, a single threshold flaps), `max.poll.records` to cap batch size and memory, manual offset commits so slow processing can't silently lose messages, consumer-lag monitoring as the detection signal, horizontal scaling as the blunt instrument. The milestone specifically calls for manual pause/resume — that's the one that teaches the most, since it forces her to own the flow-control loop rather than delegate it to config.

**Writing and Testing an Event Sourcing Microservice with Kafka and Go** (Semaphore)
<https://semaphore.io/community/tutorials/writing-and-testing-an-event-sourcing-microservice-with-kafka-and-go>
End-to-end worked example: dockerized event-sourced Go service, sarama consumer groups, scaling. Useful as a shape reference for M5. It's a tutorial, so it will hand her code — treat it as reading, not as a template to copy.

---

## The Actor Model Cluster (conceptual lens, not a dependency)

Six of the source references are about the actor model. None of the seven milestones name actors, and that's correct — **do not let her adopt an actor framework for this project.** Ergo, proto-actor, and gosiris each supply exactly the supervision, mailbox, and scheduling machinery that Milestone 2 exists to make her build by hand out of goroutines and channels. Adopting one would abstract away the engine room she deliberately dropped into.

Use them as a *lens* instead — "your worker pool is a mailbox; who owns the state?" — and as informed context for the framework-vs-primitives trade-off, which is a good architectural-audit conversation to be able to have.

| Reference | Use |
|---|---|
| [Actor model — Wikipedia](https://en.wikipedia.org/wiki/Actor_model) | Theory: asynchronous message passing as a concurrency foundation. M2 background. |
| [The Actor Model in Go (DEV)](https://dev.to/forkbikash/the-actor-model-in-go-simplifying-concurrent-programming-1j9d) | Short intro to mailboxes via goroutines + channels. M2. |
| [Kafka, distributed coordination and the actor model — Daniel Lebrero](https://danlebrero.com/2018/04/09/kafka-distributed-coordination-actor-model/) | **The most useful of the cluster.** Frames a Kafka partition as a single-threaded, partition-isolated actor — which is exactly why partition assignment gives you ordering and why consumer-group rebalancing is the hard part. Read at M5. |
| [Ergo Framework](https://docs.ergo.services/) | Erlang-style supervision trees + mutual TLS in Go. Reference for what a framework provides; not a dependency. |
| [protoactor-go/actor](https://pkg.go.dev/github.com/asynkron/protoactor-go/actor) | API-shape reference for actor types. Not a dependency. |
| [teivah/gosiris](https://github.com/teivah/gosiris) | **Archived.** Historical interest only (etcd-based actor discovery). Don't build on it. |

---

## Milestone 6–7: Scale & DevOps

**Using the Four Keys to Measure Your DevOps Performance** (Google Cloud)
<https://cloud.google.com/blog/products/devops-sre/using-the-four-keys-to-measure-your-devops-performance>

| Metric | Definition | Measures |
|---|---|---|
| Deployment Frequency | How often the org successfully releases to production | Velocity |
| Lead Time for Changes | Time from commit to production | Velocity |
| Change Failure Rate | % of deployments causing a production failure | Stability |
| Time to Restore Service | Time from incident creation to resolution | Stability |

Implementation notes for M7: Deployment Frequency counts only *successful* production deploys, so "successful" needs a defined threshold. Lead Time requires a persisted commit-SHA→deployment mapping, which is a schema decision to make before wiring the pipeline, not after. Change Failure Rate and Time to Restore both need deployments linked to incidents — meaning M7 needs *some* incident record, even a manual one, or two of the four keys are unmeasurable. Note also that DORA's 2022 report collapsed the bands from Elite/High/Medium/Low to three (High/Medium/Low); cite it accurately if it goes in the README.

---

## Video Track (side track, opened 2026-11-05)

Not a milestone. A time-boxed detour for the SimpliSafe Video interview: real-time streaming, codecs, and transport. The schedule and exit criteria live in `progress-log.md`. **All URLs below returned 200 on 2026-11-05.** Read them in this order. Each later source assumes the earlier ones.

**1. HowVideo.works** (Mux): <https://howvideo.works/>
A one-sitting overview of the whole pipeline: capture, encoding, containers, delivery protocols, adaptive bitrate, playback. Read it first so the later sources have somewhere to hang. It is a vendor's site, so read anything comparing products skeptically. The explanations themselves are sound.

**2. Digital Video Introduction** (Leandro Moreira): <https://github.com/leandromoreira/digital_video_introduction>
How a codec actually works, with images and runnable ffmpeg experiments: color models, chroma subsampling, I/P/B frames, GOPs, intra and inter prediction, transforms, quantization, entropy coding, then H.264 bitstream structure and the online-streaming section. **This is the codec vocabulary interviews assume**: keyframe interval, GOP, profile/level, bitrate vs. quality. Do the ffmpeg exercises. They make the spike's ffmpeg step make sense.

**3. WebRTC for the Curious** (Sean DuBois et al., Pion's authors): <https://webrtcforthecurious.com/>
Free, protocol-first, and written by the people who wrote the library you'll use. Read in this order:

| Chapter | Why |
|---|---|
| [1 What, Why and How](https://webrtcforthecurious.com/docs/01-what-why-and-how/) | The protocol stack in one picture |
| [2 Signaling](https://webrtcforthecurious.com/docs/02-signaling/) | SDP offer/answer. WebRTC deliberately doesn't specify signaling, so the spike's HTTP endpoint *is* yours |
| [3 Connecting](https://webrtcforthecurious.com/docs/03-connecting/) | NAT, STUN, TURN, ICE. The "why does it work at home but not on cellular" chapter |
| [4 Securing](https://webrtcforthecurious.com/docs/04-securing/) | DTLS and SRTP. Skim it |
| [5 Real-time Networking](https://webrtcforthecurious.com/docs/05-real-time-networking/) | Congestion control, NACK, FEC, jitter buffers. **The latency trade-offs live here** |
| [6 Media Communication](https://webrtcforthecurious.com/docs/06-media-communication/) | RTP/RTCP, codecs in SDP, the RTCP reports `getStats()` numbers come from |
| [9 Debugging](https://webrtcforthecurious.com/docs/09-debugging/) | `chrome://webrtc-internals`, and reading stats |

**4. Pion examples:** <https://github.com/pion/webrtc/tree/master/examples>
Read these, don't copy them. The spike is yours.
- [`save-to-disk`](https://github.com/pion/webrtc/tree/master/examples/save-to-disk): browser webcam → Pion → files on disk. The spike's closest relative. Note which codec it writes and which container (IVF/Ogg), because that decides your ffmpeg step.
- [`reflect`](https://github.com/pion/webrtc/tree/master/examples/reflect): the smallest full round trip. Read it first to see the moving parts.
- [`whip-whep`](https://github.com/pion/webrtc/tree/master/examples/whip-whep): standard HTTP signaling (WHIP is [RFC 9725](https://datatracker.ietf.org/doc/html/rfc9725)). Worth considering for the spike's signaling instead of inventing a format.
- [`stats`](https://github.com/pion/webrtc/tree/master/examples/stats) and [`rtcp-processing`](https://github.com/pion/webrtc/tree/master/examples/rtcp-processing): server-side stats in Go. Directly relevant to the getStats producer design.
- [`broadcast`](https://github.com/pion/webrtc/tree/master/examples/broadcast): one uploader fanned out to many viewers. This is a tiny SFU (selective forwarding unit), the shape of a media server's WebRTC output.

**5. Amazon Kinesis Video Streams developer guide**
- [How it works](https://docs.aws.amazon.com/kinesisvideostreams/latest/dg/how-it-works.html) and [data model](https://docs.aws.amazon.com/kinesisvideostreams/latest/dg/how-data.html): producers, streams, consumers, and fragments (MKV, time-indexed).
- Retention: `DataRetentionInHours` on [CreateStream](https://docs.aws.amazon.com/kinesisvideostreams/latest/APIReference/API_CreateStream.html). `0` means no storage, pass-through only.
- [HLS playback](https://docs.aws.amazon.com/kinesisvideostreams/latest/dg/hls-playback.html): live and on-demand from stored fragments. This is the "recording" half.
- [Producer SDK and the GStreamer `kvssink` plugin](https://docs.aws.amazon.com/kinesisvideostreams/latest/dg/examples-gstreamer-plugin.html): how cameras actually push. This is the posting's GStreamer line.
- [KVS with WebRTC](https://docs.aws.amazon.com/kinesisvideostreams-webrtc-dg/latest/devguide/what-is-kvswebrtc.html), a *separate* product inside KVS: managed signaling channels plus STUN/TURN, master/viewer roles. Also [WebRTC ingestion and storage](https://docs.aws.amazon.com/kinesisvideostreams-webrtc-dg/latest/devguide/webrtc-ingestion.html) (`JoinStorageSession`), where one WebRTC session both serves live viewers and records to a stream.
- Demos: [aws-samples/amazon-kinesis-video-streams-demos](https://github.com/aws-samples/amazon-kinesis-video-streams-demos).

**6. MistServer docs:** <https://docs.mistserver.org/> (source: [DDVTECH/mistserver](https://github.com/DDVTECH/mistserver), Unlicense, actively maintained as of 2026-11)
Read [Introduction](https://docs.mistserver.org/mistserver/introduction/), [Core Concepts](https://docs.mistserver.org/mistserver/concepts/), Inputs/Outputs, and the [protocol section](https://docs.mistserver.org/protocol/), especially [WebRTC](https://docs.mistserver.org/protocol/realtimestreaming/webrtc). The idea to take away: a media server takes **any protocol in and any protocol out** (RTMP/SRT/WHIP/RTSP in; HLS/DASH/WebRTC out), and it mostly *transmuxes* rather than transcodes. History worth knowing in one line: Livepeer acquired DDVTech in 2021, and it was spun back out to its engineers in 2023 ([announcement](https://news.mistserver.org/news/113/MistServer+independent+once+more)).

**Supporting, for lookup rather than reading through:**
- [MDN `RTCPeerConnection.getStats()`](https://developer.mozilla.org/en-US/docs/Web/API/RTCPeerConnection/getStats) and the [W3C webrtc-stats spec](https://www.w3.org/TR/webrtc-stats/): which stat lives on which report type (`inbound-rtp`, `remote-inbound-rtp`, `candidate-pair`).
- [HLS, RFC 8216](https://datatracker.ietf.org/doc/html/rfc8216): skim §3–4 for segments and playlists. [Apple's HLS docs](https://developer.apple.com/documentation/http-live-streaming) cover Low-Latency HLS.
- [ffmpeg HLS muxer options](https://ffmpeg.org/ffmpeg-formats.html#hls-2): `hls_time`, `hls_playlist_type`, `hls_segment_type` (mpegts vs. fmp4).
- [WebRTC samples](https://webrtc.github.io/samples/): browser-side `getUserMedia` and stats demos.

---

## Corrections to the Source List

The reference list in `.claude/project-context.md` has errors. Use the URLs in this file; flag these if she's citing them anywhere public (README, portfolio writeup), since a broken citation in a portfolio piece is worse than no citation.

| # | Source-doc entry | Problem |
|---|---|---|
| 3 | "ActorDB: A Unified Database Model", arXiv 2510.01234 | **Wrong paper.** That ID resolves to *"LLMRank: Understanding LLM Strengths for Model Routing"* — unrelated to actors or databases. The described paper may not exist; verify before citing. |
| 9 | "Go Kafka: Transactional Outbox Pattern (Kcode Video)" | **Dead placeholder** — the URL literally contains `youtube-outbox-placeholder`. No such link. |
| 19 | "Understanding Screenplay Pattern Part 1" | 404. Correct path inserts `/bdd/`: `cucumber.io/blog/bdd/understanding-screenplay-part-1/` |
| 5 | Educative skill path | 404. Path is singular: `/path/`, not `/paths/`. |
| — | In-text citation numbers throughout the doc | **Systematically off.** Milestone 1 cites [8] for Testcontainers (#8 is the GitHub repo; Testcontainers is #16); [10, 17] for GoMock (#17 is "Actor Model in Go"); [15] for the Red-Green-Refactor cycle (#15 is a resume docx; Beck's rule is #14). Don't trust the bracketed numbers — resolve by title. |

All other listed URLs verified reachable.
