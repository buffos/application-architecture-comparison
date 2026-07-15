# Lesson 031: Partial Returns By Line Plugin

## Objective

Make returns quantity-aware so a return request can represent a specific returned slice instead of always meaning "the whole shipped order comes back."

## Theory

Up to this point, the return workflow has assumed a simple rule:

- a return request snapshots everything that was shipped

That is useful early on, but too narrow for realistic reverse logistics.

Real systems often need:

- one returned unit from a multi-unit line
- selected lines from a larger shipment
- approval, refund, and restock only for the returned slice

In this microkernel, the important ownership split is:

- the orders plugin exposes shipped quantities in its returnable-order capability
- the returns plugin owns the explicit requested return lines
- refund, restock, and reporting use the actual requested slice rather than the full order

## Why This Matters Here

The partial-shipment lesson made fulfillment quantity-aware.

This lesson makes reverse fulfillment quantity-aware too.

That matters because:

- refund amount should match the returned quantity
- restock amount should match the returned quantity
- return-rate reporting should use the actual returned slice

Without this step, reverse workflow logic still stays coarser than forward workflow logic.

## Diagram

```mermaid
flowchart LR
    subgraph KER["Kernel"]
        direction TB
        RSP["kernel.ReturnService"]
        ROP["kernel.ReturnableOrderProvider"]
        HST["kernel.Host"]
    end

    subgraph ORP["Orders Plugin"]
        direction TB
        ORS["orders.Service<br/>GetReturnableOrder"]
        ROV["ReturnableOrder<br/>includes shipped quantities"]
    end

    subgraph RTP["Returns Plugin"]
        direction TB
        RRE["returns.Repository"]
        RWS["returns.Service<br/>RequestReturn / AcceptReturn"]
        RRQ["ReturnRequest<br/>captures requested slice"]
    end

    subgraph INF["Infrastructure"]
        direction TB
        CLI["CLI"]
    end

    CLI --> HST
    HST --> RWS
    ORS --> ROV
    RWS --> RRQ

    RSP -.used by.-> CLI
    ROP -.used by.-> RWS
    RRE -.used by.-> RWS
    ROP -.implemented by.-> ORS
    RSP -.implemented by.-> RWS

    classDef kernel fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef plugin fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef entity fill:#fff3bf,stroke:#b08900,color:#111;
    classDef framework fill:#f8f9fa,stroke:#6c757d,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class RSP,ROP,HST kernel;
    class ORS,RRE,RWS plugin;
    class ROV,RRQ entity;
    class CLI framework;
    class RSP,ROP,RRE contract;
```

Legend:

- blue: kernel-owned type or contract
- purple: plugin-owned service or registration type
- yellow: domain type
- gray: framework edge
- dashed border: contract
- dashed arrow: structural relationship such as `used by` or `implemented by`

## Implementation Focus

- explicit return line selection and quantity
- refund and restock based on the selected slice
- reporting using actual returned quantities

The code should show:

- request commands carrying return-line input
- `ReturnRequest` storing only the requested slice
- existing review and refund flow still working on that narrower model

## What To Verify

- `go test ./...` passes
- a return request can capture only part of a shipped line
- refund amount matches the returned quantity
- restock quantity matches the returned quantity
