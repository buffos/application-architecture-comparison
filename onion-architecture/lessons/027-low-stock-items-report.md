# Lesson 027: Low Stock Items Report

## Objective

Add an operational inventory report that introduces a narrow stock read boundary in the application ring.

## Theory

So far, inventory in the Onion track has only appeared as command-side behavior:

- reserve stock
- release stock
- restock returned items

That is enough for workflows, but not enough for operational visibility.

This lesson introduces a small read seam:

- infrastructure can expose stock snapshots
- the application ring decides what counts as low stock
- the report stays an application concern rather than an infrastructure query shortcut

This matters because the threshold rule is business-facing. The memory adapter should not decide what "low" means.

## Why This Matters Here

The Onion model is already showing:

- domain rules in the core
- workflow orchestration in the application ring
- storage and gateways in infrastructure

This lesson adds one more useful distinction:

- infrastructure provides raw stock state
- application turns that state into a report

That keeps the reporting rule close to the use case instead of burying it inside the repository.

## Diagram

```mermaid
flowchart LR
    subgraph DOM["Domain Ring"]
        direction TB
        STK["InventoryStockRecord"]
    end

    subgraph APP["Application Ring"]
        direction TB
        ISR["InventoryStockReader"]
        LSR["LowStockItemsReport Service"]
        LRP["Low Stock Report Rows"]
    end

    subgraph INF["Infrastructure Ring"]
        direction TB
        MIR["Memory Inventory Reservation"]
    end

    LSR --> STK
    LSR --> LRP

    ISR -.used by.-> LSR
    ISR -.implemented by.-> MIR

    classDef entity fill:#fff3bf,stroke:#b08900,color:#111;
    classDef app fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class STK entity;
    class ISR,LSR,LRP app;
    class MIR dataadapter;
    class ISR contract;
```

Legend:

- yellow: domain type
- purple: application type
- green: infrastructure data adapter
- dashed border: contract
- dashed arrow: structural relationship such as `used by` or `implemented by`

## Implementation Focus

Add:

- a stock snapshot record in the domain ring
- an application-owned stock reader contract
- a low stock report service that filters by threshold

The infrastructure adapter should only return current stock snapshots.

## What To Verify

- `go test ./...` passes
- items at or below the requested threshold are included
- the demo can print the low-stock report
