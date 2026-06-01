# Lesson 027: Low Stock Items Report

## Objective

Add an inventory-based report that introduces a narrow stock read boundary without turning the Clean track into a full inventory management slice yet.

## Theory

So far, inventory has only appeared as a command-side collaborator:

- reserve stock
- release stock
- restock items

That is enough for workflow behavior, but it does not expose stock as readable application information.

This lesson introduces a small but important idea:

- sometimes a report needs a read seam into a supporting subsystem even if that subsystem does not yet have its own full query use-case family

The report answers:

- which items are at or below a given stock threshold

This is still a Clean Architecture use case.

The application layer owns:

- the threshold input
- the low-stock selection rule
- the report output shape

The infrastructure layer only provides stock snapshots.

## Why This Matters Here

The reporting track now has workflow projections.

This lesson adds an operational projection and shows that reports can also surface infrastructure-backed operational state without giving infrastructure ownership of the report semantics.

It also prepares the ground for richer inventory lessons later if needed.

## Diagram

```mermaid
flowchart LR
    subgraph ENT["Entities"]
        direction TB
        STK["Inventory Stock Record"]
    end

    subgraph APP["Application"]
        direction TB
        RIN["LowStockItemsReport Input Boundary"]
        ROUT["LowStockItemsReport Output Boundary"]
        RUC["LowStockItemsReport Interactor"]
        ISR["Inventory Stock Reader"]
        RPT["Low Stock Report Model"]
    end

    subgraph IA["Interface Adapters"]
        direction TB
        RCTRL["LowStockItemsReport Controller"]
        RPRES["LowStockItemsReport Presenter"]
    end

    subgraph INFRA["Infrastructure / Frameworks"]
        direction TB
        CLI["CLI / HTTP Framework"]
        MIR["Memory Inventory Reservation"]
    end

    CLI --> RCTRL
    RCTRL --> RIN
    RUC --> ROUT
    RPRES --> CLI
    RUC --> STK
    RUC --> RPT

    RIN -.used by.-> RCTRL
    RIN -.implemented by.-> RUC
    ROUT -.used by.-> RUC
    ROUT -.implemented by.-> RPRES
    ISR -.used by.-> RUC
    ISR -.implemented by.-> MIR

    classDef framework fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef funcadapter fill:#ffe5d9,stroke:#bc6c25,color:#111;
    classDef app fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef entity fill:#fff3bf,stroke:#b08900,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class CLI framework;
    class MIR dataadapter;
    class RCTRL,RPRES funcadapter;
    class RIN,ROUT,RUC,ISR,RPT app;
    class STK entity;
    class RIN,ROUT,ISR contract;
```

Legend:

- blue: framework edge
- green: data adapter
- orange: translation adapter
- purple: application layer
- yellow: entity layer
- dashed border: interface / contract
- dashed arrow: structural relationship such as `used by` or `implemented by`

## Implementation Focus

Add:

- `LowStockItemsReport`

The code should show:

- a stock snapshot reader contract owned by the application layer
- threshold filtering in the interactor, not in infrastructure
- a presenter shaping low-stock results for callers

## What To Verify

- the project compiles
- `go test ./...` passes
- items at or below the threshold are included
- the demo can render the report output
