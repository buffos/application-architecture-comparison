---
name: requirements-gap-analyzer
description: Inspect discovery notes, rough specs, PRDs, canonical design artifacts, or mixed requirement sources to identify what is still missing, contradictory, underspecified, or risky. Use when artifacts already exist and Codex must determine which unanswered questions matter most, prioritize the gaps by downstream impact, and generate the next best targeted questions needed to unblock a PRD, domain model, use-case model, contract, or implementation plan. Do not use for broad greenfield discovery; use `requirements-discovery-interviewer` instead.
---

# Requirements Gap Analyzer

Find what is missing, unstable, or contradictory in the current requirement set, and convert that into the smallest set of high-value follow-up questions.

This skill is not a broad interviewer. It is a precision tool for identifying what still blocks high-quality reference documents or implementation readiness.

## Workflow

### 1. Inspect the Current Artifact Set

Read whatever requirement material exists:

- discovery notes
- interview transcripts
- PRD drafts
- domain models
- use-case models
- API/CLI contracts
- ad hoc notes

Do not assume absence means “no issue.” First determine whether something is truly unspecified, implicitly specified, or contradicted elsewhere.

### 2. Classify the Gaps

For each issue, classify the gap as one of:

- missing information
- ambiguous terminology
- conflicting statements
- underspecified workflow
- missing failure or exception path
- missing boundary or non-goal
- missing policy or invariant
- unstable status vocabulary
- missing read/report requirement

Use [references/gap-categories.md](references/gap-categories.md).

### 3. Measure Downstream Impact

Ask how each gap affects:

- PRD quality
- domain modeling
- use-case derivation
- external contract stability
- architecture comparison fairness
- testability

Use [references/impact-prioritization.md](references/impact-prioritization.md).

Not every gap matters equally. Prioritize the ones that will cause divergence or rework.

### 4. Distinguish Blocking Gaps from Deferrable Gaps

A gap is blocking if it prevents:

- stable vocabulary
- stable workflow behavior
- stable rules or invariants
- stable command/query semantics
- stable external outcomes

A gap is deferrable if:

- it is refinement rather than ambiguity
- it does not materially change behavior
- it can be safely captured as an explicit assumption

### 5. Generate the Next Best Questions

Do not dump a huge questionnaire.

Produce a small set of high-leverage questions that:

- resolve the biggest blockers first
- are specific enough to answer decisively
- reduce ambiguity rather than restating it

Prefer 3-7 questions per round.

Use [references/question-generation.md](references/question-generation.md).

### 6. Expose Acceptable Assumptions

When a gap is not worth blocking progress, say so explicitly and propose a bounded assumption.

Make clear:

- what the assumption is
- what document(s) it affects
- what risk it introduces

### 7. Produce a Gap Report

The output should include:

- confirmed strong areas
- blocking gaps
- secondary gaps
- assumptions that would allow progress
- next best questions

### 8. Run the Gap Analysis Quality Gate

Review the output with [references/gap-quality-gate.md](references/gap-quality-gate.md).

If the analysis reads like:

- a generic discovery checklist
- a restatement of the source material
- a wall of low-value questions
- or a critique with no prioritization

it is wrong. Revise it.

## Analysis Rules

### Prefer Precision Over Breadth

A good gap analysis narrows uncertainty. It does not maximize question count.

### Treat Contradictions as High Value

If two artifacts describe materially different behaviors, that usually matters more than one missing detail.

### Focus on Behavioral Gaps

Prefer questions like:

- “What makes a return valid or invalid?”
- “What event prevents cancellation?”
- “Who can approve exceptions and with what outcomes?”

over shallow questions like:

- “Anything else to add?”

### Avoid Re-asking Settled Questions

If the artifact set already answers something adequately, do not ask it again just because it appears in a checklist.

### Preserve Progress

If only a few issues remain, do not reset the conversation into full discovery mode. Push the user over the finish line with targeted clarification.

## Anti-Patterns

Do not:

- ask every possible discovery question again
- ignore artifact context and ask generic prompts
- treat all gaps as equally important
- hide contradictions inside bland wording
- block progress on trivial polish issues

## Output Contract

Produce:

- prioritized gaps
- why each gap matters
- whether each gap is blocking or deferrable
- explicit assumptions where useful
- the next best follow-up questions

If the existing material is already strong enough for the intended next artifact, say so explicitly and list only residual risks.

Persist the output to:

- `docs/requirements-gap-analysis.md`

Behavior:

- create the file if it does not exist
- overwrite or refresh it when running a new gap pass so it reflects the current blocking issues
- keep the report focused on current gaps, not historical archaeology

## Resources

Read only what you need:

- [references/gap-categories.md](references/gap-categories.md): types of gaps worth identifying
- [references/impact-prioritization.md](references/impact-prioritization.md): how to rank gap severity by downstream effect
- [references/question-generation.md](references/question-generation.md): how to generate high-value follow-up questions
- [references/gap-quality-gate.md](references/gap-quality-gate.md): review checklist before finalizing the analysis
