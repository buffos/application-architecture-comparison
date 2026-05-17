---
name: process-reference-issue
description: Implement a specific local issue or auto-select the next non-blocking agent-ready issue from the local issue registry. Use when the user wants Codex to process, implement, pick up, start, or continue an issue from the artifact-grounded issue system, or says things like "next issue". This skill is designed for issues derived from the PRD and canonical reference documents, and it must escalate back to the reference-doc pipeline when implementation exposes a real spec gap.
---

# Process Reference Issue

Implement one local issue from the artifact-grounded issue system.

This skill is for execution, not issue creation. It assumes the issue already exists in the local registry and is grounded in the reference artifacts.

## Process

### 1. Load the Issues Registry

Read the issue registry from:

- `issues/issues.md`

If it does not exist there, ask the user for the correct path.

Use [references/registry-selection.md](references/registry-selection.md).

### 2. Find the Issue to Implement

#### If the user specified an issue number

- locate that issue in the registry
- if it does not exist, inform the user and stop
- if its state is not `ready-for-agent`, inform the user of the current state and stop
- if it is blocked, inform the user which issues must be resolved first and stop

Blocked means blocked. Do not attempt implementation anyway.

#### If the user did not specify a number

Select the next issue where:

- `State` is `ready-for-agent`
- `Blocked by` is `—`

Prefer the lowest-numbered eligible issue unless the user asked for a different selection rule.

If no such issue exists, inform the user and stop.

### 3. Load the Issue File

Read the issue file from:

- `issues/pending/NNN-short-title.md`

If the file is not found, stop and inform the user.

Read:

- what to build
- acceptance criteria
- blocked-by section
- parent artifacts
- artifact anchors
- acceptance scenarios addressed

### 4. Load Context from Reference Artifacts

Read the parent artifacts listed in the issue.

Typical sources:

- `docs/prd.md`
- `docs/canonical-domain-model.md`
- `docs/canonical-use-cases.md`
- `docs/canonical-api-cli-contract.md`
- `docs/acceptance-scenarios.md`
- `docs/architecture-readiness-review.md`

The issue file is not enough by itself. Use the upstream artifacts to preserve intended behavior.

### 5. Load Context from Completed Blockers

If the issue had blockers that are now resolved, read the corresponding files from:

- `issues/done/`

Use those completed issues to understand:

- patterns already established
- interfaces already chosen
- test conventions
- behavior already implemented

This is critical for consistency.

### 6. Explore the Current Codebase

Inspect the current implementation before changing code.

Especially check:

- whether part of the issue is already implemented
- whether the code diverges from the reference docs
- whether related blockers introduced conventions this issue should follow

### 7. Implement the Issue

Implement the issue end to end.

Work from the issue and the parent artifacts, not from guesses.

As acceptance criteria are satisfied, update the issue file and mark them complete:

- `[ ]` -> `[x]`

### 8. Handle Spec Gaps Correctly

If implementation reveals a real gap in the reference-doc set:

- do not silently invent product behavior
- do not "just choose something reasonable" if that choice changes business semantics

Instead:

1. determine whether the gap is trivial and safely assumable or genuinely blocking
2. if it is genuinely blocking, stop implementation
3. update or create:
   - `docs/requirements-gap-analysis.md`
4. tell the user the issue is blocked by a spec gap
5. route back to the right upstream skill, typically:
   - `requirements-gap-analyzer`
   - `requirements-discovery-interviewer`
   - `prd-author`
   - `canonical-domain-modeler`
   - `canonical-use-case-modeler`
   - `canonical-contract-author`

Use [references/spec-gap-handling.md](references/spec-gap-handling.md).

### 9. Verify Completion

When all acceptance criteria appear met:

- run relevant tests
- run build checks
- run lint or static checks if the repo uses them

If failures occur:

- fix obvious errors first
- continue until the issue is genuinely green

Do not close the issue in a broken state.

If additional acceptance work is discovered during execution, add it to the issue as new unchecked criteria and complete it before closing.

### 10. Finalize the Issue

After all acceptance criteria are checked and verification is clean:

- move the issue file to `issues/done/`
- prefix it with `YYYYMMDD-`

Example:

- `issues/done/20260517-001-short-title.md`

Then update `issues/issues.md`:

- do not remove the historical record
- update the row state to `done`
- remove the resolved issue number from `Blocked by` fields of remaining issues

Use [references/completion-rules.md](references/completion-rules.md).

### 11. Report Back to the User

Provide a concise summary of:

- what was implemented
- what was verified
- whether any upstream artifact drift was found

Also provide a conventional commit message:

- feature: `feat(scope): description`
- fix: `fix(scope): description`

## Anti-Patterns

Do not:

- implement blocked issues
- ignore parent artifacts
- silently invent missing business behavior
- mark an issue done while tests or builds are failing
- mutate the PRD just because code was changed
- close an issue if acceptance criteria are still incomplete

## Output Contract

On a successful run:

- implemented code changes
- updated acceptance checkboxes in the issue file
- moved issue file into `issues/done/`
- updated `issues/issues.md`
- concise completion summary
- conventional commit message

On a blocked run:

- no false completion
- clear explanation of the blocker
- updated `docs/requirements-gap-analysis.md` if the blocker is a true spec gap
- recommendation for which upstream skill should run next

## Resources

Read only what you need:

- [references/registry-selection.md](references/registry-selection.md): how to choose and validate the target issue
- [references/spec-gap-handling.md](references/spec-gap-handling.md): when to escalate back to the artifact pipeline
- [references/completion-rules.md](references/completion-rules.md): how to close out the issue and update the registry safely
- [references/verification-expectations.md](references/verification-expectations.md): expected verification discipline before closing
