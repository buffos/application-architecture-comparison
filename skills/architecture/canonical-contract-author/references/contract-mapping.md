# Contract Mapping

Use these rules when mapping the canonical use-case model outward.

## Transport-Neutral First

Define:

- commands
- queries
- identifiers
- snapshots
- errors
- idempotency

before choosing specific HTTP routes or CLI command syntax.

## HTTP Mapping

When mapping to HTTP:

- use route shapes that preserve business intent
- use resource-oriented paths where natural
- use action endpoints when intent is more important than generic CRUD

Examples:

- `/quotes/{quoteId}/submit`
- `/quotes/{quoteId}/approve`
- `/quotes/{quoteId}/convert-to-order`

These are acceptable when they preserve use-case semantics clearly.

## CLI Mapping

When mapping to CLI:

- keep command groups aligned with business areas
- keep verbs aligned with use-case names
- expose actor and idempotency flags where relevant

## Preserve Naming Consistency

The same business concept should not be renamed arbitrarily across:

- PRD
- domain model
- use-case model
- HTTP contract
- CLI contract

Consistency is more important than surface cleverness.
