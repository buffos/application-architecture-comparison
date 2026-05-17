# Domain Quality Gate

Review the final domain model against this checklist.

## Language Quality

- Is the ubiquitous language stable and business-centered?
- Are important distinctions explicit?
- Would a product engineer and a domain expert recognize the concepts?

## Structural Quality

- Are subdomains justified?
- Are aggregate boundaries explained?
- Are entities and value objects meaningfully separated?
- Are policies and services used only where they help?

## Behavioral Quality

- Are invariants explicit?
- Are lifecycles and transitions visible?
- Are failure or exception paths reflected in the model?
- Are domain events meaningful business moments?

## Neutrality Quality

- Does the model avoid database-first thinking?
- Does it avoid package-structure bias?
- Could multiple architecture styles implement it faithfully?

## Downstream Usefulness

- Can application services be derived from it?
- Can canonical contracts be derived from it?
- Does it explain what business behavior must remain stable?

## Reject Conditions

Revise the model if it reads like:

- a relational schema
- a class diagram with no behavioral meaning
- a list of services without clear business ownership
- an architecture essay with thin domain grounding
