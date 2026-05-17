# Use-Case Structure

Use this structure when deriving the canonical application layer.

## Application-Layer Goals

State what the application layer is supposed to provide:

- business-capable operations
- orchestration boundaries
- stable command/query semantics

## Design Principles

Make explicit principles such as:

- intent-based naming
- command/query separation
- architecture neutrality
- domain logic staying in the right place

## Application Boundaries

Group use cases into coherent services or modules such as:

- catalog
- customer
- quote
- order
- payment
- fulfillment
- returns
- reporting
- plugins

## Use-Case Definition

For each use case, capture:

- intent
- inputs
- outputs
- preconditions
- application responsibilities
- possible outcomes

## Command and Query Shapes

Describe the contract shape conceptually so later transport contracts can map to it consistently.
