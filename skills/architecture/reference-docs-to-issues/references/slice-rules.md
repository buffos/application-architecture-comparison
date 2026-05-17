# Slice Rules

Use tracer-bullet thinking.

## Good Slice Properties

- narrow end-to-end behavior
- independently verifiable
- minimal but complete
- clear dependency story

## Vertical Slice Rule

A slice should cut through all required layers for one behavior path, such as:

- persistence changes
- domain behavior
- application orchestration
- contract exposure
- tests

Not every slice needs every layer explicitly named, but the behavior should be complete.

## Dependency Rule

Prefer:

- foundational blockers first
- feature slices after their blockers
- acceptance coverage that can be demonstrated incrementally

## HITL vs AFK

Mark as `HITL` if the slice requires:

- product-owner clarification
- human architectural decision
- human design review
- non-automatable business approval

Mark as `AFK` if an agent can implement it directly once the issue is written clearly enough.
