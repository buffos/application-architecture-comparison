# Canonical Use Cases and Application Services

## Purpose

This document defines the canonical application layer for the product described in:

- [prd.md](/c:/Users/buffo/Code/architecture/01.application.architectures/docs/prd.md)
- [canonical-domain-model.md](/c:/Users/buffo/Code/architecture/01.application.architectures/docs/canonical-domain-model.md)

Its role is to stabilize:

- the use cases the system exposes
- the command/query surface
- application service responsibilities
- orchestration boundaries
- transaction expectations
- cross-cutting concerns at the application layer

This is the reference for future architecture variants. Each implementation may organize code differently, but the same application behaviors should exist.

## Application Layer Goals

The canonical application layer should:

- expose business-capable operations, not raw persistence actions
- orchestrate domain objects, repositories, policies, and plugins
- keep UI, transport, and storage concerns out of the use-case contract
- provide a stable surface for CLI, HTTP, tests, and future adapters

The application layer is not the domain model itself. It coordinates domain work and commits results.

## Design Principles

### 1. Use Cases Are Intent-Based

Operations are named after business intent:

- `SubmitQuoteForApproval`
- `ConvertQuoteToOrder`
- `AcceptReturn`

not after CRUD mechanics such as `UpdateQuoteStatus`.

### 2. Commands Change State, Queries Read State

- commands perform business actions and may emit domain events
- queries return views or read models and do not mutate state

### 3. Application Services Orchestrate

Application services:

- load aggregates
- call domain behavior and domain services
- persist changes
- publish events or update projections

They should not own the business rules that already belong to the domain model.

### 4. Canonical Surface First, Architectural Shape Later

The same use cases should be implementable as:

- handlers calling services
- use-case interactors
- scripts
- module service methods
- command handlers

## Application Boundaries

The canonical application layer is grouped into these service areas:

- Catalog Application Service
- Customer Application Service
- Quote Application Service
- Approval Application Service
- Order Application Service
- Payment Application Service
- Fulfillment Application Service
- Return Application Service
- Plugin Application Service
- Reporting Query Service

Future implementations may collapse or split these services, but the use cases should remain equivalent.

## Canonical Commands and Queries

### Catalog Use Cases

#### CreateProduct

Intent:

- register a sellable product

Input:

- `sku`
- `name`
- `category`
- `basePrice`
- `availability`
- `isConfigurable`
- `setupFee`
- `stockShortagePolicy`
- `returnPolicy`

Output:

- created product identifier and current product snapshot

Rules:

- base price and setup fee must be non-negative
- category and policies must be valid

#### UpdateProduct

Intent:

- change commercial characteristics of a product

Input:

- `productID` or `sku`
- mutable product fields

Output:

- updated product snapshot

#### ListProducts

Intent:

- retrieve product catalog entries

Output:

- product summaries

### Inventory Use Cases

#### ReceiveStock

Intent:

- add stock into inventory

Input:

- `sku`
- `quantity`

Output:

- updated stock snapshot

Rules:

- quantity must be positive

#### AdjustReorderThreshold

Intent:

- change reporting threshold for low stock

Input:

- `sku`
- `reorderThreshold`

Output:

- updated stock snapshot

#### GetStockRecord

Intent:

- inspect current stock state

Input:

- `sku`

Output:

- stock snapshot including on-hand, reserved, and available quantities

### Customer Use Cases

#### CreateCustomer

Intent:

- register a customer

Input:

- `name`
- `tier`
- `paymentTerms`

Output:

- created customer snapshot

#### UpdateCustomerCommercialTerms

Intent:

- change tier or payment terms

Input:

- `customerID`
- `tier`
- `paymentTerms`

Output:

- updated customer snapshot

#### ListCustomers

Intent:

- retrieve customer summaries

Output:

- customer list

### Quote Use Cases

#### CreateQuote

Intent:

- create a draft quote for a customer

Input:

- `customerID`

Output:

- created quote snapshot

Preconditions:

- customer must exist and be active

#### AddQuoteLine

Intent:

- add a product request to a draft quote

Input:

- `quoteID`
- `sku`
- `quantity`
- `configurationNote`
- optional requested discount information

Output:

- updated quote snapshot with recalculated totals

Preconditions:

- quote must be in `Draft`
- product must exist and be available
- quantity must be positive

Application responsibilities:

- load quote, customer, and product data needed for pricing
- call pricing and discount policies
- recalculate totals

#### UpdateQuoteLine

Intent:

- change quantity, configuration note, or requested discount for a draft quote line

Input:

