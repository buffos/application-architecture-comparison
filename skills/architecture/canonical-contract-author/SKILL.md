---
name: canonical-contract-author
description: Derive a canonical external contract from a stable PRD, canonical domain model, and canonical use-case model. Use when Codex must define the stable API and CLI behavior of a system, including transport-neutral conventions, identifiers, payload shapes, statuses, error codes, idempotency rules, and canonical HTTP and CLI mappings without leaking internal implementation details or tying the system too tightly to a specific framework. Do not use before commands, queries, and visible outcomes are already stable.
---

# Canonical Contract Author

Produce the stable external contract that all implementations should preserve from the outside.

This skill defines external behavior, not internal architecture. It should sit on top of the canonical use-case model and make it consumable through HTTP and CLI without changing business semantics.

## Workflow

### 1. Confirm the Inputs Are Contract-Ready

Before drafting, verify the available material is strong enough to answer:

- what commands and queries exist
- what identifiers and resources matter externally
- what statuses must be exposed
- what error conditions users or clients must see
- which operations are retry-sensitive

Use [references/contract-readiness.md](references/contract-readiness.md).

If the source material is not ready, ask follow-up questions instead of fabricating a contract.

### 2. Start with Transport-Neutral Semantics

Define first:

- identifiers
- status vocabulary
- money and quantity shapes
- timestamps
- request/response conventions
- error model
- idempotency rules

Do this before choosing URL shapes or CLI verbs.

### 3. Preserve Use-Case Semantics

Map commands and queries from the canonical use-case model directly.

The contract should preserve:

- business names
- business outcomes
- business-visible errors
- stable request inputs

Do not invent endpoints or CLI commands that bypass the application model.

Use [references/contract-mapping.md](references/contract-mapping.md).

### 4. Define Resource and Payload Shapes

For each important concept exposed externally, define:

- snapshot shape
- key identifiers
- relevant statuses
- key nested structures

Focus on stability and clarity rather than exhaustive data dumping.

### 5. Define Error and Retry Semantics

Document:

- error categories
- canonical business error codes
- HTTP status mapping guidance
- CLI exit-code guidance
- idempotency behavior for sensitive commands

If retries are dangerous, say so explicitly. If retries should be safe, define what “safe” means.

### 6. Derive Canonical HTTP Mapping

After the transport-neutral contract is stable, map it to HTTP:

- base path
- route shape
- method choice
- path parameters
- query parameters
- example requests and responses

The HTTP design should be clean, but business consistency matters more than stylistic REST purity.

### 7. Derive Canonical CLI Mapping

Map the same command/query surface into CLI commands:

- stable verbs
- actor flags
- idempotency flags
- output mode guidance
- exit-code behavior

The CLI is not an afterthought. It should reflect the same use-case model as the HTTP contract.

### 8. Define Parity Rules

Document what all future implementations must preserve:

- business naming
- request meaning
- output semantics
- status vocabulary
- error codes
- retry behavior

Also document which presentation details may vary.

### 9. Run the Contract Quality Gate

Review the output with [references/contract-quality-gate.md](references/contract-quality-gate.md).

If the contract reads like:

- a framework router dump
- a random REST resource inventory
- a CLI help page without business grounding
- or an implementation detail leak

it is wrong. Revise it.

## Contract Rules

### Start from Meaning, Then Transport

External contracts should first preserve semantics, then choose representation.

Do not let route style dictate the business contract.

### Keep Status and Error Vocabularies Stable

Statuses and business-visible error codes are part of the canonical surface. Treat them as comparison-critical.

### Prefer Explicit Shapes for Ambiguous Data

For money, timestamps, and retry-sensitive commands, be explicit:

- money object instead of raw float
- RFC3339 timestamps
- explicit idempotency keys where appropriate

### Keep HTTP and CLI Semantically Aligned

If HTTP exposes a command, CLI should expose the same business operation where practical, and vice versa.

### Preserve External Neutrality

The contract should not expose:

- package names
- repository names
- framework-specific DTO quirks
- persistence-only fields without business meaning

## Anti-Patterns

Do not:

- derive the contract directly from database tables
- turn every noun into a CRUD endpoint if the use case is intent-driven
- hide important business errors behind generic “400”
- change names between HTTP and CLI without reason
- leak internal implementation terms into public contracts
- confuse a canonical contract with an OpenAPI-only artifact

## Output Contract

Produce a canonical API/CLI contract document that includes:

- purpose
- contract goals
- transport-neutral conventions
- identifier conventions
- status values
- payload shapes
- error model
- idempotency expectations
- canonical HTTP API mapping
- canonical CLI mapping
- parity rules
- minimum first-slice surface
- suggested testing contract

The result should be strong enough that multiple implementations can expose the same external behavior while differing internally.

Persist the output to:

- `docs/canonical-api-cli-contract.md`

Behavior:

- create the file if it does not exist
- update it in place when external behavior changes
- keep statuses, identifiers, and error codes aligned with upstream artifacts

## Resources

Read only what you need:

- [references/contract-readiness.md](references/contract-readiness.md): when the inputs are ready for contract definition
- [references/contract-mapping.md](references/contract-mapping.md): how to map use cases into HTTP and CLI without drifting semantically
- [references/contract-conventions.md](references/contract-conventions.md): recommended global conventions for IDs, statuses, money, errors, and idempotency
- [references/contract-quality-gate.md](references/contract-quality-gate.md): review checklist before finalizing
