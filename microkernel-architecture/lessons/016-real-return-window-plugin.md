# Lesson 016: Real Return Window Plugin

## Objective

Replace the placeholder return-eligibility rule with a real time-based return-window policy that uses shipped timestamps and per-line return-window snapshots.

## Theory

Lesson `015` separated return acceptance policy from the review workflow.

That was the right boundary, but the concrete rule was still only a placeholder:

- reject when the reason string says `outside return window`

This lesson turns that seam into a real business rule by carrying the facts the policy actually needs:

- products define `ReturnWindowDays`
- quotes snapshot that value on each line
- orders carry the snapshot forward
- orders record `ShippedAt`
- return requests record `RequestedAt`

The policy plugin can then decide eligibility from business data instead of a magic string.

This also introduces time as an explicit kernel capability, because business logic needs time facts without calling the system clock directly from plugin workflow code.

## Why This Matters Here

This lesson shows an important microkernel pressure:

- a more realistic policy often forces multiple plugins to carry richer business snapshots

The eligibility plugin did not become realistic by changing one function alone. It required:

- richer product data
- richer order data
- a clock capability
- richer return-request data

That is the kind of cross-plugin refinement that makes architecture tradeoffs visible.

## Diagram

```mermaid
flowchart LR
    subgraph KER["Kernel"]
        direction TB
        PLG["kernel.Plugin"]
        PCT["kernel.ProductCatalog"]
        CLK["kernel.Clock"]
        ROP["kernel.ReturnableOrderProvider"]
        REP["kernel.ReturnEligibilityPolicy"]
        HST["kernel.Host"]
    end

    subgraph PRP["Products Plugin"]
        direction TB
        PDT["Product<br/>ReturnWindowDays"]
        PGS["products.Service"]
        PPP["products.Plugin"]
    end

    subgraph ORP["Orders Plugin"]
        direction TB
        ORE["Order<br/>ShippedAt"]
        OSS["orders.Service"]
        OPP["orders.Plugin"]
    end

    subgraph RTP["Returns Plugin"]
        direction TB
        RRE["ReturnRequest<br/>RequestedAt"]
        RSS["returns.Service"]
        RPP["returns.Plugin"]
    end

    subgraph RLP["Return Eligibility Plugin"]
        direction TB
        RLS["returneligibility.Service<br/>Window Policy"]
        RLG["returneligibility.Plugin"]
    end

    subgraph CLP["Clock Plugin"]
        direction TB
        CLS["clock.Service"]
        CLG["clock.Plugin"]
    end

    CLI["CLI"] --> HST
    HST --> PPP
    HST --> CLG
    HST --> OPP
    HST --> RPP
    HST --> RLG
    PGS --> PDT
    OSS --> ORE
    RSS --> RRE

    PCT -.used by.-> OSS
    CLK -.used by.-> OSS
    CLK -.used by.-> RSS
    ROP -.used by.-> RSS
    REP -.used by.-> RSS
    PLG -.implemented by.-> PPP
    PLG -.implemented by.-> CLG
    PLG -.implemented by.-> OPP
    PLG -.implemented by.-> RPP
    PLG -.implemented by.-> RLG
    PCT -.implemented by.-> PGS
    CLK -.implemented by.-> CLS
    ROP -.implemented by.-> OSS
    REP -.implemented by.-> RLS

    classDef kernel fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef plugin fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef entity fill:#fff3bf,stroke:#b08900,color:#111;
    classDef framework fill:#f8f9fa,stroke:#6c757d,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class PLG,PCT,CLK,ROP,REP,HST kernel;
    class PGS,PPP,OSS,OPP,RSS,RPP,RLS,RLG,CLS,CLG plugin;
    class PDT,ORE,RRE entity;
    class CLI framework;
    class PLG,PCT,CLK,ROP,REP contract;
```

Legend:

- blue: kernel-owned type or contract
- purple: plugin-owned service or plugin registration type
- yellow: plugin-owned domain type
- gray: framework edge
- dashed border: contract
- dashed arrow: structural relationship such as `used by` or `implemented by`

## Implementation Focus

- add a kernel clock capability
- record `ReturnWindowDays` on quote, order, and return snapshots
- stamp `ShippedAt` when an order is shipped
- stamp `RequestedAt` when a return is requested
- evaluate the real return window in the eligibility plugin

Do not add reviewer metadata yet.

## What To Verify

- `go test ./...` passes
- returns inside the window can be accepted
- returns outside the window are rejected
- orders and returns get their timestamps from the clock capability
