# Lesson 008: Order Conversion With Reservation

## Objective

Add the first operational side effect to the Microkernel track by making the `orders` plugin reserve inventory through a separate plugin capability before it saves the order.

## Theory

The previous lesson proved the first cross-plugin business handoff:

- the `quotes` plugin provides an approved-quote capability
- the `orders` plugin consumes it to create an order

That is useful, but it still assumes order creation is only a data translation step.

Real workflows usually need operational coordination too.

This lesson introduces that next pressure:

- converting a quote to an order should also reserve stock

In Microkernel terms, that becomes another extension seam:

- the kernel owns an inventory reservation capability
- an `inventory` plugin implements it
- the `orders` plugin consumes it before persisting the order

This solves an important architectural problem:

- operational side effects should still flow through kernel-owned capabilities instead of being embedded as direct repository access inside another plugin

The tradeoff is that order conversion becomes a multi-step orchestration:

- load approved quote
- derive reservation items
- reserve inventory
- save order

That makes the plugin workflow more realistic and more coupled to runtime collaboration, but still through stable seams.

## Why This Matters Here

For this repository, the next Microkernel lesson should make one thing clear:

- `orders` still owns order creation
- `inventory` owns stock reservation
- the `orders` plugin does not reach into inventory storage directly

That makes the architecture show not only feature handoff, but also operational coordination through plugin capabilities.

## Diagram

```mermaid
flowchart LR
    subgraph KER["Kernel"]
        direction TB
        PLG["kernel.Plugin"]
        CDA["kernel.CustomerDirectory"]
        PCA["kernel.ProductCatalog"]
        APA["kernel.ApprovalPolicy"]
        AQP["kernel.ApprovedQuoteProvider"]
        IRA["kernel.InventoryReservation"]
        QSA["kernel.QuoteService"]
        QRA["kernel.QuoteReader"]
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
        QSS["quotes.Service"]
        QPP["quotes.Plugin"]
    end

    subgraph INP["Inventory Plugin"]
        direction TB
        STR["StockRecord"]
        INR["inventory.Repository"]
        ISS["inventory.Service<br/>Reserve"]
        IPP["inventory.Plugin"]
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
        MIR["Memory Inventory Repository"]
        MOR["Memory Order Repository"]
    end

    CLI --> HST
    HST --> CPP
    HST --> PPP
    HST --> APPP
    HST --> QPP
    HST --> IPP
    HST --> OPP
    QSS --> QTE
    OSS --> ORE
    ISS --> STR
    PSS --> PRD
    CPS --> CUS

    CDA -.used by.-> QSS
    PCA -.used by.-> QSS
    APA -.used by.-> QSS
    AQP -.used by.-> OSS
    IRA -.used by.-> OSS
    OSA -.used by.-> CLI
    QSA -.used by.-> CLI
    QRA -.used by.-> CLI
    CSR -.used by.-> CPS
    PRR -.used by.-> PSS
    QRE -.used by.-> QSS
    INR -.used by.-> ISS
    ORR -.used by.-> OSS
    PLG -.implemented by.-> CPP
    PLG -.implemented by.-> PPP
    PLG -.implemented by.-> APPP
    PLG -.implemented by.-> QPP
    PLG -.implemented by.-> IPP
    PLG -.implemented by.-> OPP
    CDA -.implemented by.-> CPS
    PCA -.implemented by.-> PSS
    APA -.implemented by.-> APS
    AQP -.implemented by.-> QSS
    IRA -.implemented by.-> ISS
    QSA -.implemented by.-> QSS
    QRA -.implemented by.-> QSS
    OSA -.implemented by.-> OSS
    CSR -.implemented by.-> MCR
    PRR -.implemented by.-> MPR
    QRE -.implemented by.-> MQR
    INR -.implemented by.-> MIR
    ORR -.implemented by.-> MOR

    classDef kernel fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef plugin fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef entity fill:#fff3bf,stroke:#b08900,color:#111;
    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef framework fill:#f8f9fa,stroke:#6c757d,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class PLG,CDA,PCA,APA,AQP,IRA,QSA,QRA,OSA,HST kernel;
    class CPS,CPP,PSS,PPP,APS,APPP,QSS,QPP,ISS,IPP,OSS,OPP,CSR,PRR,QRE,INR,ORR plugin;
    class CUS,PRD,QTE,STR,ORE entity;
    class MCR,MPR,MQR,MIR,MOR dataadapter;
    class CLI framework;
    class PLG,CDA,PCA,APA,AQP,IRA,QSA,QRA,OSA,CSR,PRR,QRE,INR,ORR contract;
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

Implement one operational flow:

- convert an approved quote to an order and reserve stock

The code should show:

- a kernel-owned inventory reservation capability
- an `inventory` plugin implementing that capability
- the `orders` plugin consuming it before saving the order
- conversion failing when inventory cannot reserve the required quantity

Do not add payment yet.

## What To Verify

- `go test ./...` passes
- the demo can convert an approved quote to an order with reservation
- insufficient stock blocks conversion in tests
- the `orders` plugin still does not access inventory storage directly
