# Lesson 002: Quote Query Through Kernel Capability

## Objective

Add the first read capability to the Microkernel track so quote lookup also flows through a kernel-owned extension contract rather than through direct repository access.

## Theory

The first lesson proved that a plugin can expose a command capability through the kernel.

That is useful, but incomplete.

A real microkernel does not only need:

- plugin registration
- command-style workflow capabilities

It also needs a way to expose read capabilities through stable kernel contracts.

This lesson introduces that next idea:

- the kernel owns a `QuoteReader` contract
- the `quotes` plugin implements it
- callers ask the kernel for quote lookup capability instead of reading plugin storage directly

This solves an important architectural problem:

- once a plugin is registered, callers should depend on stable kernel-facing capabilities, not on internal plugin repositories

The tradeoff is that the kernel contract surface grows.

That is acceptable only if the kernel keeps owning extension seams, not business details.

## Why This Matters Here

For this repository, the next Microkernel lesson should make one thing clear:

- the `quotes` plugin does not only create quotes
- it also exposes quote lookup through the kernel
- the outside world still talks to the kernel capability, not to quote storage

That keeps the plugin boundary real on both the write side and the read side.

## Diagram

```mermaid
flowchart LR
    subgraph KER["Kernel"]
        direction TB
        PLG["kernel.Plugin"]
        CDA["kernel.CustomerDirectory"]
        QSA["kernel.QuoteService"]
        QRA["kernel.QuoteReader"]
        HST["kernel.Host"]
    end

    subgraph CUP["Customers Plugin"]
        direction TB
        CUS["Customer"]
        CSR["customers.Repository"]
        CPS["customers.Service"]
        CPP["customers.Plugin"]
    end

    subgraph QUP["Quotes Plugin"]
        direction TB
        QTE["Quote"]
        QRE["quotes.Repository"]
        QSS["quotes.Service<br/>CreateDraftQuote / GetQuote"]
        QPP["quotes.Plugin"]
    end

    subgraph INF["Infrastructure"]
        direction TB
        CLI["CLI"]
        MCR["Memory Customer Repository"]
        MQR["Memory Quote Repository"]
    end

    CLI --> HST
    HST --> CPP
    HST --> QPP
    QSS --> QTE
    CPS --> CUS

    CDA -.used by.-> QSS
    QSA -.used by.-> CLI
    QRA -.used by.-> CLI
    CSR -.used by.-> CPS
    QRE -.used by.-> QSS
    PLG -.implemented by.-> CPP
    PLG -.implemented by.-> QPP
    CDA -.implemented by.-> CPS
    QSA -.implemented by.-> QSS
    QRA -.implemented by.-> QSS
    CSR -.implemented by.-> MCR
    QRE -.implemented by.-> MQR

    classDef kernel fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef plugin fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef entity fill:#fff3bf,stroke:#b08900,color:#111;
    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef framework fill:#f8f9fa,stroke:#6c757d,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class PLG,CDA,QSA,QRA,HST kernel;
    class CPS,CPP,QSS,QPP,CSR,QRE plugin;
    class CUS,QTE entity;
    class MCR,MQR dataadapter;
    class CLI framework;
    class PLG,CDA,QSA,QRA,CSR,QRE contract;
```

Legend:

- blue: kernel-owned type or contract
- purple: plugin-owned service, repository contract, or plugin registration type
- yellow: domain type
- green: data adapter
- gray: framework edge
- dashed border: contract
- dashed arrow: structural relationship such as `used by` or `implemented by`

## Implementation Focus

Implement one simple query flow:

- load a quote by id through the kernel

The code should show:

- a kernel-owned `QuoteReader` contract
- the `quotes` plugin implementing both write and read capability
- the CLI using the kernel to create and then reload the quote

Do not add quote lines, approvals, or multiple read filters yet.

## What To Verify

- `go test ./...` passes
- the demo can create a draft quote
- the demo can load that quote again through the kernel read capability
- the CLI still does not access the quote repository directly
