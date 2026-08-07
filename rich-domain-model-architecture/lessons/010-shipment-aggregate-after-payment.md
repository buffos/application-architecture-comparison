# Lesson 010: Shipment Aggregate After Payment

## Objective

Model Shipment as a separate Fulfillment aggregate and require payment before shipment creation.

## Theory

Shipment is a physical fulfillment record, not merely another Order flag. It owns shipment lines and its own Pending/Dispatched lifecycle. Order remains the commercial aggregate and records that fulfillment has started only after Shipment is dispatched.

The creation rule is explicit: a PendingPayment Order cannot create a Shipment. Payment and Shipment collaborate through guarded operations rather than one aggregate reaching into the other's internals.

## Why This Matters Here

Rich Domain Model keeps distinct business lifecycles in distinct objects. If Shipment were just a field on Order, warehouse-specific behavior would gradually make the commercial aggregate responsible for physical fulfillment as well.

The cost is a multi-aggregate workflow and explicit sequencing. The benefit is that shipment rules can grow independently and be tested without constructing the entire commerce model.

## Diagram

```mermaid
flowchart LR
    PAID["Order Paid"] --> CREATE["create shipment"]
    CREATE --> SHIPMENT["Shipment aggregate"]
    SHIPMENT --> DISPATCH["dispatch"]
    DISPATCH --> SHIPPED["Order.MarkShipped"]
    UNPAID["Order PendingPayment"] -. "blocked" .-> CREATE

    classDef order fill:#fff3bf,stroke:#b08900,color:#111;
    classDef shipment fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef command fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef blocked fill:#ffe5d9,stroke:#bc6c25,color:#111;

    class PAID,SHIPPED order;
    class SHIPMENT shipment;
    class CREATE,DISPATCH command;
    class UNPAID blocked;
```

## Implementation Focus

Implement only:

- Fulfillment Shipment identity, line snapshots, and dispatch lifecycle
- creation from a Paid Order
- Order's guarded Shipped transition
- tests for unpaid rejection and successful dispatch
- demo payment → shipment → shipped-order flow

Leave partial shipment quantities for a later lesson.

## What To Verify

- `go test ./...` passes
- unpaid Orders cannot create Shipments
- paid Orders create pending Shipments
- dispatch changes Shipment and Order state correctly
- Shipment remains a separate aggregate from Order
