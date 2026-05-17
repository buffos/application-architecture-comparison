# Pipeline Entry Points

Choose the starting point based on current artifact maturity.

## Only a Vague Idea Exists

Start with:

- `$requirements-discovery-interviewer`

Then likely:

- `$domain-glossary-extractor`
- `$requirements-gap-analyzer`

## Discovery Notes Exist but Are Messy

Start with:

- `$requirements-gap-analyzer`

Then use:

- `$domain-glossary-extractor` if terminology is unstable
- `$prd-author` when discovery is ready

## PRD Exists but Canonical Docs Do Not

Start with:

- `$domain-glossary-extractor` if terms are unstable
- otherwise `$canonical-domain-modeler`

## PRD and Domain Model Exist

Start with:

- `$canonical-use-case-modeler`

Then:

- `$canonical-contract-author`
- `$acceptance-scenario-author`

## Full Artifact Set Exists but Confidence Is Low

Start with:

- `$architecture-readiness-reviewer`

If it finds issues, route backward to the right specialist skill.
