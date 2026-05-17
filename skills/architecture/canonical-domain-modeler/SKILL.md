---
name: canonical-domain-modeler
description: Derive a canonical domain model from a stable PRD, glossary, or other sufficiently mature product requirements. Use when Codex needs to convert product requirements into a stable business model with ubiquitous language, subdomains, bounded context candidates, aggregates, entities, value objects, invariants, lifecycles, policies, and domain events without collapsing into database design or architecture-specific implementation details. Do not use while vocabulary or core workflows are still unstable.
---

# Canonical Domain Modeler

Produce the shared business model that all future implementations must preserve, regardless of architecture.

This skill exists to define the product's semantic core, not its storage schema and not its package layout.

## Workflow

### 1. Confirm the Inputs Are Domain-Model Ready

Before drafting, verify the source material is strong enough to answer:

- what the main business concepts are
- which workflows matter most
- which states and transitions matter
- which rules are hard invariants versus policies
- which concepts own behavior versus just carry data

Use [references/domain-readiness.md](references/domain-readiness.md).

If the requirements are still too fuzzy, ask follow-up questions instead of inventing a model.

### 2. Establish Ubiquitous Language

Normalize the business vocabulary first.

Extract:

- core nouns
- overloaded terms
- near-synonyms that should be collapsed
- concept pairs that must stay distinct

Do not let the model drift between product wording and technical shorthand.

### 3. Identify Subdomains and Context Candidates

Separate:

- core business areas
- supporting capabilities
- cross-cutting policy areas

Then decide whether the product naturally suggests:

- bounded contexts
- modules
- conceptual seams

Use [references/domain-structure.md](references/domain-structure.md).

### 4. Choose Aggregate Boundaries Carefully

Define aggregates only where they help preserve consistency and business invariants.

For each candidate aggregate, answer:

- what business object is the root
- what it owns
- what must change together
- what invariants it protects
- why it is one aggregate instead of many

Use [references/aggregate-heuristics.md](references/aggregate-heuristics.md).

### 5. Distinguish Entities, Value Objects, Policies, and Services

Do not treat all nouns the same.

Classify concepts as:

- entity
- value object
- aggregate root
- policy or rule object
- domain service
- event

If a concept has no independent lifecycle and is defined by value, it is a strong value-object candidate.

### 6. Model Invariants and Lifecycle Rules Explicitly

For each important concept, capture:

- creation rules
- valid states
- forbidden transitions
- quantity or threshold constraints
- consistency expectations

A domain model without invariants is just a glossary.

### 7. Capture Cross-Aggregate References and Consistency Expectations

Document:

- what references what
- where IDs or snapshots should be preferred over object graphs
- which operations are naturally single-aggregate
- which operations require coordination across aggregates

### 8. Define Domain Events and Extension Points

Capture business-significant moments that should remain visible across architectures.

Events matter when:

- state meaningfully changes
- downstream read models update
- audit trails matter
- plugins or policy engines may react

Also identify extension points where behavior may vary without changing the core business language.

### 9. Run the Domain Quality Gate

Before finalizing, review the output with [references/domain-quality-gate.md](references/domain-quality-gate.md).

If the document reads like:

- a table design
- an ORM entity dump
- a service inventory with no business semantics
- or a DDD manifesto detached from the product

it is wrong. Revise it.

## Modeling Rules

### Model the Business, Not the Database

Avoid designing:

- tables
- foreign keys
- indexes
- persistence records

unless a reference to persistence is genuinely required to explain a business constraint.

### Model the Language Before the Structure

Start by stabilizing the terms. Good aggregates emerge from the business language and invariants, not from package fantasies.

### Preserve Architecture Neutrality

The domain model should be usable by:

- layered implementations
- transaction scripts
- active record designs
- rich domain models
- clean/hexagonal/onion variants
- rules-engine and plugin variants

Do not write it so that only one of those styles makes sense.

### Make Aggregate Rationale Explicit

Whenever an aggregate is defined, explain why the boundary exists. If you cannot explain the invariant it protects, the aggregate is probably arbitrary.

### Separate Hard Invariants from Flexible Policies

Examples:

- “returned quantity cannot exceed shipped quantity” is likely invariant
- “returns are allowed within 30 days” may be policy

Keep that distinction explicit.

### Prefer Business Events Over Technical Events

Prefer:

- `QuoteApproved`
- `PaymentReviewRequired`
- `ReturnAccepted`

Avoid technical noise like:

- `RowUpdated`
- `EntityPersisted`

## Anti-Patterns

Do not:

- collapse everything into CRUD entities
- model UI screens as domain concepts
- equate aggregates with database tables
- create value objects just for style points
- declare bounded contexts without business justification
- hide missing rules behind generic “validation”
- skip failure, reversal, and exception paths

## Output Contract

Produce a canonical domain model document that includes:

- purpose
- modeling principles
- ubiquitous language
- subdomains
- bounded context candidates
- aggregate design
- entity definitions
- value objects
- enumerations or status vocabularies
- policies and rule objects
- domain services
- invariants
- lifecycle rules
- domain events
- cross-aggregate references
- consistency boundaries
- read-model expectations
- extension points
- minimum canonical scenarios
- mapping guidance for architecture variants

The result should be detailed enough that later skills can derive stable application services and external contracts from it.

Persist the output to:

- `docs/canonical-domain-model.md`

Behavior:

- create the file if it does not exist
- update it in place when domain semantics evolve
- do not silently rename core concepts without coordinating the glossary and PRD

## Resources

Read only what you need:

- [references/domain-readiness.md](references/domain-readiness.md): when the source material is ready for domain modeling
- [references/domain-structure.md](references/domain-structure.md): how to frame subdomains, contexts, and core concepts
- [references/aggregate-heuristics.md](references/aggregate-heuristics.md): how to choose aggregate boundaries without hand-waving
- [references/domain-quality-gate.md](references/domain-quality-gate.md): checklist to review the final model
