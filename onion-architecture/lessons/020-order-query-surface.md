# Lesson 020: Order Query Surface

## Objective

Add an explicit read surface for orders so the main fulfillment aggregate has the same application-owned query boundary as returns.

## Theory

The Onion track now has a meaningful order lifecycle:

- conversion
- reservation
- payment
- shipment
- cancellation

That makes order reads important enough to deserve first-class application use cases instead of ad hoc repository access.

As with other Onion query lessons:

- the application ring owns the read use cases
- infrastructure only implements storage lookup
- outer layers depend on the application surface, not on repository details

## Why This Matters Here

If callers read the order repository directly, the fulfillment workflow loses one of the main architectural lessons of the repo:

- reads should cross the same ring boundary intentionally

This lesson keeps that consistent by making order queries explicit and by letting the application ring shape the returned data.

## Diagram

```mermaid
flowchart LR
    subgraph DOM["Domain Core"]
        direction TB
        ORD["Order Entity"]
    end

    subgraph APP["Application Ring"]
        direction TB
        GOD["GetOrder Service"]
        LOD["ListOrders Service"]
        OFD["Order Finder"]
    end

    subgraph INF["Infrastructure Ring"]
        direction TB
        MOR["Memory Order Repository"]
    end

    GOD --> ORD
    LOD --> ORD

    OFD -.used by.-> GOD
    OFD -.used by.-> LOD
    OFD -.implemented by.-> MOR

    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef app fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef domain fill:#fff3bf,stroke:#b08900,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class MOR dataadapter;
    class GOD,LOD,OFD app;
    class ORD domain;
    class OFD contract;
```

## Implementation Focus

Implement two read use cases:

- get order by id
- list orders by status

The code should show:

- an order finder contract in the application ring
- query result models shaped by the application
- in-memory support for filtering by status

## What To Verify

- `go test ./...` passes
- single orders can be loaded by id
- orders can be filtered by status
- reads still flow through the application ring
