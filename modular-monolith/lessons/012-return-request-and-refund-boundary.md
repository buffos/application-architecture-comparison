# Lesson 012: Return Request And Refund Boundary

## Objective

Add the first post-shipment reverse workflow by introducing a `returns` module that requests returns against shipped orders and triggers refunds through the `payments` module.

## Theory

Cancellation already covers the pre-shipment reversal path.

That is not enough for a realistic workflow because shipped orders need a different boundary:

- the order has already moved forward
- the stock reversal will come later
- the customer-facing reversal starts with a return request and refund

This lesson keeps that path modular:

- `orders` owns whether an order is returnable
- `returns` owns the return request record
- `payments` owns the refund capability

So return processing becomes a new cross-module workflow instead of being treated like delayed cancellation.

## Why This Matters Here

This is the first place where the modular monolith needs two different reversal concepts:

- cancel before shipment
- return after shipment

That distinction matters because it proves the modules are following business boundaries, not just technical reuse. A shipped order is no longer cancellable, but it can still participate in a different workflow owned by a different module.

## Diagram

```mermaid
flowchart LR
    subgraph ORM["Orders Module"]
        direction TB
        ORO["orders.ReturnableOrderSource"]
        OMS["orders.Service"]
    end

    subgraph RTM["Returns Module"]
        direction TB
        RRE["returns.Repository"]
        RTS["returns.Service<br/>RequestReturn"]
        RTR["ReturnRequest"]
    end

    subgraph PAM["Payments Module"]
        direction TB
        PRF["payments.Refunder"]
        PMS["payments.Service"]
        PGW["payments.Gateway"]
    end

    subgraph INF["Infrastructure"]
        direction TB
        CLI["CLI"]
        MRR["Memory Return Request Repository"]
        PAG["Accept-All Payment Gateway"]
    end

    CLI --> RTS
    RTS --> RTR

    ORO -.used by.-> RTS
    RRE -.used by.-> RTS
    PRF -.used by.-> RTS
    PGW -.used by.-> PMS

    ORO -.implemented by.-> OMS
    RRE -.implemented by.-> MRR
    PRF -.implemented by.-> PMS
    PGW -.implemented by.-> PAG

    classDef module fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef entity fill:#fff3bf,stroke:#b08900,color:#111;
    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef framework fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class ORO,OMS,RRE,RTS,PRF,PMS,PGW module;
    class RTR entity;
    class MRR,PAG dataadapter;
    class CLI framework;
    class ORO,RRE,PRF,PGW contract;
```

Legend:

- yellow: domain type
- purple: module-owned service or contract
- green: data adapter
- blue: framework edge
- dashed border: contract
- dashed arrow: structural relationship such as `used by` or `implemented by`

## Implementation Focus

Implement one new post-shipment reverse workflow:

- request a return for a shipped order

The code should show:

- returnability still owned by the `orders` module
- return request storage owned by the `returns` module
- refund capability owned by the `payments` module
- shipped orders becoming returnable while non-shipped orders stay blocked

## What To Verify

- `go test ./...` passes
- only shipped orders can be returned
- successful return requests trigger a refund
- the return request is stored in the `returns` module
