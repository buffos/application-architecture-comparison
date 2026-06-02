# Lesson 011: Order Cancellation And Release

## Objective

Add the first reverse order workflow: cancel an unshipped order and release its reserved stock through the `inventory` module.

## Theory

The forward path now reaches:

- order conversion
- inventory reservation
- payment capture
- shipment creation

But a realistic order module also needs a reverse path before shipment.

This lesson keeps that reversal modular:

- `orders` owns cancellation rules
- `inventory` owns stock release
- `orders` calls the `inventory` module API to undo the reservation

That means reversal is still a workflow across modules, not just a field change inside the order record.

## Why This Matters Here

Cancellation is the first place where the modular monolith has to prove it can unwind earlier cross-module work.

That matters because it shows:

- module boundaries still hold during reversal
- releasing stock is not secretly an order responsibility
- shipped orders can be blocked by order rules before the workflow reaches inventory

This is where modules start looking less like folders and more like coordinated business capabilities.

## Diagram

```mermaid
flowchart LR
    subgraph ORM["Orders Module"]
        direction TB
        ORE["orders.Repository"]
        OMS["orders.Service<br/>CancelOrder"]
        ORD["Order"]
    end

    subgraph INM["Inventory Module"]
        direction TB
        IRL["inventory.Releaser"]
        IMS["inventory.Service"]
        IRP["inventory.Repository"]
    end

    subgraph INF["Infrastructure"]
        direction TB
        CLI["CLI"]
        MOR["Memory Order Repository"]
        MIR["Memory Inventory Repository"]
    end

    CLI --> OMS
    OMS --> ORD

    ORE -.used by.-> OMS
    IRL -.used by.-> OMS
    IRP -.used by.-> IMS

    ORE -.implemented by.-> MOR
    IRL -.implemented by.-> IMS
    IRP -.implemented by.-> MIR

    classDef module fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef entity fill:#fff3bf,stroke:#b08900,color:#111;
    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef framework fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class ORE,OMS,IRL,IMS,IRP module;
    class ORD entity;
    class MOR,MIR dataadapter;
    class CLI framework;
    class ORE,IRL,IRP contract;
```

Legend:

- yellow: domain type
- purple: module-owned service or contract
- green: data adapter
- blue: framework edge
- dashed border: contract
- dashed arrow: structural relationship such as `used by` or `implemented by`

## Implementation Focus

Implement one reverse workflow step:

- cancel an unshipped order

The code should show:

- cancellation rules owned by the `orders` module
- release capability owned by the `inventory` module
- reserved stock being released only after the order passes cancellation rules
- shipped orders staying non-cancellable

## What To Verify

- `go test ./...` passes
- unshipped orders can be cancelled
- cancellation releases reserved stock
- shipped orders are rejected
