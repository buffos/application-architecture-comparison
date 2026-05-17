---
name: architecture-readiness-reviewer
description: Review a PRD and the canonical design artifacts before architecture-specific implementation begins. Use when Codex must audit a `prd.md`, canonical domain model, canonical use-case model, canonical API/CLI contract, and related scenarios for contradictions, missing decisions, underspecified behaviors, architectural bias, testability gaps, and overall implementation readiness, then report findings in a review-first format. Do not use as a drafting skill; use it as the gate after the main artifacts exist.
---

# Architecture Readiness Reviewer

Review the reference document set as a gate before implementation begins.

This skill is not a rewriting skill by default. Its job is to find risks, gaps, contradictions, and weak assumptions across the artifact set so the user can fix them before starting architecture-specific builds.

## Workflow

### 1. Review the Full Artifact Set, Not One File in Isolation

When available, review:

- `prd.md`
- canonical domain model
- canonical use-case/application-service model
- canonical API/CLI contract

Do not evaluate any one artifact as if it stands alone. Read them as one system of references.

### 2. Use a Findings-First Review Mindset

Lead with concrete findings, not summaries.

Prioritize:

- contradictions
- missing business semantics
- underspecified workflows
- ambiguous or unstable terminology
- misaligned statuses or error vocabularies
- architectural bias that will distort comparisons
- missing testability anchors

### 3. Check Cross-Document Consistency

Compare the documents for consistency in:

- actor names
- business object names
- workflow names
- status values
- business rules
- extension points
- error and failure semantics

Use [references/cross-document-checks.md](references/cross-document-checks.md).

If one document says a behavior exists and another silently omits or weakens it, call that out.

### 4. Check Readiness for Architecture Comparison

Ask whether the artifacts actually support fair architecture comparison.

Review whether they preserve:

- one stable product baseline
- one stable domain vocabulary
- one stable use-case surface
- one stable external contract

If the documents leave too much room for architecture variants to change product behavior, flag it.

### 5. Check Missing Decisions and Unbounded Assumptions

Look for places where the documents say or imply:

- “it depends”
- “configurable” without scope
- “support reporting” without concrete read needs
- “approval” without triggers, actors, or outcomes
- “returns” without eligibility logic

Use [references/readiness-criteria.md](references/readiness-criteria.md).

If a gap will materially affect implementation shape or comparison fairness, it is review-worthy.

### 6. Check Architecture Neutrality

The artifacts should enable multiple architectures fairly.

Flag places where the documents accidentally bias toward:

- DDD-only thinking
- CRUD-only thinking
- REST-only thinking
- event-driven-only thinking
- plugin-first design without business need

Use [references/neutrality-checks.md](references/neutrality-checks.md).

### 7. Check Testability and Traceability

Review whether the artifacts support:

- acceptance tests
- command/query tests
- domain rule tests
- contract tests
- comparison tests across architectures

If a requirement or workflow cannot be tested from the existing docs, flag that gap.

### 8. Report Findings with Severity and References

For each finding, include:

- severity
- affected artifact(s)
- the issue
- why it matters
- what needs clarification or correction

Use the review format in [references/review-format.md](references/review-format.md).

### 9. State Residual Risk Even When No Findings Exist

If no major findings exist, say so explicitly and then mention:

- residual risks
- weaker areas
- recommended follow-up checks

Do not fake findings. Do not pretend zero-risk either.

## Review Rules

### Prefer Specificity Over Broad Critique

Bad:

- “The docs could be clearer.”

Good:

- “The PRD defines return eligibility broadly, but the domain model adds return-window behavior that is not grounded in the PRD. This can lead to architecture variants implementing different return rules.”

### Treat Missing Behavior as a Real Risk

An omitted failure path or state transition is often more dangerous than a minor wording issue.

### Prioritize Comparison Integrity

These artifacts exist partly to support multiple implementations of the same product. Anything that allows product behavior to drift between implementations is a high-value review target.

### Distinguish Severity

Use severity based on impact:

- `High`: likely to create inconsistent implementations, major rework, or invalid architecture comparisons
- `Medium`: likely to cause confusion, weak tests, or avoidable divergence
- `Low`: polish issue, naming weakness, or mild ambiguity

### Avoid Style Nitpicks

Do not spend findings on prose taste unless wording ambiguity materially changes behavior.

## Anti-Patterns

Do not:

- summarize the documents instead of reviewing them
- focus on formatting over semantic integrity
- report generic concerns without showing why they matter
- invent contradictions where there are only different levels of abstraction
- treat unresolved but explicit assumptions as if they were hidden defects

## Output Contract

When findings exist, output:

- findings first, ordered by severity
- open questions or assumptions next
- optional short readiness summary last

When no findings exist, output:

- explicit statement that no material findings were found
- residual risks or testing gaps
- readiness assessment

Persist the output to:

- `docs/architecture-readiness-review.md`

Behavior:

- create the file if it does not exist
- refresh it on each meaningful review pass so it reflects the latest artifact set
- keep findings ordered by severity and clearly tied to the current documents

## Resources

Read only what you need:

- [references/cross-document-checks.md](references/cross-document-checks.md): consistency checks across PRD, domain, use-case, and contract artifacts
- [references/readiness-criteria.md](references/readiness-criteria.md): what “ready for architecture implementation” means
- [references/neutrality-checks.md](references/neutrality-checks.md): architecture-neutrality review prompts
- [references/review-format.md](references/review-format.md): expected findings-first response format
