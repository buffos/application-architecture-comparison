# Lesson 024: Customer Query Surface

## Objective

Promote customers from a supporting lookup dependency into an explicit read-side application surface.

## Theory

So far, customers mostly appear as validation collaborators:

- load customer by id
- ensure the customer is active

That is enough for quote creation, but it keeps customer reads hidden inside other workflows.

Clean Architecture treats those reads as application behavior too.

So instead of letting outer layers depend directly on the customer gateway, the application layer owns:

- which customer queries are supported
- how customer filters are expressed
- how customer data is shaped for callers

This lesson uses two simple read scenarios:

- get customer by id
- list customers with an active-only filter

The tradeoff is the usual read-side ceremony:

- more interfaces
- more mapping
- more small types

## Why This Matters Here

Customers are one of the first entities touched in the workflow, so keeping them as a pure helper would leave the Clean track uneven.

This lesson balances the supporting-entity story with the catalog query lesson and makes it easier to compare how different architectures expose foundational data, not only workflow state.

## Diagram

```mermaid
flowchart LR
    subgraph ENT["Entities"]
        direction TB
        CUS["Customer Entity"]
    end

    subgraph APP["Application"]
        direction TB
        GIN["GetCustomer Input Boundary"]
        GOUT["GetCustomer Output Boundary"]
        LIN["ListCustomers Input Boundary"]
        LOUT["ListCustomers Output Boundary"]
        GUC["GetCustomer Interactor"]
        LUC["ListCustomers Interactor"]
        CDR["Customer Reader"]
        CDL["Customer Lister"]
    end

    subgraph IA["Interface Adapters"]
        direction TB
        GCTRL["GetCustomer Controller"]
        LCTRL["ListCustomers Controller"]
        GPRES["GetCustomer Presenter"]
        LPRES["ListCustomers Presenter"]
    end

    subgraph INFRA["Infrastructure / Frameworks"]
        direction TB
        CLI["CLI / HTTP Framework"]
        MCG["Memory Customer Gateway"]
    end

    CLI --> GCTRL
    CLI --> LCTRL
    GCTRL --> GIN
    LCTRL --> LIN
    GUC --> GOUT
    LUC --> LOUT
    GPRES --> CLI
    LPRES --> CLI
    GUC --> CUS
    LUC --> CUS

    GIN -.used by.-> GCTRL
    GIN -.implemented by.-> GUC
    GOUT -.used by.-> GUC
    GOUT -.implemented by.-> GPRES
    LIN -.used by.-> LCTRL
    LIN -.implemented by.-> LUC
    LOUT -.used by.-> LUC
    LOUT -.implemented by.-> LPRES
    CDR -.used by.-> GUC
    CDL -.used by.-> LUC
    CDR -.implemented by.-> MCG
    CDL -.implemented by.-> MCG

    classDef framework fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef funcadapter fill:#ffe5d9,stroke:#bc6c25,color:#111;
    classDef app fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef entity fill:#fff3bf,stroke:#b08900,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class CLI framework;
    class MCG dataadapter;
    class GCTRL,LCTRL,GPRES,LPRES funcadapter;
    class GIN,GOUT,LIN,LOUT,GUC,LUC,CDR,CDL app;
    class CUS entity;
    class GIN,GOUT,LIN,LOUT,CDR,CDL contract;
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

- `GetCustomer`
- `ListCustomers`

The code should show:

- a single-customer query use case
- an active-only list query use case
- the customer gateway implementing reader and lister contracts
- presenters shaping customer read models for callers

## What To Verify

- the project compiles
- `go test ./...` passes
- a customer can be loaded through a query interactor
- active customers can be listed
