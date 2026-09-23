# athenas-telemetry-svc

[![Validation](https://github.com/angiebrr/athenas-telemetry-svc/actions/workflows/validation.yaml/badge.svg)](https://github.com/angiebrr/athenas-telemetry-svc/actions/workflows/validation.yaml)

- [athenas-telemetry-svc](#athenas-telemetry-svc)
  - [About](#about)
    - [Architecture](#architecture)
  - [Running tests](#running-tests)
  - [More details](#more-details)
    - [Self-study](#self-study)
    - [TDD](#tdd)
    - [Deployment](#deployment)
  - [License](#license)

## About

`athenas-telemetry-svc` is intended to be a distributed telemetry ingestion service written in Go that I'm using as a self-study mechanism to learn more about distributed systems in addition to keeping my skills sharp while job hunting.

> Currently on Milestone 2 of 7, see [self-study.md](docs/self-study.md)

### Architecture

We have 2 Golang modules, `athenas-telemetry-svc/` and `athenas-acceptance-tests/`:

| Path | Role |
|---|---|
| `athenas-telemetry-svc/` | The service. `cmd/athenas` → `internal/api` (Gin transport) → `internal/telemetry` (domain rules) → `internal/data` (storage port), plus `models` and the exported `specifications/` package |
| `athenas-acceptance-tests/` | Black-box acceptance test suite that runs the service locally and runs the exported specifications against it |

Some details worth noting:
- `athenas-telemetry-svc/specifications/` lives in the service module and is intentionally exported rather than in `internal/` so the acceptance tests can access them
- Acceptance tests build `deploy/Dockerfile` so the tests run against the same image that would be shipped

## Running tests

Tooling is managed with [mise](https://mise.jdx.dev/), and `mise install` sets up the necessary golang toolchain.

```bash
mise run test                                # every project, including the container suite (needs Docker)
mise run test --short                        # every project, skipping the acceptance suite
mise run test telemetry-svc                  # one project
mise run test telemetry-svc . TestFoo        # a single test by name or regex
mise run test telemetry-svc --race           # with the data race detector
mise run test telemetry-svc --bench --memprofile   # benchmarks, allocations, memory profile
```

`mise` also has other tasks that were defined in [`mise.toml`](mise.toml), and can be seen with `mise tasks ls`.

## More details

### Self-study

I'm building this as a self-study project alongside [a distributed systems course in Educative](https://www.educative.io/path/become-a-distributed-systems-professional) and **[Growing Object-Oriented Software, Guided by Tests](https://growing-object-oriented-software.com/)** (i.e. **GOOS**). Claude acts as a mentor rather than writing the code for me: it handles docs, tooling and CI upkeep, and I write the internal code, the spec and the tests.

See [self-study.md](docs/self-study.md) for how the self-study works, the milestone roadmap, and the reading list.

### TDD

When working on this service, I followed TDD principles influenced by **[GOOS](https://growing-object-oriented-software.com/)** and **[Learn Go with Tests](https://quii.gitbook.io/learn-go-with-tests)**, and more information on why can be found here: [docs/why-tdd.md](docs/why-tdd.md)

### Deployment

Since this is a toy service that I don't want to pay cloud costs for, it isn't actually deployed anywhere. 

Does today:

- run automated acceptance tests in CI/CD against the service's Docker image
- run dependency vuln checks in CI/CD
- run automated unit tests in CI/CD

Planned:

- use IaC + Minikube + LocalStack to emulate K8s clusters and cloud services locally

## License

MIT - see [LICENSE](LICENSE).
