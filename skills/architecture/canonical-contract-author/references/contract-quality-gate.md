# Contract Quality Gate

Review the final contract against this checklist.

## Semantic Quality

- Does the contract preserve the canonical use-case semantics?
- Are resource and command names business-centered?
- Are statuses and errors stable and meaningful?

## Transport Quality

- Are global conventions explicit?
- Are payload shapes clear enough for multiple implementations?
- Are HTTP and CLI mappings semantically aligned?

## Retry and Failure Quality

- Are idempotency-sensitive operations identified?
- Are business-visible failures represented explicitly?
- Are transport-level mappings reasonable but not overbearing?

## Neutrality Quality

- Does the contract avoid leaking internal package or persistence details?
- Could multiple implementations expose the same contract faithfully?

## Downstream Usefulness

- Can client tests be written from this contract?
- Can acceptance tests target both HTTP and CLI from it?
- Does it preserve a stable comparison surface across architectures?

## Reject Conditions

Revise the contract if it reads like:

- a framework-specific route file
- a random CRUD endpoint list
- a CLI syntax sheet detached from business semantics
- or an implementation-specific DTO dump
