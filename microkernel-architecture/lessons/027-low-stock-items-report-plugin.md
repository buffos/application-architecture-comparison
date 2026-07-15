# Lesson 027: Low Stock Items Report Plugin

## Objective

Add an operational inventory report and introduce a narrow stock read seam without turning the inventory plugin into a full query track yet.

## Theory

So far, the `inventory` plugin has only appeared as a command-side collaborator:

- reserve stock
- release stock
- restock returned items

That is enough for workflow behavior, but it does not yet expose stock as readable business information.

This lesson adds a small but important idea:

- some reports need a read seam into a supporting plugin
- that does not mean the plugin must immediately grow a full query family

The report answers:

- which items are at or below a given threshold

The inventory plugin contributes:

- stock snapshots

The reporting plugin owns:

- the threshold rule
- the low-stock projection
- the report output shape

## Why This Matters Here

Operational reporting is another place where microkernel systems often lose discipline and go straight to storage.

This lesson keeps the design honest:

- inventory still owns stock data access
- reporting still owns the meaning of the low-stock report
- infrastructure does not decide what "low stock" means

## Diagram

```mermaid
flowchart LR
    subgraph KER["Kernel"]
        direction TB
        RPT["kernel.Reporting"]
        INR["kernel.InventoryReader"]
        HST["kernel.Host"]
    end

    subgraph RPP["Reporting Plugin"]
        direction TB
        LSR["reporting.Service<br/>LowStockItemsReport"]
        LRP["LowStockItemsReport"]
        LRR["LowStockItemsReportRow"]
        RPG["reporting.Plugin"]
    end

    subgraph INP["Inventory Plugin"]
        direction TB
        IQS["inventory.Service<br/>ListStock"]
    end

    subgraph INF["Infrastructure"]
        direction TB
        CLI["CLI"]
        MIR["Memory Inventory Repository"]
    end

    CLI --> HST
    HST --> RPG
    LSR --> LRP
    LRP --> LRR

    INR -.used by.-> LSR
    RPT -.used by.-> CLI
    INR -.implemented by.-> IQS
    RPT -.implemented by.-> LSR
    MIR -.used by.-> IQS

    classDef kernel fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef plugin fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef entity fill:#fff3bf,stroke:#b08900,color:#111;
    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef framework fill:#f8f9fa,stroke:#6c757d,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class RPT,INR,HST kernel;
    class LSR,RPG,IQS plugin;
    class LRP,LRR entity;
    class MIR dataadapter;
    class CLI framework;
    class RPT,INR contract;
```

Legend:

- blue: kernel-owned type or contract
- purple: plugin-owned service or registration type
- yellow: report model
- green: adapter or technical implementation
- gray: framework edge
- dashed border: contract
- dashed arrow: structural relationship such as `used by` or `implemented by`

## Implementation Focus

- add a narrow `InventoryReader` capability with stock snapshots
- add `LowStockItemsReport`
- keep the threshold filtering inside the reporting plugin

Do not add a broader inventory query surface yet.

## What To Verify

- `go test ./...` passes
- items at or below the threshold are included
- the demo can render the low-stock output
