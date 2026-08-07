# Lesson 021: Shipment Query Surface

## Objective

Expose Shipment details and summaries through an application query surface without putting read concerns on the Fulfillment aggregate.

## Theory

Shipment commands protect dispatch state and line invariants. Operational screens usually need a projection keyed by shipment ID: order reference, dispatch status, and line counts.

The query reader translates the rich aggregate into details and summaries. It can later be replaced with a database-backed reader without changing Shipment behavior.

## Why This Matters Here

Fulfillment queries have different shape and sorting needs from Order commands. Keeping the read model outside Shipment avoids a growing aggregate API designed for screens rather than business transitions.

## Diagram

```mermaid
flowchart LR
    SHIPMENT["Shipment aggregate\nPending/Dispatched"] --> READER["application Shipment Reader"]
    READER --> DETAILS["ShipmentDetails"]
    READER --> SUMMARY["ShipmentSummary list"]

    classDef domain fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef application fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef view fill:#e8eefc,stroke:#3559b5,color:#111;

    class SHIPMENT domain;
    class READER application;
    class DETAILS,SUMMARY view;
```

## Implementation Focus

Implement only:

- Shipment detail and summary read types
- an application `Reader` contract
- an in-memory projection with status filtering and copied line views
- tests for save, get, list, and missing shipments
- demo query output

Leave database adapters and fulfillment search concerns for later work.

## What To Verify

- `go test ./...` passes
- Shipment details contain order, status, and lines
- list filtering and deterministic ordering work
- missing shipments return a query-specific error
- the Shipment aggregate remains command-focused
