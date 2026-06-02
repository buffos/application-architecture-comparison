# Lesson 008: Order Conversion With Reservation

## Objective

Extend quote-to-order conversion so the `orders` module reserves stock through an `inventory` module before it saves the new order.

## Theory

Lesson `007` proved that one module can hand a business snapshot to another.

That is useful, but still incomplete for a real workflow.

Order creation usually has an operational consequence:

- stock must be reserved

In a modular monolith, the important question is not just _whether_ stock is reserved, but _through which module boundary_ that happens.

This lesson makes that explicit:

- `orders` still owns order creation
- `inventory` owns stock reservation rules and storage
- `orders` calls the public reservation capability of `inventory`

So the workflow crosses modules, but each module still owns its own business area.

## Why This Matters Here

Without this step, order conversion is only a document transformation.

Reservation makes the workflow operational:

- an approved quote becomes an order
- the order consumes inventory capacity
- failure in one module can stop a workflow in another

That is where modular boundaries start to matter more than simple code grouping.

## Diagram

```mermaid
flowchart LR
    subgraph QTM["Quotes Module"]
        direction TB
        QAP["quotes.ApprovedQuoteSource"]
        QMS["quotes.Service"]
    end

    subgraph INM["Inventory Module"]
        direction TB
        IRS["inventory.Reserver"]
        IMS["inventory.Service"]
        IRP["inventory.Repository"]
    end

    subgraph ORM["Orders Module"]
        direction TB
        ORE["orders.Repository"]
        OMS["orders.Service<br/>ConvertQuoteToOrder"]
        ORD["Order"]
    end

    subgraph INF["Infrastructure"]
        direction TB
        CLI["CLI"]
        MOR["Memory Order Repository"]
        MIR["Memory Inventory Repository"]
    end

    CLI --> OMS
    OMS --> ORD

    QAP -.used by.-> OMS
    IRS -.used by.-> OMS
    ORE -.used by.-> OMS
    IRP -.used by.-> IMS

    QAP -.implemented by.-> QMS
    IRS -.implemented by.-> IMS
    ORE -.implemented by.-> MOR
    IRP -.implemented by.-> MIR

    classDef module fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef entity fill:#fff3bf,stroke:#b08900,color:#111;
    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef framework fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class QAP,QMS,IRS,IMS,IRP,ORE,OMS module;
    class ORD entity;
    class MOR,MIR dataadapter;
    class CLI framework;
    class QAP,IRS,IRP,ORE contract;
```

Legend:

- yellow: domain type
- purple: module-owned service or contract
- green: data adapter
- blue: framework edge
- dashed border: contract
- dashed arrow: structural relationship such as `used by` or `implemented by`

## Implementation Focus

Implement one operational extension:

- reserve stock during quote-to-order conversion

The code should show:

- a new `inventory` module
- a reservation capability published by that module
- `orders` depending on the inventory module API instead of inventory storage
- conversion stopping if reservation fails

## What To Verify

- `go test ./...` passes
- approved quotes reserve stock before order save
- insufficient stock stops conversion
- `orders` depends on the `inventory` module API, not on memory storage