- `quoteID`
- `quoteLineID`
- mutable fields

Output:

- updated quote snapshot with recalculated totals

#### RemoveQuoteLine

Intent:

- remove a line from a draft quote

Input:

- `quoteID`
- `quoteLineID`

Output:

- updated quote snapshot with recalculated totals

#### RepriceQuote

Intent:

- explicitly recompute commercial terms for a draft quote

Input:

- `quoteID`

Output:

- repriced quote snapshot

Usefulness:

- important for architectures that separate editing from policy evaluation
- useful when plugins or policy configuration change

#### SubmitQuoteForApproval

Intent:

- finalize a draft quote for approval evaluation

Input:

- `quoteID`
- `submittedBy`

Output:

- quote snapshot
- approval summary
- policy findings

Application responsibilities:

- ensure quote has at least one line
- evaluate approval policy
- create or update approval request if needed
- move quote to `Approved` or `PendingApproval` or `Rejected`

#### ApproveQuote

Intent:

- approve a pending quote

Input:

- `quoteID`
- `reviewedBy`
- optional `decisionComment`

Output:

- approved quote snapshot

Preconditions:

- quote must be pending approval

#### RejectQuote

Intent:

- reject a pending quote

Input:

- `quoteID`
- `reviewedBy`
- `decisionComment`

Output:

- rejected quote snapshot

#### GetQuote

Intent:

- read a quote with lines, totals, status, and policy findings

Input:

- `quoteID`

Output:

- quote detail view

#### ListQuotes

Intent:

- list quotes by status or customer

Input:

- optional filters

Output:

- quote summaries

### Order Use Cases

#### ConvertQuoteToOrder

Intent:

- convert an approved quote into an order and attempt reservation

Input:

- `quoteID`
- `requestedBy`

Output:

- created order snapshot
- reservation outcome
- policy findings

Preconditions:

- quote must be approved or otherwise not require approval
- quote must not already be converted

Application responsibilities:

- load quote, customer, products, and stock records
- create order from quote snapshot
- attempt stock reservation
- apply inventory shortage policy
- persist order and inventory changes consistently
- mark quote as converted

Possible outcomes:

- order created and reserved
- order created and backordered
- conversion rejected with business reason

#### CancelOrder

Intent:

- cancel an order before shipment

Input:

- `orderID`
- `cancelledBy`
- `reason`

Output:

- cancelled order snapshot

Application responsibilities:

- verify shipment has not started
- cancel order
- release reservations

#### GetOrder

Intent:

- retrieve order detail, including payment, reservation, and fulfillment state

Input:

- `orderID`

Output:

- order detail view

#### ListOrders

Intent:

- list orders by status, customer, or fulfillment/payment state

Input:

- optional filters

Output:

- order summaries

### Payment Use Cases

#### CapturePayment

Intent:

- simulate payment processing for an order

Input:

- `orderID`
- optional simulated outcome or scenario flag

Output:

- updated order payment state

Application responsibilities:

- load order and customer terms
- apply payment review policy
- decide whether payment is accepted, failed, or requires manual review
- persist payment state

#### ApprovePaymentReview

Intent:

- resolve an order in manual payment review

Input:

- `orderID`
- `reviewedBy`
- decision

Output:

- updated order snapshot

#### GetPaymentStatus

Intent:

- inspect order payment state

Input:

- `orderID`

Output:

- payment status view

### Fulfillment Use Cases

#### CreateShipment

Intent:

- allocate and ship part or all of an order

Input:

- `orderID`
- list of `{orderLineID, quantity}`
- `shippedBy`

Output:

- shipment snapshot
- updated order fulfillment snapshot

Preconditions:

- order must be eligible for shipment
- payment must satisfy payment policy
- quantities must not exceed remaining shippable quantities

Application responsibilities:

- load order and relevant stock/inventory data
- validate shipment policy
- create shipment
- update shipped quantities
- consume reserved stock

#### GetShipment

Intent:

- read shipment details

Input:

- `shipmentID`

Output:

- shipment detail view

#### ListShipments

Intent:

- list shipments by order or status

Input:

- optional filters

Output:

- shipment summaries

### Return Use Cases

#### RequestReturn

Intent:

- create a return request for shipped items

Input:

- `orderID`
- list of `{orderLineID, quantity}`
- `reason`
- `requestedBy`

Output:

- return request snapshot
- eligibility findings

Application responsibilities:

- load order and shipped quantities
- evaluate return eligibility
- accept immediately or mark rejected

#### AcceptReturn

Intent:

- accept a pending or eligible return request

Input:

