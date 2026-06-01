# Lesson 032: Plugin Pricing Extension Point

## Objective

Add a real extension seam so pricing behavior can change by enabling plugins without changing the quote workflow use case itself.

## Theory

Replaceability and extensibility are not the same thing.

Earlier Clean lessons already showed replaceable boundaries:

- approval policy
- payment gateway
- refund gateway

But a plugin extension point is a different architectural idea.

The question is no longer just:

- which one implementation do we inject?

It becomes:

- how does the application decide which optional behaviors are active?
- how can business behavior change without rewriting the core use case?

This lesson keeps the extension point deliberately narrow:

- quote line pricing

The shape is:

- the application owns plugin registration use cases
- infrastructure stores plugin registrations
- a plugin-aware pricing adapter composes enabled pricing plugins
- `AddQuoteLine` depends only on a pricing contract

That keeps the quote workflow stable while pricing behavior becomes extensible.

## Why This Matters Here

This is one of the strongest “why architecture matters” lessons in the Clean track.

Without an extension seam, every new pricing experiment would push more conditional logic into:

- `AddQuoteLine`
- `Quote`
- or some framework-specific service

With the seam, the use case keeps its orchestration role and the extension mechanism stays outside the entity model while still remaining explicit in the application layer.

## Diagram

```mermaid
flowchart LR
    subgraph ENT["Entities"]
        direction TB
        QTE["Quote Entity"]
        PLG["Plugin Registration Entity"]
    end

    subgraph APP["Application"]
        direction TB
        AIN["AddQuoteLine Input Boundary"]
        AOUT["AddQuoteLine Output Boundary"]
        RIN["RegisterPricingPlugin Input Boundary"]
        ROUT["RegisterPricingPlugin Output Boundary"]
        EIN["EnablePlugin Input Boundary"]
        EOUT["EnablePlugin Output Boundary"]
        LIN["ListPlugins Input Boundary"]
        LOUT["ListPlugins Output Boundary"]
        AUC["AddQuoteLine Interactor"]
        RUC["RegisterPricingPlugin Interactor"]
        EUC["EnablePlugin Interactor"]
        LUC["ListPlugins Interactor"]
        QED["Quote Editor"]
        PGR["Product Reader"]
        PRC["Pricing Policy"]
        PIR["Plugin Repository"]
    end

    subgraph IA["Interface Adapters"]
        direction TB
        ACTRL["AddQuoteLine Controller"]
        RCTRL["RegisterPricingPlugin Controller"]
        ECTRL["EnablePlugin Controller"]
        LCTRL["ListPlugins Controller"]
        APRES["AddQuoteLine Presenter"]
        RPRES["RegisterPricingPlugin Presenter"]
        EPRES["EnablePlugin Presenter"]
        LPRES["ListPlugins Presenter"]
    end

    subgraph INFRA["Infrastructure / Frameworks"]
        direction TB
        CLI["CLI / HTTP Framework"]
        MQG["Memory Quote Gateway"]
        MPG["Memory Product Gateway"]
        MPL["Memory Plugin Gateway"]
        FPP["Fixed Pricing Policy"]
        PPP["Plugin-Aware Pricing Policy"]
    end

    CLI --> ACTRL
    CLI --> RCTRL
    CLI --> ECTRL
    CLI --> LCTRL
    ACTRL --> AIN
    RCTRL --> RIN
    ECTRL --> EIN
    LCTRL --> LIN
    AUC --> AOUT
    RUC --> ROUT
    EUC --> EOUT
    LUC --> LOUT
    APRES --> CLI
    RPRES --> CLI
    EPRES --> CLI
    LPRES --> CLI
    AUC --> QTE
    RUC --> PLG
    EUC --> PLG
    LUC --> PLG

    AIN -.used by.-> ACTRL
    AIN -.implemented by.-> AUC
    AOUT -.used by.-> AUC
    AOUT -.implemented by.-> APRES
    RIN -.used by.-> RCTRL
    RIN -.implemented by.-> RUC
    ROUT -.used by.-> RUC
    ROUT -.implemented by.-> RPRES
    EIN -.used by.-> ECTRL
    EIN -.implemented by.-> EUC
    EOUT -.used by.-> EUC
    EOUT -.implemented by.-> EPRES
    LIN -.used by.-> LCTRL
    LIN -.implemented by.-> LUC
    LOUT -.used by.-> LUC
    LOUT -.implemented by.-> LPRES
    QED -.used by.-> AUC
    PGR -.used by.-> AUC
    PRC -.used by.-> AUC
    PIR -.used by.-> RUC
    PIR -.used by.-> EUC
    PIR -.used by.-> LUC
    QED -.implemented by.-> MQG
    PGR -.implemented by.-> MPG
    PIR -.implemented by.-> MPL
    PRC -.implemented by.-> FPP
    PRC -.implemented by.-> PPP

    classDef framework fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef funcadapter fill:#ffe5d9,stroke:#bc6c25,color:#111;
    classDef app fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef entity fill:#fff3bf,stroke:#b08900,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class CLI framework;
    class MQG,MPG,MPL dataadapter;
    class FPP,PPP,ACTRL,RCTRL,ECTRL,LCTRL,APRES,RPRES,EPRES,LPRES funcadapter;
    class AIN,AOUT,RIN,ROUT,EIN,EOUT,LIN,LOUT,AUC,RUC,EUC,LUC,QED,PGR,PRC,PIR app;
    class QTE,PLG entity;
    class AIN,AOUT,RIN,ROUT,EIN,EOUT,LIN,LOUT,QED,PGR,PRC,PIR contract;
```

Legend:

- blue: framework edge
- green: data adapter
- orange: service or translation adapter
- purple: application layer
- yellow: entity layer
- dashed border: interface / contract
- dashed arrow: structural relationship such as `used by` or `implemented by`

## Implementation Focus

Add:

- plugin registration, enable, and list use cases
- a plugin repository
- a pricing contract for `AddQuoteLine`
- a plugin-aware pricing adapter with one sample plugin: `seasonal-pricing`

The code should show:

- the quote use case stays stable
- pricing changes only because the enabled plugin set changes
- plugin registration and activation are application behavior, not framework magic

## What To Verify

- the project compiles
- `go test ./...` passes
- a pricing plugin can be registered and enabled
- enabling `seasonal-pricing` changes the quote line unit price and total
