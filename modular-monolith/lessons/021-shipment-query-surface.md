# Lesson 021: Shipment Query Surface

## Objective

Give the `shipments` module an explicit read surface so callers load shipments through the module API instead of treating the repository as the public interface.

## Theory

The `shipments` module already owns creation:

- receive a shipment request
- build a shipment record
- persist it

But without explicit queries, outside code still has an easy shortcut:

- read shipment storage directly

That weakens the modular boundary because the repository starts to look like the real public API.

This lesson closes that gap:

- `shipments` still owns persistence
- the module now publishes `GetShipment`
- the module now publishes `ListShipments`

So both write and read access go through the module surface.

## Why This Matters Here

In a modular monolith, even small modules should own their public read shape.

If callers create through the module but read through storage, the architecture quietly drifts toward:

- command methods on services
- queries on repositories

That makes repositories shared access points again. An explicit query surface keeps the module boundary visible:

- the repository remains internal plumbing
- the module owns the read model it exposes
- callers depend on shipment capabilities, not storage details

## Diagram

```mermaid
flowchart LR
    subgraph SHM["Shipments Module"]
        direction TB
        SRE["shipments.Repository"]
        SQS["shipments.Service<br/>GetShipment / ListShipments"]
        SDT["ShipmentDetails"]
    end

    subgraph INF["Infrastructure"]
        direction TB
        CLI["CLI"]
        MSR["Memory Shipment Repository"]
    end

    CLI --> SQS
    SQS --> SDT

    SRE -.used by.-> SQS
    SRE -.implemented by.-> MSR

    classDef module fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef entity fill:#fff3bf,stroke:#b08900,color:#111;
    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef framework fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class SRE,SQS module;
    class SDT entity;
    class MSR dataadapter;
    class CLI framework;
    class SRE contract;
```

Legend:

- yellow: query model or business-facing read shape
- purple: module-owned service or contract
- green: adapter or technical implementation
- blue: framework edge
- dashed border: contract
- dashed arrow: structural relationship such as `used by` or `implemented by`

## Implementation Focus

Implement one explicit read boundary:

- query shipments through the `shipments` module

The code should show:

- `GetShipment`
- `ListShipments`
- repository support for list-by-order-id
- callers reading through the module service, not the repository directly

## What To Verify

- `go test ./...` passes
- a stored shipment can be loaded through the module API
- shipments can be listed for an order
- the demo can load and list shipments without direct repository access
