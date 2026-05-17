# Scenario Structure

Use a stable, implementation-neutral structure.

## Recommended Form

- Title
- Purpose
- Given
- When
- Then
- Optional rule/coverage note

## Writing Guidance

### Given

Capture the minimum meaningful precondition.

### When

State the business action, not the transport mechanism.

### Then

State the externally visible outcome in plain business language.

## Good Example Shape

- Given an approved quote for an in-stock product
- When the quote is converted to an order
- Then stock is reserved and the order becomes ready for payment

## Bad Example Shape

- Given I call POST /api/v1/quotes/x/convert-to-order
- Then the repository contains a row

That is contract or implementation detail, not canonical acceptance behavior.
