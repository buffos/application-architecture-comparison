# Lesson 025: Quote Conversion Report Plugin

## Objective

Introduce the first projection-style report in the microkernel track and make it explicit that cross-plugin reporting should still depend on published read capabilities, not on storage access.

## Theory

Up to this point, the microkernel read side has focused on:

- get one thing
- list a group of things

Reports are different.

A report often does not belong to one entity or one plugin record.

Instead, it:

- reads from multiple plugin capabilities
- combines or aggregates those results
- produces a report model with its own meaning

This lesson uses a simple quote conversion report:

- total quotes
- approved quotes
- converted quotes
- conversion rate

The important architectural point is that the report still does not read repositories directly. It depends on the existing `QuoteReader` and `OrderReader` kernel capabilities.

## Why This Matters Here

Cross-plugin reporting is one of the easiest places for a microkernel to lose discipline.

Without a clear home, teams often jump straight to:

- direct repository reads
- storage-shaped reporting code
- adapters that bypass plugin boundaries

This lesson keeps the design honest:

- reporting is a plugin of its own
- it depends on published read capabilities from other plugins
- repositories remain internal to their owning plugins

## Diagram

```mermaid
flowchart LR
    subgraph KER["Kernel"]
        direction TB
        RPT["kernel.Reporting"]
        QRR["kernel.QuoteReader"]
        ORR["kernel.OrderReader"]
        HST["kernel.Host"]
    end

    subgraph RPP["Reporting Plugin"]
        direction TB
        RQS["reporting.Service<br/>QuoteConversionReport"]
        RPM["QuoteConversionReport"]
        RPG["reporting.Plugin"]
    end

    subgraph QTP["Quotes Plugin"]
        direction TB
        QQS["quotes.Service<br/>ListQuotes"]
    end

    subgraph ORP["Orders Plugin"]
        direction TB
        OQS["orders.Service<br/>ListOrders"]
    end

    subgraph INF["Infrastructure"]
        direction TB
        CLI["CLI"]
    end

    CLI --> HST
    HST --> RPG
    RQS --> RPM

    QRR -.used by.-> RQS
    ORR -.used by.-> RQS
    RPT -.used by.-> CLI
    QRR -.implemented by.-> QQS
    ORR -.implemented by.-> OQS
    RPT -.implemented by.-> RQS

    classDef kernel fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef plugin fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef entity fill:#fff3bf,stroke:#b08900,color:#111;
    classDef framework fill:#f8f9fa,stroke:#6c757d,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class RPT,QRR,ORR,HST kernel;
    class RQS,RPG,QQS,OQS plugin;
    class RPM entity;
    class CLI framework;
    class RPT,QRR,ORR contract;
```

Legend:

- blue: kernel-owned type or contract
- purple: plugin-owned service or plugin registration type
- yellow: report model
- gray: framework edge
- dashed border: contract
- dashed arrow: structural relationship such as `used by` or `implemented by`

## Implementation Focus

- add a dedicated reporting plugin
- expose `QuoteConversionReport`
- depend on `QuoteReader` and `OrderReader`
- keep repositories out of the reporting plugin

Do not add the other reports yet.

## What To Verify

- `go test ./...` passes
- the report combines quote and order counts correctly
- the demo can render the report output
