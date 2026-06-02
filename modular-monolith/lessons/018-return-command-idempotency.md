# Lesson 018: Return Command Idempotency

## Objective

Make return review commands retry-safe so duplicate accept and reject requests do not replay refund or restock side effects.

## Theory

The return workflow is now:

- requested
- reviewed
- policy-checked
- refunded and restocked on acceptance

That means `AcceptReturn` is no longer a harmless command to retry.

If the same command is sent twice after a timeout or retry, it could otherwise:

- refund twice
- restock twice
- produce inconsistent audit history

This lesson introduces a separate idempotency module so retry handling stays outside the return entity and outside the payment or inventory modules.

## Why This Matters Here

This is a good modular-monolith example because retry safety is a cross-cutting concern, but it still needs a clear home.

Putting idempotency directly inside:

- the payment module
- the inventory module
- or the return entity

would blur responsibilities.

A separate module lets the `returns` workflow ask one focused capability:

- has this command already succeeded?
- if so, return the stored result instead of replaying side effects

## Diagram

```mermaid
flowchart LR
    subgraph RTM["Returns Module"]
        direction TB
        RRE["returns.Repository"]
        RAS["returns.Service<br/>AcceptReturn / RejectReturn"]
        RTR["ReturnRequest"]
    end

    subgraph IDM["Idempotency Module"]
        direction TB
        IDS["idempotency.Store"]
        IDM1["idempotency.Service"]
    end

    subgraph PAM["Payments Module"]
        direction TB
        PRF["payments.Refunder"]
        PMS["payments.Service"]
    end

    subgraph INM["Inventory Module"]
        direction TB
        IRS["inventory.Restocker"]
        IMS["inventory.Service"]
    end

    subgraph INF["Infrastructure"]
        direction TB
        CLI["CLI"]
        MRR["Memory Return Request Repository"]
        MIS["Memory Idempotency Store"]
    end

    CLI --> RAS
    RAS --> RTR

    RRE -.used by.-> RAS
    IDS -.used by.-> RAS
    PRF -.used by.-> RAS
    IRS -.used by.-> RAS

    RRE -.implemented by.-> MRR
    IDS -.implemented by.-> MIS
    IDS -.implemented by.-> IDM1
    PRF -.implemented by.-> PMS
    IRS -.implemented by.-> IMS

    classDef module fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef entity fill:#fff3bf,stroke:#b08900,color:#111;
    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef framework fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class RRE,RAS,IDS,IDM1,PRF,PMS,IRS,IMS module;
    class RTR entity;
    class MRR,MIS dataadapter;
    class CLI framework;
    class RRE,IDS,PRF,IRS contract;
```

Legend:

- yellow: domain type or workflow record
- purple: module-owned service or contract
- green: adapter or technical implementation
- blue: framework edge
- dashed border: contract
- dashed arrow: structural relationship such as `used by` or `implemented by`

## Implementation Focus

Implement one retry-safety layer:

- accept and reject should be idempotent

The code should show:

- a separate `idempotency` module
- review commands carrying an idempotency key
- stored review results being returned on retries
- refunds and restocks not replaying when a result is already stored

## What To Verify

- `go test ./...` passes
- repeated accept commands reuse the stored result
- repeated reject commands reuse the stored result
- idempotent replays do not trigger refund or restock again
