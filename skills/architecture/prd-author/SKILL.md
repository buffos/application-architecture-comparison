---
name: prd-author
description: Turn structured discovery material, rough application notes, interview transcripts, or partial requirement fragments into a high-quality product requirements document. Use when discovery is already strong enough to define the product clearly and Codex needs to produce a rigorous, scoped, architecture-neutral `prd.md` that captures goals, actors, workflows, rules, constraints, non-goals, and acceptance scenarios. Do not use when the main problem is unresolved ambiguity or contradictions; use `requirements-gap-analyzer` or `requirements-discovery-interviewer` first.
---

# PRD Author

Write a PRD that is concrete enough to guide design and implementation, but neutral enough not to hard-code one architectural solution prematurely.

Do not write glossy product fluff. Write a planning-grade product document.

## Workflow

### 1. Confirm the Inputs Are PRD-Ready

Before drafting, confirm the input material answers these questions well enough:

- what problem the product solves
- who uses it
- what workflows matter most
- what is in and out of scope
- what business rules and constraints matter
- what success looks like

Use [references/prd-readiness.md](references/prd-readiness.md).

If the material is not ready, stop and ask targeted follow-up questions instead of fabricating certainty.

### 2. Normalize the Discovery Material

Distill the source material into:

- problem statement
- actors
- scope
- workflows
- core business objects
- policies and rules
- reporting needs
- open issues and assumptions

Resolve obvious terminology drift before drafting the PRD.

### 3. Write the PRD in a Stable Structure

Use the section model in [references/prd-sections.md](references/prd-sections.md).

At minimum, cover:

- title
- purpose
- product summary
- goals
- non-goals
- user/actor definitions
- domain scope
- business capabilities
- required workflows
- business rules
- functional requirements
- non-functional requirements
- acceptance scenarios
- rationale for fit, if the exercise is comparative or educational

### 4. Keep the PRD Product-Centric

The PRD should describe:

- what the system must do
- who it serves
- what constraints matter
- what outcomes define success

The PRD should not drift into:

- package layout
- adapter design
- repository choices
- domain model implementation details
- API endpoint exhaustiveness
- database schema design

Those belong in later artifacts.

### 5. Make Requirements Testable

Every important feature or rule should be expressible as:

- a workflow
- a constraint
- an acceptance scenario
- or a comparison-relevant requirement

Vague statements like “support reporting” are insufficient. Replace them with specific operational or managerial questions the system must answer.

### 6. Separate Hard Requirements from Preferences

If the user mixes business needs with preferences like:

- “Use plugins”
- “Support microservices later”
- “I want DDD”

Capture them accurately but classify them correctly:

- product requirement
- architectural preference
- future direction
- non-goal for current scope

### 7. Surface Assumptions

If the source material is incomplete but draftable, make assumptions explicit. Do not smuggle them in as settled facts.

### 8. Run the PRD Quality Gate

Before finalizing, review the document using [references/prd-quality-gate.md](references/prd-quality-gate.md).

If the PRD reads like:

- a feature wishlist
- a UI outline
- a technical design spec
- or a marketing brief

it is wrong. Fix it.

## Writing Rules

### Write in Business Terms First

Prefer:

- actors
- workflows
- capabilities
- policies
- constraints

Avoid premature technical detail unless it changes the requirement itself.

### Be Concrete

Replace:

- “users can manage orders”

with:

- “sales clerks create quotes, submit discount exceptions for approval, and convert approved quotes into orders”

### Include Failure and Exception Paths

A PRD without failure, rejection, exception, or approval flows is incomplete.

For every important happy path, try to include:

- at least one blocked path
- at least one exception path
- at least one policy-driven branch

### Bound the Scope Aggressively

Strong PRDs explicitly state what is out of scope. If the product can expand infinitely, the PRD is weak.

### Prefer Comparability When the Exercise Is Architectural

If the user intends to build the same app multiple ways, the PRD must define a stable product baseline:

- same actors
- same workflows
- same business rules
- same visible behavior

This prevents later architecture comparisons from becoming meaningless.

## Anti-Patterns

Do not:

- write generic corporate fluff
- hide missing decisions behind broad language
- treat CRUD as sufficient product definition
- collapse business rules into “validation”
- omit operational reporting needs
- over-specify internal technical design
- confuse PRD content with canonical domain model content

## Output Contract

Produce:

- a complete PRD draft
- a short list of explicit assumptions
- a short list of unresolved questions, only if they materially matter

The PRD should be good enough that the next skills can derive:

- a canonical domain model
- a canonical use-case/application-service model
- a canonical API/CLI contract

Persist the output to:

- `docs/prd.md`

Behavior:

- create the file if it does not exist
- update it in place when refining scope or requirements
- preserve stable product semantics unless upstream discovery changed them intentionally

## Resources

Read only what you need:

- [references/prd-sections.md](references/prd-sections.md): recommended PRD structure and section intent
- [references/prd-readiness.md](references/prd-readiness.md): when discovery is sufficient to draft
- [references/prd-quality-gate.md](references/prd-quality-gate.md): review checklist before finalizing
