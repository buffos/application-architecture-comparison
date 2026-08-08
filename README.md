# Application Architecture Comparison Repository

## Purpose

This repository implements the same business application across multiple
architectural styles.

The goal is not to declare one architecture universally best. The goal is to
keep the problem space stable while changing the internal design, so the
trade-offs become visible in code:

- the same business domain
- the same broad workflows
- the same canonical reference documents
- different boundaries, dependencies, and places for business logic

The repository is organized as a sequence of small, tagged lessons. Each
lesson adds one architectural or workflow idea and leaves the code runnable.

## Domain and Reference Documents

Most architecture tracks use the quote/order domain described by the product
requirements. The canonical material is kept in [docs](docs/):

- [Product requirements](docs/prd.md)
- [Canonical domain model](docs/canonical-domain-model.md)
- [Canonical use cases](docs/canonical-use-cases.md)
- [Canonical API/CLI contract](docs/canonical-api-cli-contract.md)
- [Domain glossary](docs/domain-glossary.md)
- [Architecture catalogue](docs/architectures.md)
- [Git lesson workflow](docs/git-how-to.md)

The [Blackboard Architecture](blackboard-architecture/) track is the
intentional exception: it uses the smaller Smart Invoice-to-Payment Matcher
example so that the shared blackboard, knowledge sources, controller, and
concurrent evaluation remain easy to see.

## Implemented Architecture Tracks

The table below is the current lesson inventory. “Start” and “Latest” link to
the first and final lesson currently present in each track.

| Architecture | Lessons | Start | Latest |
| --- | ---: | --- | --- |
| [Layered](layered-architecture/) | 14 (001–014) | [001 — skeleton](layered-architecture/lessons/001-layered-skeleton.md) | [014 — conclusion](layered-architecture/lessons/014-why-not-stop-at-layered.md) |
| [Hexagonal / Ports and Adapters](hexagonal-architecture/) | 33 (001–033) | [001 — skeleton](hexagonal-architecture/lessons/001-ports-and-adapters-skeleton.md) | [033 — conclusion](hexagonal-architecture/lessons/033-why-not-stop-at-hexagonal.md) |
| [Clean](clean-architecture/) | 34 (000–033) | [000 — transition](clean-architecture/lessons/000-from-hexagonal-to-clean.md) | [033 — conclusion](clean-architecture/lessons/033-why-not-stop-at-clean.md) |
| [Onion](onion-architecture/) | 34 (000–033) | [000 — transition](onion-architecture/lessons/000-from-clean-to-onion.md) | [033 — conclusion](onion-architecture/lessons/033-why-not-stop-at-onion.md) |
| [Modular Monolith](modular-monolith/) | 34 (000–033) | [000 — transition](modular-monolith/lessons/000-from-onion-to-modular-monolith.md) | [033 — conclusion](modular-monolith/lessons/033-why-not-stop-at-modular-monolith.md) |
| [Microkernel / Plugin](microkernel-architecture/) | 34 (000–033) | [000 — transition](microkernel-architecture/lessons/000-from-modular-monolith-to-microkernel.md) | [033 — conclusion](microkernel-architecture/lessons/033-why-not-stop-at-microkernel.md) |
| [Component-Based](component-based-architecture/) | 34 (000–033) | [000 — transition](component-based-architecture/lessons/000-from-microkernel-to-component-based.md) | [033 — conclusion](component-based-architecture/lessons/033-why-not-stop-at-component-based.md) |
| [Domain-Driven Design](domain-driven-design-architecture/) | 34 (000–033) | [000 — ubiquitous language](domain-driven-design-architecture/lessons/000-ubiquitous-language-and-quote-aggregate.md) | [033 — conclusion](domain-driven-design-architecture/lessons/033-why-not-stop-at-ddd.md) |
| [Transaction Script](transaction-script-architecture/) | 34 (000–033) | [000 — transition](transaction-script-architecture/lessons/000-from-ddd-to-transaction-script.md) | [033 — conclusion](transaction-script-architecture/lessons/033-why-not-stop-at-transaction-script.md) |
| [Active Record](active-record-architecture/) | 34 (000–033) | [000 — transition](active-record-architecture/lessons/000-from-transaction-script-to-active-record.md) | [033 — conclusion](active-record-architecture/lessons/033-why-not-stop-at-active-record.md) |
| [Rich Domain Model](rich-domain-model-architecture/) | 34 (000–033) | [000 — transition](rich-domain-model-architecture/lessons/000-from-active-record-to-rich-domain-model.md) | [033 — conclusion](rich-domain-model-architecture/lessons/033-why-not-stop-at-rich-domain-model.md) |
| [Blackboard](blackboard-architecture/) | 6 (001–006) | [001 — skeleton](blackboard-architecture/lessons/001-blackboard-skeleton.md) | [006 — concurrent knowledge sources](blackboard-architecture/lessons/006-concurrent-knowledge-sources.md) |
| [Rules Engine / Knowledge-Based](rules-engine-architecture/) | 34 (000–033) | [000 — transition](rules-engine-architecture/lessons/000-from-rich-domain-model-to-rule-engine.md) | [033 — final trade-offs](rules-engine-architecture/lessons/033-plugin-rule-extension-and-final-tradeoffs.md) |

The canonical ordering and the architectural summaries are maintained in
[docs/architectures.md](docs/architectures.md).

## How To Study a Track

For each architecture:

1. Read the lessons in numeric order.
2. Inspect the code introduced by the current lesson.
3. Run the tests and the demo.
4. Compare the result with the previous architecture.

Every lesson has a corresponding commit and annotated tag. Tags use the
architecture slug and lesson number, for example:

~~~powershell
git tag --list "rules-engine-*"
git switch --detach rules-engine-033
~~~

Return to the development branch when finished inspecting a tag:

~~~powershell
git switch main
~~~

The exact commit and tag workflow is documented in
[docs/git-how-to.md](docs/git-how-to.md).

## Running the Examples

The architecture implementations are Go modules with tests and executable
demos. From the repository root, choose a track and run:

~~~powershell
Set-Location rules-engine-architecture
go test ./...
go run ./cmd/quote-demo
~~~

The Blackboard track has its own matcher demo:

~~~powershell
Set-Location blackboard-architecture
go test ./...
go run ./cmd/matcher-demo
~~~

For the other tracks, replace the directory name and keep the same
go test ./... / go run ./cmd/quote-demo commands.

## What To Compare

As you move through the tracks, pay attention to:

- where business rules live
- which direction dependencies point
- how persistence and external systems are isolated
- how workflows are composed and tested
- how new policies, modules, plugins, or components are added
- what complexity the architecture makes easier and what complexity it
  introduces

The Rules Engine track returns to the canonical quote/order domain and makes
facts, rules, working memory, inference cycles, derived facts, stage
fulfilments, approvals, and extension points explicit. The Blackboard track
focuses on opportunistic reasoning: independent knowledge sources contribute
evidence to shared working memory until a match converges.

## Maintaining This README

When a new architecture or lesson is added:

- update [docs/architectures.md](docs/architectures.md) first
- add or update its row in the architecture table above
- verify the start/latest lesson links and count
- keep the demo commands and git workflow current
