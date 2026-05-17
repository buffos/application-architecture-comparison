---
name: domain-glossary-extractor
description: Extract, normalize, and stabilize domain vocabulary from discovery notes, a PRD, interviews, or other product artifacts. Use when Codex needs to identify the business terms that should become the project's canonical language, resolve synonyms and overloaded terminology, separate actor names from domain objects and policies, and produce a glossary that downstream PRD, domain-model, use-case, and contract work can rely on consistently. Use this early whenever terminology drift is the main blocker; do not substitute it for full domain modeling.
---

# Domain Glossary Extractor

Build the canonical vocabulary for the product before terminology drift hardens into the design.

This skill exists to stabilize names, meanings, distinctions, and aliases. It is not a domain model by itself, but it is often the fastest way to prevent the domain model from becoming inconsistent later.

## Workflow

### 1. Gather Raw Terms from the Source Material

Extract candidate terms from:

- discovery notes
- user answers
- PRDs
- rough specs
- existing docs

Look for:

- business nouns
- actor names
- workflow names
- status labels
- policy labels
- error or exception terms

Use [references/term-categories.md](references/term-categories.md).

### 2. Separate Term Types Early

Do not mix all nouns together.

Classify each term as one of:

- actor
- business object
- workflow/action
- policy/rule
- state/status
- metric/reporting term
- external system or integration term

This prevents early confusion like treating a workflow as an entity or a UI label as a domain object.

### 3. Detect Synonyms, Near-Synonyms, and Collisions

Look for cases where:

- two words mean the same thing
- one word means different things in different places
- a UI term conflicts with a business term
- a technical convenience label is replacing a better business name

Use [references/conflict-patterns.md](references/conflict-patterns.md).

Flag examples like:

- `quote` vs `proposal`
- `order` vs `purchase`
- `approval` vs `review`
- `return` vs `refund`

### 4. Choose Canonical Terms

For each concept, choose one primary name.

Document:

- canonical term
- definition
- aliases or discouraged synonyms
- why the distinction matters, if relevant

Prefer business-facing language over technical shorthand unless the technical term is already the stable business term.

### 5. Define Critical Distinctions

A strong glossary does more than list terms. It preserves distinctions that later modeling depends on.

Examples:

- `Quote` is not `Order`
- `ReturnRequest` is not `Refund`
- `Reservation` is not `Shipment`
- `Policy` is not `Invariant`

Capture those distinctions explicitly.

### 6. Identify Unstable or Missing Vocabulary

If the source material keeps changing words for the same concept, or cannot name a concept clearly, flag it.

This is often a sign that:

- the product understanding is still fuzzy
- the workflow is underspecified
- the domain model is not ready yet

Use [references/glossary-readiness.md](references/glossary-readiness.md).

### 7. Produce a Glossary That Downstream Skills Can Reuse

The output should be structured enough that later skills can rely on it for:

- PRD wording
- domain modeling
- use-case naming
- API/CLI naming
- architecture review

### 8. Run the Glossary Quality Gate

Review the output with [references/glossary-quality-gate.md](references/glossary-quality-gate.md).

If the glossary reads like:

- a random noun list
- a class diagram substitute
- a UI label dump
- or a weak synonym table with no semantic distinctions

it is wrong. Revise it.

## Extraction Rules

### Prefer Business Terms Over UI Labels

If the source says “order screen” or “manage products page,” look through to the underlying business concept.

### Prefer One Canonical Name Per Concept

Allow aliases to be recorded, but do not let every synonym remain equally valid in the reference set.

### Keep Definitions Short but Precise

Definitions should be one to three sentences, enough to remove ambiguity without turning into mini-specs.

### Preserve Critical Negative Distinctions

Good glossary entries often say what a concept is not.

Examples:

- `Quote`: negotiable commercial proposal, not a committed order
- `Refund`: financial consequence of an accepted return, not the return request itself

### Do Not Force Artificial Precision

If two terms are genuinely unresolved, note that clearly rather than inventing certainty.

## Anti-Patterns

Do not:

- treat every noun as a domain concept
- preserve all synonyms as equally valid
- define concepts in purely technical language
- ignore actor terminology
- ignore state and status vocabulary
- jump straight from glossary work into aggregate design unless asked

## Output Contract

Produce a glossary with:

- canonical term
- category
- concise definition
- aliases or discouraged terms
- important distinctions where needed
- open terminology issues, if any

If terminology is too unstable for confident glossary creation, stop and surface the conflicts rather than pretending the language is settled.

Persist the output to:

- `docs/domain-glossary.md`

Behavior:

- create the file if it does not exist
- update it when terminology is refined
- keep canonical terms, aliases, and distinctions stable unless there is explicit reason to rename them

## Resources

Read only what you need:

- [references/term-categories.md](references/term-categories.md): categories of terms worth extracting
- [references/conflict-patterns.md](references/conflict-patterns.md): common synonym and overload patterns
- [references/glossary-readiness.md](references/glossary-readiness.md): when vocabulary is stable enough to normalize
- [references/glossary-quality-gate.md](references/glossary-quality-gate.md): review checklist before finalizing
