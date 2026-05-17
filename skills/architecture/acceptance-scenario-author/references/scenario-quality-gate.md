# Scenario Quality Gate

Review the scenario set against this checklist.

## Behavioral Quality

- Are the scenarios externally observable?
- Do they express meaningful business outcomes?
- Are important rules and branches covered?

## Coverage Quality

- Is there more than a happy path?
- Are failure, approval, or exception paths included where relevant?
- Are read/report scenarios included when product value depends on them?

## Neutrality Quality

- Do the scenarios avoid API- or CLI-specific wording?
- Could multiple architectures be tested against them unchanged?

## Efficiency Quality

- Are the scenarios distinct and non-redundant?
- Does each scenario prove something useful?

## Reject Conditions

Revise the set if it reads like:

- a UI walkthrough
- a route-level test list
- a long repetition of nearly identical cases
- a feature checklist with no real outcomes