- `returnRequestID`
- `reviewedBy`

Output:

- accepted return snapshot
- refund status

Application responsibilities:

- confirm eligibility still holds
- update returned quantities
- restock inventory
- create refund workflow state

#### RejectReturn

Intent:

- reject a return request

Input:

- `returnRequestID`
- `reviewedBy`
- `reason`

Output:

- rejected return snapshot

#### CompleteRefund

Intent:

- mark refund as completed for an accepted return

Input:

- `returnRequestID`
- `processedBy`

Output:

- updated return/refund snapshot

#### GetReturnRequest

Intent:

- read return request and refund state

Input:

- `returnRequestID`

Output:

- return detail view

### Plugin and Policy Use Cases

#### RegisterPlugin

Intent:

- register an available plugin capability

Input:

- `pluginKey`
- `pluginType`
- `version`
- `config`

Output:

- plugin registration snapshot

#### EnablePlugin

Intent:

- enable a registered plugin for participation in business flows

Input:

- `pluginKey`

Output:

- enabled plugin snapshot

#### DisablePlugin

Intent:

- disable a plugin

Input:

- `pluginKey`

Output:

- disabled plugin snapshot

#### UpdatePluginConfiguration

Intent:

- change plugin behavior inputs

Input:

- `pluginKey`
- `config`

Output:

- updated plugin snapshot

#### ListPlugins

Intent:

- retrieve plugin registrations and states

Output:

- plugin summaries

### Reporting Queries

#### GetOrdersAwaitingApproval

Output:

- quote or approval work queue view

#### GetLowStockItems

Output:

- low stock report view

#### GetQuoteConversionReport

Output:

- quote-to-order conversion metrics

#### GetTopDiscountedProducts

Output:

- ranked discounted product report

#### GetReturnRateByCategory

Output:

- return metrics grouped by product category

## Canonical Application Services

The following service decomposition is canonical. Implementations may use different names, but the responsibilities should map cleanly.

### CatalogApplicationService

Owns:

- `CreateProduct`
- `UpdateProduct`
- `ListProducts`
- `ReceiveStock`
- `AdjustReorderThreshold`
- `GetStockRecord`

Dependencies:

- product repository
- stock record repository

### CustomerApplicationService

Owns:

- `CreateCustomer`
- `UpdateCustomerCommercialTerms`
- `ListCustomers`

Dependencies:

- customer repository

### QuoteApplicationService

Owns:

- `CreateQuote`
- `AddQuoteLine`
- `UpdateQuoteLine`
- `RemoveQuoteLine`
- `RepriceQuote`
- `SubmitQuoteForApproval`
- `ApproveQuote`
- `RejectQuote`
- `GetQuote`
- `ListQuotes`

Dependencies:

- quote repository
- customer repository
- product repository
- pricing policy/service
- discount policy
- approval policy
- plugin registry/provider

### OrderApplicationService

Owns:

- `ConvertQuoteToOrder`
- `CancelOrder`
- `GetOrder`
- `ListOrders`

Dependencies:

- quote repository
- order repository
- stock record repository
- order creation service
- inventory reservation service
- inventory policy

### PaymentApplicationService

Owns:

- `CapturePayment`
- `ApprovePaymentReview`
- `GetPaymentStatus`

Dependencies:

- order repository
- customer repository
- payment policy

### FulfillmentApplicationService

Owns:

- `CreateShipment`
- `GetShipment`
- `ListShipments`

Dependencies:

- order repository
- shipment repository
- stock record repository
- shipment policy

### ReturnApplicationService

Owns:

- `RequestReturn`
- `AcceptReturn`
- `RejectReturn`
- `CompleteRefund`
- `GetReturnRequest`

Dependencies:

- order repository
- return repository
- stock record repository
- return policy evaluator

### PluginApplicationService

Owns:

- `RegisterPlugin`
- `EnablePlugin`
- `DisablePlugin`
- `UpdatePluginConfiguration`
- `ListPlugins`

Dependencies:

- plugin repository
- plugin registry/provider

### ReportingQueryService

Owns:

- reporting queries only

Dependencies:

- read model store, projections, or repositories optimized for queries

## Command Contract Shape

Each command should conceptually follow this shape:

- command name
- actor or initiator
- target aggregate identity
- required business inputs
- optional idempotency key

Example command shape:

```text
ConvertQuoteToOrderCommand
- QuoteID
- RequestedBy
- IdempotencyKey
```

Application responses should include:

- resulting aggregate ID
- resulting status
- policy findings or business warnings when relevant
- enough snapshot data for adapter layers to render a response

