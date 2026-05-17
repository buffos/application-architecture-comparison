# Aggregate Heuristics

Use these heuristics when choosing aggregate boundaries.

## Start from Invariants

Ask:

- what must remain true after every command
- what data must change together to preserve that truth

That is the best reason to form an aggregate.

## Prefer Small, Defensible Boundaries

If an aggregate grows because “it is related,” the boundary is probably weak.

Prefer boundaries that can be justified by:

- transactional consistency
- lifecycle ownership
- policy enforcement

## Root Selection

An aggregate root should:

- have a recognizable business identity
- control access to internal consistency
- serve as the natural command target

## Warning Signs

Reconsider the boundary if:

- the aggregate contains several independent lifecycles
- every child concept wants to be loaded or updated separately
- the only reason for grouping is table convenience
- no invariant clearly requires the grouping

## Cross-Aggregate Coordination

If two concepts participate in the same workflow but do not need one transaction to preserve one invariant, they may be separate aggregates coordinated by an application or domain service.
