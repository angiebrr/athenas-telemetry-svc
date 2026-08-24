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
