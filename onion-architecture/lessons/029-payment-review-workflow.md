# Lesson 029: Payment Review Workflow

## Objective

Introduce a payment review state so payment capture can produce a business outcome other than immediate success.

## Theory

Until now, payment capture in the Onion track has been modeled as:

- success, or
- technical failure

That is too narrow for many real payment integrations.

A gateway may report:

- approve immediately
- send to manual review
- fail technically

The important point is that manual review is not just an error. It is a business outcome that changes workflow state.

In Onion Architecture, that means:

- the application ring owns the capture outcome contract
- the domain ring owns the `PaymentReview -> Paid` state rule
- infrastructure only reports the external capture outcome

## Why This Matters Here

This adds a real branch to the order lifecycle:

- `PendingPayment`
- `PaymentReview`
- `Paid`
- `Shipped`

That makes the Onion workflow more realistic and shows how application services translate external outcomes into domain state transitions without pushing gateway semantics into the domain model itself.

## Diagram

```mermaid
flowchart LR
    subgraph DOM["Domain Ring"]
        direction TB
        ORD["Order"]
    end

    subgraph APP["Application Ring"]
        direction TB
        ORR["OrderRepository"]
        PG["PaymentGateway"]
        CPS["CapturePayment Service"]
        APR["ApprovePaymentReview Service"]
    end

    subgraph INF["Infrastructure Ring"]
        direction TB
        MOR["Memory Order Repository"]
        AAG["AcceptAllGateway"]
        MRG["ManualReviewGateway"]
    end

    CPS --> ORD
    APR --> ORD

    ORR -.used by.-> CPS
    ORR -.used by.-> APR
    PG -.used by.-> CPS
    ORR -.implemented by.-> MOR
    PG -.implemented by.-> AAG
    PG -.implemented by.-> MRG

    classDef entity fill:#fff3bf,stroke:#b08900,color:#111;
    classDef app fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef funcadapter fill:#ffe5d9,stroke:#bc6c25,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class ORD entity;
    class ORR,PG,CPS,APR app;
    class MOR dataadapter;
    class AAG,MRG funcadapter;
    class ORR,PG contract;
```

Legend:

- yellow: domain type
- purple: application type
- green: infrastructure data adapter
- orange: infrastructure behavior adapter
- dashed border: contract
- dashed arrow: structural relationship such as `used by` or `implemented by`

## Implementation Focus

Add:

- `PaymentReview` as an order state
- payment capture outcomes for approved vs review
- `ApprovePaymentReview`

The code should show:

- the gateway returning a business outcome instead of only `error`
- capture moving some orders into `PaymentReview`
- shipment remaining blocked until review is approved

## What To Verify

- `go test ./...` passes
- capture can move an order to `PaymentReview`
- approving payment review moves the order to `Paid`
- shipment is rejected while the order is still in review
