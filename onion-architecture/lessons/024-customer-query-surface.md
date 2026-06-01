# Lesson 024: Customer Query Surface

## Objective

Add an explicit customer read surface through the application ring so the main supporting entities now all have first-class query use cases.

## Theory

Customers have been present since the first Onion lesson, but only as a supporting dependency for quote creation.

At this point, they should be queryable through the same architectural path as products, quotes, orders, returns, and shipments.

That keeps the repository consistent:

- application ring owns the query use cases
- infrastructure only implements lookup and filtering
- outer layers depend on application, not directly on storage

## Why This Matters Here

Customer reads often support:

- quote creation
- order review
- support workflows

If they bypass the application ring, one of the core teaching points of the repo becomes inconsistent on a basic entity.

## Diagram

```mermaid
flowchart LR
    subgraph DOM["Domain Core"]
        direction TB
        CUS["Customer Entity"]
    end

    subgraph APP["Application Ring"]
        direction TB
        GCU["GetCustomer Service"]
        LCU["ListCustomers Service"]
        CLK["Customer Lookup"]
    end

    subgraph INF["Infrastructure Ring"]
        direction TB
        MCR["Memory Customer Repository"]
    end

    GCU --> CUS
    LCU --> CUS

    CLK -.used by.-> GCU
    CLK -.used by.-> LCU
    CLK -.implemented by.-> MCR

    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef app fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef domain fill:#fff3bf,stroke:#b08900,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class MCR dataadapter;
    class GCU,LCU,CLK app;
    class CUS domain;
    class CLK contract;
```

## Implementation Focus

Implement two read use cases:

- get customer by id
- list customers by active status

The code should show:

- a customer lookup contract in the application ring
- application-shaped customer query results
- in-memory filtering by active status

## What To Verify

- `go test ./...` passes
- customers can be loaded by id
- customers can be filtered by active status
- customer reads now cross the application ring explicitly
