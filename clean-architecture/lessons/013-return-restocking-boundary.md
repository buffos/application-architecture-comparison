# Lesson 013: Return Restocking Boundary

## Objective

Extend the return workflow so refunded returns also restock inventory through a dedicated contract.

## Theory

The previous lesson established the money side of a return:

- shipped order
- refund gateway
- return request record

But one side effect was still missing:

- the stock that left during fulfillment has not been put back

That makes restocking the natural next step.

This is useful architecturally because it shows that one use case can coordinate multiple external side effects without pushing those details into the entity:

- refund money
- restock inventory
- save the return request

The return request entity still owns its own local meaning.

The restock operation stays behind an application-owned contract.

The interactor owns the sequencing.

The tradeoff is another boundary and another operational dependency in the same workflow.

## Why This Matters Here

Without restocking, the reverse workflow is incomplete from an inventory perspective.

This lesson closes that gap and makes the post-shipment reversal more realistic:

- cancellation releases reservation
- returns refund money and restock inventory

That distinction is important for understanding how the architecture treats similar-but-different workflows.

## Diagram

```mermaid
flowchart TD
    subgraph INFRA[Infrastructure / Frameworks]
        CLI[CLI Framework]
        MOG[Memory Order Gateway]
        MRG[Memory Return Request Gateway]
        RFG[Accept-All Refund Gateway]
        MIR[Memory Inventory Reservation]
    end

    subgraph IA[Interface Adapters]
        CTRL[RequestReturn Controller]
        PRES[RequestReturn Presenter]
    end

    subgraph APP[Application]
        IN[RequestReturn Input Model]
        INB[RequestReturn Input Boundary]
        UC[RequestReturn Interactor]
        OG[Order Reader Gateway]
        RRG[Return Request Writer Gateway]
        REF[Refund Gateway Contract]
        RES[Inventory Restock Contract]
        OUT[RequestReturn Output Model]
        OUTB[RequestReturn Output Boundary]
    end

    subgraph ENT[Entities]
        ORDER[Order Entity]
        RET[Return Request Entity]
    end

    CLI --> CTRL
    CTRL --> IN
    IN --> INB
    INB --> UC
    OG --> UC
    RRG --> UC
    REF --> UC
    RES --> UC
    UC --> ORDER
    UC --> RET
    UC --> OUTB
    OUTB --> OUT
    OUT --> PRES
    INB -.used by.-> CTRL
    INB -.implements.-> UC
    OG -.used by.-> UC
    OG -.implemented by.-> MOG
    RRG -.used by.-> UC
    RRG -.implemented by.-> MRG
    REF -.used by.-> UC
    REF -.implemented by.-> RFG
    RES -.used by.-> UC
    RES -.implemented by.-> MIR
    OUTB -.used by.-> UC
    OUTB -.implemented by.-> PRES

    classDef framework fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef funcadapter fill:#ffe5d9,stroke:#bc6c25,color:#111;
    classDef app fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef entity fill:#fff3bf,stroke:#b08900,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class CLI framework;
    class MOG,MRG,MIR dataadapter;
    class RFG,CTRL,PRES funcadapter;
    class IN,INB,UC,OG,RRG,REF,RES,OUT,OUTB app;
    class ORDER,RET entity;
    class INB,OG,RRG,REF,RES,OUTB contract;
```

Legend:

- blue: framework edge
- green: data adapter
- orange: functionality / policy / translation adapter
- purple: application layer
- yellow: entity layer
- dashed border: interface / contract
- dashed arrow: structural relationship

## Implementation Focus

Extend one existing use case:

- request a return, refund the order, and restock inventory

The code should show:

- an inventory restock contract
- restock items derived from the order lines
- the existing in-memory inventory adapter implementing restock
- the return use case coordinating refund, restock, and persistence

Do not add return review or partial returns yet.

## What To Verify

- the project compiles
- `go test ./...` passes
- a shipped order return refunds and restocks
- a non-shipped order still cannot be returned
