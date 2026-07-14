# Lesson 009: Payment Gateway And Order Capture

## Objective

Add the first order-side business integration seam by making the `orders` plugin capture payment through a separate plugin capability before it marks an order as paid.

## Theory

The previous lesson made order creation operational:

- the `orders` plugin consumes an approved quote
- it reserves inventory through another plugin capability
- it then saves the order

That makes order creation realistic, but the order lifecycle is still incomplete.

A real order does not stop at `PendingPayment`.

This lesson introduces the next pressure:

- payment capture should be a separate integration seam

In Microkernel terms, that becomes another kernel-owned capability:

- the kernel owns a payment capture contract
- a `payments` plugin implements it
- the `orders` plugin consumes it and then applies its own order transition

That distinction matters because:

- external payment execution is not the same thing as order lifecycle ownership

The payment plugin decides:

- whether payment can be captured successfully

The `orders` plugin still decides:

- whether an order is payable
- how its own status changes after a successful capture

This solves an important architectural problem:

- integration with an external business service should still pass through a stable kernel seam instead of being embedded directly in the `orders` plugin

The tradeoff is that another workflow step now depends on runtime plugin collaboration, but the boundary remains explicit.

## Why This Matters Here

For this repository, the next Microkernel lesson should make one thing clear:

- `orders` still owns order state
- `payments` owns payment capture integration
- `orders` becomes `Paid` only after the payment capability succeeds

That makes the first order-side integration seam visible in the architecture.

## Diagram

```mermaid
flowchart LR
    subgraph KER["Kernel"]
        direction TB
        PLG["kernel.Plugin"]
        CDA["kernel.CustomerDirectory"]
        PCA["kernel.ProductCatalog"]
        APA["kernel.ApprovalPolicy"]
        AQP["kernel.ApprovedQuoteProvider"]
        IRA["kernel.InventoryReservation"]
        PMA["kernel.PaymentCapture"]
        QSA["kernel.QuoteService"]
        QRA["kernel.QuoteReader"]
        OSA["kernel.OrderService"]
        HST["kernel.Host"]
    end

    subgraph QUP["Quotes Plugin"]
        direction TB
        QTE["Quote"]
        QRE["quotes.Repository"]
        QSS["quotes.Service"]
        QPP["quotes.Plugin"]
    end

    subgraph INP["Inventory Plugin"]
        direction TB
        STR["StockRecord"]
        INR["inventory.Repository"]
        ISS["inventory.Service<br/>Reserve"]
        IPP["inventory.Plugin"]
    end

    subgraph PMP["Payments Plugin"]
        direction TB
        PGS["payments.Service<br/>Capture"]
        PPP["payments.Plugin"]
    end

    subgraph ORP["Orders Plugin"]
        direction TB
        ORE["Order<br/>MarkPaid()"]
        ORR["orders.Repository"]
        OSS["orders.Service<br/>ConvertQuoteToOrder / CapturePayment"]
        OPP["orders.Plugin"]
    end

    subgraph INF["Infrastructure"]
        direction TB
        CLI["CLI"]
        MQR["Memory Quote Repository"]
        MIR["Memory Inventory Repository"]
        MOR["Memory Order Repository"]
    end

    CLI --> HST
    HST --> QPP
    HST --> IPP
    HST --> PPP
    HST --> OPP
    QSS --> QTE
    OSS --> ORE
    ISS --> STR

    AQP -.used by.-> OSS
    IRA -.used by.-> OSS
    PMA -.used by.-> OSS
    OSA -.used by.-> CLI
    QSA -.used by.-> CLI
    QRA -.used by.-> CLI
    QRE -.used by.-> QSS
    INR -.used by.-> ISS
    ORR -.used by.-> OSS
    PLG -.implemented by.-> QPP
    PLG -.implemented by.-> IPP
    PLG -.implemented by.-> PPP
    PLG -.implemented by.-> OPP
    AQP -.implemented by.-> QSS
    IRA -.implemented by.-> ISS
    PMA -.implemented by.-> PGS
    QSA -.implemented by.-> QSS
    QRA -.implemented by.-> QSS
    OSA -.implemented by.-> OSS
    QRE -.implemented by.-> MQR
    INR -.implemented by.-> MIR
    ORR -.implemented by.-> MOR

    classDef kernel fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef plugin fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef entity fill:#fff3bf,stroke:#b08900,color:#111;
    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef framework fill:#f8f9fa,stroke:#6c757d,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class PLG,CDA,PCA,APA,AQP,IRA,PMA,QSA,QRA,OSA,HST kernel;
    class QSS,QPP,ISS,IPP,PGS,PPP,OSS,OPP,QRE,INR,ORR plugin;
    class QTE,STR,ORE entity;
    class MQR,MIR,MOR dataadapter;
    class CLI framework;
    class PLG,AQP,IRA,PMA,QSA,QRA,OSA,QRE,INR,ORR contract;
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

Implement one payment flow:

- capture payment for an order

The code should show:

- a kernel-owned payment capture capability
- a `payments` plugin implementing that capability
- the `orders` plugin consuming it and then marking the order paid
- payment capture rejected when the order is not payable

Do not add shipment yet.

## What To Verify

- `go test ./...` passes
- the demo can capture payment for the converted order
- capturing payment twice is rejected in tests
- the `orders` plugin still does not own payment integration directly
