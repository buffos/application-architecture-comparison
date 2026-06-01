# Lesson 015: Return Eligibility Policy

## Objective

Make return acceptance policy-aware so not every requested return is accepted automatically.

## Theory

The previous lesson added a review boundary:

- request
- accept or reject
- refund and restock only on acceptance

That still leaves one simplifying assumption:

- every requested return is acceptable

Onion Architecture handles that by adding another inward-facing policy contract owned by the application ring.

This keeps the responsibilities clear:

- the domain core owns the return state machine
- the application ring asks whether the request is eligible
- infrastructure provides the concrete policy implementation

## Why This Matters Here

If acceptance is unconditional, the review boundary is only procedural.

Adding an eligibility policy makes the review decision substantive:

- some returns are accepted
- some returns are refused by policy

This is exactly the kind of rule that often changes independently from core lifecycle rules, which makes it a strong candidate for an application-owned boundary.

## Diagram

```mermaid
flowchart LR
    subgraph DOM["Domain Core"]
        direction TB
        RET["Return Request Entity"]
        ACC["Accept() Transition"]
    end

    subgraph APP["Application Ring"]
        direction TB
        ARS["AcceptReturn Service"]
        POL["Return Eligibility Policy"]
        RFD["Refund Gateway"]
        RSK["Inventory Restock"]
    end

    subgraph INF["Infrastructure Ring"]
        direction TB
        RPL["Reason Policy"]
        ARG["Accept-All Refund Gateway"]
        MIR["Memory Inventory Reservation"]
    end

    ARS --> RET
    RET --> ACC

    POL -.used by.-> ARS
    RFD -.used by.-> ARS
    RSK -.used by.-> ARS
    POL -.implemented by.-> RPL
    RFD -.implemented by.-> ARG
    RSK -.implemented by.-> MIR

    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef app fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef domain fill:#fff3bf,stroke:#b08900,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class RPL,ARG,MIR dataadapter;
    class ARS,POL,RFD,RSK app;
    class RET,ACC domain;
    class POL,RFD,RSK contract;
```

## Implementation Focus

Implement one rule refinement:

- acceptance depends on a return eligibility policy

The code should show:

- a policy contract in the application ring
- a simple reason-based policy in infrastructure
- acceptance blocked before refund and restock when policy rejects the request

## What To Verify

- `go test ./...` passes
- eligible requested returns can still be accepted
- policy-blocked returns stay in `Requested`
- blocked returns do not refund or restock
