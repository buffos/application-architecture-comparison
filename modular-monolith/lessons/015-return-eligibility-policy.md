# Lesson 015: Return Eligibility Policy

## Objective

Make return acceptance policy-aware by moving the acceptance rule into a dedicated `returneligibility` module.

## Theory

Lesson `014` introduced a real review workflow for returns:

- request
- accept
- reject

But acceptance was still unconditional.

This lesson makes one more boundary explicit:

- `returns` owns the review workflow
- `returneligibility` owns the acceptance rule
- `payments` and `inventory` still own side effects after acceptance

That keeps policy separate from workflow orchestration.

## Why This Matters Here

Without a separate policy seam, the `returns` module would accumulate both:

- workflow state management
- business acceptance rules

That is manageable for one rule, but it gets muddy quickly as policy grows. A dedicated module keeps the rule replaceable and easier to reason about.

## Diagram

```mermaid
flowchart LR
    subgraph RTM["Returns Module"]
        direction TB
        RRE["returns.Repository"]
        RAS["returns.Service<br/>AcceptReturn"]
        RTR["ReturnRequest"]
    end

    subgraph RLM["Return Eligibility Module"]
        direction TB
        RLE["returneligibility.Evaluator"]
        RLS["returneligibility.Service"]
    end

    subgraph PAM["Payments Module"]
        direction TB
        PRF["payments.Refunder"]
        PMS["payments.Service"]
    end

    subgraph INM["Inventory Module"]
        direction TB
        IRS["inventory.Restocker"]
        IMS["inventory.Service"]
    end

    subgraph INF["Infrastructure"]
        direction TB
        CLI["CLI"]
        MRR["Memory Return Request Repository"]
    end

    CLI --> RAS
    RAS --> RTR

    RRE -.used by.-> RAS
    RLE -.used by.-> RAS
    PRF -.used by.-> RAS
    IRS -.used by.-> RAS

    RRE -.implemented by.-> MRR
    RLE -.implemented by.-> RLS
    PRF -.implemented by.-> PMS
    IRS -.implemented by.-> IMS

    classDef module fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef entity fill:#fff3bf,stroke:#b08900,color:#111;
    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef framework fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class RRE,RAS,RLE,RLS,PRF,PMS,IRS,IMS module;
    class RTR entity;
    class MRR dataadapter;
    class CLI framework;
    class RRE,RLE,PRF,IRS contract;
```

Legend:

- yellow: domain type
- purple: module-owned service or contract
- green: data adapter
- blue: framework edge
- dashed border: contract
- dashed arrow: structural relationship such as `used by` or `implemented by`

## Implementation Focus

Implement one policy seam:

- accepting a return should consult a separate eligibility capability

The code should show:

- a `returneligibility` module
- `returns` asking that module before refund and restock
- blocked returns becoming `Rejected` without side effects

## What To Verify

- `go test ./...` passes
- eligible returns still refund and restock
- policy-blocked returns are rejected
- blocked returns do not refund or restock
