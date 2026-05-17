# Impact Prioritization

Rank gaps by the damage they can cause downstream.

## High Impact

Usually includes gaps that can cause:

- divergent implementations
- contradictory domain models
- unstable use-case naming
- incompatible contracts
- broken acceptance tests

## Medium Impact

Usually includes gaps that can cause:

- avoidable confusion
- extra review cycles
- minor implementation drift
- weak but still recoverable tests

## Low Impact

Usually includes gaps that affect:

- polish
- wording clarity with little behavioral effect
- non-blocking future refinements

## Prioritization Questions

- Will this gap change system behavior?
- Will it change who can do what and when?
- Will it change state transitions or failure semantics?
- Will different architects likely make different choices if it stays unresolved?
- Can we safely proceed with a bounded assumption?
