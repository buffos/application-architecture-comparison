# Artifact Readiness

Use this checklist before concluding discovery.

## Ready for PRD

A PRD is draftable when:

- the product goal is clear
- target actors are known
- core workflows are named
- scope and non-goals are bounded
- business rules are partially concrete
- acceptance scenarios can be described

## Ready for Canonical Domain Model

A domain model is draftable when:

- core business nouns are stable
- at least one main lifecycle is clear
- important invariants are visible
- approval, exception, and failure flows are not hidden
- entity versus policy boundaries are becoming clear

## Ready for Canonical Use Cases

A use-case model is draftable when:

- user intents are clear
- major commands and queries can be named
- orchestration-heavy flows are known
- read concerns and write concerns are both visible

## Ready for Canonical API/CLI Contract

A contract is draftable when:

- command names are stable
- result states are stable
- major request inputs are known
- major failure modes are known
- core resources and identifiers are known

## Not Ready Signals

Do not proceed if:

- the user only described UI screens
- the business objects still have fuzzy names
- state transitions are hand-wavy
- policies and invariants are mixed together
- there is no unhappy path
- no one can say what success looks like
