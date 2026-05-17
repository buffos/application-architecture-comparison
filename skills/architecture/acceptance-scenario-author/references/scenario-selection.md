# Scenario Selection

Choose scenarios that carry the most behavioral value.

## High-Value Scenario Types

- primary happy path
- approval or review branch
- policy rejection or failure path
- stock or resource shortage path
- cancellation or reversal path
- return or refund path
- extensibility or plugin variation path
- operational read/report path

## Selection Questions

- Does this scenario prove a rule that matters?
- Would two architects be likely to implement this differently if it were underspecified?
- Does the scenario exercise a meaningful branch?
- Does the scenario expose a business-visible failure?

## Avoid

- trivial CRUD-only scenarios unless the CRUD action is itself business-significant
- duplicates that differ only cosmetically
