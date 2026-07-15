# Lesson 030: Partial Shipment Support Plugin

## Objective

Make fulfillment quantity-aware so an order can be shipped in multiple steps instead of only as an all-or-nothing transition.

## Theory

Up to this point, shipment creation has assumed a simple rule:

- once an order is payable, one shipment ships everything

That is useful early on, but too narrow for a realistic fulfillment workflow.

Real systems often need:

- a first shipment for currently available quantity
- later shipments for the remaining quantity

In this microkernel, the important ownership split is:

- the orders plugin tracks shipped progress per line and owns fulfillment state
- the shipments plugin records the shipped slice
- the order service decides whether to ship explicit quantities or all remaining quantities

The key lifecycle change is the new intermediate state:

- `PartiallyShipped`

## Why This Matters Here

The payment review lesson added a branch before fulfillment.

This lesson adds incremental fulfillment inside fulfillment itself.

The shipping workflow is no longer a one-time state flip. It becomes progress over time, and that affects other application behavior too:

- cancellation must now reject partially shipped orders
- later shipment commands must continue from remaining quantity
- return logic can now distinguish ordered quantity from shipped quantity

## Diagram

```mermaid
flowchart LR
    subgraph KER["Kernel"]
        direction TB
        OST["kernel.OrderService"]
        SCT["kernel.ShipmentCreation"]
        HST["kernel.Host"]
    end

    subgraph ORP["Orders Plugin"]
        direction TB
        OWS["orders.Service<br/>CreateShipment"]
        ORD["Order<br/>tracks shipped quantity"]
        ORE["orders.Repository"]
    end

    subgraph SHP["Shipments Plugin"]
        direction TB
        SWS["shipments.Service"]
        SHI["Shipment<br/>captures shipped slice"]
    end

    subgraph INF["Infrastructure"]
        direction TB
        CLI["CLI"]
    end

    CLI --> HST
    HST --> OWS
    OWS --> ORD
    OWS --> SHI

    OST -.used by.-> CLI
    ORE -.used by.-> OWS
    SCT -.used by.-> OWS
    OST -.implemented by.-> OWS
    SCT -.implemented by.-> SWS

    classDef kernel fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef plugin fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef entity fill:#fff3bf,stroke:#b08900,color:#111;
    classDef framework fill:#f8f9fa,stroke:#6c757d,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class OST,SCT,HST kernel;
    class OWS,ORE,SWS plugin;
    class ORD,SHI entity;
    class CLI framework;
    class OST,SCT contract;
```

Legend:

- blue: kernel-owned type or contract
- purple: plugin-owned service or registration type
- yellow: domain type
- gray: framework edge
- dashed border: contract
- dashed arrow: structural relationship such as `used by` or `implemented by`

## Implementation Focus

- `PartiallyShipped` as an order state
- shipped quantity tracking on order lines
- explicit shipment line input for partial shipment
- default "ship all remaining" behavior for the existing full-shipment path

The code should show:

- the orders plugin updating shipment progress
- later shipments continuing from remaining quantity
- cancellation treating partial shipment as already fulfilled

## What To Verify

- `go test ./...` passes
- a partial shipment stores only the requested quantity
- a later shipment can ship the remaining quantity
- partially shipped orders cannot be cancelled
