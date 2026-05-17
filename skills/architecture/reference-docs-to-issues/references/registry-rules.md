# Registry Rules

Use this format for `issues/issues.md`.

```markdown
# Issues Registry

| # | Title | Category | State | Blocked by |
|---|---|---|---|---|
| 001 | Short issue title | feature | ready-for-agent | — |
| 002 | Another issue | feature | ready-for-human | 001 |
```

## Column Rules

- **#**: zero-padded issue number
- **Title**: short descriptive title matching the file
- **Category**: usually `feature` or `bug`
- **State**:
  - `needs-triage`
  - `needs-info`
  - `ready-for-agent`
  - `ready-for-human`
  - `done`
  - `wontfix`
- **Blocked by**: comma-separated issue numbers or `—`

## State Assignment

- fully specified `AFK` issues -> `ready-for-agent`
- `HITL` issues -> `ready-for-human`
- issues with unresolved details -> `needs-info`
- externally reported unclear work -> `needs-triage`
- completed issues -> `done`

## Max Issue ID Section

Always maintain:

```markdown
# Current Max Issue ID

NNN
```
