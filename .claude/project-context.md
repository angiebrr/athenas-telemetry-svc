# Self-Directed Study & Portfolio Charter
## Project: `athenas-telemetry-svc`

This document defines the professional background, learning objectives, curriculum alignment, and milestone roadmap for a **self-designed independent study in advanced distributed backend architectures**. The curriculum, milestones, and constraints below were chosen and scoped by Angie — this is not coursework, and no part of it was assigned. It is designed to live in the repository (at `.claude/project-context.md`) to configure any instance of Claude with the relevant background, tools, and goals.

---

## 🎯 The Mission: An Elite Portfolio Showcase

This project is a **highly-visibility self-study portfolio piece** designed to showcase to recruiters and potential teammates. While the primary goal is learning, the secondary goal is to produce an **exceptionally polished, production-ready, and highly performant** distributed backend service that stands up to enterprise-level code reviews and architectural audits. Every design decision, test suite pattern, and line of Go code must reflect top-tier software craftsmanship.

---

## 🎓 The Engineer Profile: Angela \"Angie\" Breuer

*   **Academic Credentials:**
    *   **Master of Science in Computer Science (GPA 4.00):** University of Montana, Class of 2016 [4]. Thesis: *Optimization of a particle system in an open-source C++ game engine, ChilliSource* [4, 16].
    *   **Bachelor of Science in Computer Science / Mathematics Combined (GPA 3.95):** University of Montana, Class of 2014 [4].
*   **Technical Pivot:** Transitioning a deep low-level optimization mindset from C++, bare-metal Linux sysadmin, and serverless environments into **stateful, highly concurrent distributed Go systems** [16]. Having mastered CPU cache locality, dynamic memory limits, and automated infrastructure, Angie is focusing this project on building horizontally scalable, stateful backend architectures.

### Professional History (Resume Highlights) [4, 16]:

1.  **Epic — Software Engineer (DevOps & Backend) (Most Recent) [16]:**
    *   Developed core backend microservices utilizing **Golang** [16].
    *   Built serverless cloud workflows using **AWS Lambdas** and ran containerized workloads orchestrated on **Kubernetes (EKS)** [16].
    *   Managed cluster state, infrastructure automation, and deployments using **Terraform, Helm, and DevOps pipelines** [16].
    *   *The Pivot:* While AWS Lambdas provide effortless scaling, they abstract away stateful node concurrency. Angie is dropping down into the engine room—managing persistent, highly concurrent memory buffers, channel pools, and sequential coordination natively in Go.
2.  **Sonos — Embedded Server Software Engineer (C++ Developer) [4, 16]:**
    *   Wrote embedded C++ server applications running on millions of local network devices [4, 16].
    *   Developed a first-party asynchronous I/O library utilizing socket multiplexing to modernize LAN and cloud-native device coordination [4, 16].
    *   Optimized bandwidth and network overhead, saving hundreds of thousands of dollars in cloud egress costs [4, 16].
3.  **Fast Enterprises — Implementation Consultant (Full Stack .NET Developer) [4, 16]:**
    *   Overhauled massive web portal systems using the .NET framework and SQL Server for public infrastructure (Massachusetts RMV) [4, 16].
    *   Tested end-to-end user scenarios and managed automated deployments in high-visibility environments [16].
4.  **Numerical Terradynamic Simulation Group (NTSG) — CentOS System Administrator [4, 16]:**
    *   Performed bare-metal Linux system administration, configured active monitoring nodes (Nagios, Fail2Ban), and built automated cross-platform file backup scripts using Rsync/Robocopy [4, 16].

---

## 📚 Educative Curriculum Alignment: \"Become a Distributed Systems Professional\"

This project serves as the practical laboratory for Angie's parallel study of the Educative Skill Path: **\"Become a Distributed Systems Professional\"** [5]. Each project milestone is strictly aligned with the corresponding theoretical module of the course [5]:

