# Neutrality Checks

Review whether the artifacts stay neutral enough for fair comparison.

## DDD Bias

Flag if:

- the PRD assumes aggregates or bounded contexts as product requirements
- the domain model presumes rich modeling where the requirements do not justify it

## CRUD Bias

Flag if:

- important workflows are flattened into generic create/update/delete language
- approval, policy, and lifecycle behaviors disappear behind CRUD wording

## REST Bias

Flag if:

- external contract semantics drive use-case naming
- route shape seems to define the product more than the use-case model does

## Event-Driven Bias

Flag if:

- domain events are treated as mandatory implementation style rather than preserved business moments

## Plugin Bias

Flag if:

- plugin architecture is over-specified without business justification
- extension points outnumber real variable behaviors
