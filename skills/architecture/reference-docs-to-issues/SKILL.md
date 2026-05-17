---
name: reference-docs-to-issues
description: Break a PRD and the canonical reference artifact set into independently workable local issues. Use when the product specs are stable enough that Codex should translate them into dependency-aware vertical slices, propose the breakdown for approval, and then write local markdown issue files plus an issue registry. Prefer this over PRD-only issue generation when any of the canonical artifacts exist because it can ground issue boundaries in domain rules, use cases, contracts, acceptance scenarios, and readiness findings.
---

# Reference Docs to Issues

Break the reference artifact set into independently grabbable local issues using vertical slices.

This skill is the implementation bridge between the reference documents and actual work items. It should produce thin, demoable, dependency-aware issues that an agent can process one at a time.

## Source of Truth

Use whichever of these artifacts exist, in this priority:

1. `docs/prd.md`
2. `docs/canonical-domain-model.md`
3. `docs/canonical-use-cases.md`
4. `docs/canonical-api-cli-contract.md`
5. `docs/acceptance-scenarios.md`
6. `docs/architecture-readiness-review.md`
7. `docs/domain-glossary.md`
8. `docs/discovery-notes.md`
9. `docs/requirements-gap-analysis.md`
10. `docs/reference-doc-orchestration-status.md`

If only some artifacts exist, use the strongest available set. If the artifacts are too weak or contradictory to support issue breakdown, stop and tell the user to route back through the missing authoring or review skill first.

Use [references/artifact-priority.md](references/artifact-priority.md).

## Process

### 1. Inventory the Artifact Set

Check which source artifacts exist and treat them as the current reference set.

If the user explicitly names a different source file or folder, use that instead.

Before creating issues, determine:

- which artifacts are present
- which are authoritative
- whether the set is implementation-ready enough

If `docs/architecture-readiness-review.md` contains unresolved high-severity findings, do not blindly generate agent-ready issues from a broken spec set. Surface that problem first.

### 2. Explore the Codebase When Relevant

If the user wants issues for an existing codebase, inspect the current implementation so the issue breakdown reflects what already exists versus what still needs to be built.

This is especially important for:

- brownfield systems
- partially implemented products
- gaps between spec and code

### 3. Derive Vertical Slices

Break the reference set into **tracer bullet** issues.

Each issue should be a thin vertical slice that cuts through all relevant layers needed to make the behavior real and verifiable.

Use [references/slice-rules.md](references/slice-rules.md).

Rules:

- each slice must deliver a narrow but complete path
- a completed slice must be demoable or verifiable on its own
- prefer many thin slices over few thick ones
- do not create horizontal layer-only issues unless the work is genuinely cross-cutting and cannot be sliced vertically

Possible slice types:

- `AFK`: can be implemented and merged without human interaction
- `HITL`: requires human input, decision, review, or approval

Prefer `AFK` over `HITL` where possible.

### 4. Ground the Slices in the Artifact Set

Each proposed issue should be anchored to whichever of these exist:

- PRD sections
- domain concepts or invariants
- use cases or application services
- contract endpoints or CLI commands
- acceptance scenarios
- readiness review findings, if the issue exists to resolve one

Do not create issues that float free from the reference documents.

### 5. Present the Breakdown for Approval

Before writing files, present the proposed breakdown as a numbered list.

For each slice, show:

- **Title**
- **Type**: `AFK` or `HITL`
- **Blocked by**
- **Artifacts covered**
- **Acceptance scenarios covered** when available

Ask the user:

- Is the granularity right?
- Are the dependency relationships correct?
- Should any slices be split or merged?
- Are the right slices marked `AFK` and `HITL`?

Iterate until the user approves the breakdown.

### 6. Create the Issue Files

After approval, write one markdown file per issue in:

- `issues/pending/`

Use the filename pattern:

- `issues/pending/NNN-short-title.md`

Create files in dependency order so blockers are written first.

Read the current max issue ID from:

- `issues/issues.md`

Look for the section:

```markdown
# Current Max Issue ID

NNN
```

If that section does not exist, fall back to scanning:

- `issues/pending/`
- `issues/done/`

for the highest existing issue number.

Do not use GitHub issues or `gh issue create`. These are local issue files only.

### 7. Use the Issue Template

Use the template in [references/issue-template.md](references/issue-template.md).

Each issue file must include:

- parent artifacts
- what to build
- acceptance criteria
- blocked by
- artifact anchors
- acceptance scenarios addressed

If some artifact types do not exist, omit only the irrelevant references.

### 8. Create or Update the Issues Registry

Create or update:

- `issues/issues.md`

Append new issues to the existing table. Do not overwrite existing rows.

Registry convention:

- keep active and completed issues in the same registry
- completed issues remain in the table with state `done`
- issue file history lives in `issues/done/`, but the registry remains the index of all issued work

Use the registry rules in [references/registry-rules.md](references/registry-rules.md).

### 9. Update the Max Issue ID

After appending the new rows, update:

```markdown
# Current Max Issue ID

NNN
```

so future runs continue numbering correctly.

## Brownfield Guidance

For brownfield projects, the issue set should reflect what is already implemented.

Do not generate issues as if the codebase were empty.

Instead:

- compare the artifact set to the existing code
- generate only the missing or corrective slices
- create explicit spec-alignment issues when implementation and reference docs differ

If the artifact set is weak, route back to:

- `requirements-gap-analyzer`
- `architecture-readiness-reviewer`

before generating implementation issues.

## Anti-Patterns

Do not:

- generate horizontal layer-only issues by default
- generate one giant “implement the whole feature” issue
- generate issues from the PRD only when richer artifacts exist
- ignore acceptance scenarios
- ignore readiness-review findings
- create issue files before the user approves the breakdown

## Output Contract

Produce:

- an approved issue breakdown
- local issue markdown files in `issues/pending/`
- an updated `issues/issues.md` registry
- an updated max issue ID

If the artifact set is not strong enough, stop and report which upstream artifact must be fixed first.

## Resources

Read only what you need:

- [references/artifact-priority.md](references/artifact-priority.md): which artifacts to trust and how strongly
- [references/slice-rules.md](references/slice-rules.md): vertical-slice rules and dependency guidance
- [references/issue-template.md](references/issue-template.md): local issue file template
- [references/registry-rules.md](references/registry-rules.md): issue registry format and state rules