## Query Contract Shape

Queries should conceptually return read models, not mutable domain objects.

Examples:

- `QuoteDetailView`
- `OrderSummaryView`
- `LowStockItemView`
- `ReturnMetricsView`

This matters because some architecture variants will use dedicated projections while others query directly from the write model.

## Transaction Boundaries

Canonical expectations:

### Single-Aggregate Transactions

These should normally complete atomically inside one application command:

- create customer
- create product
- create quote
- edit draft quote
- approve or reject quote
- enable or disable plugin

### Coordinated Multi-Aggregate Transactions

These involve more than one aggregate and need explicit orchestration:

- convert quote to order and reserve stock
- cancel order and release reservations
- create shipment and consume reserved stock
- accept return and restock inventory

Implementations may choose:

- one database transaction
- application-level orchestration with compensating logic
- sequential consistency with explicit failure handling

The business behavior must remain the same even if technical transaction mechanics differ.

## Idempotency Expectations

The following commands should ideally support idempotent behavior because they are sensitive to retries:

- `ConvertQuoteToOrder`
- `CapturePayment`
- `CreateShipment`
- `AcceptReturn`
- `CompleteRefund`

Canonical behavior under retry:

- do not create duplicate orders, shipments, or refunds
- return the already-established business outcome when the same request is retried safely

## Cross-Cutting Application Concerns

Application services are responsible for:

- authorization hook points, if introduced later
- validation of command completeness
- loading required aggregates
- coordinating repositories and domain services
- transaction handling
- idempotency handling
- logging/auditing hook points
- event publication or projection triggering

Application services are not responsible for:

- HTTP request parsing
- CLI flag parsing
- SQL details
- serialization formats

## Failure Model

Canonical failures should be expressed as business-visible outcomes, not only technical exceptions.

Examples:

- `QuoteNotApprovable`
- `QuoteAlreadyConverted`
- `InsufficientStock`
- `ShipmentNotAllowedUntilPaymentAccepted`
- `ReturnNotEligible`
- `OrderAlreadyCancelled`

Each command should distinguish:

- validation failures
- business rule failures
- missing resource failures
- concurrency or consistency failures
- infrastructure failures

This distinction is important for comparing architectures fairly.

## Domain Events at the Application Boundary

Commands may produce domain events that application services publish or hand off for projection work.

Typical examples:

- `SubmitQuoteForApproval` publishes quote submission and approval-required outcomes
- `ConvertQuoteToOrder` publishes order creation and reservation outcome events
- `CreateShipment` publishes shipment and stock-consumption events
- `AcceptReturn` publishes return acceptance and restocking events
- `EnablePlugin` publishes plugin state change events

Simple architectures may process these inline. More advanced architectures may route them through explicit event dispatch.

## Canonical End-to-End Use Case Chains

### Standard Sales Flow

1. `CreateCustomer`
2. `CreateProduct`
3. `ReceiveStock`
4. `CreateQuote`
5. `AddQuoteLine`
6. `SubmitQuoteForApproval`
7. `ConvertQuoteToOrder`
8. `CapturePayment`
9. `CreateShipment`

### Discount Approval Flow

1. `CreateQuote`
2. `AddQuoteLine` with discount
3. `SubmitQuoteForApproval`
4. `ApproveQuote` or `RejectQuote`
5. `ConvertQuoteToOrder`

### Stock Shortage Flow

1. `ConvertQuoteToOrder`
2. inventory policy decides `Backordered` or rejection

### Return Flow

1. `RequestReturn`
2. `AcceptReturn` or `RejectReturn`
3. `CompleteRefund`

### Plugin Variation Flow

1. `RegisterPlugin`
2. `EnablePlugin`
3. `RepriceQuote` or `AddQuoteLine`
4. observe changed pricing or approval outcome

## Mapping Guidance For Architecture Variants

The canonical use cases should map cleanly across architectural styles:

- Layered Architecture: service layer methods with repositories underneath
- Hexagonal / Clean / Onion: input ports/use cases with adapters around them
- Modular Monolith: module services with explicit contracts
- Transaction Script: command procedures centered on each use case
- Active Record: thin application coordination around persistence-aware models
- Rich Domain Model / DDD: application services orchestrating expressive aggregates and domain services
- Rules Engine: application services invoking a rule evaluation boundary
- Microkernel: application services invoking plugin extension points during orchestration

## Deliverable Role

This document is the canonical application/use-case reference for every future implementation of the system. If two architecture variants expose different use-case semantics, this document is the baseline used to detect and explain that divergence.
