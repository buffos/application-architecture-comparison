# Lesson 030: Support Partial Shipments

## Objective

Extend shipment processing so a paid order can ship selected quantities while retaining the remaining fulfillment state.

## Theory

The original `Order.CreateShipment` assumed every reserved line shipped together. This lesson adds `Order.CreatePartialShipment`, which:

1. requires a fulfillable order and shipper;
2. derives all remaining lines when no selection is supplied;
3. validates selected line quantities against remaining reservations;
4. creates a shipment for only the selected quantities;
5. consumes the matching stock;
6. sets the order to `PartiallyShipped` or `Shipped` based on what remains.

`CreateShipment` becomes a convenience operation that delegates to the same Active Record behavior.

## Why This Matters Here

The change is small in code but meaningful in the model: line-level bookkeeping now influences aggregate status. Active Record keeps the coordination on `Order`, making the growing fulfillment responsibility easy to see.

## Diagram

```mermaid
flowchart LR
    FULL["Order.CreateShipment"] --> PARTIAL["Order.CreatePartialShipment"]
    PARTIAL -.reads.-> ORDER["Order Active Record"]
    PARTIAL -.checks.-> STOCK["StockRecord"]
    PARTIAL --> SHIPMENT["Shipment Active Record"]
    PARTIAL --> STATE["PartiallyShipped\nor Shipped"]

    classDef operation fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef record fill:#fff3bf,stroke:#b08900,color:#111;

    class FULL,PARTIAL operation;
    class ORDER,STOCK,SHIPMENT,STATE record;
```

Legend:

- purple: Active Record fulfillment operations;
- yellow: persisted records and resulting state;
- dashed arrows: reads and preflight checks;
- solid arrows: shipment and status writes.

## Implementation Focus

Implement only:

- `CreatePartialShipment`;
- line-level quantity validation;
- partial versus complete order status;
- full-shipment delegation through the same model operation;
- tests for partial, complete, and invalid shipments.

Leave partial returns for the next lesson.

## What To Verify

- `go test ./...` passes from `active-record-architecture/`;
- a subset of reserved quantities can be shipped;
- the order becomes `PartiallyShipped` when work remains;
- a later shipment can finish the order;
- stock is consumed only for selected quantities.
