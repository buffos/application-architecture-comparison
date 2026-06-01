# Lesson 007: Convert Quote To Order

## Objective

Add the first cross-entity workflow in the Clean Architecture track by converting an approved quote into an order through a dedicated use case.

## Theory

Until now, every lesson has stayed inside the quote workflow.

That has been useful for making boundaries visible, but it does not yet show how Clean Architecture handles a workflow that spans two business concepts.

Converting a quote to an order is the first good example.

The use case must:

- load the quote
- verify the quote is in a convertible state
- create an order entity from the quote data
- save the new order

This is where the value of the application layer becomes more concrete.

The interactor is not just calling one repository method.

It is coordinating a business workflow across boundaries while leaving entity-specific rules inside the entities themselves.

The tradeoff is more gateways, more models, and more orchestration code in the use case layer.

## Why This Matters Here

The sample application is not only about quote management.

It eventually needs:

- orders
- payment
- shipment
- returns

Before reaching those later workflows, the architecture needs one clean example of crossing from one aggregate or entity concept into another.

Quote-to-order conversion is the natural next step.

## Diagram

```mermaid
flowchart TD
    subgraph INFRA[Infrastructure / Frameworks]
        CLI[CLI Framework]
        MQG[Memory Quote Gateway]
        MOG[Memory Order Gateway]
    end

    subgraph IA[Interface Adapters]
        CTRL[ConvertQuoteToOrder Controller]
        PRES[ConvertQuoteToOrder Presenter]
    end

    subgraph APP[Application]
        IN[ConvertQuoteToOrder Input Model]
        INB[ConvertQuoteToOrder Input Boundary]
        UC[ConvertQuoteToOrder Interactor]
        QG[Quote Reader Gateway]
        OG[Order Writer Gateway]
        OUT[ConvertQuoteToOrder Output Model]
        OUTB[ConvertQuoteToOrder Output Boundary]
    end

    subgraph ENT[Entities]
        QUOTE[Quote Entity]
        ORDER[Order Entity]
    end

    CLI --> CTRL
    CTRL --> IN
    IN --> INB
    INB --> UC
    QG --> UC
    OG --> UC
    UC --> QUOTE
    UC --> ORDER
    UC --> OUTB
    OUTB --> OUT
    OUT --> PRES
    INB -.used by.-> CTRL
    INB -.implements.-> UC
    QG -.used by.-> UC
    QG -.implemented by.-> MQG
    OG -.used by.-> UC
    OG -.implemented by.-> MOG
    OUTB -.used by.-> UC
    OUTB -.implemented by.-> PRES

    classDef framework fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef funcadapter fill:#ffe5d9,stroke:#bc6c25,color:#111;
    classDef app fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef entity fill:#fff3bf,stroke:#b08900,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class CLI framework;
    class MQG,MOG dataadapter;
    class CTRL,PRES funcadapter;
    class IN,INB,UC,QG,OG,OUT,OUTB app;
    class QUOTE,ORDER entity;
    class INB,QG,OG,OUTB contract;
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

- convert an approved quote into an order

The code should show:

- an `Order` entity
- entity validation that only approved quotes can be converted
- an order gateway contract and in-memory adapter
- a `ConvertQuoteToOrder` interactor
- a controller and presenter for the new workflow
- the CLI demo creating, editing, submitting, and converting a quote

Do not add inventory reservation or payment yet.

## What To Verify

- the project compiles
- `go test ./...` passes
- an approved quote can become an order
- a non-approved quote cannot be converted
