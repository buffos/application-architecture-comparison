# Completion Rules

Use these rules when closing an issue.

## Issue File

- all acceptance criteria must be checked
- if extra work was discovered, add new criteria and complete them too

## Move to Done

Move the issue from:

- `issues/pending/NNN-short-title.md`

to:

- `issues/done/YYYYMMDD-NNN-short-title.md`

## Registry Update

After closing the issue:

- update its row state to `done`
- remove the issue number from blockers of dependent issues

## No Silent Closure

Do not close the issue if:

- tests fail
- build fails
- lint/static checks fail in relevant scope
- acceptance criteria remain incomplete
