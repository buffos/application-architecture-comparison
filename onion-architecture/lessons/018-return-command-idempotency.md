# Lesson 018: Return Command Idempotency

## Objective

Make return review commands safe to retry without replaying refund or restock side effects.

## Theory

The return workflow now has real branching, policy, time, and actor metadata.

That increases one operational risk:

- the same accept or reject command may be retried by a caller

If retries run the full workflow again, the system can refund twice or restock twice.

Onion Architecture handles that by adding another application-owned contract:

- the application ring owns an idempotency store
- infrastructure provides the storage implementation
- the domain core remains unchanged

This is the right boundary because idempotency is workflow protection, not a domain invariant on the return aggregate itself.

## Why This Matters Here

Without idempotency, the return review workflow is correct only under ideal delivery conditions.

Adding an idempotency boundary makes the workflow safer under retries:

- the first successful command records the result
- a duplicate command returns that result
- refund and restock do not happen twice

## Diagram

```mermaid
flowchart LR
    subgraph DOM["Domain Core"]
        direction TB
        RET["Return Request Entity"]
    end

    subgraph APP["Application Ring"]
        direction TB
        ARS["AcceptReturn Service"]
        RJS["RejectReturn Service"]
        IDS["Idempotency Store"]
        RFD["Refund Gateway"]
        RSK["Inventory Restock"]
    end

    subgraph INF["Infrastructure Ring"]
        direction TB
        MID["Memory Idempotency Store"]
        ARG["Accept-All Refund Gateway"]
        MIR["Memory Inventory Reservation"]
    end

    ARS --> RET
    RJS --> RET

    IDS -.used by.-> ARS
    IDS -.used by.-> RJS
    RFD -.used by.-> ARS
    RSK -.used by.-> ARS
    IDS -.implemented by.-> MID
    RFD -.implemented by.-> ARG
    RSK -.implemented by.-> MIR

    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef app fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef domain fill:#fff3bf,stroke:#b08900,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class MID,ARG,MIR dataadapter;
    class ARS,RJS,IDS,RFD,RSK app;
    class RET domain;
    class IDS,RFD,RSK contract;
```

## Implementation Focus

Implement one retry-safety refinement:

- idempotent accept and reject return commands

The code should show:

- idempotency keys on review commands
- an application-owned idempotency store contract
- an in-memory idempotency store
- tests proving duplicate retries do not replay side effects

## What To Verify

- `go test ./...` passes
- duplicate accept commands do not refund or restock twice
- duplicate reject commands return the same result
- idempotency stays outside the domain core
