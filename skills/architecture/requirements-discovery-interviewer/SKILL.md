---
name: requirements-discovery-interviewer
description: Interrogate vague or incomplete software product ideas until they are concrete enough to support high-quality planning artifacts. Use when a user wants to build an application but does not yet have a stable PRD or reference design documents and Codex must lead broad, phased discovery by asking targeted questions, surfacing assumptions, pressure-testing workflows, and producing organized discovery notes. Do not use for narrow follow-up on an already mature artifact set; use `requirements-gap-analyzer` instead.
---

# Requirements Discovery Interviewer

Drive the user from a fuzzy idea to a decision-ready problem statement. Ask questions in phases, close ambiguities before moving on, and leave behind structured discovery notes that another skill can convert into formal documents.

Do not rush into solutioning. The purpose of this skill is to extract the missing information, challenge weak assumptions, and make hidden decisions explicit.

## Workflow

### 1. Establish the Engagement Mode

- Confirm whether the user wants:
- a brand-new application definition
- a refinement of an existing rough idea
- a rescue of a contradictory or underspecified spec set

If the user already has documents, interrogate them instead of starting from zero.

### 2. Stay in Discovery Until Exit Criteria Are Met

Do not move to PRD writing while major unknowns remain in any of these areas:

- user and actor model
- core business workflow
- domain vocabulary
- critical state transitions
- business rules and constraints
- integration or boundary assumptions
- non-goals
- success criteria

Use the phase checklist in [references/interview-phases.md](references/interview-phases.md).

### 3. Ask Questions in Batches, Not Dumps

Ask 3-7 high-value questions at a time. Prefer:

- one anchor question
- a few narrowing questions
- one challenge question that tests an assumption or exposes a tradeoff

Do not ask twenty generic questions in one burst. Make each round responsive to what the user already answered.

### 4. Prefer Gap-Closing Questions

Use questions that eliminate future ambiguity:

- “Who is the primary actor and who only reads data?”
- “What event turns a draft record into a committed business object?”
- “What must never happen even if the UI tries to allow it?”
- “What changes if this flow fails halfway through?”
- “Which rules are policy choices versus hard invariants?”

Use the deeper prompts in [references/question-bank.md](references/question-bank.md).

### 5. Pressure-Test the User's Mental Model

Actively look for:

- conflicting goals
- hand-wavy workflows
- overloaded terminology
- fake requirements that are really implementation preferences
- “CRUD disguise” hiding actual business workflows
- missing failure states
- missing approval or exception paths
- requirements that collapse under scale, retries, or policy changes

When a user gives a broad statement, convert it into a testable question.

Example:

User says: “It should support returns.”

Ask:

- “Who can initiate a return?”
- “What makes a return valid or invalid?”
- “Does accepted return restock inventory, create a refund, or both?”
- “What quantities or time windows constrain it?”

### 6. Keep a Running Discovery Ledger

Maintain these sections as you go:

- confirmed facts
- open questions
- assumptions you are making
- contradictions needing resolution
- likely future document sections

Do not keep this only in hidden reasoning. Surface it back to the user regularly so they can correct drift early.

### 7. Stop Only When the Outputs Are Draftable

This skill is complete when you can produce discovery notes that are sufficient for:

- a PRD
- a canonical domain model
- a canonical use-case/application-service document
- a canonical API/CLI contract

Use the artifact readiness checklist in [references/artifact-readiness.md](references/artifact-readiness.md).

## Interview Strategy

### Start Broad, Then Narrow

Open with the smallest set of questions that identifies:

- what the application is for
- who uses it
- what core workflow matters most
- why the user wants the system built

Then narrow into:

- stateful business objects
- core lifecycle transitions
- important policies
- edge cases and failure paths
- external boundaries

### Prefer Domain Language Over Technical Language

Ask:

- “What is the business object called?”
- “What happens before it is considered final?”
- “What can invalidate it?”

Avoid early questions like:

- “Do you want microservices?”
- “Which database?”
- “Should we use DDD?”

Unless the user's problem explicitly depends on those concerns.

### Separate Business Requirements from Solution Preferences

If the user says:

- “I want event sourcing”
- “I want a plugin architecture”
- “I want CQRS”

Ask what business pressure motivates that preference. Capture both:

- the underlying requirement
- the preferred solution direction

This prevents architecture bias from polluting the PRD.

### Always Pull Out the Negative Space

Ask what the system will not do:

- excluded actors
- excluded workflows
- postponed features
- fake nice-to-haves
- forbidden states
- unacceptable outcomes

Non-goals are as important as goals.

## Questioning Rules

### Ask for Examples

Whenever the user gives an abstract answer, ask for one realistic example.

Examples:

- a typical happy path
- an approval exception
- a failure scenario
- a report someone would actually read

### Ask for Thresholds and Triggers

Whenever a rule sounds vague, ask what triggers it.

Examples:

- “What discount requires approval?”
- “What amount triggers manual review?”
- “When is cancellation no longer allowed?”

### Ask About State Changes

For each important business object, determine:

- how it is created
- what states it goes through
- what transitions are valid
- what transitions are forbidden
- who performs each transition

### Ask About Policy Flexibility

For each rule, determine whether it is:

- a hard invariant
- a configurable policy
- an environment-specific behavior
- a likely plugin or extension point

### Ask About Failure and Recovery

For each major workflow, ask:

- what can fail
- what happens next
- what gets rolled back
- what gets partially preserved
- what the user must see

## Anti-Patterns

Do not:

- accept vague nouns without definitions
- confuse screens with use cases
- confuse database tables with domain concepts
- assume happy path implies sufficient requirements
- skip reporting and read-model questions
- skip approval, exception, and retry scenarios
- let the user answer everything with “it depends” without forcing a decision or bounded assumption

## Output Contract

At the end of the interview, produce a discovery summary with these sections:

- problem statement
- target users and actors
- core workflows
- domain vocabulary
- major business entities
- business rules and constraints
- critical states and transitions
- external integrations or system boundaries
- reporting/read needs
- explicit non-goals
- open questions that still block drafting
- assumptions accepted for drafting

If major blockers remain, stop there and keep interviewing. Do not pretend the system is specified when it is not.

Persist the output to:

- `docs/discovery-notes.md`

Behavior:

- create the file if it does not exist
- update the file if discovery is continuing
- keep prior confirmed facts unless they were explicitly corrected

## Resources

Read only what you need:

- [references/interview-phases.md](references/interview-phases.md): phased workflow and exit criteria
- [references/question-bank.md](references/question-bank.md): high-signal question prompts by topic
- [references/artifact-readiness.md](references/artifact-readiness.md): readiness checks for PRD and canonical documents

---

**Not every skill requires all three types of resources.**
