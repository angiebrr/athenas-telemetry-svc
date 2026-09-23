# Self-study

- [Self-study](#self-study)
  - [How this self-study works](#how-this-self-study-works)
  - [Self-study roadmap](#self-study-roadmap)
  - [Self-study curriculum](#self-study-curriculum)

## How this self-study works

This self-study is driven by a skill I created with Claude (Opus 5) and can be found in [.claude/skills/self-study-mentor/SKILL.md](../.claude/skills/self-study-mentor/SKILL.md). The idea was to work on this project while I worked on [a distributed systems course in Educative](https://www.educative.io/path/become-a-distributed-systems-professional) and read **[Growing Object-Oriented Software, Guided by Tests](https://growing-object-oriented-software.com/)** (i.e. **GOOS**).

Claude and I both came up with the curriculum, and Claude helps me learn new concepts as I work through it. The agent is instructed to not write the code for me to maximize knowledge gain, but can help me do menial tasks such as update progress logs, stale documentation, etc.

So, the repo is a bit of a mix of AI-generated code and my own hand-written code (i.e. Claude handled docs, tooling and CI upkeep, and I wrote the internal code, the spec and the tests).

It's worth noting that, in a production project, I would encourage more code generation (especially when using a good model like Opus). However, since I wanted to *truly understand* what I was doing so I could replicate it later manually or via prompts, I opted for Claude to not write it all for me.

## Self-study roadmap

The plan is 7 milestones. Each one pairs some theory from the [Educative distributed systems path](https://www.educative.io/path/become-a-distributed-systems-professional) with the next few chapters of **GOOS**, and adds one new distributed systems piece to the service.

| # | Milestone | Adds | Status |
|---|---|---|---|
| 1 | Walking Skeleton | HTTP ingest endpoint, validation, containerized E2E harness | ✅ Complete |
| 2 | In-Memory Concurrency | Thread-safe store and query endpoint, then a channel-based dispatch ring and worker pools; allocation profiling under load | 🔨 In progress |
| 3 | Queue Integration | Buffer flush to AWS SQS, exercised offline against LocalStack | ⬜ Planned |
| 4 | Transactional Outbox | Atomic Postgres write of event + outbox row, with poll-and-publish workers | ⬜ Planned |
| 5 | Event Streaming | Kafka consumer groups with explicit consumer-side backpressure | ⬜ Planned |
| 6 | Distributed Coordination | Multi-pod Kubernetes deployment; Redis-backed global rate limiting | ⬜ Planned |
| 7 | Delivery & Observability | GitHub Actions CI/CD and DORA four-keys instrumentation | ⬜ Planned |

The milestones go in order, and, within each one, we want to make it work first, then make it right, and only *then* make it fast ([Make It Work, Make It Right, Make It Fast](https://wiki.c2.com/?MakeItWorkMakeItRightMakeItFast)). In order to prevent premature optimization, we don't optimize until a benchmark says it's actually a problem.

## Self-study curriculum

The main things I'm working through or leaning on:

- **[Growing Object-Oriented Software, Guided by Tests](https://growing-object-oriented-software.com/)** by Freeman & Pryce - the backbone of the whole project. I'm reading it cover to cover alongside the build, and it's where the walking skeleton, the spec/driver split, and the outside-in cycle come from.
- **[Become a Distributed Systems Professional](https://www.educative.io/path/become-a-distributed-systems-professional)** on Educative - three courses covering distributed systems fundamentals, protocols, and production concerns. This repo is basically my lab for it.
- **[Learn Go with Tests](https://quii.gitbook.io/learn-go-with-tests)** by Chris James - my go-to for doing TDD the Go way, and where I got the structure for the acceptance tests from.
- **[The Practical Test Pyramid](https://martinfowler.com/articles/practical-test-pyramid.html)** by Ham Vocke - a good reminder to keep the slow, container-backed tests few and push most of the testing down to fast unit tests.
