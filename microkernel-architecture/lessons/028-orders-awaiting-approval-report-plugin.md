# Lesson 028: Orders Awaiting Approval Report Plugin

## Objective

Add an approval-queue style report that exposes pending approval work as a reporting-plugin projection.

## Theory

The name "orders awaiting approval" is slightly imperfect in this microkernel track.

The current model does not have:

- a separate approval aggregate
- an order that exists before quote approval

What it does have is:

- quotes in `PendingApproval`

So the honest projection is:

- an approval queue over pending-approval quotes

This is still a useful lesson because operational reports do not need to mirror aggregate names mechanically. The reporting plugin can speak in the language of work queues while still being explicit about the underlying model it reads.

## Why This Matters Here

The reporting track already includes:

- conversion metrics
- return analysis
- operational stock visibility

This lesson adds a human workflow queue.

That broadens the reporting story without inventing domain structures the current model does not actually own.

## Diagram

```mermaid
flowchart LR
    subgraph KER["Kernel"]
        direction TB
        RPT["kernel.Reporting"]
        QRR["kernel.QuoteReader"]
        HST["kernel.Host"]
    end

    subgraph RPP["Reporting Plugin"]
        direction TB
        AQS["reporting.Service<br/>OrdersAwaitingApprovalReport"]
        AQP["OrdersAwaitingApprovalReport"]
        AQR["OrdersAwaitingApprovalRow"]
        RPG["reporting.Plugin"]
    end

    subgraph QTP["Quotes Plugin"]
        direction TB
        QQS["quotes.Service<br/>ListQuotes"]
    end

    subgraph INF["Infrastructure"]
        direction TB
        CLI["CLI"]
    end

    CLI --> HST
    HST --> RPG
    AQS --> AQP
    AQP --> AQR

    QRR -.used by.-> AQS
    RPT -.used by.-> CLI
    QRR -.implemented by.-> QQS
    RPT -.implemented by.-> AQS

    classDef kernel fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef plugin fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef entity fill:#fff3bf,stroke:#b08900,color:#111;
    classDef framework fill:#f8f9fa,stroke:#6c757d,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class RPT,QRR,HST kernel;
    class AQS,RPG,QQS plugin;
    class AQP,AQR entity;
    class CLI framework;
    class RPT,QRR contract;
```

Legend:

- blue: kernel-owned type or contract
- purple: plugin-owned service or registration type
- yellow: report model
- gray: framework edge
- dashed border: contract
- dashed arrow: structural relationship such as `used by` or `implemented by`

## Implementation Focus

- add `OrdersAwaitingApprovalReport`
- enrich the quote read model with `TotalAmount`
- build the queue projection from pending-approval quotes

Do not add a separate approval aggregate or workflow yet.

## What To Verify

- `go test ./...` passes
- pending approval quotes appear in the queue
- line counts and total amounts are surfaced correctly
- the demo can render the approval queue output
