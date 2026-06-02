# Lesson 007: Convert Quote To Order

## Objective

Add the first cross-module workflow that turns an approved quote into an order owned by a new `orders` module.

## Theory

The Modular Monolith track now has a meaningful quote lifecycle:

- draft
- pending approval
- approved

The next step is to use an approved quote as input for another module.

This is an important modular-monolith lesson because it shows a different kind of module interaction:

- `quotes` still owns quote lifecycle
- `orders` owns order creation and order storage
- `orders` depends on a narrow quote-read capability instead of reaching into quote persistence directly

The business workflow crosses modules, but module ownership stays clear.

## Why This Matters Here

Without a cross-module workflow, the architecture still only proves that each module can manage its own local behavior.

Conversion makes the module boundaries more realistic:

- one approved business document becomes another
- the `orders` module coordinates the handoff
- the `quotes` module provides a stable, narrow API

That is the first point where the modular monolith starts to show why module APIs matter as much as internal design.

## Diagram

```mermaid
flowchart LR
    subgraph QTM["Quotes Module"]
        direction TB
        QAP["quotes.ApprovedQuoteSource"]
        QMS["quotes.Service"]
    end

    subgraph ORM["Orders Module"]
        direction TB
        ORD["Order"]
        ORE["orders.Repository"]
        OMS["orders.Service<br/>ConvertQuoteToOrder"]
    end

    subgraph INF["Infrastructure"]
        direction TB
        CLI["CLI"]
        MOR["Memory Order Repository"]
        MQR["Memory Quote Repository"]
    end

    CLI --> OMS
    OMS --> ORD

    QAP -.used by.-> OMS
    ORE -.used by.-> OMS
    QAP -.implemented by.-> QMS
    ORE -.implemented by.-> MOR

    classDef module fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef entity fill:#fff3bf,stroke:#b08900,color:#111;
    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef framework fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class QAP,QMS,ORE,OMS module;
    class ORD entity;
    class MOR,MQR dataadapter;
    class CLI framework;
    class QAP,ORE contract;
```

Legend:

- yellow: domain type
- purple: module-owned service or contract
- green: data adapter
- blue: framework edge
- dashed border: contract
- dashed arrow: structural relationship such as `used by` or `implemented by`

## Implementation Focus

Implement one workflow:

- convert approved quote to order

The code should show:

- a new `orders` module
- a narrow approved-quote API exposed by `quotes`
- an order created as a business snapshot from the approved quote
- the demo reaching conversion after approval

## What To Verify

- `go test ./...` passes
- approved quotes can be converted
- non-approved quotes cannot be converted
- the `orders` module depends on a `quotes` capability, not on quote storage
