# Cross-Document Checks

Use this checklist to compare the reference artifacts.

## Vocabulary Consistency

Check whether the same concepts use the same names across:

- PRD
- domain model
- use-case model
- contract

Look for drift in:

- actor names
- entity names
- workflow names
- status labels

## Rule Consistency

Check whether business rules:

- appear in the PRD
- are modeled in the domain model
- influence the use-case layer
- surface through the contract where externally relevant

Flag gaps where a rule exists in one layer but disappears in another.

## Workflow Consistency

Check whether major workflows:

- exist in the PRD
- have matching domain semantics
- have corresponding commands and queries
- have contract-level visibility where needed

## Failure Consistency

Check whether failures and blocked paths:

- are described in the PRD
- are modeled as rules or invariants
- appear as business-visible outcomes in use cases
- are reflected in the contract error model

## Extension Consistency

Check whether plugins, policies, or extensibility points:

- are justified by the PRD
- are modeled in the domain
- influence application services
- appear meaningfully in the contract
