# Lesson 013: Return Restocking Boundary

## Objective

Complete the stock-side reversal after returns by restocking inventory as part of the return workflow.

## Theory

The previous lesson added the first post-shipment reverse path:

- shipped order can request a return
- refund is executed through an external gateway

That still leaves one missing operational concern:

- returned stock should go back into inventory

Onion Architecture handles this the same way it handles reservation and release:

- the application ring owns the workflow coordination
- the domain core stays focused on business concepts
- infrastructure implements the stock operation

The important separation is:

- refund is a money-side boundary
- restock is an inventory-side boundary

They belong to the same workflow, but they are not the same dependency.

## Why This Matters Here

If returns stop at refund, the reverse workflow is only financially complete.

Adding restocking makes it operationally complete as well:

- refund compensates the customer
- restock compensates inventory

This makes the Onion application ring more realistic because it now coordinates multiple external boundaries in the same use case.

## Diagram

```mermaid
flowchart LR
    subgraph DOM["Domain Core"]
        direction TB
        ORD["Order Entity"]
        RET["Return Request Entity"]
    end

    subgraph APP["Application Ring"]
        direction TB
        RRS["RequestReturn Service"]
        OFD["Order Finder"]
        RST["Return Store"]
        RFG["Refund Gateway"]
        RSK["Inventory Restock"]
    end

    subgraph INF["Infrastructure Ring"]
        direction TB
        CLI["CLI Framework"]
        MOR["Memory Order Repository"]
        MRR["Memory Return Repository"]
        ARG["Accept-All Refund Gateway"]
        MIR["Memory Inventory Reservation"]
    end

    CLI --> RRS
    RRS --> ORD
    RRS --> RET

    OFD -.used by.-> RRS
    RST -.used by.-> RRS
    RFG -.used by.-> RRS
    RSK -.used by.-> RRS
    OFD -.implemented by.-> MOR
    RST -.implemented by.-> MRR
    RFG -.implemented by.-> ARG
    RSK -.implemented by.-> MIR

    classDef framework fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef app fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef domain fill:#fff3bf,stroke:#b08900,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class CLI framework;
    class MOR,MRR,ARG,MIR dataadapter;
    class RRS,OFD,RST,RFG,RSK app;
    class ORD,RET domain;
    class OFD,RST,RFG,RSK contract;
```

Legend:

- blue: framework edge
- green: data adapter
- purple: application ring
- yellow: domain core
- dashed border: interface / contract
- dashed arrow: structural relationship

## Implementation Focus

Implement one stock-side extension:

- restock inventory during return request processing

The code should show:

- a distinct restock item type
- an inventory restock contract in the application ring
- in-memory restock support in the inventory adapter
- tests proving stock rises after a refunded return

## What To Verify

- `go test ./...` passes
- shipped orders still produce refunded return requests
- return processing restocks inventory
- restocking stays outside the domain core
