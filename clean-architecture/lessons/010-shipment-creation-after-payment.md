# Lesson 010: Shipment Creation After Payment

## Objective

Add shipment creation after payment, so the Clean Architecture track now shows the first end-to-end order fulfillment gate: reservation, payment, then shipment.

## Theory

By this point the order can be:

- created from an approved quote
- reserved against inventory
- marked as paid

The next natural step is shipment.

This lesson is useful because shipment is not just another order status update.

It introduces:

- a second entity on the order side
- a new persistence boundary
- an additional state rule on the order

The use case now coordinates:

- loading the order
- asking the order entity whether shipment is allowed
- creating a shipment entity
- saving the shipment
- updating the order state

This is exactly the kind of flow the application layer exists to coordinate in Clean Architecture.

The entity owns shipment eligibility.

The shipment entity owns shipment data.

The interactor owns the sequencing across those concepts and boundaries.

The tradeoff is another gateway and another use case to wire.

## Why This Matters Here

This lesson completes the first narrow happy path of the sample application:

- draft quote
- add line
- submit / approve
- convert to order
- reserve stock
- capture payment
- create shipment

That makes later reverse flows like cancellation and returns easier to introduce because the forward lifecycle is now visible.

## Diagram

```mermaid
flowchart TD
    subgraph INFRA[Infrastructure / Frameworks]
        CLI[CLI Framework]
        MOG[Memory Order Gateway]
        MSG[Memory Shipment Gateway]
    end

    subgraph IA[Interface Adapters]
        CTRL[CreateShipment Controller]
        PRES[CreateShipment Presenter]
    end

    subgraph APP[Application]
        IN[CreateShipment Input Model]
        INB[CreateShipment Input Boundary]
        UC[CreateShipment Interactor]
        OG[Order Editor Gateway]
        SG[Shipment Writer Gateway]
        OUT[CreateShipment Output Model]
        OUTB[CreateShipment Output Boundary]
    end

    subgraph ENT[Entities]
        ORDER[Order Entity]
        SHIP[Shipment Entity]
    end

    CLI --> CTRL
    CTRL --> IN
    IN --> INB
    INB --> UC
    OG --> UC
    SG --> UC
    UC --> ORDER
    UC --> SHIP
    UC --> OUTB
    OUTB --> OUT
    OUT --> PRES
    INB -.used by.-> CTRL
    INB -.implements.-> UC
    OG -.used by.-> UC
    OG -.implemented by.-> MOG
    SG -.used by.-> UC
    SG -.implemented by.-> MSG
    OUTB -.used by.-> UC
    OUTB -.implemented by.-> PRES

    classDef framework fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef funcadapter fill:#ffe5d9,stroke:#bc6c25,color:#111;
    classDef app fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef entity fill:#fff3bf,stroke:#b08900,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class CLI framework;
    class MOG,MSG dataadapter;
    class CTRL,PRES funcadapter;
    class IN,INB,UC,OG,SG,OUT,OUTB app;
    class ORDER,SHIP entity;
    class INB,OG,SG,OUTB contract;
```

Legend:

- blue: framework edge
- green: data adapter
- orange: functionality / translation adapter
- purple: application layer
- yellow: entity layer
- dashed border: interface / contract
- dashed arrow: structural relationship

## Implementation Focus

Implement one use case:

- create a shipment for a paid order

The code should show:

- a `Shipment` entity
- a shipped order status
- entity validation that only paid orders can be shipped
- a shipment gateway contract and in-memory adapter
- the CLI demo creating a shipment after payment

Do not add partial shipment or shipment queries yet.

## What To Verify

- the project compiles
- `go test ./...` passes
- a paid order can be shipped
- an unpaid order cannot be shipped
