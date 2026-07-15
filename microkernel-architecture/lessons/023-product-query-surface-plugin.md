# Lesson 023: Product Query Surface Plugin

## Objective

Promote products from a supporting quote dependency into an explicit read surface with product queries through the plugin boundary.

## Theory

The products plugin already exposes one narrow capability:

- `GetProductForQuote`

That is useful for the quotes workflow, but it is not the same as saying the plugin has a real public read API for product browsing or lookup.

This lesson adds that missing surface:

- the products plugin still supports quote pricing and category lookup
- the plugin now exposes `GetProduct`
- the plugin now exposes `ListProducts`

So the plugin has:

- one specialized capability for quote creation
- one general read surface for product access

## Why This Matters Here

Without explicit product queries, the products plugin remains a helper instead of a visible business boundary.

That encourages a common drift:

- quote workflows use the products plugin
- everything else reads product storage directly

Adding product queries keeps the architecture consistent:

- the repository remains internal plumbing
- the products plugin owns the read shapes it exposes
- callers depend on product capabilities, not storage details

## Diagram

```mermaid
flowchart LR
    subgraph KER["Kernel"]
        direction TB
        PCT["kernel.ProductCatalog"]
        PRD["kernel.ProductReader"]
        HST["kernel.Host"]
    end

    subgraph PRP["Products Plugin"]
        direction TB
        PRE["products.Repository"]
        PSS["products.Service<br/>GetProductForQuote / GetProduct / ListProducts"]
        PPP["products.Plugin"]
    end

    subgraph INF["Infrastructure"]
        direction TB
        CLI["CLI"]
        MPR["Memory Product Repository"]
    end

    CLI --> HST
    HST --> PPP

    PRE -.used by.-> PSS
    PCT -.used by.-> CLI
    PRD -.used by.-> CLI
    PCT -.implemented by.-> PSS
    PRD -.implemented by.-> PSS
    PRE -.implemented by.-> MPR

    classDef kernel fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef plugin fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef framework fill:#f8f9fa,stroke:#6c757d,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class PCT,PRD,HST kernel;
    class PRE,PSS,PPP plugin;
    class MPR dataadapter;
    class CLI framework;
    class PCT,PRD,PRE contract;
```

Legend:

- blue: kernel-owned type or contract
- purple: plugin-owned service or plugin registration type
- green: data adapter
- gray: framework edge
- dashed border: contract
- dashed arrow: structural relationship such as `used by` or `implemented by`

## Implementation Focus

- keep `GetProductForQuote`
- add `GetProduct`
- add `ListProducts`
- support category and active filtering in the repository-backed read surface

Do not add customer query surfaces yet.

## What To Verify

- `go test ./...` passes
- a stored product can be loaded through the kernel capability
- products can be listed by category and activity through the kernel capability
- the demo can load and list products without direct repository access
