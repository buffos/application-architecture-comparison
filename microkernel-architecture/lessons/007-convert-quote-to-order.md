# Lesson 007: Convert Quote To Order

## Objective

Add the first cross-plugin fulfillment handoff by letting an `orders` plugin consume an approved-quote capability from the `quotes` plugin through the kernel.

## Theory

The previous lessons established the full quote-side workflow:

- draft quote creation
- quote editing
- policy-aware submission
- explicit approval

That makes the `quotes` plugin meaningful on its own.

The next architectural step is a real plugin-to-plugin business handoff:

- one plugin finishes its workflow
- another plugin begins its own workflow from that result

This lesson introduces that next idea:

- the kernel owns a narrow approved-quote capability
- the `quotes` plugin implements that capability
- a new `orders` plugin consumes it to create an order

This matters because a microkernel should not make plugins depend on each other's repositories or internal entities directly.

The handoff should happen through a kernel-owned extension seam.

That solves an important architectural problem:

- order creation depends on quote approval, but it still should not couple directly to quote storage or quote internals

The tradeoff is that the kernel now owns another cross-plugin handoff contract.

That is acceptable only if the contract stays narrow and capability-oriented rather than exposing entire plugin internals.

## Why This Matters Here

For this repository, the next Microkernel lesson should make one thing clear:

- `quotes` owns quote approval and quote conversion readiness
- `orders` owns order creation
- the handoff between them happens through a kernel capability for approved quotes

That makes the first true multi-plugin business workflow visible in the code.

## Diagram

```mermaid
flowchart LR
    subgraph KER["Kernel"]
        direction TB
        PLG["kernel.Plugin"]
        CDA["kernel.CustomerDirectory"]
        PCA["kernel.ProductCatalog"]
        APA["kernel.ApprovalPolicy"]
        QSA["kernel.QuoteService"]
        QRA["kernel.QuoteReader"]
        AQP["kernel.ApprovedQuoteProvider"]
        OSA["kernel.OrderService"]
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

    subgraph APP["Approvals Plugin"]
        direction TB
        APS["approvals.Service"]
        APPP["approvals.Plugin"]
    end

    subgraph QUP["Quotes Plugin"]
        direction TB
        QTE["Quote"]
        QRE["quotes.Repository"]
        QSS["quotes.Service<br/>... / GetApprovedQuoteForOrder"]
        QPP["quotes.Plugin"]
    end

    subgraph ORP["Orders Plugin"]
        direction TB
        ORE["Order"]
        ORR["orders.Repository"]
        OSS["orders.Service<br/>ConvertQuoteToOrder"]
        OPP["orders.Plugin"]
    end

    subgraph INF["Infrastructure"]
        direction TB
        CLI["CLI"]
        MCR["Memory Customer Repository"]
        MPR["Memory Product Repository"]
        MQR["Memory Quote Repository"]
        MOR["Memory Order Repository"]
    end

    CLI --> HST
    HST --> CPP
    HST --> PPP
    HST --> APPP
    HST --> QPP
    HST --> OPP
    QSS --> QTE
    OSS --> ORE
    PSS --> PRD
    CPS --> CUS

    CDA -.used by.-> QSS
    PCA -.used by.-> QSS
    APA -.used by.-> QSS
    QSA -.used by.-> CLI
    QRA -.used by.-> CLI
    AQP -.used by.-> OSS
    OSA -.used by.-> CLI
    CSR -.used by.-> CPS
    PRR -.used by.-> PSS
    QRE -.used by.-> QSS
    ORR -.used by.-> OSS
    PLG -.implemented by.-> CPP
    PLG -.implemented by.-> PPP
    PLG -.implemented by.-> APPP
    PLG -.implemented by.-> QPP
    PLG -.implemented by.-> OPP
    CDA -.implemented by.-> CPS
    PCA -.implemented by.-> PSS
    APA -.implemented by.-> APS
    QSA -.implemented by.-> QSS
    QRA -.implemented by.-> QSS
    AQP -.implemented by.-> QSS
    OSA -.implemented by.-> OSS
    CSR -.implemented by.-> MCR
    PRR -.implemented by.-> MPR
    QRE -.implemented by.-> MQR
    ORR -.implemented by.-> MOR

    classDef kernel fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef plugin fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef entity fill:#fff3bf,stroke:#b08900,color:#111;
    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef framework fill:#f8f9fa,stroke:#6c757d,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class PLG,CDA,PCA,APA,QSA,QRA,AQP,OSA,HST kernel;
    class CPS,CPP,PSS,PPP,APS,APPP,QSS,QPP,OSS,OPP,CSR,PRR,QRE,ORR plugin;
    class CUS,PRD,QTE,ORE entity;
    class MCR,MPR,MQR,MOR dataadapter;
    class CLI framework;
    class PLG,CDA,PCA,APA,QSA,QRA,AQP,OSA,CSR,PRR,QRE,ORR contract;
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

Implement one handoff flow:

- convert an approved quote to an order

The code should show:

- a kernel-owned approved-quote capability
- the `quotes` plugin implementing it
- a new `orders` plugin using that capability to create an order
- the demo creating and then converting an approved quote

Do not add reservation or payment yet.

## What To Verify

- `go test ./...` passes
- the demo can convert an approved quote to an order
- converting a non-approved quote is rejected in tests
- the `orders` plugin does not access quote storage directly
