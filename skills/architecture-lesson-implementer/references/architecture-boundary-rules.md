# Architecture Boundary Rules

The goal of this repository is to make differences visible.

## Preserve Architectural Shape

Each solution should visibly express its target architecture.

Examples:

- Layered should make layering visible.
- Hexagonal should make ports and adapters visible.
- Clean should make dependency inversion and use-case boundaries visible.
- Transaction Script should keep orchestration procedural instead of pretending to be rich domain modeling.
- Active Record should make persistence-aware model behavior visible.

## Avoid Cross-Contamination

Do not import patterns from another architecture unless the lesson is explicitly about comparison.

Examples:

- avoid rich aggregates in a transaction-script lesson
- avoid unnecessary ports in a simple layered lesson
- avoid generic repositories everywhere just because other architectures use them

## Show the Tradeoff

Where useful, make the lesson reveal the tradeoff:

- simpler but more coupled
- more decoupled but more indirect
- richer domain semantics but more modeling overhead

## Make Boundaries Concrete In Diagrams

If a lesson includes a diagram, the diagram should reinforce the real architectural boundaries expressed by the code.

Avoid diagrams that:

- place interfaces in the wrong layer for the actual implementation
- hide whether an element is a contract or a concrete type
- collapse different adapter responsibilities into one vague box when the distinction matters

Prefer diagrams that make it easy to see:

- where each element belongs
- what depends on what
- what implements what
- which arrows represent runtime flow versus structural relationships
