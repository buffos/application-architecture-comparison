# Lesson 014: Return Review Plugin

## Objective

Insert an explicit review step into the return workflow so a return request is created first, and refund plus restock happen only when the request is accepted.

## Theory

Lesson `013` completed the technical reversal for a shipped order:

- load a returnable order
- refund payment
- restock inventory
- store the return request

That works mechanically, but it assumes every request should be accepted immediately.

This lesson introduces a real workflow boundary:

- requesting a return
- accepting a return
- rejecting a return

In microkernel architecture, the important point is not only adding more states.

It is keeping ownership explicit:

- the `returns` plugin owns the return-request state machine
- the `payments` plugin still owns refunds
- the `inventory` plugin still owns restocking
- the kernel exposes the capabilities that let those plugins collaborate

That makes return review a first-class plugin workflow instead of an implicit side effect of request creation.

## Why This Matters Here

This is the first lesson where the `returns` plugin becomes more than a thin orchestration wrapper.

It now owns:

- the persistent return request
- the `Requested` state
- the acceptance path
- the rejection path

That makes the plugin boundary more meaningful, because the workflow state now lives in the plugin that owns the process.

## Diagram

```mermaid
flowchart LR
    subgraph KER["Kernel"]
        direction TB
        PLG["kernel.Plugin"]
        ROP["kernel.ReturnableOrderProvider"]
        RFC["kernel.PaymentRefund"]
        IRS["kernel.InventoryRestock"]
        RSA["kernel.ReturnService"]
        HST["kernel.Host"]
    end

    subgraph RTP["Returns Plugin"]
        direction TB
        RRE["ReturnRequest<br/>Request / Accept / Reject"]
        RTR["returns.Repository"]
        RSS["returns.Service<br/>RequestReturn / AcceptReturn / RejectReturn"]
        RPP["returns.Plugin"]
    end

    subgraph ORP["Orders Plugin"]
        direction TB
        OSS["orders.Service<br/>GetReturnableOrder"]
        OPP["orders.Plugin"]
    end

    subgraph PMP["Payments Plugin"]
        direction TB
        PGS["payments.Service<br/>Refund"]
        PPP["payments.Plugin"]
    end

    subgraph INP["Inventory Plugin"]
        direction TB
        ISS["inventory.Service<br/>Restock"]
        IPP["inventory.Plugin"]
    end

    subgraph INF["Infrastructure"]
        direction TB
        CLI["CLI"]
        MRR["Memory Return Repository"]
    end

    CLI --> HST
    HST --> RPP
    HST --> OPP
    HST --> PPP
    HST --> IPP
    RSS --> RRE

    ROP -.used by.-> RSS
    RTR -.used by.-> RSS
    RFC -.used by.-> RSS
    IRS -.used by.-> RSS
    RSA -.used by.-> CLI
    PLG -.implemented by.-> RPP
    PLG -.implemented by.-> OPP
    PLG -.implemented by.-> PPP
    PLG -.implemented by.-> IPP
    ROP -.implemented by.-> OSS
    RFC -.implemented by.-> PGS
    IRS -.implemented by.-> ISS
    RSA -.implemented by.-> RSS
    RTR -.implemented by.-> MRR

    classDef kernel fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef plugin fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef entity fill:#fff3bf,stroke:#b08900,color:#111;
    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef framework fill:#f8f9fa,stroke:#6c757d,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class PLG,ROP,RFC,IRS,RSA,HST kernel;
    class RTR,RSS,RPP,OSS,OPP,PGS,PPP,ISS,IPP plugin;
    class RRE entity;
    class MRR dataadapter;
    class CLI framework;
    class PLG,ROP,RFC,IRS,RSA,RTR contract;
```

Legend:

- blue: kernel-owned type or contract
- purple: plugin-owned service, repository contract, or plugin registration type
- yellow: plugin-owned domain type
- green: data adapter
- gray: framework edge
- dashed border: contract
- dashed arrow: structural relationship such as `used by` or `implemented by`

## Implementation Focus

- make `RequestReturn` create only a requested return
- add explicit `AcceptReturn` and `RejectReturn` operations to the return capability
- move refund and restock to the acceptance path
- keep rejection side-effect free

Do not add return policy yet.

## What To Verify

- `go test ./...` passes
- request creation stores only a requested return
- accepting a return triggers refund and restock
- rejecting a return does not trigger refund or restock
