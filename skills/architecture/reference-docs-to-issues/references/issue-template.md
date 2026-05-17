# Issue Template

Use this local markdown template.

```markdown
## Parent Artifacts

- `docs/prd.md`
- `docs/canonical-domain-model.md`
- `docs/canonical-use-cases.md`
- `docs/canonical-api-cli-contract.md`
- `docs/acceptance-scenarios.md`

Include only the artifacts that actually exist and are relevant.

## What to build

A concise description of the vertical slice. Describe the end-to-end behavior, not a layer-by-layer task list.

## Acceptance criteria

- [ ] Criterion 1
- [ ] Criterion 2
- [ ] Criterion 3

## Blocked by

- Blocked by `issues/pending/NNN-title.md`

Or:

`None - can start immediately`

## Artifact anchors

- PRD: section or requirement anchor
- Domain model: concept, invariant, or lifecycle anchor
- Use cases: command/query or service anchor
- Contract: endpoint/CLI or external behavior anchor

Include only the anchors that exist and matter.

## Acceptance scenarios addressed

- Scenario title 1
- Scenario title 2

If no acceptance-scenarios artifact exists yet, say:

`Not yet defined`
```
