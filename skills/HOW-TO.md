# Skills How-To

## Purpose

This document explains how to use the skills in [skills](</c:/Users/buffo/Code/architecture/01.application.architectures/skills>) to build and maintain a coherent reference-document set for an application.

The skills are useful for both:

- greenfield projects
- brownfield projects

They are not limited to new projects. In fact, they are often more valuable in brownfield work because existing systems usually contain terminology drift, undocumented rules, contradictory behavior, and hidden assumptions.

## Standard Output Files

With the tightened skill definitions, the pipeline now has deterministic artifact targets under `docs/`:

- `discovery-notes.md`
- `domain-glossary.md`
- `requirements-gap-analysis.md`
- `prd.md`
- `canonical-domain-model.md`
- `canonical-use-cases.md`
- `canonical-api-cli-contract.md`
- `acceptance-scenarios.md`
- `architecture-readiness-review.md`
- `reference-doc-orchestration-status.md`

## The Skills

- `requirements-discovery-interviewer`
- `domain-glossary-extractor`
- `requirements-gap-analyzer`
- `prd-author`
- `canonical-domain-modeler`
- `canonical-use-case-modeler`
- `canonical-contract-author`
- `acceptance-scenario-author`
- `architecture-readiness-reviewer`
- `reference-doc-orchestrator`
- `reference-docs-to-issues`
- `process-reference-issue`

## Core Principle

These skills are not a one-time ceremony.

They form a pipeline for:

1. discovering the product
2. documenting the product
3. stabilizing the design reference
4. reviewing readiness before implementation
5. revisiting the reference set when the product changes
6. converting the reference set into actionable local issues
7. implementing those issues one at a time

That last point matters. New requirements, production discoveries, edge cases, and changed business rules are normal. When that happens, you re-enter the pipeline at the right point instead of starting over.

## Recommended Default

If you are unsure which skill to use first, use:

- `reference-doc-orchestrator`

That skill decides which specialist skill should run next.

If the reference docs are already stable and you want to start delivery work:

- use `reference-docs-to-issues` to create the issue queue
- use `process-reference-issue` to execute one issue at a time

## Greenfield

### When to Use the Full Pipeline

Use the full pipeline when:

- the app idea is still vague
- there is no real PRD
- workflows are not yet stable
- terminology is still loose
- architecture work has not started yet

### Recommended Order

1. `requirements-discovery-interviewer`
   Use this first when the idea is fuzzy. Its job is to grill the user until the product problem, actors, workflows, rules, and boundaries are explicit enough to draft serious artifacts.

2. `domain-glossary-extractor`
   Use this early if the same concept is being named in multiple ways or if actor, object, and workflow terms are still muddled.

3. `requirements-gap-analyzer`
   Use this after early discovery to identify what still blocks a strong PRD or stable design artifacts.

4. `prd-author`
   Use this when discovery is strong enough to write the product requirements document.

5. `canonical-domain-modeler`
   Use this when the PRD is stable enough to derive the core business model.

6. `canonical-use-case-modeler`
   Use this when the domain model is stable enough to derive commands, queries, application services, orchestration, and failure behavior.

7. `canonical-contract-author`
   Use this when the use-case model is stable enough to define external HTTP and CLI behavior.

8. `acceptance-scenario-author`
   Use this when the system behavior is clear enough to define reusable acceptance scenarios.

9. `architecture-readiness-reviewer`
   Use this last before implementation starts. It checks contradictions, missing decisions, architecture bias, and testability gaps across the artifact set.

10. `reference-docs-to-issues`
    Use this after the reference set is stable enough to break it into dependency-aware implementation slices.

11. `process-reference-issue`
    Use this after issue generation to implement one approved issue at a time.

### Greenfield Shortcut

If the product idea is already unusually well specified, you may skip parts of early discovery:

- if terminology is already stable, skip directly past `domain-glossary-extractor`
- if the PRD already exists and is solid, start at `canonical-domain-modeler`
- if the PRD and domain model already exist, start at `canonical-use-case-modeler`

## Brownfield

### Can These Skills Be Used in Brownfield?

Yes.

They are absolutely usable in brownfield work.

