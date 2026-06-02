# Lesson 016: Real Return Window Policy

## Objective

Replace the placeholder return-eligibility rule with a real time-based return-window policy that uses shipped timestamps and per-product return-window snapshots.

## Theory

Lesson `015` separated policy from the return-review workflow.

But the policy itself was still only a placeholder:

- reject if reason is `outside return window`

That is not a real business rule.

This lesson turns it into one by threading the necessary data through the workflow:

- products define `ReturnWindowDays`
- quotes snapshot that value on each line
- orders carry the snapshot forward
- shipments stamp `ShippedAt`
- return requests stamp `RequestedAt`

That lets the `returneligibility` module evaluate an actual business window instead of reading a magic string.

## Why This Matters Here

This lesson is important because it shows a common modular-monolith pressure: a policy module often forces upstream modules to capture more truthful business data.

The policy did not become more realistic by changing one function alone. It required:

- richer product data
- richer order data
- a time source boundary
- a better return-request record

That is the kind of cross-module refinement that makes architectural tradeoffs visible.

## Diagram

```mermaid
flowchart LR
    subgraph PRM["Products Module"]
        direction TB
        PCT["products.Catalog"]
        PDS["products.Service"]
        PDT["Product<br/>ReturnWindowDays"]
    end

    subgraph ORM["Orders Module"]
        direction TB
        ORO["orders.ReturnableOrderSource"]
        OCL["orders.Clock"]
        OMS["orders.Service"]
        ODT["Order<br/>ShippedAt"]
    end

    subgraph RTM["Returns Module"]
        direction TB
        RCL["returns.Clock"]
        RTS["returns.Service"]
        RRT["ReturnRequest<br/>RequestedAt"]
    end

    subgraph RLM["Return Eligibility Module"]
        direction TB
        RLE["returneligibility.Evaluator"]
        RLS["returneligibility.Service"]
    end

    subgraph INF["Infrastructure"]
        direction TB
        CLI["CLI"]
        CLK["System Clock"]
    end

    CLI --> OMS
    CLI --> RTS
    OMS --> ODT
    RTS --> RRT

    PCT -.used by.-> OMS
    OCL -.used by.-> OMS
    ORO -.used by.-> RTS
    RCL -.used by.-> RTS
    RLE -.used by.-> RTS

    PCT -.implemented by.-> PDS
    OCL -.implemented by.-> CLK
    RCL -.implemented by.-> CLK
    RLE -.implemented by.-> RLS

    classDef module fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef entity fill:#fff3bf,stroke:#b08900,color:#111;
    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef framework fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class PCT,PDS,ORO,OCL,OMS,RCL,RTS,RLE,RLS module;
    class PDT,ODT,RRT entity;
    class CLK dataadapter;
    class CLI framework;
    class PCT,ORO,OCL,RCL,RLE contract;
```

Legend:

- yellow: domain type or business snapshot
- purple: module-owned service or contract
- green: adapter or technical implementation
- blue: framework edge
- dashed border: contract
- dashed arrow: structural relationship such as `used by` or `implemented by`

## Implementation Focus

Implement one real policy upgrade:

- return eligibility should depend on actual shipment and request timing

The code should show:

- `ReturnWindowDays` on products
- return-window snapshots carried from quote to order
- `ShippedAt` recorded when shipment is created
- `RequestedAt` recorded when the return is requested
- the `returneligibility` module evaluating those values

## What To Verify

- `go test ./...` passes
- in-window returns are allowed
- out-of-window returns are rejected
- the time source is abstracted behind a clock boundary
