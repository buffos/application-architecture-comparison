# Lesson 021: Shipment Query Surface

## Objective

Add explicit read-side use cases for shipments so the fulfillment slice has the same query boundary shape as quotes, returns, and orders.

## Theory

Shipments already exist on the write side:

- an order is paid
- a shipment is created
- the order moves to shipped

But without read use cases, shipment access still lives only inside infrastructure.

Clean Architecture treats shipment queries as application behavior too.

That means even a simple shipment lookup should still pass through:

- input boundary
- interactor
- output boundary
- presenter

This lesson also shows that query filters do not have to copy the write-side status pattern exactly.

For shipments, the natural list filter is:

- by `OrderID`

The benefit is that the application layer still owns:

- which shipment queries are allowed
- how shipment data is shaped
- how callers depend on shipment reads

The tradeoff is more code around a small read path.

## Why This Matters Here

Orders and returns already have explicit query seams.

Shipments are the missing workflow object in that same read-side pattern.

Adding this lesson makes the fulfillment path easier to compare across architectures because shipment reads are now visible as first-class Clean use cases instead of hidden infrastructure access.

## Diagram

```mermaid
flowchart LR
    subgraph ENT["Entities"]
        direction TB
        SHP["Shipment Entity"]
    end

    subgraph APP["Application"]
        direction TB
        GIN["GetShipment Input Boundary"]
        GOUT["GetShipment Output Boundary"]
        LIN["ListShipments Input Boundary"]
        LOUT["ListShipments Output Boundary"]
        GUC["GetShipment Interactor"]
        LUC["ListShipments Interactor"]
        SHR["Shipment Reader"]
        SHL["Shipment Lister"]
    end

    subgraph IA["Interface Adapters"]
        direction TB
        GCTRL["GetShipment Controller"]
        LCTRL["ListShipments Controller"]
        GPRES["GetShipment Presenter"]
        LPRES["ListShipments Presenter"]
    end

    subgraph INFRA["Infrastructure / Frameworks"]
        direction TB
        CLI["CLI / HTTP Framework"]
        MSG["Memory Shipment Gateway"]
    end

    CLI --> GCTRL
    CLI --> LCTRL
    GCTRL --> GIN
    LCTRL --> LIN
    GUC --> GOUT
    LUC --> LOUT
    GPRES --> CLI
    LPRES --> CLI
    GUC --> SHP
    LUC --> SHP

    GIN -.used by.-> GCTRL
    GIN -.implemented by.-> GUC
    GOUT -.used by.-> GUC
    GOUT -.implemented by.-> GPRES
    LIN -.used by.-> LCTRL
    LIN -.implemented by.-> LUC
    LOUT -.used by.-> LUC
    LOUT -.implemented by.-> LPRES
    SHR -.used by.-> GUC
    SHL -.used by.-> LUC
    SHR -.implemented by.-> MSG
    SHL -.implemented by.-> MSG

    classDef framework fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef funcadapter fill:#ffe5d9,stroke:#bc6c25,color:#111;
    classDef app fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef entity fill:#fff3bf,stroke:#b08900,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class CLI framework;
    class MSG dataadapter;
    class GCTRL,LCTRL,GPRES,LPRES funcadapter;
    class GIN,GOUT,LIN,LOUT,GUC,LUC,SHR,SHL app;
    class SHP entity;
    class GIN,GOUT,LIN,LOUT,SHR,SHL contract;
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

- `GetShipment`
- `ListShipments`

The code should show:

- a single-shipment query use case
- a list-by-order query use case
- the shipment gateway implementing reader and lister contracts
- presenters shaping shipment read models for callers

## What To Verify

- the project compiles
- `go test ./...` passes
- a shipment can be loaded through a query interactor
- shipments can be listed by order id
