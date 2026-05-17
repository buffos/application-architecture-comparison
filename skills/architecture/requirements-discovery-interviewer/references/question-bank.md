# Question Bank

Use these prompts selectively. Do not dump them all at once.

## Product and User

- What job is the user trying to get done?
- Who initiates the workflow and who only observes it?
- Who approves exceptions?
- Who is accountable when something goes wrong?

## Workflow

- Walk me through the most common happy path.
- What happens immediately before the core record is created?
- What event makes it committed or final?
- What are the common exception paths?

## Domain Concepts

- What are the nouns the business already uses?
- Which of those nouns are distinct objects versus just attributes?
- Which objects have their own lifecycle?

## Rules

- What must never be allowed?
- What is always allowed?
- What depends on thresholds, customer tier, category, or time window?
- Which rules are likely to change often?

## State and Lifecycle

- Which states matter to the business, not just to the UI?
- Who can move an object from one state to another?
- What prevents reversal of a transition?

## Integrations and Boundaries

- Which external systems provide data?
- Which systems consume outcomes?
- What can be mocked for the first implementation?

## Reporting and Read Models

- What dashboard or report would a manager ask for on day one?
- Which queues require operational visibility?
- What summaries help detect business problems early?

## Failure and Recovery

- What happens if a workflow succeeds halfway?
- What must be rolled back versus retained?
- What must be visible to the user after failure?
- Which actions must be safe to retry?

## Extensibility

- Which behaviors vary by customer, region, category, or policy?
- Which rules would likely become configurable?
- Where would future plugins or policy modules attach?