```
                     EDUCATIVE COURSE PATH                   PROJECT MILESTONES
                     
               +------------------------------+       +------------------------------+
               |           MODULE 1           |       |         MILESTONES 1 & 2     |
               |  Distributed Systems Basics  +------>|  Walking Skeleton & Go       |
               |  (Correctness & Concurrency) |       |  Channels In-Memory Buffers  |
               +--------------+---------------+       +--------------+---------------+
                              |                                      |
                              v                                      v
               +--------------+---------------+       +--------------+---------------+
               |           MODULE 2           |       |        MILESTONES 3, 4 & 5   |
               |  Algorithms & Protocols      +------>|  LocalStack SQS, Postgres    |
               |  (Queues, Outbox, streaming) |       |  Outbox, Kafka backpressure   |
               +--------------+---------------+       +--------------+---------------+
                              |                                      |
                              v                                      v
               +--------------+---------------+       +--------------+---------------+
               |           MODULE 3           |       |         MILESTONES 6 & 7     |
               |  Real-World Software Building +------>|  Multi-Pod Kubernetes, Redis |
               |  (Scale, clustering, DevOps) |       |  Global Rate-Limiter, DORA   |
               +------------------------------+       +------------------------------+
```

### Module 1: Introduction to Distributed Systems [5]
*   **Theoretical Focus:** Core tenets of distributed processing, message formatting, and localized correctness [5].
*   **Project Milestones:**
    *   **Milestone 1: The Walking Skeleton (ACTIVE):** Establish a framework-decoupled, physical module boundary using Go workspaces (Path A). Wire up a Gin HTTP handler returning `202 Accepted` on `/v1/telemetry`, and write a black-box E2E test using Testcontainers-go to verify basic transport [8].
    *   **Milestone 2: In-Memory Concurrency & Buffering:** Design an unbuffered memory dispatch ring inside `internal/engine` utilizing native **Go channels** and worker pools. Optimize memory allocation to handle 100,000 RPS by parsing metric slices directly on a concrete parent struct, bypassing dynamic maps and heap boxing.

### Module 2: Distributed Systems for Practitioners [5]
*   **Theoretical Focus:** Basic protocols, distributed queues, data consistency, and event-sourced state replication [5].
*   **Project Milestones:**
    *   **Milestone 3: SQS Integration with LocalStack:** Flush in-memory buffers into AWS SQS. To ensure Angie can test these AWS resources completely offline without paying cloud costs, integrate **LocalStack** inside the Testcontainers E2E harness. Use GoMock to generate clean interfaces for local unit tests [10, 17].
    *   **Milestone 4: Transactional Outbox Pattern:** Maintain strong transactional consistency by saving telemetry events and system outbox logs in a single atomic PostgreSQL write operation, then spawn background workers to poll and publish [9].
    *   **Milestone 5: Event-Sourced Streaming with Apache Kafka:** Stream outbox records into Kafka partitions. Programmatically coordinate consumer groups and implement consumer-side backpressure (manually pausing and resuming partition ingestion) to protect resources under spiky loads [12, 22].

### Module 3: Distributed Systems: Building Software for the Real World [5]
*   **Theoretical Focus:** Horizontal scaling, failure domains, clustering, and cluster coordination algorithms [5].
*   **Project Milestones:**
    *   **Milestone 6: Distributed Coordination & Kubernetes:** Deploy multiple instances of your service on a multi-pod Kubernetes cluster (Minikube). Set up Redis as a shared distributed state store to coordinate alert thresholds and enforce global rate-limiting across pods.
    *   **Milestone 7: CD Automation & DevOps Observability:** Build an automated GitHub Actions CI/CD pipeline and integrate DORA Four Keys metrics trackers (Deployment Frequency, Lead Time, CFR, MTTR) to monitor engineering velocity and reliability [21].

---

## 🤖 Interaction Directives for Local Claude

When serving as Angie's mentor, follow these strict pedagogical guidelines:

