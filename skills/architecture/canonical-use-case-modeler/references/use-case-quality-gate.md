# Use-Case Quality Gate

Review the final use-case document against this checklist.

## Intent Quality

- Are the use cases named after business intent?
- Can a reader understand what each command is for without reading implementation details?

## Structural Quality

- Are commands and queries clearly separated?
- Are application services grouped coherently?
- Are orchestration-heavy workflows described explicitly?

## Behavioral Quality

- Are preconditions and outcomes visible?
- Are failure modes represented?
- Are transaction and retry expectations surfaced where needed?

## Neutrality Quality

- Does the document avoid HTTP or CLI leakage?
- Does it avoid framework or package-layout bias?
- Could multiple architectures implement the same use-case surface faithfully?

## Downstream Usefulness

- Can an API/CLI contract be derived from this document?
- Can tests be mapped to the command/query surface?
- Does it explain what the application layer must preserve across implementations?

## Reject Conditions

Revise the document if it reads like:

- a route map
- a thin controller inventory
- a service class list without business semantics
- a CRUD checklist with renamed verbs
