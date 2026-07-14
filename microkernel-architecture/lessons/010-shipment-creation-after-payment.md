# Lesson 010: Shipment Creation After Payment

## Objective

Add the first shipping seam by making the `orders` plugin create a shipment through a separate plugin capability after payment succeeds.

## Theory

The previous lesson made the order lifecycle advance through payment:

- order conversion created `PendingPayment`
- payment capture moved the order to `Paid`

That is important, but a fulfillment workflow still needs a shipping step.

This lesson introduces that next architectural pressure:

- shipment creation should be a separate plugin capability
- but order shippability should still belong to the `orders` plugin

So this lesson introduces:

- a kernel-owned shipment creation capability
- a `shipments` plugin that implements it
- an order-side `MarkShipped()` transition inside the `orders` plugin

That distinction matters because:

- shipment persistence and creation are not the same thing as order lifecycle ownership

The shipments plugin decides:

- how shipment records are created and stored

The `orders` plugin still decides:

- whether an order is shippable
- when order status becomes `Shipped`

This solves an important architectural problem:

- fulfillment integration should still pass through a kernel seam instead of being collapsed into direct order-side storage

The tradeoff is that the order workflow now coordinates yet another plugin capability, but the boundary remains explicit and capability-oriented.

## Why This Matters Here

For this repository, the next Microkernel lesson should make one thing clear:

- `orders` still owns order state
- `shipments` owns shipment record creation
- an order becomes `Shipped` only after shipment creation succeeds

That completes the first narrow forward fulfillment path in the Microkernel track.

## Diagram

```mermaid
flowchart LR
    subgraph KER["Kernel"]
        direction TB
        PLG["kernel.Plugin"]
        AQP["kernel.ApprovedQuoteProvider"]
        IRA["kernel.InventoryReservation"]
        PMA["kernel.PaymentCapture"]
        SHA["kernel.ShipmentCreation"]
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

    subgraph SHP["Shipments Plugin"]
        direction TB
        SHT["Shipment"]
        SHR["shipments.Repository"]
        SHS["shipments.Service<br/>CreateShipment"]
        SPP["shipments.Plugin"]
    end

    subgraph ORP["Orders Plugin"]
        direction TB
        ORE["Order<br/>MarkPaid()<br/>MarkShipped()"]
        ORR["orders.Repository"]
        OSS["orders.Service<br/>ConvertQuoteToOrder / CapturePayment / CreateShipment"]
        OPP["orders.Plugin"]
    end

    subgraph INF["Infrastructure"]
        direction TB
        CLI["CLI"]
        MQR["Memory Quote Repository"]
        MIR["Memory Inventory Repository"]
        MOR["Memory Order Repository"]
        MSH["Memory Shipment Repository"]
    end

    CLI --> HST
    HST --> QPP
    HST --> IPP
    HST --> PPP
    HST --> SPP
    HST --> OPP
    QSS --> QTE
    OSS --> ORE
    ISS --> STR
    SHS --> SHT

    AQP -.used by.-> OSS
    IRA -.used by.-> OSS
    PMA -.used by.-> OSS
    SHA -.used by.-> OSS
    OSA -.used by.-> CLI
    QRE -.used by.-> QSS
    INR -.used by.-> ISS
    SHR -.used by.-> SHS
    ORR -.used by.-> OSS
    PLG -.implemented by.-> QPP
    PLG -.implemented by.-> IPP
    PLG -.implemented by.-> PPP
    PLG -.implemented by.-> SPP
    PLG -.implemented by.-> OPP
    AQP -.implemented by.-> QSS
    IRA -.implemented by.-> ISS
    PMA -.implemented by.-> PGS
    SHA -.implemented by.-> SHS
    OSA -.implemented by.-> OSS
    QRE -.implemented by.-> MQR
    INR -.implemented by.-> MIR
    SHR -.implemented by.-> MSH
    ORR -.implemented by.-> MOR

    classDef kernel fill:#e8eefc,stroke:#3559b5,color:#111;
    classDef plugin fill:#f3e8ff,stroke:#7b2cbf,color:#111;
    classDef entity fill:#fff3bf,stroke:#b08900,color:#111;
    classDef dataadapter fill:#d8f3dc,stroke:#2d6a4f,color:#111;
    classDef framework fill:#f8f9fa,stroke:#6c757d,color:#111;
    classDef contract stroke-dasharray: 6 4;

    class PLG,AQP,IRA,PMA,SHA,OSA,HST kernel;
    class QSS,QPP,ISS,IPP,PGS,PPP,SHS,SPP,OSS,OPP,QRE,INR,SHR,ORR plugin;
    class QTE,STR,SHT,ORE entity;
    class MQR,MIR,MSH,MOR dataadapter;
    class CLI framework;
    class PLG,AQP,IRA,PMA,SHA,OSA,QRE,INR,SHR,ORR contract;
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

Implement one shipment flow:

- create a shipment for a paid order

The code should show:

- a kernel-owned shipment creation capability
- a `shipments` plugin implementing that capability
- the `orders` plugin consuming it and then marking the order shipped
- shipment creation rejected when the order is not shippable

Do not add cancellation yet.

## What To Verify

- `go test ./...` passes
- the demo can create a shipment after payment
- attempting to ship a non-paid order is rejected in tests
- the `orders` plugin still does not own shipment persistence directly
