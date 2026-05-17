# Interview Phases

## Phase 1: Problem Framing

Goals:

- understand what the application is
- identify why it exists
- identify the primary user or buyer

Minimum questions:

- What problem is this application solving?
- Who is the primary user?
- What is the most important workflow?
- What is the business outcome of success?

Exit criteria:

- one-sentence product definition exists
- at least one primary actor is known
- one primary workflow is named

## Phase 2: Scope and Boundaries

Goals:

- separate goals from non-goals
- determine system boundary
- identify what is intentionally out of scope

Minimum questions:

- What must the first version do?
- What explicitly does not belong in this system?
- What other systems or manual processes touch this workflow?

Exit criteria:

- first-slice scope is bounded
- at least three non-goals or exclusions are known
- boundary assumptions are visible

## Phase 3: Workflow and States

Goals:

- identify major business objects
- map key lifecycle transitions
- uncover exception and approval flows

Minimum questions:

- What are the key records or business objects?
- When is each object created?
- What states does it pass through?
- What transitions are forbidden?
- What approvals or exceptions exist?

Exit criteria:

- at least one core workflow is mapped end to end
- at least one object lifecycle is explicit
- at least one unhappy path is explicit

## Phase 4: Rules and Constraints

Goals:

- capture hard invariants
- identify configurable policies
- identify thresholds, triggers, and edge cases

Minimum questions:

- What must never happen?
- What rules are configurable?
- What thresholds trigger review, rejection, or escalation?
- What edge cases matter in practice?

Exit criteria:

- hard rules and flexible policies are separated
- at least three business constraints are explicit
- at least one threshold-driven rule is explicit

## Phase 5: Read Needs and Operational Concerns

Goals:

- identify reporting needs
- capture visibility requirements
- expose retry, audit, and traceability concerns

Minimum questions:

- What summaries or reports do users need?
- What decisions depend on those reports?
- What actions need to be auditable or explainable?
- What happens when commands are retried or partially fail?

Exit criteria:

- at least two read/reporting needs are explicit
- auditability expectations are known
- at least one failure/recovery concern is explicit

## Phase 6: Drafting Readiness

Goals:

- validate there is enough information to write reference documents
- identify residual risks and bounded assumptions

Minimum questions:

- What is still ambiguous?
- Can we proceed with explicit assumptions?
- Which unresolved items block drafting versus implementation?

Exit criteria:

- open questions are listed
- assumptions are listed
- drafting feasibility is explicit
