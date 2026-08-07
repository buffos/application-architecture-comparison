# Lesson 020: Add Order Query Scripts

## Objective

Expose order reads through `GetOrder` and `ListOrders` procedures rather than direct map access.

## Theory

Orders now carry quote snapshots, reservation state, payment state, fulfillment state, and cancellation metadata. Query procedures provide a stable read surface for that growing record.

The scripts will:

- load one order by ID;
- list orders with an optional status filter;
- sort by order ID;
- copy line slices before returning results.

They remain read-only and do not reuse command procedures for convenience.

## Why This Matters Here

The Transaction Script shape works for queries as well as commands, but the distinction is worth naming. Keeping reads separate prevents a caller from accidentally invoking a state-changing workflow just to inspect an order.

## Diagram

```mermaid
flowchart LR
    subgraph SCRIPT["internal/scripts"]
        GET["GetOrder"]
        LIST["ListOrders"]
    end

    subgraph DATA["internal/data"]
        STORE["Store.Orders"]
        SNAPSHOT["Order snapshots"]
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
- yellow: passive storage and snapshots;
- dashed arrows: non-mutating reads;
- solid arrows: result shaping.

## Implementation Focus

Implement only:

- `GetOrder`;
- `ListOrders` with optional status filtering;
- deterministic order-ID sorting and defensive line copies;
- query tests for found, missing, filtered, and unfiltered results.

Leave shipment and quote query surfaces for later lessons.

## What To Verify

- `go test ./...` passes from `transaction-script-architecture/`;
- one order can be read by ID;
- missing orders return a business error;
- status filtering and sorting work;
- modifying a returned line slice does not mutate the store.
