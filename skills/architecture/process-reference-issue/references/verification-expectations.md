# Verification Expectations

Before closing an issue, verify the implementation in a way that matches the repo and issue scope.

## Minimum Expectations

- relevant tests pass
- relevant build succeeds
- obvious lint or static failures in touched scope are resolved

## If Verification Fails

- fix obvious breakage first
- keep working until the issue is genuinely green
- add missing acceptance criteria if new required work is discovered

## Report Clearly

When reporting back, state:

- what was tested
- what passed
- what could not be run, if anything
