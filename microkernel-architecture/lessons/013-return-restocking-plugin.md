# Lesson 013: Return Restocking Plugin

## Objective

Extend the return workflow so refunded returns also restock inventory through a distinct inventory capability exposed by the microkernel.

## Theory

Lesson `012` introduced a new `returns` plugin and showed that post-shipment reversal can be implemented by composing kernel capabilities from other plugins.

But that flow still reversed only the money side:

- load a returnable order
- refund the payment
- store the return request

It did not reverse stock.

This lesson makes stock reversal explicit by adding a separate kernel-owned `InventoryRestock` capability.

That matters because:

- `orders` still owns whether an order is returnable
- `payments` still owns refund behavior
- `inventory` owns stock changes
- `returns` orchestrates the workflow, but does not own inventory rules

The microkernel value here is not that one plugin does everything.

It is that a new workflow can be assembled by combining stable capabilities from other plugins while keeping ownership where it belongs.

## Why This Matters Here

This is the first microkernel lesson where one plugin coordinates three distinct capabilities from the rest of the system:

- order lookup for a returnable order
- payment refund
- inventory restock

That makes the host-plus-capability model more concrete than the earlier forward-only flow.

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

    subgraph ORP["Orders Plugin"]
        direction TB
        ORE["Order"]
        OSS["orders.Service<br/>... / GetReturnableOrder"]
        OPP["orders.Plugin"]
    end

    subgraph PMP["Payments Plugin"]
        direction TB
        PGS["payments.Service<br/>Capture / Refund"]
        PPP["payments.Plugin"]
    end

    subgraph INP["Inventory Plugin"]
        direction TB
        STR["StockRecord"]
        INR["inventory.Repository"]
        ISS["inventory.Service<br/>Reserve / Release / Restock"]
        IPP["inventory.Plugin"]
    end

    subgraph RTP["Returns Plugin"]
        direction TB
        RRE["ReturnRequest"]
        RTR["returns.Repository"]
        RSS["returns.Service<br/>RequestReturn"]
        RPP["returns.Plugin"]
    end

    subgraph INF["Infrastructure"]
        direction TB
        CLI["CLI"]
        MIR["Memory Inventory Repository"]
        MRR["Memory Return Repository"]
    end

    CLI --> HST
    HST --> OPP
    HST --> PPP
    HST --> IPP
    HST --> RPP
    OSS --> ORE
    ISS --> STR
    RSS --> RRE

    ROP -.used by.-> RSS
    RFC -.used by.-> RSS
    IRS -.used by.-> RSS
    RSA -.used by.-> CLI
    INR -.used by.-> ISS
    RTR -.used by.-> RSS
    PLG -.implemented by.-> OPP
    PLG -.implemented by.-> PPP
    PLG -.implemented by.-> IPP
    PLG -.implemented by.-> RPP
    ROP -.implemented by.-> OSS
    RFC -.implemented by.-> PGS
    IRS -.implemented by.-> ISS
    RSA -.implemented by.-> RSS
    INR -.implemented by.-> MIR
    RTR -.implemented by.-> MRR

    classDef kernel fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef plugin fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef entity fill:#fff3bf,stroke:#b08900,color:#111;
    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef framework fill:#f8f9fa,stroke:#6c757d,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class PLG,ROP,RFC,IRS,RSA,HST kernel;
    class OSS,OPP,PGS,PPP,ISS,IPP,RSS,RPP,INR,RTR plugin;
    class ORE,STR,RRE entity;
    class MIR,MRR dataadapter;
    class CLI framework;
    class PLG,ROP,RFC,IRS,RSA,INR,RTR contract;
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

- add a kernel-owned `InventoryRestock` capability
- expose that capability from the `inventory` plugin
- make the `returns` plugin refund and restock in the same workflow
- keep return persistence inside the `returns` plugin

Do not add return review yet.

## What To Verify

- `go test ./...` passes
- successful return requests trigger both refund and restock
- a restock failure stops the workflow
- restocking still happens through an inventory capability rather than direct repository access from `returns`
