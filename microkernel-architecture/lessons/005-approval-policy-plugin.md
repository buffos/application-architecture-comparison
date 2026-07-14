# Lesson 005: Approval Policy Plugin

## Objective

Introduce the first external policy seam in the Microkernel track so quote submission still belongs to the `quotes` plugin, but the decision about whether approval is required comes from a separate plugin capability.

## Theory

The previous lesson made submission a real lifecycle transition inside the `quotes` plugin.

That was important, but it still assumed the submission outcome was fully local:

- every valid draft became `Submitted`

This lesson introduces a more realistic architectural pressure:

- some quotes should go straight through
- some quotes should stop for approval

In Microkernel terms, that is a useful distinction because it separates:

- plugin-owned lifecycle behavior

from:

- plugin-external policy evaluation

So this lesson introduces:

- a kernel-owned `ApprovalPolicy` contract
- an `approvals` plugin that implements it
- a policy-aware `Quote.Submit(requiresApproval bool)` rule inside the `quotes` plugin

This solves an important architectural problem:

- the `quotes` plugin should still own the transition
- but the approval decision itself should be replaceable through a kernel extension seam

The tradeoff is that submission now coordinates one more capability, which makes plugin registration and capability discovery more central to the workflow.

## Why This Matters Here

For this repository, the next Microkernel lesson should make one thing clear:

- submission is still a `quotes` plugin capability
- but whether submission ends in `Approved` or `PendingApproval` is no longer hard-wired inside `quotes`
- the approval rule comes from a separate plugin

That keeps the kernel seam honest and makes extension more than a storage or query story.

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
        QTE["Quote<br/>Submit(requiresApproval)"]
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

Implement one policy-aware submission flow:

- submit a quote with an external approval rule

The code should show:

- a kernel-owned `ApprovalPolicy` contract
- an `approvals` plugin implementing that contract
- the `quotes` plugin consulting that capability through the kernel
- a quote ending in `Approved` or `PendingApproval` based on plugin-provided policy

Do not add explicit approval action yet.

## What To Verify

- `go test ./...` passes
- the demo can submit a standard quote straight to `Approved`
- a custom-build quote becomes `PendingApproval`
- the `quotes` plugin still does not own the approval rule directly
