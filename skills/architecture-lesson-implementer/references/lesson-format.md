# Lesson Format

Use a short, repeatable lesson shape.

## Recommended Sections

- `# Lesson NNN: Title`
- `## Objective`
- `## Theory`
- `## Why This Matters Here`
- `## Diagram`
- `## Implementation Focus`
- `## What To Verify`

## Theory Guidance

Keep theory brief.

Answer:

- what is this concept
- why are we using it here
- what problem does it solve

## Diagram Guidance

Use Mermaid when it improves understanding.

Good uses:

- dependency direction
- module boundaries
- request flow
- adapter relationships

Prefer diagrams that match the actual code boundaries of the lesson, not only a generic version of the architecture.

When useful, make these distinctions visible:

- layer or boundary ownership
- contracts/interfaces versus concrete types
- runtime flow versus implementation relationships
- different adapter responsibilities such as persistence/data adapters versus translation/behavioral adapters

Useful conventions:

- subgraphs for layers
- dashed borders for contracts/interfaces
- dashed arrows for implementation or dependency-to-contract links
- solid arrows for request/runtime flow
- color with a short legend when it adds clarity

## Implementation Focus

State exactly what this lesson will implement and what it will deliberately leave for later.
