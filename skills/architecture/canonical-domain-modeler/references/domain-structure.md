# Domain Structure

Use this framing when deriving the canonical model.

## Ubiquitous Language

Start with the business terms the product actually uses.

Look for:

- terms that recur in workflows
- words that carry business meaning
- terms that are distinct in business language but easy to conflate technically

## Subdomains

Group concepts into:

- core domain
- supporting subdomains
- policy or extension areas

Use this to explain where the real business complexity lives.

## Bounded Context Candidates

A bounded context candidate is justified when:

- the same term means different things in different areas
- a workflow cluster is cohesive and relatively independent
- rules and language naturally cluster together

Do not invent contexts only because the architecture list includes DDD-oriented styles.

## Cross-Cutting Areas

Some concepts are better treated as policy spaces or supporting areas:

- pricing
- approvals
- reporting
- plugin registration

These may influence multiple workflows without becoming central aggregates themselves.
