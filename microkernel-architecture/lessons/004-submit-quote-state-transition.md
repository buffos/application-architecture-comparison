# Lesson 004: Submit Quote State Transition

## Objective

Make quote submission an explicit lifecycle rule inside the `quotes` plugin instead of treating status as only passive data.

## Theory

The previous lesson proved that the `quotes` plugin can coordinate with another plugin through the kernel to edit a draft quote.

That is useful, but it still leaves an important modeling gap:

- when does a draft stop being editable?
- what makes a quote submittable?
- where should that lifecycle rule live?

This lesson introduces the next architectural idea:

- the `quotes` plugin still exposes submission through the kernel
- but the actual submission rule belongs to the `Quote` entity inside the plugin

That matters because Microkernel Architecture is not only about plugin registration.

It still needs each plugin to protect its own business invariants.

So the kernel should own:

- extension seams
- capability discovery

while the plugin should still own:

- its own lifecycle rules
- its own state transition behavior

This solves an important architectural problem:

- exposing a capability through the kernel should not flatten plugin business behavior into procedural status updates

The tradeoff is that plugin services become coordinators around richer plugin-owned behavior instead of being only thin wrappers over storage.

## Why This Matters Here

For this repository, the next Microkernel lesson should make one thing clear:

- a quote can only be submitted from `Draft`
- a quote with no lines cannot be submitted
- once submitted, it is no longer editable
- the CLI still reaches that behavior through the kernel contract

That keeps the microkernel structure and the business rule both visible at the same time.

## Diagram

```mermaid
flowchart LR
    subgraph KER["Kernel"]
        direction TB
        PLG["kernel.Plugin"]
        CDA["kernel.CustomerDirectory"]
        PCA["kernel.ProductCatalog"]
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

    subgraph PRP["Products Plugin"]
        direction TB
        PRD["Product"]
        PRR["products.Repository"]
        PSS["products.Service"]
        PPP["products.Plugin"]
    end

    subgraph QUP["Quotes Plugin"]
        direction TB
        QTE["Quote<br/>Submit()"]
        QRE["quotes.Repository"]
        QSS["quotes.Service<br/>CreateDraftQuote / AddQuoteLine / SubmitQuote / GetQuote"]
        QPP["quotes.Plugin"]
    end

    subgraph INF["Infrastructure"]
        direction TB
        CLI["CLI"]
        MCR["Memory Customer Repository"]
        MPR["Memory Product Repository"]
        MQR["Memory Quote Repository"]
    end

    CLI --> HST
    HST --> CPP
    HST --> PPP
    HST --> QPP
    QSS --> QTE
    PSS --> PRD
    CPS --> CUS

    CDA -.used by.-> QSS
    PCA -.used by.-> QSS
    QSA -.used by.-> CLI
    QRA -.used by.-> CLI
    CSR -.used by.-> CPS
    PRR -.used by.-> PSS
    QRE -.used by.-> QSS
    PLG -.implemented by.-> CPP
    PLG -.implemented by.-> PPP
    PLG -.implemented by.-> QPP
    CDA -.implemented by.-> CPS
    PCA -.implemented by.-> PSS
    QSA -.implemented by.-> QSS
    QRA -.implemented by.-> QSS
    CSR -.implemented by.-> MCR
    PRR -.implemented by.-> MPR
    QRE -.implemented by.-> MQR

    classDef kernel fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef plugin fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef entity fill:#fff3bf,stroke:#b08900,color:#111;
    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef framework fill:#f8f9fa,stroke:#6c757d,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class PLG,CDA,PCA,QSA,QRA,HST kernel;
    class CPS,CPP,PSS,PPP,QSS,QPP,CSR,PRR,QRE plugin;
    class CUS,PRD,QTE entity;
    class MCR,MPR,MQR dataadapter;
    class CLI framework;
    class PLG,CDA,PCA,QSA,QRA,CSR,PRR,QRE contract;
```

Legend:

- blue: kernel-owned type or contract
- purple: plugin-owned service, repository contract, or plugin registration type
- yellow: plugin-owned domain type
- green: data adapter
- gray: framework edge
- dashed border: contract
- dashed arrow: structural relationship such as `used by` or `implemented by`

## Implementation Focus

Implement one lifecycle flow:

- submit a quote

The code should show:

- a kernel-level submit command on `QuoteService`
- a `Quote.Submit()` rule inside the `quotes` plugin
- submission blocked when the quote has no lines
- adding more lines blocked after submission

Do not add approval policy yet.

## What To Verify

- `go test ./...` passes
- the demo can create a draft quote
- the demo can add a line and submit the quote
- reloading the quote shows the submitted status
- trying to edit a submitted quote is rejected in tests
