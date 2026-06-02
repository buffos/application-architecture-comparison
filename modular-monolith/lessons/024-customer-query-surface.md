# Lesson 024: Customer Query Surface

## Objective

Promote the `customers` module from a supporting validation dependency into an explicit read surface with customer queries through the module boundary.

## Theory

The `customers` module already exposes one narrow capability:

- `RequireActiveCustomer`

That is useful for quote creation, but it does not yet make the module a visible public read boundary for customer lookup or browsing.

This lesson adds that missing surface:

- `customers` still supports active-customer validation
- the module now publishes `GetCustomer`
- the module now publishes `ListCustomers`

So the module has:

- one specialized capability for validation
- one general read surface for customer access

## Why This Matters Here

Without explicit customer queries, the customer module stays a helper and storage becomes the natural place to read customers from.

That weakens the modular-monolith story because the system drifts toward:

- module services for workflow checks
- repositories for ordinary reads

Adding customer queries keeps the boundary consistent:

- the repository remains internal plumbing
- the `customers` module owns the read shapes it exposes
- callers depend on customer capabilities, not storage details

## Diagram

```mermaid
flowchart LR
    subgraph CUM["Customers Module"]
        direction TB
        CRE["customers.Repository"]
        CQS["customers.Service<br/>RequireActiveCustomer / GetCustomer / ListCustomers"]
        CDT["CustomerDetails"]
    end

    subgraph INF["Infrastructure"]
        direction TB
        CLI["CLI"]
        MCR["Memory Customer Repository"]
    end

    CLI --> CQS
    CQS --> CDT

    CRE -.used by.-> CQS
    CRE -.implemented by.-> MCR

    classDef module fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef entity fill:#fff3bf,stroke:#b08900,color:#111;
    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef framework fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class CRE,CQS module;
    class CDT entity;
    class MCR dataadapter;
    class CLI framework;
    class CRE contract;
```

Legend:

- yellow: query model or business-facing read shape
- purple: module-owned service or contract
- green: adapter or technical implementation
- blue: framework edge
- dashed border: contract
- dashed arrow: structural relationship such as `used by` or `implemented by`

## Implementation Focus

Implement one explicit read surface:

- query customers through the `customers` module

The code should show:

- `GetCustomer`
- `ListCustomers`
- repository support for active filtering
- existing validation behavior still available through `RequireActiveCustomer`

## What To Verify

- `go test ./...` passes
- a stored customer can be loaded through the module API
- customers can be listed with active-only filtering through the module API
- the demo can load and list customers without direct repository access
