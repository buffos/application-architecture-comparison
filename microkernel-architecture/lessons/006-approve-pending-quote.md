# Lesson 006: Approve Pending Quote

## Objective

Turn `PendingApproval` into a real workflow state by adding an explicit approval action inside the `quotes` plugin.

## Theory

The previous lesson introduced an external approval policy.

That created an important new branch:

- some quotes become `Approved`
- some quotes become `PendingApproval`

But that state is still incomplete unless there is also a real action that moves a pending quote forward.

This lesson introduces that next idea:

- the approval rule still comes from the external `approvals` plugin
- but the action that moves a quote from `PendingApproval` to `Approved` belongs to the `quotes` plugin itself

That matters because policy and workflow are not the same thing.

The approval policy decides:

- whether approval is required

The quote workflow still decides:

- whether a pending quote can transition to approved
- what state changes are valid

This solves an important architectural problem:

- external policy should not absorb the quote lifecycle itself

The tradeoff is that the `quotes` plugin now exposes more than one command on the same kernel capability, which makes the capability richer but also more central.

## Why This Matters Here

For this repository, the next Microkernel lesson should make one thing clear:

- `PendingApproval` is not just a passive status label
- there is a specific `ApproveQuote` action
- that action is still quote-owned behavior even though the approval requirement came from another plugin

That makes the workflow explicit and keeps the plugin boundary honest.

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
        QTE["Quote<br/>Submit(requiresApproval)<br/>Approve()"]
        QRE["quotes.Repository"]
        QSS["quotes.Service<br/>CreateDraftQuote / AddQuoteLine / SubmitQuote / ApproveQuote / GetQuote"]
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
    HST --> APPP
    HST --> QPP
    QSS --> QTE
    PSS --> PRD
    CPS --> CUS

    CDA -.used by.-> QSS
    PCA -.used by.-> QSS
    APA -.used by.-> QSS
    QSA -.used by.-> CLI
    QRA -.used by.-> CLI
    CSR -.used by.-> CPS
    PRR -.used by.-> PSS
    QRE -.used by.-> QSS
    PLG -.implemented by.-> CPP
    PLG -.implemented by.-> PPP
    PLG -.implemented by.-> APPP
    PLG -.implemented by.-> QPP
    CDA -.implemented by.-> CPS
    PCA -.implemented by.-> PSS
    APA -.implemented by.-> APS
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

    class PLG,CDA,PCA,APA,QSA,QRA,HST kernel;
    class CPS,CPP,PSS,PPP,APS,APPP,QSS,QPP,CSR,PRR,QRE plugin;
    class CUS,PRD,QTE entity;
    class MCR,MPR,MQR dataadapter;
    class CLI framework;
    class PLG,CDA,PCA,APA,QSA,QRA,CSR,PRR,QRE contract;
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

Implement one approval flow:

- approve a pending quote

The code should show:

- a kernel-level `ApproveQuote` command on `QuoteService`
- a `Quote.Approve()` rule inside the `quotes` plugin
- approval only valid from `PendingApproval`
- the demo moving a custom quote from `PendingApproval` to `Approved`

Do not convert quotes to orders yet.

## What To Verify

- `go test ./...` passes
- the demo can submit a standard quote directly to `Approved`
- the demo can submit a custom quote to `PendingApproval`
- the demo can explicitly approve that pending quote
- approving a non-pending quote is rejected in tests
