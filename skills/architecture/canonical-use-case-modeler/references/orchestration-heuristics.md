# Orchestration Heuristics

Use these heuristics to decide what belongs in the canonical application layer.

## A Good Application Use Case

Usually:

- starts with an intent
- loads the required aggregates
- coordinates domain behavior
- persists results
- returns a business-meaningful outcome

## Single-Aggregate Commands

These usually have:

- one clear aggregate root
- one local consistency boundary
- minimal orchestration

## Coordinated Commands

These usually:

- involve multiple aggregates
- need explicit transaction expectations
- may need retry or compensation thinking

Examples include:

- converting a quote to an order while reserving stock
- shipping while consuming reservations
- accepting returns while restoring inventory

## Idempotency Signals

Consider idempotency when:

- duplicates would be harmful
- external callers may retry
- the command creates new records as part of processing

## Event Awareness

Application services may publish or hand off events when:

- read models must update
- auditability matters
- plugins or policy subsystems may react

The use-case model should preserve those moments even if implementations differ technically.
