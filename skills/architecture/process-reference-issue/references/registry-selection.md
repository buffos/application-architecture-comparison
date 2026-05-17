# Registry Selection

Use these rules to choose the target issue.

## Explicit Issue Number

If the user names an issue number:

- find that exact issue
- verify it exists
- verify its state is `ready-for-agent`
- verify it is not blocked

If any of those checks fail, stop.

## Implicit Selection

If the user says:

- "next issue"
- "pick one"
- "continue"

select the next issue where:

- `State` is `ready-for-agent`
- `Blocked by` is `—`

Prefer the lowest-numbered eligible issue unless the user indicates a different ordering rule.

## Registry Truth

The registry is the selection source of truth.

The issue file is the implementation source of truth once selection succeeds.
