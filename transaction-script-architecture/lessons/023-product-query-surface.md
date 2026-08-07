# Lesson 023: Add Product Query Scripts

## Objective

Expose catalog reads through `GetProduct` and `ListProducts` procedures.

## Theory

Products are used by quote scripts, stock scripts, pricing rules, and return eligibility. A small read surface lets callers inspect catalog data without reaching directly into `Store.Products`.

`ListProducts` supports optional category and active-only filters and sorts by SKU.

## Why This Matters Here

The same passive record is now a shared input to many procedures. A query script creates a named seam without pretending that the product is a rich aggregate. The tradeoff is another small procedure and continued coupling to the record shape.

## Diagram

```mermaid
flowchart LR
    subgraph SCRIPT["internal/scripts"]
        GET["GetProduct"]
        LIST["ListProducts\ncategory + active filters"]
    end

    subgraph DATA["internal/data"]
        STORE["Store.Products"]
        SNAPSHOT["Product snapshots"]
    end

    GET -.reads.-> STORE
    LIST -.filters and reads.-> STORE
    GET --> SNAPSHOT
    LIST --> SNAPSHOT

    classDef script fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef data fill:#fff3bf,stroke:#b08900,color:#111;
    class GET,LIST script;
    class STORE,SNAPSHOT data;
```

Legend:

- purple: query procedures;
- yellow: passive catalog data and views;
- dashed arrows: reads and filters;
- solid arrows: result shaping.

## Implementation Focus

Implement only:

- `GetProduct`;
- `ListProducts` with category and active-only filters;
- deterministic SKU sorting;
- product query tests.

Leave customer queries for the next lesson.

## What To Verify

- `go test ./...` passes from `transaction-script-architecture/`;
- products can be read by SKU;
- category and active-only filters work;
- unavailable products can still be read when the filter is disabled.