The key difference is that you do not treat the existing system as automatically correct or complete. In brownfield work, the codebase, API, docs, and team knowledge are all sources of truth, but they may disagree.

The goal is usually:

- reconstruct missing reference docs
- formalize undocumented behavior
- identify inconsistencies between implementation and intended behavior
- extend the system without breaking conceptual integrity

### Typical Brownfield Starting Points

#### Case 1: There Is Existing Code But Weak or Missing Docs

Recommended order:

1. `requirements-gap-analyzer`
2. `domain-glossary-extractor`
3. `requirements-discovery-interviewer`
4. `prd-author`
5. `canonical-domain-modeler`
6. `canonical-use-case-modeler`
7. `canonical-contract-author`
8. `acceptance-scenario-author`
9. `architecture-readiness-reviewer`

Why:

- first identify what is known and what is missing
- then stabilize terminology
- then interrogate the product owner or team on contradictions and hidden rules
- then write the missing reference docs

After the reference set is stable:

10. `reference-docs-to-issues`
11. `process-reference-issue`

#### Case 2: There Is a PRD, But It Is Weak or Outdated

Recommended order:

1. `requirements-gap-analyzer`
2. `requirements-discovery-interviewer`
3. `prd-author`
4. then continue with the canonical documents

Why:

- the gap analyzer identifies what the current PRD fails to define
- the interviewer closes the gaps
- the PRD author rewrites or strengthens the document instead of patching it blindly

#### Case 3: PRD Exists, But the Domain Language Is a Mess

Recommended order:

1. `domain-glossary-extractor`
2. `requirements-gap-analyzer`
3. `canonical-domain-modeler`

Why:

- vocabulary instability will infect all later documents if not fixed first

#### Case 4: Core Docs Exist, but New Features Must Be Added

Recommended order:

1. `requirements-gap-analyzer`
2. `requirements-discovery-interviewer` if the new feature is still vague
3. `domain-glossary-extractor` if the feature introduces new language
4. update `prd-author`
5. update `canonical-domain-modeler`
6. update `canonical-use-case-modeler`
7. update `canonical-contract-author` if external behavior changes
8. update `acceptance-scenario-author`
9. rerun `architecture-readiness-reviewer`

Why:

- new features usually change the reference set, not just the code

After the updated docs are stable:

10. `reference-docs-to-issues`
11. `process-reference-issue`

#### Case 5: Implementation Exists, and You Need to Compare It Against Intended Design

Recommended order:

1. `architecture-readiness-reviewer`
2. `requirements-gap-analyzer`
3. reroute to the specific missing-authoring skill

Why:

- the reviewer tells you where the artifact set is weak or contradictory
- the gap analyzer turns that into focused follow-up work

## How to Re-Enter the Pipeline

Sooner or later, the programmer will discover:

- a missing rule
- a missing exception path
- a forgotten actor
- a new approval policy
- a new integration constraint
- a new report requirement
- a feature that changes existing business behavior

When that happens, do not just patch the code and move on.

Re-enter the pipeline at the highest affected level.

### Re-Entry Rule

If the change affects:

- wording only: update the relevant document directly
- terminology: start with `domain-glossary-extractor`
- product scope or goals: start with `prd-author` after gap analysis/interview as needed
- business objects, invariants, or lifecycle rules: start with `canonical-domain-modeler`
- commands, queries, orchestration, or failure semantics: start with `canonical-use-case-modeler`
- external API/CLI behavior: start with `canonical-contract-author`
- expected behavior checks: start with `acceptance-scenario-author`
- overall coherence: finish with `architecture-readiness-reviewer`

### Practical Example

Suppose a brownfield app gains partial refunds after return approval.

Likely flow:

1. `requirements-gap-analyzer`
   Determine what is missing: refund rules, actor permissions, thresholds, states.

2. `requirements-discovery-interviewer`
   Ask the product owner the missing questions if the feature is still vague.

3. `prd-author`
   Update product requirements and workflows.

4. `canonical-domain-modeler`
   Add or revise `Refund`, return lifecycle, and related invariants.

5. `canonical-use-case-modeler`
   Add or revise refund-related commands and failure modes.

