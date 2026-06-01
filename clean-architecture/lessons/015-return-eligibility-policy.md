# Lesson 015: Return Eligibility Policy

## Objective

Make return acceptance depend on a dedicated eligibility policy contract instead of allowing any requested return to be accepted unconditionally.

## Theory

The previous lesson introduced review as a separate step.

That was useful because it split:

- creating a return request
- accepting a return
- rejecting a return

But the acceptance path was still too permissive.

The interactor would accept any request that was still in the `Requested` state.

In many systems, return review depends on additional policy questions:

- is the reason acceptable?
- is the product category returnable?
- is the request inside the allowed window?

Those rules do not belong inside the review controller, and they do not have to live directly inside the entity either.

They are a good application-layer policy seam.

So the acceptance flow now becomes:

- load the return request
- load the order
- ask the eligibility policy whether this return may be accepted
- if allowed, continue with refund and restock

This is a good Clean Architecture lesson because it shows that a review use case can depend on both:

- entity lifecycle rules
- replaceable business policy boundaries

The tradeoff is another contract and another adapter.

## Why This Matters Here

Without an eligibility policy, the review step is mostly manual ceremony.

With a policy seam, the architecture now shows a more realistic distinction:

- reviewers act inside a policy framework
- the use case enforces that policy through a contract

This also prepares the next lesson naturally if you want a real time-based return window later.

## Diagram

```mermaid
flowchart TD
    subgraph INFRA[Infrastructure / Frameworks]
        CLI[CLI Framework]
        MOG[Memory Order Gateway]
        MRG[Memory Return Request Gateway]
        RFG[Accept-All Refund Gateway]
        MIR[Memory Inventory Reservation]
        REP[Reason Eligibility Policy]
    end

    subgraph IA[Interface Adapters]
        CTRL[AcceptReturn Controller]
        PRES[AcceptReturn Presenter]
    end

    subgraph APP[Application]
        IN[AcceptReturn Input Model]
        INB[AcceptReturn Input Boundary]
        UC[AcceptReturn Interactor]
        OG[Order Reader Gateway]
        RRG[Return Request Gateway]
        REF[Refund Gateway Contract]
        RES[Inventory Restock Contract]
        ELI[Return Eligibility Policy]
        OUT[AcceptReturn Output Model]
        OUTB[AcceptReturn Output Boundary]
    end

    subgraph ENT[Entities]
        ORDER[Order Entity]
        RET[Return Request Entity]
    end

    CLI --> CTRL --> IN --> INB --> UC
    OG --> UC
    RRG --> UC
    REF --> UC
    RES --> UC
    ELI --> UC
    UC --> ORDER
    UC --> RET
    UC --> OUTB --> OUT --> PRES

    INB -.used by.-> CTRL
    INB -.implements.-> UC
    OG -.used by.-> UC
    OG -.implemented by.-> MOG
    RRG -.used by.-> UC
    RRG -.implemented by.-> MRG
    REF -.used by.-> UC
    REF -.implemented by.-> RFG
    RES -.used by.-> UC
    RES -.implemented by.-> MIR
    ELI -.used by.-> UC
    ELI -.implemented by.-> REP
    OUTB -.used by.-> UC
    OUTB -.implemented by.-> PRES

    classDef framework fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef funcadapter fill:#ffe5d9,stroke:#bc6c25,color:#111;
    classDef app fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef entity fill:#fff3bf,stroke:#b08900,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class CLI framework;
    class MOG,MRG,MIR dataadapter;
    class RFG,REP,CTRL,PRES funcadapter;
    class IN,INB,UC,OG,RRG,REF,RES,ELI,OUT,OUTB app;
    class ORDER,RET entity;
    class INB,OG,RRG,REF,RES,ELI,OUTB contract;
```

Legend:

- blue: framework edge
- green: data adapter
- orange: functionality / policy / translation adapter
- purple: application layer
- yellow: entity layer
- dashed border: interface / contract
- dashed arrow: structural relationship

## Implementation Focus

Extend one existing use case:

- accept a return only when the eligibility policy allows it

The code should show:

- a return eligibility policy contract
- a simple concrete policy adapter
- `AcceptReturn` consulting that policy before refund and restock
- tests for both allowed and blocked acceptance

Do not add a real date-based return window yet.

## What To Verify

- the project compiles
- `go test ./...` passes
- eligible returns can still be accepted
- policy-blocked returns stay `Requested`
