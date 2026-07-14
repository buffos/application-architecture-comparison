# Lesson 015: Return Eligibility Plugin

## Objective

Make return acceptance policy-aware by moving the acceptance rule into a dedicated eligibility plugin exposed through a kernel capability.

## Theory

Lesson `014` introduced a real review workflow:

- request a return
- accept a return
- reject a return

But acceptance was still unconditional.

This lesson makes one more boundary explicit:

- the `returns` plugin owns the return review workflow
- a separate `returneligibility` plugin owns the acceptance rule
- the `payments` and `inventory` plugins still own the side effects after acceptance

That matters because workflow orchestration and business policy change for different reasons.

In a microkernel design, the kernel provides the stable capability seam and a plugin supplies the current rule implementation.

## Why This Matters Here

Without a separate policy seam, the `returns` plugin would accumulate both:

- return state management
- acceptance policy logic

That is manageable for one `if` statement, but it gets muddy quickly when policy grows. A separate plugin keeps the rule replaceable and easier to evolve.

## Diagram

```mermaid
flowchart LR
    subgraph KER["Kernel"]
        direction TB
        PLG["kernel.Plugin"]
        REP["kernel.ReturnEligibilityPolicy"]
        RFC["kernel.PaymentRefund"]
        IRS["kernel.InventoryRestock"]
        RSA["kernel.ReturnService"]
        HST["kernel.Host"]
    end

    subgraph RTP["Returns Plugin"]
        direction TB
        RRE["ReturnRequest<br/>Request / Accept / Reject"]
        RTR["returns.Repository"]
        RSS["returns.Service<br/>AcceptReturn consults policy"]
        RPP["returns.Plugin"]
    end

    subgraph RLP["Return Eligibility Plugin"]
        direction TB
        RLS["returneligibility.Service"]
        RLPG["returneligibility.Plugin"]
    end

    subgraph PMP["Payments Plugin"]
        direction TB
        PGS["payments.Service<br/>Refund"]
        PPP["payments.Plugin"]
    end

    subgraph INP["Inventory Plugin"]
        direction TB
        ISS["inventory.Service<br/>Restock"]
        IPP["inventory.Plugin"]
    end

    subgraph INF["Infrastructure"]
        direction TB
        CLI["CLI"]
        MRR["Memory Return Repository"]
    end

    CLI --> HST
    HST --> RPP
    HST --> RLPG
    HST --> PPP
    HST --> IPP
    RSS --> RRE

    RTR -.used by.-> RSS
    REP -.used by.-> RSS
    RFC -.used by.-> RSS
    IRS -.used by.-> RSS
    RSA -.used by.-> CLI
    PLG -.implemented by.-> RPP
    PLG -.implemented by.-> RLPG
    PLG -.implemented by.-> PPP
    PLG -.implemented by.-> IPP
    REP -.implemented by.-> RLS
    RFC -.implemented by.-> PGS
    IRS -.implemented by.-> ISS
    RSA -.implemented by.-> RSS
    RTR -.implemented by.-> MRR

    classDef kernel fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef plugin fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef entity fill:#fff3bf,stroke:#b08900,color:#111;
    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef framework fill:#f8f9fa,stroke:#6c757d,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class PLG,REP,RFC,IRS,RSA,HST kernel;
    class RTR,RSS,RPP,RLS,RLPG,PGS,PPP,ISS,IPP plugin;
    class RRE entity;
    class MRR dataadapter;
    class CLI framework;
    class PLG,REP,RFC,IRS,RSA,RTR contract;
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

- add a kernel-owned return eligibility capability
- implement it with a dedicated `returneligibility` plugin
- make `AcceptReturn` ask that policy before refund and restock
- auto-reject blocked returns without side effects

Do not add a real date-based return window yet.

## What To Verify

- `go test ./...` passes
- eligible returns still refund and restock
- policy-blocked returns are rejected
- blocked returns do not refund or restock
