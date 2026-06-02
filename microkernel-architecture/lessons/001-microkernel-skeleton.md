# Lesson 001: Microkernel Skeleton

## Objective

Build the first runnable slice of the application in Microkernel / Plugin Architecture and make the stable kernel plus plugin boundary visible through a `customers` plugin and a `quotes` plugin.

## Theory

Microkernel Architecture keeps a small stable core and lets features grow around it as plugins.

The key idea is:

- the kernel owns plugin registration and stable extension contracts
- plugins implement business capabilities
- plugins discover other capabilities through the kernel instead of acting like ordinary peer modules

This solves a different problem from the Modular Monolith skeleton.

The Modular Monolith lesson asked:

- how do we make business modules explicit?

This lesson asks:

- what belongs in the stable core?
- what should be a plugin?
- how does one plugin consume another capability without collapsing back into direct ownership?

The tradeoff is that the kernel must stay disciplined.

If every business concept moves into the kernel, the architecture loses its point.

## Why This Matters Here

For this repository, the first Microkernel lesson should make one thing unmistakable:

- the kernel owns registration and capability discovery
- the `customers` plugin provides customer validation capability
- the `quotes` plugin provides draft quote creation capability
- the `quotes` plugin gets customer validation through the kernel, not from customer storage and not from direct module wiring in `main`

That is the first meaningful difference from the Modular Monolith baseline.

## Diagram

```mermaid
flowchart LR
    subgraph KER["Kernel"]
        direction TB
        PLG["kernel.Plugin"]
        CDA["kernel.CustomerDirectory"]
        QSA["kernel.QuoteService"]
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
        QSS["quotes.Service<br/>CreateDraftQuote"]
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
    CSR -.used by.-> CPS
    QRE -.used by.-> QSS
    PLG -.implemented by.-> CPP
    PLG -.implemented by.-> QPP
    CDA -.implemented by.-> CPS
    QSA -.implemented by.-> QSS
    CSR -.implemented by.-> MCR
    QRE -.implemented by.-> MQR

    classDef kernel fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef plugin fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef entity fill:#fff3bf,stroke:#b08900,color:#111;
    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef framework fill:#f8f9fa,stroke:#6c757d,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class PLG,CDA,QSA,HST kernel;
    class CPS,CPP,QSS,QPP,CSR,QRE plugin;
    class CUS,QTE entity;
    class MCR,MQR dataadapter;
    class CLI framework;
    class PLG,CDA,QSA,CSR,QRE contract;
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

Implement one simple flow:

- create a draft quote

The code should show:

- a kernel that owns plugin registration and capability discovery
- a `customers` plugin that exposes customer validation
- a `quotes` plugin that exposes draft quote creation
- in-memory repositories wired from the outside
- one CLI demo that boots the kernel, registers plugins, and exercises the plugin capability

Do not add quote lines, approvals, or reporting yet.

## What To Verify

- the project compiles
- `go test ./...` passes
- the demo can create a draft quote
- the `quotes` plugin gets customer validation capability from the kernel rather than from direct storage access
