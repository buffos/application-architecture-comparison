# Lesson 008: Order Conversion With Reservation

## Objective

Extend quote-to-order conversion so the application ring also reserves inventory before persisting the order.

## Theory

The previous lesson showed a cross-aggregate workflow:

- load approved quote
- create order
- save order

Real workflows often need one more thing:

- a side effect in an external operational subsystem

Onion Architecture handles that by keeping the domain unchanged and letting the application ring coordinate an additional inward-facing contract.

In this lesson:

- the domain still creates the order from the quote
- the application ring translates order lines into reservation items
- infrastructure provides the inventory implementation

This keeps the responsibility split clean:

- the domain decides what an order is
- the application ring decides which collaborators the workflow needs
- infrastructure performs the external stock operation

## Why This Matters Here

Without reservation, conversion is still only a document handoff.

Reservation makes the workflow operational:

- an order now claims stock
- failure in stock reservation blocks order creation
- the application ring coordinates both business and operational steps

That makes the Onion boundary around external services more visible.

## Diagram

```mermaid
flowchart LR
    subgraph DOM["Domain Core"]
        direction TB
        QTE["Quote Entity"]
        ORD["Order Entity"]
    end

    subgraph APP["Application Ring"]
        direction TB
        CTO["ConvertQuoteToOrder Service"]
        QFD["Quote Finder"]
        OST["Order Store"]
        RSV["Inventory Reservation"]
    end

    subgraph INF["Infrastructure Ring"]
        direction TB
        CLI["CLI Framework"]
        MQR["Memory Quote Repository"]
        MOR["Memory Order Repository"]
        MIR["Memory Inventory Reservation"]
    end

    CLI --> CTO
    CTO --> QTE
    CTO --> ORD

    QFD -.used by.-> CTO
    OST -.used by.-> CTO
    RSV -.used by.-> CTO
    QFD -.implemented by.-> MQR
    OST -.implemented by.-> MOR
    RSV -.implemented by.-> MIR

    classDef framework fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef app fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef domain fill:#fff3bf,stroke:#b08900,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class CLI framework;
    class MQR,MOR,MIR dataadapter;
    class CTO,QFD,OST,RSV app;
    class QTE,ORD domain;
    class QFD,OST,RSV contract;
```

Legend:

- blue: framework edge
- green: data adapter
- purple: application ring
- yellow: domain core
- dashed border: interface / contract
- dashed arrow: structural relationship

## Implementation Focus

Implement one operational extension:

- reserve inventory during quote-to-order conversion

The code should show:

- reservation item types in the domain layer
- an inventory reservation contract in the application ring
- an in-memory reservation adapter
- conversion failing when stock is insufficient

## What To Verify

- `go test ./...` passes
- approved quotes reserve stock when converted
- insufficient stock blocks conversion
- the order is only saved after successful reservation
