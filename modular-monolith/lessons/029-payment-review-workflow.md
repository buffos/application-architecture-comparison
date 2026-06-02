# Lesson 029: Payment Review Workflow

## Objective

Introduce a payment review state so payment capture can produce a business outcome other than immediate success.

## Theory

Until now, payment capture in the Modular Monolith track has been modeled as:

- success, or
- technical failure

That is too narrow for many real payment integrations.

A gateway may report:

- approve immediately
- send to manual review
- fail technically

The important point is that manual review is not just an error. It is a business outcome that changes workflow state.

In this modular monolith, that means:

- the `payments` module owns the capture outcome contract
- the `orders` module owns the `PendingPayment -> PaymentReview -> Paid` lifecycle
- infrastructure only reports the external capture outcome

## Why This Matters Here

This adds a real branch to the order lifecycle:

- `PendingPayment`
- `PaymentReview`
- `Paid`
- `Shipped`

That makes the workflow more realistic and shows how one module can translate another module's business outcome into its own state transition without leaking gateway details into the order model.

## Diagram

```mermaid
flowchart LR
    subgraph ODM["Orders Module"]
        direction TB
        ORE["orders.Repository"]
        OWS["orders.Service<br/>CapturePayment / ApprovePaymentReview"]
        ORD["Order"]
    end

    subgraph PAM["Payments Module"]
        direction TB
        PPR["payments.Processor"]
        PWS["payments.Service"]
    end

    subgraph INF["Infrastructure"]
        direction TB
        CLI["CLI"]
        MOR["Memory Order Repository"]
        AAG["AcceptAllGateway"]
        MRG["ManualReviewGateway"]
    end

    CLI --> OWS
    OWS --> ORD

    ORE -.used by.-> OWS
    PPR -.used by.-> OWS
    ORE -.implemented by.-> MOR
    PPR -.implemented by.-> PWS
    PPR -.implemented by.-> AAG
    PPR -.implemented by.-> MRG

    classDef module fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef entity fill:#fff3bf,stroke:#b08900,color:#111;
    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef funcadapter fill:#ffe5d9,stroke:#bc6c25,color:#111;
    classDef framework fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class ORE,OWS,PPR,PWS module;
    class ORD entity;
    class MOR dataadapter;
    class AAG,MRG funcadapter;
    class CLI framework;
    class ORE,PPR contract;
```

Legend:

- yellow: domain type or workflow record
- purple: module-owned service or contract
- green: data adapter
- orange: behavior adapter
- blue: framework edge
- dashed border: contract
- dashed arrow: structural relationship such as `used by` or `implemented by`

## Implementation Focus

Implement one branching payment workflow:

- payment capture returns a business outcome
- review moves the order into `PaymentReview`
- a separate command approves review and moves the order to `Paid`

The code should show:

- a capture outcome contract in the `payments` module
- `PaymentReview` as an order state
- `ApprovePaymentReview` in the `orders` module
- shipment remaining blocked while the order is still in review

## What To Verify

- `go test ./...` passes
- capture can move an order to `PaymentReview`
- approving payment review moves the order to `Paid`
- shipment is rejected while the order is still in review
