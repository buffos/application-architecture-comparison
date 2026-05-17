# Artifact Priority

Use the strongest available artifacts first.

## Highest Value

- `prd.md`
- `canonical-domain-model.md`
- `canonical-use-cases.md`
- `canonical-api-cli-contract.md`
- `acceptance-scenarios.md`

These should define the intended product behavior.

## Supporting Value

- `architecture-readiness-review.md`
- `domain-glossary.md`
- `discovery-notes.md`
- `requirements-gap-analysis.md`
- `reference-doc-orchestration-status.md`

These help explain weaknesses, vocabulary, and process state.

## Rule

If richer canonical artifacts exist, do not fall back to PRD-only slicing.

Use the whole reference set so issue boundaries reflect:

- domain rules
- use-case boundaries
- external contract shape
- acceptance behavior

## Stop Condition

If the current artifacts are contradictory or obviously incomplete, do not generate agent-ready issues until the gap is resolved.
