# Conflict Patterns

Look for these terminology problems.

## One Concept, Many Names

Examples:

- quote / proposal / estimate
- order / purchase / booking

Choose one canonical term and list the others as aliases or discouraged terms.

## One Name, Many Meanings

Examples:

- approval meaning business approval in one place and payment authorization in another
- return meaning both request and inventory restock

Split the meanings explicitly.

## UI Language Replacing Business Language

Examples:

- “manage order screen”
- “inventory page”

Translate the UI phrasing into business concepts.

## Technical Language Replacing Business Language

Examples:

- entity
- record
- transaction

These may be useful later, but they are not automatically business terms.

## Policy vs Invariant Confusion

Examples:

- “cannot return clearance items” may be policy
- “returned quantity cannot exceed shipped quantity” is closer to invariant

Do not flatten both into one generic “rule” unless the distinction truly does not matter.
