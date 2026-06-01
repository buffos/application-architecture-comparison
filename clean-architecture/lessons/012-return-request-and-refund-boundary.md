# Lesson 012: Return Request And Refund Boundary

## Objective

Add the first post-shipment reverse workflow by introducing return requests and a refund gateway boundary.

## Theory

Cancellation and returns are not the same thing.

Cancellation happens before fulfillment is complete.

Returns happen after shipment and usually involve a different business path:

- the order has already moved forward
- the customer is asking to reverse part of that outcome
- refund is not an entity-only decision because it crosses a payment boundary

That makes returns a useful Clean Architecture lesson because the workflow now has to coordinate:

- loading the order
- checking that the order is in a returnable state
- creating a return request entity
- calling a refund gateway
- saving the return request

This reinforces a key Clean point:

- the application layer owns workflow sequencing
- the entities own local validity rules
- external money movement stays behind a contract

The tradeoff is another entity, another gateway, and another write path.

## Why This Matters Here

The order lifecycle is now complete enough that the next useful step is not another forward action.

It is a post-fulfillment reverse path with different constraints from cancellation.

That makes the difference between:

- pre-shipment reversal
- post-shipment reversal

architecturally visible instead of only conceptually mentioned.

## Diagram

```mermaid
flowchart TD
    subgraph INFRA[Infrastructure / Frameworks]
        CLI[CLI Framework]
        MOG[Memory Order Gateway]
        MRG[Memory Return Request Gateway]
        RFG[Accept-All Refund Gateway]
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
    OUTB -.used by.-> UC
    OUTB -.implemented by.-> PRES

    classDef framework fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef funcadapter fill:#ffe5d9,stroke:#bc6c25,color:#111;
    classDef app fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef entity fill:#fff3bf,stroke:#b08900,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class CLI framework;
    class MOG,MRG dataadapter;
    class RFG,CTRL,PRES funcadapter;
    class IN,INB,UC,OG,RRG,REF,OUT,OUTB app;
    class ORDER,RET entity;
    class INB,OG,RRG,REF,OUTB contract;
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

Implement one use case:

- request a return for a shipped order and issue a refund

The code should show:

- a `ReturnRequest` entity
- entity validation that only shipped orders can be returned
- a refund gateway contract
- a return request gateway contract and in-memory adapter
- the CLI demo path staying unchanged while tests cover the return workflow

Do not add return review or restocking yet.

## What To Verify

- the project compiles
- `go test ./...` passes
- a shipped order can produce a return request
- a non-shipped order cannot be returned
