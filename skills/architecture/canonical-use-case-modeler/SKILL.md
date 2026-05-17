---
name: canonical-use-case-modeler
description: Derive a canonical application/use-case model from a stable PRD and canonical domain model. Use when Codex must define the stable application-layer surface of a system, including commands, queries, application services, orchestration responsibilities, transaction expectations, idempotency-sensitive operations, failure categories, and end-to-end use-case chains without collapsing into controllers, transport contracts, or implementation-specific package layouts. Do not use before the domain semantics are stable enough to name intents cleanly.
---

# Canonical Use-Case Modeler

Produce the stable application-layer contract that sits between the domain model and external interfaces.

This skill defines what the system does through application services and use cases. It does not define HTTP endpoints, CLI flags, or internal package layout.

## Workflow

### 1. Confirm the Inputs Are Use-Case Ready

Before drafting, verify the available material is strong enough to answer:

- what user or system intents exist
- which workflows change state
- which workflows only read data
- which domain concepts participate in each workflow
- which rules and failures materially affect orchestration

Use [references/use-case-readiness.md](references/use-case-readiness.md).

If the material is not ready, ask follow-up questions instead of inventing application behavior.

### 2. Extract Intents, Not Screens

Start from business intents such as:

- create quote
- submit for approval
- convert to order
- capture payment
- accept return

Do not start from UI pages or menu items.

If the input is screen-oriented, translate it into business-capable actions first.

### 3. Separate Commands from Queries

For each use case, determine whether it:

- changes state
- reads state
- coordinates multiple aggregates
- publishes events or triggers projection work

Commands and queries should be stable enough that later HTTP and CLI contracts can map onto them cleanly.

Use [references/use-case-structure.md](references/use-case-structure.md).

### 4. Define Application Service Boundaries

Group use cases into coherent application services or modules.

A good grouping usually reflects:

- cohesive workflow ownership
- shared dependencies
- shared orchestration concerns

Do not group services only by CRUD resource naming.

### 5. Make Orchestration Explicit

For each non-trivial command, document:

- inputs
- outputs
- preconditions
- application responsibilities
- possible outcomes

The application layer should:

- load aggregates
- coordinate domain behavior
- persist changes
- enforce transaction boundaries
- handle idempotency where needed

It should not absorb domain logic that belongs in the model or policies.

### 6. Define Transaction and Consistency Expectations

For each important command, determine whether it is:

- single-aggregate
- coordinated multi-aggregate
- retry-sensitive
- event-emitting

Use [references/orchestration-heuristics.md](references/orchestration-heuristics.md).

### 7. Model Failures as Business Outcomes

Document failure classes such as:

- validation failure
- missing resource
- business rule violation
- state conflict
- infrastructure failure

If a command can fail in a meaningful business way, make that visible in the use-case model.

### 8. Capture End-to-End Use-Case Chains

Describe canonical flows that tie the commands together:

- happy path
- approval path
- shortage path
- return path
- plugin or policy variation path

This makes later contract and testing work far more stable.

### 9. Run the Use-Case Quality Gate

Review the output with [references/use-case-quality-gate.md](references/use-case-quality-gate.md).

If the document reads like:

- a controller list
- a REST route outline
- a service layer with no orchestration semantics
- or CRUD wrappers pretending to be use cases

it is wrong. Revise it.

## Modeling Rules

### Use Intent-Based Names

Prefer:

- `SubmitQuoteForApproval`
- `ConvertQuoteToOrder`
- `ApprovePaymentReview`

Avoid:

- `UpdateQuote`
- `SetOrderStatus`
- `ProcessEntity`

unless the business intent is truly that generic.

### Keep Commands and Queries Distinct

Commands mutate or coordinate state.

Queries return views, summaries, or snapshots.

Do not hide mutations inside read-sounding operations.

### Keep the Layer Architecture-Neutral

This model should work for:

- layered service methods
- clean/hexagonal input ports
- transaction scripts
- rich domain model orchestration
- modular monolith modules

Do not bake in HTTP semantics or framework assumptions.

### Model Application Responsibilities, Not Everything

A use-case document should explain:

- what the application layer coordinates
- which dependencies it uses
- which outcomes it returns

It should not restate the full domain model or jump ahead to transport contracts.

### Call Out Idempotency Where Retries Matter

If a command could create duplicates or replay dangerous work, mark it as retry-sensitive and describe the expected canonical behavior.

### Preserve Read/Write Tension

Some queries will later come from projections or specialized read models. The use-case model should allow for that without forcing it.

## Anti-Patterns

Do not:

- derive use cases from tables
- mirror UI navigation directly
- collapse everything into CRUD verbs
- hide business failure modes
- put domain logic into application responsibilities
- drift into endpoint design
- confuse repositories with application services

## Output Contract

Produce a canonical use-case/application-service document that includes:

- purpose
- application-layer goals
- design principles
- application boundaries
- canonical commands and queries
- application service groupings
- command contract shape
- query contract shape
- transaction boundaries
- idempotency expectations
- cross-cutting application concerns
- failure model
- domain events at the application boundary
- end-to-end use-case chains
- mapping guidance for different architectures

The result should be stable enough that later skills can derive a canonical API/CLI contract from it.

Persist the output to:

- `docs/canonical-use-cases.md`

Behavior:

- create the file if it does not exist
- update it in place when commands, queries, or orchestration semantics change
- keep command and query names stable unless upstream artifacts changed materially

## Resources

Read only what you need:

- [references/use-case-readiness.md](references/use-case-readiness.md): when the inputs are ready for application modeling
- [references/use-case-structure.md](references/use-case-structure.md): how to structure commands, queries, and service boundaries
- [references/orchestration-heuristics.md](references/orchestration-heuristics.md): how to think about orchestration, transactions, and retries
- [references/use-case-quality-gate.md](references/use-case-quality-gate.md): review checklist before finalizing
