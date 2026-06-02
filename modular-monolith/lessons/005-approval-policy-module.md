# Lesson 005: Approval Policy Module

## Objective

Introduce the first external business-policy seam in the Modular Monolith track by making quote submission depend on an `approvals` module capability.

## Theory

The previous lesson moved submission into the `Quote` entity.

That solved one problem:

- lifecycle rules no longer lived in the caller

But another problem remains:

- some submission outcomes depend on business policy that may change more often than the core lifecycle rule

In a modular monolith, that is a good place for another business module:

- `quotes` still owns quote lifecycle
- `approvals` owns approval-decision policy
- `quotes` depends on a narrow approval capability instead of hard-coding policy rules

This keeps two concerns separate:

- the entity owns how submission changes state
- the approval module decides which submission path applies

## Why This Matters Here

If the `Quote` entity hard-codes category-specific approval rules immediately, the `quotes` module becomes too coupled to one policy variant.

If the caller hard-codes the whole decision, the `quotes` module becomes too weak again.

The modular-monolith answer is:

- keep the transition in `quotes`
- keep the changing approval rule in `approvals`
- connect them through a narrow module API

## Diagram

```mermaid
flowchart LR
    subgraph APM["Approvals Module"]
        direction TB
        AEV["approvals.Evaluator"]
        AMS["approvals.Service"]
    end

    subgraph QTM["Quotes Module"]
        direction TB
        QTE["Quote"]
        SUB["Quote.Submit(requiresApproval)"]
        QSR["quotes.Repository"]
        QMS["quotes.Service<br/>SubmitQuote"]
    end

    subgraph INF["Infrastructure"]
        direction TB
        CLI["CLI"]
        MQR["Memory Quote Repository"]
    end

    CLI --> QMS
    QMS --> QTE
    QTE --> SUB

    QSR -.used by.-> QMS
    AEV -.used by.-> QMS
    QSR -.implemented by.-> MQR
    AEV -.implemented by.-> AMS

    classDef module fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef entity fill:#fff3bf,stroke:#b08900,color:#111;
    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef framework fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class AEV,AMS,QSR,QMS module;
    class QTE,SUB entity;
    class MQR dataadapter;
    class CLI framework;
    class AEV,QSR contract;
```

Legend:

- yellow: domain type
- purple: module-owned service or contract
- green: data adapter
- blue: framework edge
- dashed border: contract
- dashed arrow: structural relationship such as `used by` or `implemented by`

## Implementation Focus

Implement one policy-aware workflow:

- submit quote with approval decision

The code should show:

- quote lines carrying enough information for approval evaluation
- an `approvals` module with a narrow evaluator API
- `quotes` submission ending in either `Approved` or `PendingApproval`

## What To Verify

- `go test ./...` passes
- standard quotes become `Approved`
- custom-build quotes become `PendingApproval`
- the submission transition still belongs to the `Quote` entity