1.  **Do Not Code for Her (Strict Non-Prescriptive Guidance):** Angie is an experienced backend and systems engineer. Do not spit out full-file templates. Act as an advisor: outline architectural trade-offs (Option A vs. Option B), analyze memory structures, discuss garbage collector behavior under load, and let her write the code.
2.  **Teach & Reinforce Key Concepts Actively:** For every milestone, proactively explain, teach, and review relevant system patterns:
    *   *Distributed Systems Concepts:* Eventual consistency, CAP Theorem, linearizability, outbox/inbox guarantees, backpressure patterns, distributed rate-limiting, and transactional delivery guarantees.
    *   *TDD & BDD Concepts:* The Outside-In TDD cycle (GOOS) [6]. Ensure she writes or updates the abstract *Specification* (The Rules) and *Driver* (The Plumbing) to fail first (Red) before writing any production code in `internal/` (Green), and only then optimizing structures (Refactor) [15].
    *   *Go Concurrency Patterns:* The elegant coordination of goroutines, channels (buffered vs. unbuffered), select statements, close signals, and mutex-free state ownership.
3.  **Honor Her Preferred Tooling:** Always support her choice of repo orchestration tools. Favour `mise.toml` for managing directory-agnostic tool configurations and task executions over legacy Makefiles.

---

## 📖 References & Links

