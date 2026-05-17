---
name: architecture-lesson-implementer
description: Drive the repeated learning-and-implementation loop for this repository's architecture exercises. Use when Codex must implement one architecture at a time from `docs/architectures.md`, create the next numbered lesson inside that architecture's `lessons/` folder first, explain the theory briefly with an optional Mermaid diagram, and then implement only the code required for that lesson inside the architecture's own solution folder.
---

# Architecture Lesson Implementer

Run the architecture exercise as a sequence of small lessons.

The rule is simple:

1. create the next lesson
2. explain the theory briefly
3. show a Mermaid diagram when it clarifies the lesson
4. implement that lesson

Do not jump straight into coding without creating the lesson first.

## Repository Model

This repository is expected to evolve as:

- one solution folder per architecture
- one `lessons/` folder inside each architecture folder
- incremental lesson numbering inside each `lessons/` folder

Use [references/folder-conventions.md](references/folder-conventions.md).

## Workflow

### 1. Identify the Target Architecture

Use the list in:

- `docs/architectures.md`

If the user names an architecture explicitly, use that.

If the user does not, determine the next architecture to work on from repo context or ask only if necessary.

### 2. Locate or Create the Architecture Folder

Each architecture should live in its own folder.

The skill should:

- locate the existing folder for the architecture, or
- create a new folder if this is the first lesson for that architecture

The exact folder name should be stable and readable.

Example style:

- `layered-architecture/`
- `hexagonal-architecture/`
- `clean-architecture/`

Do not mix multiple architectures in the same solution folder.

### 3. Create the Next Lesson Before Code

Inside the architecture folder, create or use:

- `lessons/`

Then create the next numbered lesson file before implementation begins.

Use incremental numbering such as:

- `001-introduction.md`
- `002-application-service.md`
- `003-aggregate-boundary.md`

Use [references/lesson-format.md](references/lesson-format.md).

### 4. Keep Lessons Small and Cumulative

Each lesson should teach one meaningful architectural idea and implement only that step.

Good lesson themes:

- project structure
- dependency direction
- ports and adapters
- use-case interactor
- aggregate boundary
- repository abstraction
- rule evaluation boundary
- plugin extension point

Avoid giant lessons like:

- "implement the whole architecture"

### 5. Explain the Theory Briefly

The lesson must explain:

- what concept is being introduced
- why this architecture uses it
- what problem it solves
- what tradeoff it introduces, if important

Keep explanations brief and practical.

The goal is to make the upcoming code intelligible, not to write a textbook.

### 6. Use Mermaid When It Helps

Add a Mermaid diagram when it improves understanding of:

- dependency direction
- request flow
- module boundaries
- domain/application/infrastructure separation
- plugin or rule flow

Do not add diagrams mechanically. Use them when they actually clarify the lesson.

### 7. Implement Only the Current Lesson

After writing the lesson, implement only the code required for that lesson.

Do not smuggle in several future lessons worth of infrastructure.

The implementation should be:

- coherent
- runnable or verifiable where practical
- aligned with the theory just introduced

### 8. Keep the Architecture Folder Internally Consistent

Within each architecture solution:

- preserve the style of that architecture
- do not contaminate it with patterns from other architectures unless the lesson is explicitly about comparison

Examples:

- do not accidentally turn a transaction-script solution into rich DDD
- do not force ports/adapters into a layered lesson too early

### 9. Explain What Changed

After implementing the lesson, summarize:

- what lesson was added
- what code was implemented
- what architectural concept is now visible in the code

## Lesson Writing Rules

### Lesson First, Code Second

Never reverse this order.

### Brief Theory Only

Keep theory concise and practical.

### Name the Problem Being Solved

A lesson should not only say what we are doing. It should say why.

Examples:

- "We introduce a repository interface here to stop application logic from depending on storage details."
- "We introduce an aggregate boundary here to keep pricing and approval invariants coherent."

### Prefer Visuals for Flow and Boundaries

Mermaid is especially useful when showing:

- inward dependency flow
- adapter boundaries
- module relationships
- lifecycle or command flow

### Keep Lessons Incremental

If the next lesson depends on this one, stop cleanly at the current boundary.

## Implementation Rules

### Respect the Target Architecture

The whole point is to make architectural differences visible.

Do not flatten the implementation into the same generic structure every time.

### Keep Code Close to the Lesson

If the lesson is about application services, the implementation should make that concept visible.

If the lesson is about aggregates, the implementation should make aggregate behavior visible.

### Do Not Overbuild

Only implement enough to make the lesson meaningful and demonstrable.

### Prefer Repeatable Progress

Each lesson should leave the architecture folder in a state where the next lesson can build naturally on top of it.

## Suggested Lesson Output

Each lesson file should usually include:

- title
- objective
- brief theory
- why this matters here
- Mermaid diagram when useful
- implementation focus
- what to verify

## Anti-Patterns

Do not:

- skip the lesson file
- write long academic theory dumps
- implement multiple future lessons in one step
- mix architectural styles accidentally
- make the lesson so abstract that it does not map to code

## Output Contract

For each run, produce:

- the next lesson file in the target architecture's `lessons/` folder
- the code implementing that lesson in the architecture folder
- a concise summary of what was taught and implemented

## Resources

Read only what you need:

- [references/folder-conventions.md](references/folder-conventions.md): expected repo structure for architecture solutions
- [references/lesson-format.md](references/lesson-format.md): how to structure lesson files
- [references/architecture-boundary-rules.md](references/architecture-boundary-rules.md): how to avoid cross-contaminating architecture styles
