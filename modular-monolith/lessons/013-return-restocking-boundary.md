# Lesson 013: Return Restocking Boundary

## Objective

Extend the return workflow so refunded returns also restock inventory through a distinct inventory capability.

## Theory

Lesson `012` made returns a separate post-shipment workflow.

But it still only handled the customer-facing side:

- verify the order is returnable
- refund the payment
- store the return request

The stock-side reversal was still missing.

This lesson makes that explicit:

- `returns` still owns the return-request workflow
- `payments` still owns refund capability
- `inventory` now owns restocking capability

That matters because returns are not just reversed payments. They also change stock state, and that stock state should still belong to the inventory module.

## Why This Matters Here

This is the first post-shipment workflow that fans out into multiple collaborating modules:

- `orders` decides whether the order is eligible
- `payments` handles the money reversal
- `inventory` handles the stock reversal
- `returns` orchestrates and records the workflow

That is a more realistic modular-monolith example than a simple two-module call chain.

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

    subgraph INM["Inventory Module"]
        direction TB
        IRS["inventory.Restocker"]
        IMS["inventory.Service"]
        IRP["inventory.Repository"]
    end

    subgraph INF["Infrastructure"]
        direction TB
        CLI["CLI"]
        MRR["Memory Return Request Repository"]
        MIR["Memory Inventory Repository"]
        PAG["Accept-All Payment Gateway"]
    end

    CLI --> RTS
    RTS --> RTR

    ORO -.used by.-> RTS
    RRE -.used by.-> RTS
    PRF -.used by.-> RTS
    IRS -.used by.-> RTS
    PGW -.used by.-> PMS
    IRP -.used by.-> IMS

    ORO -.implemented by.-> OMS
    RRE -.implemented by.-> MRR
    PRF -.implemented by.-> PMS
    PGW -.implemented by.-> PAG
    IRS -.implemented by.-> IMS
    IRP -.implemented by.-> MIR

    classDef module fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef entity fill:#fff3bf,stroke:#b08900,color:#111;
    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef framework fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class ORO,OMS,RRE,RTS,PRF,PMS,PGW,IRS,IMS,IRP module;
    class RTR entity;
    class MRR,MIR,PAG dataadapter;
    class CLI framework;
    class ORO,RRE,PRF,PGW,IRS,IRP contract;
```

Legend:

- yellow: domain type
- purple: module-owned service or contract
- green: data adapter
- blue: framework edge
- dashed border: contract
- dashed arrow: structural relationship such as `used by` or `implemented by`

## Implementation Focus

Implement one missing stock-side reversal:

- restock inventory when a return is refunded

The code should show:

- restocking as a distinct inventory capability
- `returns` orchestrating refund and restock together
- return requests still being stored only by the `returns` module

## What To Verify

- `go test ./...` passes
- successful returns trigger both refund and restock
- restock failures stop the workflow
- the restock logic stays behind the `inventory` module boundary
