# Lesson 031: Partial Returns By Line

## Objective

Make returns quantity-aware so a return request can cover only some shipped units instead of always implying the whole order.

## Theory

Partial shipment support created an asymmetry:

- fulfillment can now happen in slices
- returns still behave like all shipped quantity must come back together

That is too coarse.

A realistic reverse flow often needs:

- return 1 of 3 shipped units
- return one SKU but not another
- make another return later for the remaining shipped quantity

Clean Architecture handles this by keeping the quantity rules in the inner layers:

- the request use case accepts explicit return lines
- the return request entity snapshots those lines
- the order entity tracks how much has already been returned
- acceptance updates order return progress and restocks only the accepted lines

The important change is that a return request is no longer just:

- order id
- reason

It is now also:

- the specific line quantities being returned

## Why This Matters Here

This is the natural counterpart to partial shipment.

Without it, the forward flow becomes more realistic while the reverse flow stays artificially simple.

This lesson brings the two halves back into alignment and makes the order entity carry more of the lifecycle truth around what has shipped versus what has been returned.

## Diagram

```mermaid
flowchart LR
    subgraph ENT["Entities"]
        direction TB
        ORD["Order Entity<br/>tracks shipped and returned quantity"]
        RET["Return Request Entity<br/>tracks requested return lines"]
    end

    subgraph APP["Application"]
        direction TB
        RIN["RequestReturn Input Boundary"]
        ROUT["RequestReturn Output Boundary"]
        AIN["AcceptReturn Input Boundary"]
        AOUT["AcceptReturn Output Boundary"]
        RUC["RequestReturn Interactor"]
        AUC["AcceptReturn Interactor"]
        OED["Order Editor"]
        RRE["Return Request Editor"]
        REF["Refund Gateway"]
        RES["Inventory Restock"]
    end

    subgraph IA["Interface Adapters"]
        direction TB
        RCTRL["RequestReturn Controller"]
        ACTRL["AcceptReturn Controller"]
        RPRES["RequestReturn Presenter"]
        APRES["AcceptReturn Presenter"]
    end

    subgraph INFRA["Infrastructure / Frameworks"]
        direction TB
        CLI["CLI / HTTP Framework"]
        MOG["Memory Order Gateway"]
        MRG["Memory Return Request Gateway"]
        RFG["Refund Gateway"]
        MIR["Memory Inventory Reservation"]
    end

    CLI --> RCTRL
    CLI --> ACTRL
    RCTRL --> RIN
    ACTRL --> AIN
    RUC --> ROUT
    AUC --> AOUT
    RPRES --> CLI
    APRES --> CLI
    RUC --> ORD
    RUC --> RET
    AUC --> ORD
    AUC --> RET

    RIN -.used by.-> RCTRL
    RIN -.implemented by.-> RUC
    ROUT -.used by.-> RUC
    ROUT -.implemented by.-> RPRES
    AIN -.used by.-> ACTRL
    AIN -.implemented by.-> AUC
    AOUT -.used by.-> AUC
    AOUT -.implemented by.-> APRES
    OED -.used by.-> RUC
    OED -.used by.-> AUC
    RRE -.used by.-> RUC
    RRE -.used by.-> AUC
    REF -.used by.-> AUC
    RES -.used by.-> AUC
    OED -.implemented by.-> MOG
    RRE -.implemented by.-> MRG
    REF -.implemented by.-> RFG
    RES -.implemented by.-> MIR

    classDef framework fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef funcadapter fill:#ffe5d9,stroke:#bc6c25,color:#111;
    classDef app fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef entity fill:#fff3bf,stroke:#b08900,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class CLI framework;
    class MOG,MRG,MIR dataadapter;
    class RFG,RCTRL,ACTRL,RPRES,APRES funcadapter;
    class RIN,ROUT,AIN,AOUT,RUC,AUC,OED,RRE,REF,RES app;
    class ORD,RET entity;
    class RIN,ROUT,AIN,AOUT,OED,RRE,REF,RES contract;
```

Legend:

- blue: framework edge
- green: data adapter
- orange: translation or service adapter
- purple: application layer
- yellow: entity layer
- dashed border: interface / contract
- dashed arrow: structural relationship such as `used by` or `implemented by`

## Implementation Focus

Add:

- explicit return line input
- return request line snapshots
- returned quantity tracking on order lines
- restock and return accounting based only on the accepted return lines

The code should show:

- requesting only some shipped quantity
- accepting a partial return updates order return progress
- later returns cannot exceed what has already shipped minus what was already returned

## What To Verify

- the project compiles
- `go test ./...` passes
- a partial return request stores only the requested line quantities
- accepting a partial return restocks only those quantities
- returned quantity on the order increases correctly
