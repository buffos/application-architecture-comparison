# Lesson 032: Plugin Pricing Extension Point Plugin

## Objective

Add a real extension seam so enabled plugins can change quote-line pricing without changing the `quotes` workflow structure.

## Theory

Up to now, quote pricing in the microkernel track has just been the product's stored unit price.

That proves a simple workflow, but not extensibility.

This lesson adds a different architectural idea:

- the quotes plugin keeps its `AddQuoteLine` use case stable
- the kernel owns a narrow pricing capability
- pricing plugins can publish or decorate that capability

The extension stays deliberately narrow:

- quote-line unit price

The quote workflow does not change, but the effective unit price can change because the enabled plugin set changes.

## Why This Matters Here

A microkernel is not only about workflow slicing. It also needs disciplined places for optional behavior to grow.

Without this seam, every pricing experiment would push conditionals into:

- `quotes.Service`
- `Quote`
- or random infrastructure helpers

With it:

- the quotes plugin still owns quote editing
- the kernel owns the pricing contract
- pricing plugins own price calculation behavior

## Diagram

```mermaid
flowchart LR
    subgraph KER["Kernel"]
        direction TB
        QPR["kernel.QuotePricer"]
        QSV["kernel.QuoteService"]
        HST["kernel.Host"]
    end

    subgraph QTP["Quotes Plugin"]
        direction TB
        QWS["quotes.Service<br/>AddQuoteLine"]
    end

    subgraph PRP["Pricing Plugins"]
        direction TB
        BPR["pricing.Service<br/>passthrough"]
        SPR["seasonalpricing.Service<br/>decorator"]
    end

    subgraph INF["Infrastructure"]
        direction TB
        CLI["CLI"]
    end

    CLI --> HST
    HST --> QWS

    QPR -.used by.-> QWS
    QSV -.used by.-> CLI
    QPR -.implemented by.-> BPR
    QPR -.implemented by.-> SPR
    QSV -.implemented by.-> QWS

    classDef kernel fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef plugin fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef framework fill:#f8f9fa,stroke:#6c757d,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class QPR,QSV,HST kernel;
    class QWS,BPR,SPR plugin;
    class CLI framework;
    class QPR,QSV contract;
```

Legend:

- blue: kernel-owned type or contract
- purple: plugin-owned service or registration type
- gray: framework edge
- dashed border: contract
- dashed arrow: structural relationship such as `used by` or `implemented by`

## Implementation Focus

- add a kernel `QuotePricer` capability
- add a base pricing plugin
- add one sample decorator plugin: `seasonalpricing`

The code should show:

- the quote use case stays structurally stable
- pricing changes only because a pricing plugin is registered
- the extension seam lives at the kernel boundary, not inside the quote entity

## What To Verify

- `go test ./...` passes
- quote pricing works with the base pricer
- registering `seasonalpricing` changes the quote line unit price
- the demo can show the pricing impact
