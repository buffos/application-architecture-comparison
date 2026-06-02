# Lesson 017: Return Actor Metadata

## Objective

Make the return workflow auditable by recording who requested, reviewed, and processed each return.

## Theory

The return workflow now has:

- request
- policy-aware review
- refund and restock on acceptance
- rejection when the review fails

But without actor metadata, the workflow still lacks operational accountability.

This lesson adds that missing part:

- the requester is recorded when the return is opened
- the reviewer is recorded when the request is accepted or rejected
- the processor is recorded when acceptance triggers refund and restock
- review notes travel with the decision

That keeps the auditing concern inside the `returns` workflow where it belongs, instead of scattering actor fields across payment or inventory modules.

## Why This Matters Here

This is the point where the return module stops being only a state machine and starts to look like a real business record.

Operational workflows often need answers to questions like:

- who asked for the return?
- who approved it?
- who processed the financial and stock reversal?
- what note was attached to the decision?

Those questions are part of the return workflow itself, not side effects of payments or inventory.

## Diagram

```mermaid
flowchart LR
    subgraph RTM["Returns Module"]
        direction TB
        RRE["returns.Repository"]
        RQS["returns.Service<br/>RequestReturn"]
        RAS["returns.Service<br/>AcceptReturn / RejectReturn"]
        RTR["ReturnRequest<br/>RequestedBy / ReviewedBy / ProcessedBy"]
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

    CLI --> RQS
    CLI --> RAS
    RQS --> RTR
    RAS --> RTR

    RRE -.used by.-> RQS
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

    class RRE,RQS,RAS,RLE,RLS,PRF,PMS,IRS,IMS module;
    class RTR entity;
    class MRR dataadapter;
    class CLI framework;
    class RRE,RLE,PRF,IRS contract;
```

Legend:

- yellow: domain type or business record
- purple: module-owned service or contract
- green: adapter or technical implementation
- blue: framework edge
- dashed border: contract
- dashed arrow: structural relationship such as `used by` or `implemented by`

## Implementation Focus

Implement one auditability upgrade:

- record actors and notes on return request and review

The code should show:

- `RequestedBy` stored when the return is created
- `ReviewedBy` stored for both acceptance and rejection
- `ProcessedBy` stored when acceptance triggers refund and restock
- missing actors being rejected as invalid input

## What To Verify

- `go test ./...` passes
- return requests require a requester
- accepting a return records reviewer and processor
- rejecting a return records reviewer without processing side effects
