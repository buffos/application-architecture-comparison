# Question Generation

Generate follow-up questions that reduce uncertainty quickly.

## Good Question Properties

- specific
- answerable in one reply
- tied to a concrete gap
- high downstream leverage

## Good Examples

- What event makes a quote no longer editable?
- Which discounts require approval, and which are outright rejected?
- Is shipment blocked until payment is accepted for all customers, or only some?
- What summaries does a manager need to see daily?

## Weak Examples

- Can you tell me more?
- What else should the system do?
- Any special cases?

## Batch Size

Ask 3-7 questions at a time unless the user explicitly asks for exhaustive questioning.

## Ordering

Order questions by:

1. blocker severity
2. behavioral impact
3. vocabulary stabilization
4. non-blocking refinements
