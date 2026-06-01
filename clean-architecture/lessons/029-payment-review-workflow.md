# Lesson 029: Payment Review Workflow

## Objective

Introduce a payment review state so payment capture can produce a business outcome other than immediate success.

## Theory

Until now, payment capture has been modeled as:

- success, or
- technical failure

That is often too narrow.

Real payment integrations may say:

- approve immediately
- send to manual review
- fail technically

The important point is that "manual review" is not just an error.

It is a business outcome that changes the workflow state.

Clean Architecture handles this by letting the application boundary own:

- the outcome contract returned by the payment gateway
- the workflow transition for `PaymentReview`
- the explicit command that approves a reviewed payment

The infrastructure layer only reports the capture outcome.

The application layer decides what that means for the order lifecycle.

## Why This Matters Here

This is a stronger lesson than another report or list query because it adds a genuinely new state boundary.

The order workflow is no longer linear:

- pending payment
- paid
- shipped

It now has a real branch:

- pending payment
- payment review
- paid
- shipped

That makes the Clean use cases more representative of real business processes and shows how an interactor can coordinate business-state transitions from an external outcome without leaking gateway details inward.

## Diagram

```mermaid
flowchart LR
    subgraph ENT["Entities"]
        direction TB
        ORD["Order Entity"]
    end

    subgraph APP["Application"]
        direction TB
        CIN["CapturePayment Input Boundary"]
        COUT["CapturePayment Output Boundary"]
        AIN["ApprovePaymentReview Input Boundary"]
        AOUT["ApprovePaymentReview Output Boundary"]
        CUC["CapturePayment Interactor"]
        AUC["ApprovePaymentReview Interactor"]
        OED["Order Editor"]
        PGW["Payment Gateway"]
    end

    subgraph IA["Interface Adapters"]
        direction TB
        CCTRL["CapturePayment Controller"]
        ACTRL["ApprovePaymentReview Controller"]
        CPRES["CapturePayment Presenter"]
        APRES["ApprovePaymentReview Presenter"]
    end

    subgraph INFRA["Infrastructure / Frameworks"]
        direction TB
        CLI["CLI / HTTP Framework"]
        MOG["Memory Order Gateway"]
        APG["Accept-All Payment Gateway"]
        MPG["Manual-Review Payment Gateway"]
    end

    CLI --> CCTRL
    CLI --> ACTRL
    CCTRL --> CIN
    ACTRL --> AIN
    CUC --> COUT
    AUC --> AOUT
    CPRES --> CLI
    APRES --> CLI
    CUC --> ORD
    AUC --> ORD

    CIN -.used by.-> CCTRL
    CIN -.implemented by.-> CUC
    COUT -.used by.-> CUC
    COUT -.implemented by.-> CPRES
    AIN -.used by.-> ACTRL
    AIN -.implemented by.-> AUC
    AOUT -.used by.-> AUC
    AOUT -.implemented by.-> APRES
    OED -.used by.-> CUC
    OED -.used by.-> AUC
    PGW -.used by.-> CUC
    OED -.implemented by.-> MOG
    PGW -.implemented by.-> APG
    PGW -.implemented by.-> MPG

    classDef framework fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef funcadapter fill:#ffe5d9,stroke:#bc6c25,color:#111;
    classDef app fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef entity fill:#fff3bf,stroke:#b08900,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class CLI framework;
    class MOG dataadapter;
    class APG,MPG,CCTRL,ACTRL,CPRES,APRES funcadapter;
    class CIN,COUT,AIN,AOUT,CUC,AUC,OED,PGW app;
    class ORD entity;
    class CIN,COUT,AIN,AOUT,OED,PGW contract;
```

Legend:

- blue: framework edge
- green: data adapter
- orange: translation or service adapter
- purple: application layer
- yellow: entity layer
- dashed border: interface / contract
- dashed arrow: structural relationship such as `used by` or `implemented by`

## Implementation Focus

Add:

- `PaymentReview` as an order state
- payment capture outcomes for approved vs review
- `ApprovePaymentReview`

The code should show:

- the payment gateway returning a business outcome instead of only `error`
- capture moving some orders into `PaymentReview`
- shipment remaining blocked until review is approved

## What To Verify

- the project compiles
- `go test ./...` passes
- capture can move an order to `PaymentReview`
- approving payment review moves the order to `Paid`
- shipment is rejected while the order is still in review
