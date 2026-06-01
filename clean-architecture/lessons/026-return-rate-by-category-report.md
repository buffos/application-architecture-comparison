# Lesson 026: Return Rate By Category Report

## Objective

Add a second projection-style report that combines shipped orders, refunded returns, and product category lookup into one application-owned read model.

## Theory

The first report lesson showed that reports are not just "big queries."

They are application projections with their own boundaries.

This lesson goes one step further:

- it reads from multiple workflow sources
- it resolves supporting catalog data
- it groups the result into a category-based metric

The report answers:

- how many units were shipped per category
- how many units were later returned per category
- what the resulting return rate is

That is still a Clean Architecture use case.

The report model does not belong to any one entity.

It belongs to the application layer because the application defines what the metric means.

The tradeoff is broader reader dependencies:

- orders
- returns
- products

But that dependency breadth is explicit and still points inward to application-owned contracts.

## Why This Matters Here

The repository already has entity-centric reads.

This lesson shows a more realistic analytics-style use case, where the application layer must coordinate several data sources to express a business metric that no single aggregate owns by itself.

That makes the reporting story more complete and easier to compare with the other architecture tracks.

## Diagram

```mermaid
flowchart LR
    subgraph APP["Application"]
        direction TB
        RIN["ReturnRateByCategoryReport Input Boundary"]
        ROUT["ReturnRateByCategoryReport Output Boundary"]
        RUC["ReturnRateByCategoryReport Interactor"]
        ORR["Order Report Reader"]
        RRR["Return Report Reader"]
        PRR["Product Report Reader"]
        RPT["Return Rate Report Model"]
    end

    subgraph IA["Interface Adapters"]
        direction TB
        RCTRL["ReturnRateByCategoryReport Controller"]
        RPRES["ReturnRateByCategoryReport Presenter"]
    end

    subgraph INFRA["Infrastructure / Frameworks"]
        direction TB
        CLI["CLI / HTTP Framework"]
        MOG["Memory Order Gateway"]
        MRG["Memory Return Request Gateway"]
        MPG["Memory Product Gateway"]
    end

    CLI --> RCTRL
    RCTRL --> RIN
    RUC --> ROUT
    RPRES --> CLI
    RUC --> RPT

    RIN -.used by.-> RCTRL
    RIN -.implemented by.-> RUC
    ROUT -.used by.-> RUC
    ROUT -.implemented by.-> RPRES
    ORR -.used by.-> RUC
    RRR -.used by.-> RUC
    PRR -.used by.-> RUC
    ORR -.implemented by.-> MOG
    RRR -.implemented by.-> MRG
    PRR -.implemented by.-> MPG

    classDef framework fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef funcadapter fill:#ffe5d9,stroke:#bc6c25,color:#111;
    classDef app fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class CLI framework;
    class MOG,MRG,MPG dataadapter;
    class RCTRL,RPRES funcadapter;
    class RIN,ROUT,RUC,ORR,RRR,PRR,RPT app;
    class RIN,ROUT,ORR,RRR,PRR contract;
```

Legend:

- blue: framework edge
- green: data adapter
- orange: translation adapter
- purple: application layer
- dashed border: interface / contract
- dashed arrow: structural relationship such as `used by` or `implemented by`

## Implementation Focus

Add:

- `ReturnRateByCategoryReport`

The code should show:

- a report interactor that reads from orders, returns, and products
- category resolution as an explicit supporting dependency
- a presenter shaping the aggregated category metrics for callers

## What To Verify

- the project compiles
- `go test ./...` passes
- shipped and refunded quantities are grouped by category correctly
- the demo can render the report output
