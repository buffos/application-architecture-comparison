# Lesson 020: Order Query Surface

## Objective

Add explicit read-side use cases for orders so the order slice has the same Clean Architecture query seam as quotes and return requests.

## Theory

Orders already have a substantial write workflow:

- conversion from quote
- payment capture
- shipment creation
- cancellation

But without dedicated query use cases, reads still happen only as internal support for write-side interactors.

Clean Architecture treats order reads as application behavior too.

That means a caller should not jump straight from a controller into a concrete order gateway.

Instead, the read path should still go through:

- input boundary
- interactor
- output boundary
- presenter

The benefit is consistency.

The application layer decides:

- which order data is exposed
- which filters exist
- which read scenarios are officially supported

The tradeoff is more code for operations that may seem simple.

## Why This Matters Here

The return query lesson showed that reads are first-class in Clean Architecture.

Orders are the next natural slice because they sit in the center of the workflow and are already reused by many write-side use cases.

This lesson makes that central concept readable through the same boundary structure rather than only through internal collaborators.

## Diagram

```mermaid
flowchart LR
    subgraph ENT["Entities"]
        direction TB
        ORD["Order Entity"]
    end

    subgraph APP["Application"]
        direction TB
        GIN["GetOrder Input Boundary"]
        GOUT["GetOrder Output Boundary"]
        LIN["ListOrders Input Boundary"]
        LOUT["ListOrders Output Boundary"]
        GUC["GetOrder Interactor"]
        LUC["ListOrders Interactor"]
        ODR["Order Reader"]
        ODL["Order Lister"]
    end

    subgraph IA["Interface Adapters"]
        direction TB
        GCTRL["GetOrder Controller"]
        LCTRL["ListOrders Controller"]
        GPRES["GetOrder Presenter"]
        LPRES["ListOrders Presenter"]
    end

    subgraph INFRA["Infrastructure / Frameworks"]
        direction TB
        CLI["CLI / HTTP Framework"]
        MOG["Memory Order Gateway"]
    end

    CLI --> GCTRL
    CLI --> LCTRL
    GCTRL --> GIN
    LCTRL --> LIN
    GUC --> GOUT
    LUC --> LOUT
    GPRES --> CLI
    LPRES --> CLI
    GUC --> ORD
    LUC --> ORD

    GIN -.used by.-> GCTRL
    GIN -.implemented by.-> GUC
    GOUT -.used by.-> GUC
    GOUT -.implemented by.-> GPRES
    LIN -.used by.-> LCTRL
    LIN -.implemented by.-> LUC
    LOUT -.used by.-> LUC
    LOUT -.implemented by.-> LPRES
    ODR -.used by.-> GUC
    ODL -.used by.-> LUC
    ODR -.implemented by.-> MOG
    ODL -.implemented by.-> MOG

    classDef framework fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef funcadapter fill:#ffe5d9,stroke:#bc6c25,color:#111;
    classDef app fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef entity fill:#fff3bf,stroke:#b08900,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class CLI framework;
    class MOG dataadapter;
    class GCTRL,LCTRL,GPRES,LPRES funcadapter;
    class GIN,GOUT,LIN,LOUT,GUC,LUC,ODR,ODL app;
    class ORD entity;
    class GIN,GOUT,LIN,LOUT,ODR,ODL contract;
```

Legend:

- blue: framework edge
- green: data adapter
- orange: translation adapter
- purple: application layer
- yellow: entity layer
- dashed border: interface / contract
- dashed arrow: structural relationship such as `used by` or `implemented by`

## Implementation Focus

Add:

- `GetOrder`
- `ListOrders`

The code should show:

- a single-order query use case
- a list-by-status query use case
- the order gateway implementing reader and lister contracts
- presenters shaping the result for callers instead of exposing raw entities

## What To Verify

- the project compiles
- `go test ./...` passes
- an order can be loaded through a query interactor
- paid orders can be listed by status
