# Spec Gap Handling

Use this rule when code work exposes a real artifact problem.

## What Counts as a Real Spec Gap

Examples:

- the issue depends on a business rule not defined anywhere
- the parent artifacts contradict one another
- a required failure outcome is unspecified
- the contract implies behavior the use-case model does not define
- a new edge case would materially change business semantics

## What Does Not Automatically Count

Examples:

- trivial naming cleanup
- minor implementation detail choices
- obvious local refactors with no product-semantics impact

## Escalation Rule

If the gap changes business behavior:

- stop pretending the issue is implementable
- record the blocker in `docs/requirements-gap-analysis.md`
- tell the user which upstream artifact is insufficient
- recommend the correct upstream skill

## Assumption Rule

Only proceed with a local assumption if:

- it is clearly bounded
- it does not materially alter the product behavior
- it is consistent with the reference set

If in doubt, escalate.
