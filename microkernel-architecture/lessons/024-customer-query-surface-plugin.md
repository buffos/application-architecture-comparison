# Lesson 024: Customer Query Surface Plugin

## Objective

Promote customers from a supporting validation dependency into an explicit read surface with customer queries through the plugin boundary.

## Theory

The customers plugin already exposes one narrow capability:

- `RequireActiveCustomer`

That is useful for quote creation, but it does not yet make the plugin a visible public read boundary for customer lookup or browsing.

This lesson adds that missing surface:

- the customers plugin still supports active-customer validation
- the plugin now exposes `GetCustomer`
- the plugin now exposes `ListCustomers`

So the plugin has:

- one specialized capability for validation
- one general read surface for customer access

## Why This Matters Here

Without explicit customer queries, the customers plugin stays a helper and storage becomes the natural place to read customers from.

That weakens the microkernel story because the system drifts toward:

- plugin capabilities for workflow checks
- repositories for ordinary reads

Adding customer queries keeps the boundary consistent:

- the repository remains internal plumbing
- the customers plugin owns the read shapes it exposes
- callers depend on customer capabilities, not storage details

## Diagram

```mermaid
flowchart LR
    subgraph KER["Kernel"]
        direction TB
        CDR["kernel.CustomerDirectory"]
        CRD["kernel.CustomerReader"]
        HST["kernel.Host"]
    end

    subgraph CUP["Customers Plugin"]
        direction TB
        CRE["customers.Repository"]
        CSS["customers.Service<br/>RequireActiveCustomer / GetCustomer / ListCustomers"]
        CPP["customers.Plugin"]
    end

    subgraph INF["Infrastructure"]
        direction TB
        CLI["CLI"]
        MCR["Memory Customer Repository"]
    end

    CLI --> HST
    HST --> CPP

    CRE -.used by.-> CSS
    CDR -.used by.-> CLI
    CRD -.used by.-> CLI
    CDR -.implemented by.-> CSS
    CRD -.implemented by.-> CSS
    CRE -.implemented by.-> MCR

    classDef kernel fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef plugin fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef framework fill:#f8f9fa,stroke:#6c757d,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class CDR,CRD,HST kernel;
    class CRE,CSS,CPP plugin;
    class MCR dataadapter;
    class CLI framework;
    class CDR,CRD,CRE contract;
```

Legend:

- blue: kernel-owned type or contract
- purple: plugin-owned service or plugin registration type
- green: data adapter
- gray: framework edge
- dashed border: contract
- dashed arrow: structural relationship such as `used by` or `implemented by`

## Implementation Focus

- keep `RequireActiveCustomer`
- add `GetCustomer`
- add `ListCustomers`
- support active filtering in the repository-backed read surface

Do not add reporting yet.

## What To Verify

- `go test ./...` passes
- a stored customer can be loaded through the kernel capability
- customers can be listed with active-only filtering through the kernel capability
- the demo can load and list customers without direct repository access
