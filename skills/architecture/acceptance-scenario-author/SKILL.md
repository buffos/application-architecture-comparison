---
name: acceptance-scenario-author
description: Derive stable, high-value acceptance scenarios from a PRD and canonical design artifacts. Use when Codex must turn already-modeled requirements into behavior-focused scenarios that capture happy paths, exceptions, policy branches, failure modes, and comparison-critical outcomes so the same product behavior can later be tested across architectures, application services, APIs, and CLIs. Do not use as a substitute for writing the requirements themselves.
---

# Acceptance Scenario Author

Turn requirements and design artifacts into stable behavior scenarios that later implementations can all be tested against.

This skill produces canonical scenarios, not framework-specific tests. Its job is to define what must be true from the outside when the system behaves correctly or incorrectly.

## Workflow

### 1. Review the Behavior Sources

Read the available sources, especially:

- PRD
- canonical domain model
- canonical use-case model
- canonical API/CLI contract

Extract:

- major workflows
- key business rules
- state transitions
- exception paths
- externally visible outcomes

### 2. Identify Scenario-Worthy Behaviors

Not every requirement deserves its own canonical scenario.

Prioritize behaviors that:

- define the core happy path
- express important invariants or policies
- branch on approval or exception logic
- involve failure or blocked outcomes
- exercise retry-sensitive or state-sensitive behavior
- matter for architecture comparison fairness

Use [references/scenario-selection.md](references/scenario-selection.md).

### 3. Balance Happy and Unhappy Paths

A strong scenario set should include:

- standard success path
- approval or review path
- invalid or rejected path
- shortage or failure path
- reversal or compensation path
- extensibility or policy variation path, where relevant

If the scenario set only proves the happy path, it is weak.

### 4. Write Scenarios in Stable Behavioral Language

Prefer concise scenario statements using structures like:

- Given / When / Then
- Preconditions / Action / Outcome

Keep the language:

- business-centered
- implementation-neutral
- externally observable

Do not write assertions that depend on package layout or persistence internals.

Use [references/scenario-structure.md](references/scenario-structure.md).

### 5. Tie Scenarios to Rules and Outcomes

Each scenario should make clear:

- what initial condition matters
- what action occurs
- what rule or policy is being exercised
- what outcome must be visible

If a scenario exists only as a vague workflow summary, sharpen it.

### 6. Preserve Comparison Stability

For architecture-comparison use cases, scenario wording should be stable enough that:

- every implementation can be tested against the same scenario
- failure and success semantics stay equivalent
- differences in internal design do not alter the expected outcome

### 7. Include Read-Side and Operational Scenarios Where Needed

Acceptance scenarios are not only about commands. Some should cover:

- reports
- queues
- visibility of approvals
- low-stock or operational summaries

If these are part of the product behavior, scenario coverage should include them.

### 8. Run the Scenario Quality Gate

Review the final scenario set with [references/scenario-quality-gate.md](references/scenario-quality-gate.md).

If the scenarios read like:

- a feature checklist
- low-value UI steps
- technical integration tests
- or repetitive variants with no new behavioral signal

they are wrong. Revise them.

## Writing Rules

### Prefer Behavioral Outcomes

Good:

- “Then shipment is rejected because payment is not accepted”

Weak:

- “Then the shipment service returns false”

### Keep Preconditions Sharp

The “Given” should isolate the thing that matters:

- customer tier
- stock level
- approval status
- payment state
- product category
- plugin enabled state

### One Main Behavioral Point per Scenario

A scenario can have multiple effects, but it should test one central behavioral idea.

### Cover Policy Branches Explicitly

If behavior changes based on:

- thresholds
- categories
- roles
- terms
- plugin contributions

then at least some scenarios should isolate those branches.

### Avoid Redundant Scenario Variants

Do not create five scenarios that prove the same rule with minor wording changes.

## Anti-Patterns

Do not:

- write only happy-path scenarios
- encode API routes or CLI flags into the canonical scenario text
- depend on internal object structures
- mistake long step-by-step UI procedures for acceptance behavior
- omit reports or read models if they are part of product value

## Output Contract

Produce a scenario set that includes:

- scenario title
- short purpose or target behavior
- precondition(s)
- action
- expected outcome
- optional note on what rule or workflow branch it covers

The output should be reusable later for:

- application-service tests
- contract tests
- architecture comparison checks
- implementation readiness review

Persist the output to:

- `docs/acceptance-scenarios.md`

Behavior:

- create the file if it does not exist
- update it in place as behaviors are added or clarified
- keep scenario titles and meanings stable so later tests can map to them consistently

## Resources

Read only what you need:

- [references/scenario-selection.md](references/scenario-selection.md): how to choose which behaviors deserve canonical scenarios
- [references/scenario-structure.md](references/scenario-structure.md): how to write stable behavioral scenarios
- [references/scenario-coverage.md](references/scenario-coverage.md): coverage dimensions to include
- [references/scenario-quality-gate.md](references/scenario-quality-gate.md): review checklist before finalizing
