# Lesson 032: Plugin Pricing Extension Point

## Objective

Add a real extension seam so enabled plugins can change quote-line pricing without changing the `quotes` workflow structure.

## Theory

Up to now, quote pricing in the Modular Monolith track has just been the product's stored unit price.

That proves a simple workflow, but not extensibility.

This lesson adds a different architectural idea:

- the `quotes` module keeps its `AddQuoteLine` use case stable
- a `plugins` module owns registration, enablement, and listing
- a `pricing` module owns the pricing capability that `quotes` depends on

The extension stays deliberately narrow:

- quote-line unit price

The quote workflow does not change, but the effective unit price can change because the enabled plugin set changes.

## Why This Matters Here

A modular monolith is not only about workflow orchestration. It also needs disciplined places for optional behavior to grow.

Without this seam, every pricing experiment would push conditionals into:

- `quotes.Service`
- `Quote`
- or random infrastructure helpers

With it:

- the `quotes` module still owns quote editing
- the `plugins` module owns activation state
- the `pricing` module owns price calculation behavior

## Diagram

```mermaid
flowchart LR
    subgraph QTM["Quotes Module"]
        direction TB
        QRE["quotes.Repository"]
        QWS["quotes.Service<br/>AddQuoteLine"]
    end

    subgraph PGM["Plugins Module"]
        direction TB
        PRE["plugins.Repository"]
        PWS["plugins.Service<br/>Register / Enable / List"]
        PLG["PluginRegistration"]
    end

    subgraph PRM["Pricing Module"]
        direction TB
        PPR["pricing.QuotePricer"]
        PRS["pricing.Service"]
    end

    subgraph INF["Infrastructure"]
        direction TB
        CLI["CLI"]
        MPR["Memory Plugin Repository"]
    end

    CLI --> QWS
    CLI --> PWS
    QWS --> PLG

    PPR -.used by.-> QWS
    PRE -.used by.-> PWS
    PRE -.used by.-> PRS
    PPR -.implemented by.-> PRS
    PRE -.implemented by.-> MPR

    classDef module fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef entity fill:#fff3bf,stroke:#b08900,color:#111;
    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef framework fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class QRE,QWS,PRE,PWS,PPR,PRS module;
    class PLG entity;
    class MPR dataadapter;
    class CLI framework;
    class PRE,PPR contract;
```

Legend:

- yellow: workflow record or domain-facing state
- purple: module-owned service or contract
- green: data adapter
- blue: framework edge
- dashed border: contract
- dashed arrow: structural relationship such as `used by` or `implemented by`

## Implementation Focus

Implement one extension seam:

- plugin registration, enable, and list services
- a pricing capability that `quotes` depends on
- one sample plugin: `seasonal-pricing`

The code should show:

- the quote use case stays structurally stable
- pricing changes only because the enabled plugin set changes
- plugin registration and activation are explicit module behavior

## What To Verify

- `go test ./...` passes
- a pricing plugin can be registered and enabled
- enabling `seasonal-pricing` changes the quote line unit price
- the demo can show plugin registration, activation, and pricing impact
