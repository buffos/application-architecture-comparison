# Lesson 010: Shipment Creation After Payment

## Objective

Complete the first narrow fulfillment path by creating a shipment only after payment has been captured.

## Theory

The Onion track now reaches a paid order.

The next step is fulfillment:

- create a shipment
- mark the order as shipped

This is another useful Onion lesson because it adds a second aggregate around the order workflow while still keeping the responsibilities clean:

- the order owns whether shipment is allowed
- the shipment captures the fulfillment snapshot
- the application ring coordinates repository access

The shipment itself is not created by infrastructure.

It is created in the domain core and then persisted by infrastructure from the outside.

## Why This Matters Here

If shipment creation happens without a rule on the order, then fulfillment eligibility leaks out of the core.

If the application service mutates order state directly, the domain becomes weak again.

The Onion pattern stays the same:

- domain owns the lifecycle transition
- application orchestrates
- infrastructure stores the result

## Diagram

```mermaid
flowchart LR
    subgraph DOM["Domain Core"]
        direction TB
        ORD["Order Entity"]
        SHP["Shipment Entity"]
        FUL["MarkShipped() Transition"]
    end

    subgraph APP["Application Ring"]
        direction TB
        CSS["CreateShipment Service"]
        OFD["Order Finder"]
        OST["Order Store"]
        SST["Shipment Store"]
    end

    subgraph INF["Infrastructure Ring"]
        direction TB
        CLI["CLI Framework"]
        MOR["Memory Order Repository"]
        MSR["Memory Shipment Repository"]
    end

    CLI --> CSS
    CSS --> ORD
    CSS --> SHP
    ORD --> FUL

    OFD -.used by.-> CSS
    OST -.used by.-> CSS
    SST -.used by.-> CSS
    OFD -.implemented by.-> MOR
    OST -.implemented by.-> MOR
    SST -.implemented by.-> MSR

    classDef framework fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef app fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef domain fill:#fff3bf,stroke:#b08900,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class CLI framework;
    class MOR,MSR dataadapter;
    class CSS,OFD,OST,SST app;
    class ORD,SHP,FUL domain;
    class OFD,OST,SST contract;
```

Legend:

- blue: framework edge
- green: data adapter
- purple: application ring
- yellow: domain core
- dashed border: interface / contract
- dashed arrow: structural relationship

## Implementation Focus

Implement one fulfillment workflow:

- create shipment for paid order

The code should show:

- shipment as a domain concept
- an order transition from `Paid` to `Shipped`
- an application service that creates and stores a shipment
- in-memory shipment storage

## What To Verify

- `go test ./...` passes
- paid orders can be shipped
- unpaid orders cannot be shipped
- the order and shipment are both persisted through the application ring