6. `canonical-contract-author`
   Update API and CLI behavior if refund operations are externally visible.

7. `acceptance-scenario-author`
   Add scenarios for accepted partial refund, invalid refund attempt, and visibility outcomes.

8. `architecture-readiness-reviewer`
   Confirm the updated reference set is still coherent.

9. `reference-docs-to-issues`
   Generate the new implementation slices required by the change.

10. `process-reference-issue`
    Execute those slices one at a time.

## Which Skill to Use When

Use `requirements-discovery-interviewer` when:

- the product or feature is still fuzzy
- you need broad discovery
- there are too many unknowns for drafting

Use `requirements-gap-analyzer` when:

- artifacts already exist
- the problem is not “start from zero”
- you need the next best clarifying questions

Use `domain-glossary-extractor` when:

- vocabulary is drifting
- synonyms or overloaded terms are causing confusion

Use `prd-author` when:

- discovery is strong enough to write or revise the PRD

Use `canonical-domain-modeler` when:

- business semantics need to be stabilized

Use `canonical-use-case-modeler` when:

- application-layer behavior needs to be stabilized

Use `canonical-contract-author` when:

- external API/CLI behavior needs to be stabilized

Use `acceptance-scenario-author` when:

- you need stable behavior scenarios for testing and comparison

Use `architecture-readiness-reviewer` when:

- the document set appears complete
- or before implementation of a new architecture or major feature begins

Use `reference-doc-orchestrator` when:

- you are not sure which step comes next
- multiple artifacts are missing or weak
- you want the whole workflow coordinated

Use `reference-docs-to-issues` when:

- the reference docs are stable enough to be converted into implementation work
- you want local issue files in `issues/pending/`
- you want a dependency-aware issue registry in `issues/issues.md`

Use `process-reference-issue` when:

- you want to implement a specific local issue
- you want the next available `ready-for-agent` issue
- you want execution to stay grounded in the reference docs
- you want implementation to escalate back to the artifact pipeline if a true spec gap is discovered

## Minimum Brownfield Discipline

In brownfield work, the minimum safe loop is usually:

1. `requirements-gap-analyzer`
2. the specific authoring skill for the affected artifact
3. `acceptance-scenario-author` if behavior changed
4. `architecture-readiness-reviewer`
5. `reference-docs-to-issues` for newly actionable work
6. `process-reference-issue` to implement it

This is the smallest loop that still preserves coherence.

## Recommended Habit

For both greenfield and brownfield:

- do not let reference docs lag far behind behavior
- do not patch only the lowest-level artifact
- when behavior changes, walk upward to the right reference layer
- always end meaningful changes with readiness review

## Suggested Usage Patterns

### Pattern A: New Product From Scratch

Use:

1. `reference-doc-orchestrator`
2. follow the routed sequence
3. end with `architecture-readiness-reviewer`
4. use `reference-docs-to-issues`
5. then use `process-reference-issue`

### Pattern B: Existing Product, Missing Docs

Use:

1. `reference-doc-orchestrator`
2. or manually start with `requirements-gap-analyzer`
3. then reconstruct the missing artifacts
4. then use `reference-docs-to-issues`
5. then use `process-reference-issue`

### Pattern C: Existing Product, New Major Feature

Use:

1. `requirements-gap-analyzer`
2. `requirements-discovery-interviewer` if needed
3. update the affected authoring skills’ artifacts
4. `acceptance-scenario-author`
5. `architecture-readiness-reviewer`
6. `reference-docs-to-issues`
7. `process-reference-issue`

### Pattern D: Pre-Implementation Gate

Use:

1. `architecture-readiness-reviewer`

If it finds issues, reroute to the correct authoring skill and review again.

### Pattern E: Stable Specs to Execution

Use:

1. `reference-docs-to-issues`
2. approve the issue breakdown
3. `process-reference-issue`
4. repeat until the `ready-for-agent` queue is empty
5. if implementation exposes a true spec gap, route backward through the appropriate authoring skill

## Final Guidance

These skills are not only for creating the first version of the docs.

They are for maintaining a stable, explicit reference model for the system as it evolves.

Greenfield uses them to create clarity.

Brownfield uses them to recover clarity and preserve it while the software changes.
