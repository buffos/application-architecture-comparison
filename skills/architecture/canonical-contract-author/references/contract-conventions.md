# Contract Conventions

Use these conventions unless the user explicitly needs a different contract style.

## Identifiers

- use opaque string identifiers externally
- keep business-stable keys like `sku` explicit where relevant

## Money

- serialize money as an object with `amount` and `currency`
- avoid float-only shapes

## Timestamps

- use RFC3339 UTC strings

## Statuses

- expose stable enum-like strings
- keep them aligned with the canonical domain model

## Errors

- define machine-readable error codes
- distinguish validation, not found, conflict, business rule, and infrastructure failures

## Idempotency

- expose idempotency support for commands where replay can create duplicates or repeat dangerous work

## CLI Output

- support a machine-readable mode such as JSON
- keep CLI result shapes aligned with HTTP `data` payloads where practical
