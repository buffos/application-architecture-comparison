# Lesson 027: Low-Stock Items Report

## Objective

Build an application report that identifies StockRecord aggregates at or below a chosen availability threshold.

## Theory

Low-stock reporting is a collection-level question. Each StockRecord owns its availability and low-stock rule, while the report selects and sorts many records for an operational view.

The report consumes the aggregate's public query methods and does not recalculate or mutate stock state.

## Why This Matters Here

The report demonstrates the boundary between aggregate behavior and application composition. StockRecord remains authoritative for available quantity; the report chooses the threshold and presentation shape.

## Diagram

```mermaid
flowchart LR
    STOCK["StockRecord collection"] --> REPORT["LowStockItemsReport"]
    THRESHOLD["application threshold"] --> REPORT
    REPORT --> ITEMS["SKU + available items"]

    classDef source fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef report fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef result fill:#e8eefc,stroke:#3559b5,color:#111;

    class STOCK,THRESHOLD source;
    class REPORT report;
    class ITEMS result;
```

## Implementation Focus

Implement only:

- low-stock report row and report types
- threshold-based selection using StockRecord.Available
- deterministic item ordering
- tests and demo output

Leave replenishment commands and notification delivery for later work.

## What To Verify

- `go test ./...` passes
- records at or below the threshold are included
- records above the threshold are excluded
- report ordering is deterministic
- report generation does not mutate stock aggregates
