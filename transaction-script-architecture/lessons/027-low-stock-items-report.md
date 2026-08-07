# Lesson 027: Build A Low-Stock Report

## Objective

Add a report that identifies products whose available quantity is at or below their reorder threshold.

## Theory

Inventory records keep `OnHand`, `Reserved`, and `ReorderThreshold`. The report calculates `Available = OnHand - Reserved` and returns rows at or below the threshold, optionally enriching them with product names from the catalog.

This is a read-time query over the write model; it does not reserve, receive, or alter stock.

## Why This Matters Here

The procedure is compact and useful for a small application. Its direct joins across `Stocks` and `Products` also show the cost of keeping all coordination in scripts: a catalog rename or a future inventory projection changes the report's assumptions directly.

## Diagram

```mermaid
flowchart LR
    subgraph SCRIPT["internal/scripts"]
        REPORT["GetLowStockItemsReport"]
    end

    subgraph DATA["internal/data"]
        STOCK["Store.Stocks\non-hand / reserved / threshold"]
        PRODUCTS["Store.Products\nname"]
        VIEW["Low-stock rows"]
    end

    REPORT -.reads.-> STOCK
    REPORT -.looks up names.-> PRODUCTS
    REPORT --> VIEW

    classDef script fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef data fill:#fff3bf,stroke:#b08900,color:#111;
    class REPORT script;
    class STOCK,PRODUCTS,VIEW data;
```

Legend:

- purple: report procedure;
- yellow: passive source data and read model;
- dashed arrows: read-only scans/lookups;
- solid arrow: report shaping.

## Implementation Focus

Implement only:

- a low-stock report row;
- `GetLowStockItemsReport`;
- available-quantity calculation and deterministic SKU sorting;
- tests for threshold boundaries and missing product names.

Leave approval queue reporting for the next lesson.

## What To Verify

- `go test ./...` passes from `transaction-script-architecture/`;
- available quantity uses on-hand minus reserved;
- equal-to-threshold items are included;
- healthy stock is excluded;
- the report does not mutate inventory.
