# Lesson 022: Quote List Query Surface

## Objective

Add list-by-status querying for quotes so the quote slice has the same explicit read-side surface as orders, returns, and shipments.

## Theory

Quotes already introduced the first Clean Architecture query lesson through:

- `GetQuote`

But that still leaves the quote read side narrower than the other main workflow objects.

Clean Architecture treats query breadth as application behavior too.

The question is not only:

- can we load one quote?

It is also:

- which quote listing scenarios does the application officially support?

This lesson keeps the same boundary pattern:

- controller
- input boundary
- interactor
- gateway contract
- output boundary
- presenter

The benefit is consistency across the architecture.

The application layer, not infrastructure, decides that listing quotes by status is a supported use case and how that result is shaped for callers.

## Why This Matters Here

Quotes start the whole workflow, so they should not end up with a weaker read model than downstream objects.

Adding `ListQuotes` also completes the comparison with the recent order, shipment, and return query lessons and makes the quote slice feel like a full application surface instead of a one-off example.

## Diagram

```mermaid
flowchart LR
    subgraph ENT["Entities"]
        direction TB
        QTE["Quote Entity"]
    end

    subgraph APP["Application"]
        direction TB
        GIN["GetQuote Input Boundary"]
        GOUT["GetQuote Output Boundary"]
        LIN["ListQuotes Input Boundary"]
        LOUT["ListQuotes Output Boundary"]
        GUC["GetQuote Interactor"]
        LUC["ListQuotes Interactor"]
        QTR["Quote Reader"]
        QTL["Quote Lister"]
    end

    subgraph IA["Interface Adapters"]
        direction TB
        GCTRL["GetQuote Controller"]
        LCTRL["ListQuotes Controller"]
        GPRES["GetQuote Presenter"]
        LPRES["ListQuotes Presenter"]
    end

    subgraph INFRA["Infrastructure / Frameworks"]
        direction TB
        CLI["CLI / HTTP Framework"]
        MQG["Memory Quote Gateway"]
    end

    CLI --> GCTRL
    CLI --> LCTRL
    GCTRL --> GIN
    LCTRL --> LIN
    GUC --> GOUT
    LUC --> LOUT
    GPRES --> CLI
    LPRES --> CLI
    GUC --> QTE
    LUC --> QTE

    GIN -.used by.-> GCTRL
    GIN -.implemented by.-> GUC
    GOUT -.used by.-> GUC
    GOUT -.implemented by.-> GPRES
    LIN -.used by.-> LCTRL
    LIN -.implemented by.-> LUC
    LOUT -.used by.-> LUC
    LOUT -.implemented by.-> LPRES
    QTR -.used by.-> GUC
    QTL -.used by.-> LUC
    QTR -.implemented by.-> MQG
    QTL -.implemented by.-> MQG

    classDef framework fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef funcadapter fill:#ffe5d9,stroke:#bc6c25,color:#111;
    classDef app fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef entity fill:#fff3bf,stroke:#b08900,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class CLI framework;
    class MQG dataadapter;
    class GCTRL,LCTRL,GPRES,LPRES funcadapter;
    class GIN,GOUT,LIN,LOUT,GUC,LUC,QTR,QTL app;
    class QTE entity;
    class GIN,GOUT,LIN,LOUT,QTR,QTL contract;
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

- `ListQuotes`

The code should show:

- list-by-status as an explicit quote read use case
- the quote gateway implementing a lister contract in addition to single-item lookup
- a presenter shaping quote list results for callers

## What To Verify

- the project compiles
- `go test ./...` passes
- approved quotes can be listed by status
- the existing single-quote query flow still works