1.  **Actor Model | Ergo Framework documentation**  
    *Enforces actor model isolation, supervisor trees, and mutual TLS in Go.*  
    [https://docs.ergo.services/](https://docs.ergo.services/)
2.  **Actor model - Wikipedia**  
    *Historical and theoretical foundation of concurrent computation via asynchronous message passing.*  
    [https://en.wikipedia.org/wiki/Actor_model](https://en.wikipedia.org/wiki/Actor_model)
3.  **ActorDB: A Unified Database Model (arXiv PDF)**  
    *Integrating single-writer actor persistence, stream processing, and zero-trust security.*  
    [https://arxiv.org/abs/2510.01234](https://arxiv.org/abs/2510.01234)
4.  **Angela Breuer Resume - 06-01-26 - Initial Draft.docx**  
    *Archived professional curriculum vitae highlighting embedded and systems administration credentials.*  
    [Local Source Reference / Word Document]
5.  **Become a Distributed Systems Professional - Educative**  
    *Curated graduate Skill Path focusing on distributed protocols, algorithms, and practical architecture.*  
    [https://www.educative.io/paths/become-a-distributed-systems-professional](https://www.educative.io/paths/become-a-distributed-systems-professional)
6.  **Growing Object-Oriented Software Guided by Tests (GOOS PDF)**  
    *The pioneering guide on TDD outside-in, using mock objects to describe component relationships.*  
    [https://github.com/GunterMueller/Books-3/blob/master/Growing%20Object%20Oriented%20Software%20Guided%20by%20Tests.pdf](https://github.com/GunterMueller/Books-3/blob/master/Growing%20Object%20Oriented%20Software%20Guided%20by%20Tests.pdf)
7.  **Growing Object-Oriented Software Guided by Tests: About the Book**  
    *Methodological patterns for establishing stable, evolutionary software designs.*  
    [https://growing-object-oriented-software.com/](https://growing-object-oriented-software.com/)
8.  **GitHub - angiebrr/athenas-telemetry-svc · GitHub**  
    *The target course repository housing the isolated multi-module Gin telemetry skeleton.*  
    [https://github.com/angiebrr/athenas-telemetry-svc](https://github.com/angiebrr/athenas-telemetry-svc)
9.  **Go Kafka: Transactional Outbox Pattern in Microservices (Kcode Video)**  
    *Architectural patterns for guaranteeing database and message broker consistency.*  
    [https://www.youtube.com/watch?v=youtube-outbox-placeholder](https://www.youtube.com/watch?v=youtube-outbox-placeholder)
10. **GoMock - Uber Go Mock Framework**  
    *Google/Uber toolset for generating typesafe mock implementations for Go interface structures.*  
    [https://github.com/uber-go/mock](https://github.com/uber-go/mock)
11. **How to Handle Backpressure in Kafka Consumers - OneUptime**  
    *Mitigation strategies for consumer lag including partition pausing and flow control loops.*  
    [https://oneuptime.com/blog/post/2026-01-24-handle-backpressure-kafka-consumers/view](https://oneuptime.com/blog/post/2026-01-24-handle-backpressure-kafka-consumers/view)
12. **Kafka, distributed coordination and the actor model - Daniel Lebrero**  
    *Conceptualizing event-driven streaming components as partition-isolated single-threaded actors.*  
    [https://danlebrero.com/2018/04/09/kafka-distributed-coordination-actor-model/](https://danlebrero.com/2018/04/09/kafka-distributed-coordination-actor-model/)
13. **Learn Go with Tests | Learn Go with tests (Quii GitBook)**  
    *Golang-centric TDD guide introducing shared specs, adapters, and mock-free design.*  
    [https://quii.gitbook.io/learn-go-with-tests](https://quii.gitbook.io/learn-go-with-tests)
14. **Make It Work, Make It Right, Make It Fast - C2 Wiki**  
    *Kent Beck's developmental pattern prioritizing correct design over premature tuning.*  
    [https://wiki.c2.com/?MakeItWorkMakeItRightMakeItFast](https://wiki.c2.com/?MakeItWorkMakeItRightMakeItFast)
15. **Resume Notes 06-01-26.docx**  
    *Retrospective engineering journal documenting professional challenges, re-orgs, and technical wins.*  
    [Local Source Reference / Word Document]
16. **Testcontainers for Go Quickstart**  
    *Guide to programmatically spinning up databases, queues, and LocalStack containers in Go.*  
    [https://golang.testcontainers.org/quickstart/](https://golang.testcontainers.org/quickstart/)
17. **The Actor Model in Go: DEV Community**  
    *An introduction to applying erlang-style mailboxes and scalability patterns using lightweight Go actors.*  
    [https://dev.to/forkbikash/the-actor-model-in-go-simplifying-concurrent-programming-1j9d](https://dev.to/forkbikash/the-actor-model-in-go-simplifying-concurrent-programming-1j9d)
18. **The Practical Test Pyramid - Martin Fowler (Ham Vocke)**  
    *Establishing an automated testing mix that prioritizes feedback speed and regression stability.*  
    [https://martinfowler.com/articles/practical-test-pyramid.html](https://martinfowler.com/articles/practical-test-pyramid.html)
19. **Understanding Screenplay Pattern Part 1 - Cucumber Blog**  
    *Decoupling automation complexity by expressing assertions as actor-driven, composite tasks.*  
    [https://cucumber.io/blog/understanding-screenplay-part-1](https://cucumber.io/blog/understanding-screenplay-part-1)
20. **Using the Four Keys to Measure Your DevOps Performance - Google Cloud**  
    *Evaluating software engineering maturity through automated metrics pipeline dashboards.*  
    [https://cloud.google.com/blog/products/devops-sre/using-the-four-keys-to-measure-your-devops-performance](https://cloud.google.com/blog/products/devops-sre/using-the-four-keys-to-measure-your-devops-performance)
21. **Writing and Testing an Event Sourcing Microservice - Semaphore**  
    *Tutorial on dockerizing event-sourced Go systems, cluster scaling, and sarama consumer groups.*  
    [https://semaphore.io/community/tutorials/writing-and-testing-an-event-sourcing-microservice-with-kafka-and-go](https://semaphore.io/community/tutorials/writing-and-testing-an-event-sourcing-microservice-with-kafka-and-go)
22. **actor package - github.com/asynkron/protoactor-go/actor**  
    *Official package documentation for types representing asynchronous actor-based processing.*  
    [https://pkg.go.dev/github.com/asynkron/protoactor-go/actor](https://pkg.go.dev/github.com/asynkron/protoactor-go/actor)
23. **teivah/gosiris: An actor framework for Go - GitHub**  
    *An archived, open-source framework demonstrating remote deployment and discoverability using etcd.*  
    [https://github.com/teivah/gosiris](https://github.com/teivah/gosiris)
