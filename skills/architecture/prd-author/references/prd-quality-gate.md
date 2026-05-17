# PRD Quality Gate

Review the draft against this checklist before finalizing.

## Product Clarity

- Is the product understandable in a few paragraphs?
- Is the core workflow obvious?
- Are the primary actors clearly defined?

## Scope Discipline

- Are non-goals explicit?
- Is the first-slice scope bounded?
- Does the document avoid pretending everything is in scope?

## Requirement Quality

- Are requirements concrete rather than generic?
- Are important rules and thresholds explicit?
- Are failure and exception paths represented?
- Are reporting and read needs represented?

## Architectural Neutrality

- Does the PRD avoid premature internal design choices?
- Are business needs separated from architectural preferences?
- Would multiple architectures still be able to implement the same product?

## Downstream Usefulness

- Can a domain model be derived from this PRD?
- Can use cases be named from it?
- Can acceptance scenarios be extracted from it?
- Would an implementation team know what behavior must exist?

## Reject Conditions

Revise the PRD if it reads like:

- a sprint backlog
- a UI wireframe description
- an architecture proposal
- a database schema draft
- a vague product vision statement without operational detail
