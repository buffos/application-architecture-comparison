# Canonical Domain Model

## Purpose

This document defines the canonical domain model for the product described in [prd.md](/c:/Users/buffo/Code/architecture/01.application.architectures/docs/prd.md).

Its purpose is to keep all future Go implementations behaviorally aligned while allowing architectural freedom in:

- package structure
- dependency direction
- persistence approach
- application boundary style
- modeling richness

This is the reference source for:

- ubiquitous language
- subdomains
- aggregate boundaries
- entities
- value objects
- policies and rules
- state transitions
- domain events
- invariants

## Modeling Principles

The canonical model is intentionally designed to be:

- rich enough for DDD, Rich Domain Model, Rules Engine, and Plugin Architecture
- still mappable to Layered, Transaction Script, and Active Record implementations
- stable across all architecture variants

This means:

- domain concepts are defined independently of storage or transport
- behavior is described even if some implementations place it in services instead of entities
- aggregate boundaries are explicit so transaction decisions can be compared fairly

## Ubiquitous Language

### Core Terms

- `Customer`: a buyer with tier and payment terms
- `Product`: a sellable catalog item with pricing and return behavior
- `StockRecord`: inventory state for a product
- `Quote`: a pre-order commercial proposal for a customer
- `QuoteLine`: one requested product and quantity inside a quote
- `ApprovalRequest`: a request to review a quote or policy exception
- `Order`: a committed commercial transaction created from an approved quote
- `OrderLine`: one purchased product and quantity inside an order
- `Reservation`: stock held for an order
- `Payment`: simulated payment attempt or approval state for an order
- `Shipment`: allocation and dispatch record for part or all of an order
- `ReturnRequest`: request to return previously shipped goods
- `Refund`: simulated refund status for accepted returns
- `Policy`: business rule used to allow, reject, or require review
- `Plugin`: an extension unit that contributes behavior to pricing, shipping, or approval flows

### Important Distinctions

- A `Quote` is negotiable. An `Order` is committed.
- `Approval` is a business decision, not a technical authorization concern.
- `Reservation` is inventory intent. `Shipment` is physical fulfillment.
- `ReturnRequest` evaluates eligibility. `Refund` reflects financial follow-up.
- `Policy` is core business logic. `Plugin` is an extension mechanism that may provide or modify policy behavior.

## Subdomains

### Core Domain

- quoting
- order management
- pricing and discount decisions
- approval decisions
- inventory reservation
- shipping eligibility
- return eligibility

These are the parts that should carry the most attention in rich architectural variants.

### Supporting Subdomains

- catalog management
- customer management
- payment simulation
- refund simulation
- reporting/read models
- plugin registration and configuration

## Bounded Context Candidates

The application can be implemented as one deployable unit, but the canonical model recognizes these conceptual boundaries:

### Catalog Context

- Product
- Product pricing attributes
- Product returnability attributes

### Customer Context

- Customer
- Customer tier
- Payment terms

### Quoting Context

- Quote
- QuoteLine
- quote pricing
- quote approval requirements

### Ordering Context

- Order
- OrderLine
- Reservation
- cancellation decisions

### Fulfillment Context

- Shipment
- allocation decisions

### Returns Context

- ReturnRequest
- Refund

### Policy Context

- policy definitions
- evaluation outcomes
- plugin-contributed rules

Implementations may collapse or separate these, but the terminology should remain stable.

## Aggregate Design

The following aggregate roots are the canonical default.

### 1. Customer Aggregate

Root:

- `Customer`

Owns:

- customer identity
- tier
- payment terms
- active/inactive flag

Rationale:

- customer commercial traits affect pricing, approval, and shipping eligibility
- changes are infrequent and transactional scope is small

### 2. Product Aggregate

Root:

- `Product`

Owns:

- SKU
- name
- category
- base price
- availability flag
- configurability flag
- setup fee policy
- stock shortage policy
- return policy attributes

Rationale:

- product commercial rules should be kept coherent
- inventory quantity itself is separated to avoid turning product updates into high-contention transactions

### 3. StockRecord Aggregate

Root:

- `StockRecord`

Owns:

- SKU
- on-hand quantity
- reserved quantity
- reorder threshold

Rationale:

- inventory mutations are frequent and should be isolated
- reservation consistency is central to the domain

### 4. Quote Aggregate

Root:

- `Quote`

Owns:

- quote status
- customer snapshot or customer reference
- quote lines
- calculated totals
- requested discounts
- policy findings
- approval state

Rationale:

- quote edits, recalculation, submission, and approval requirements belong together
- quote is the natural unit for pricing and pre-order validation

### 5. Order Aggregate

Root:

- `Order`

Owns:

- order status
- customer snapshot or customer reference
- order lines
- reservation summary
- payment summary
- shipping summary
- cancellation state

Rationale:

- the order is the center of commitment, payment readiness, and fulfillment progression
- line-level and order-level status transitions must remain consistent

### 6. Shipment Aggregate

Root:

- `Shipment`

Owns:

- order reference
- shipped lines and quantities
- shipment status
- shipment timestamps

Rationale:

- partial shipments benefit from explicit lifecycle and auditability
- shipment history should not overload the order aggregate with mutable fulfillment records

### 7. ReturnRequest Aggregate

Root:

- `ReturnRequest`

Owns:

- order reference
- returned lines and quantities
- return reason
- eligibility outcome
- approval decision
- refund status

Rationale:

- returns have their own workflow and constraints
- they must preserve shipped-vs-returned consistency and refund traceability

### 8. PluginRegistration Aggregate

Root:

- `PluginRegistration`

Owns:

- plugin key
- plugin type
- enabled flag
- configuration payload
- version or capability metadata

Rationale:

- enables Microkernel/Plugin implementations without contaminating the core order model

## Entity Definitions

### Customer

Attributes:

- `CustomerID`
- `Name`
- `Tier`
- `PaymentTerms`
- `Status`

Behavior:

- change tier
- change payment terms
- deactivate customer

### Product

Attributes:

- `ProductID`
- `SKU`
- `Name`
- `Category`
- `BasePrice`
- `Availability`
- `IsConfigurable`
- `SetupFee`
- `StockShortagePolicy`
- `ReturnPolicy`

Behavior:

- change price
- change availability
- determine whether discount is allowed
- determine whether return is allowed by product category

### StockRecord

Attributes:

- `SKU`
- `OnHand`
- `Reserved`
- `ReorderThreshold`

Derived values:

- `Available = OnHand - Reserved`

Behavior:

- increase stock
- reserve stock
- release reservation
- consume shipped stock
- restock returned stock

### Quote

Attributes:

- `QuoteID`
- `CustomerID`
- `Status`
- `Lines`
- `Totals`
- `ApprovalState`
- `PolicyFindings`
- `CreatedAt`
- `SubmittedAt`
- `ApprovedAt`

Behavior:

- add line
- update line quantity
- remove line
- recalculate totals
- submit for approval
- mark approved
- mark rejected
- convert eligibility check

### QuoteLine

Attributes:

- `QuoteLineID`
- `SKU`
- `ProductNameSnapshot`
- `Quantity`
- `BaseUnitPrice`
- `AdjustedUnitPrice`
- `SetupFee`
- `Discounts`
- `LineTotal`
- `ConfigurationNote`

Behavior:

- change quantity
- apply discount result
- recalculate line total

### ApprovalRequest

Attributes:

- `ApprovalRequestID`
- `TargetType`
- `TargetID`
- `ReasonCodes`
- `Status`
- `RequestedBy`
- `ReviewedBy`
- `DecisionComment`

Behavior:

- approve
- reject
- request more info

Note:

- some implementations may persist approval inside `Quote`
- others may model it as a distinct aggregate or workflow record

### Order

Attributes:

- `OrderID`
- `SourceQuoteID`
- `CustomerID`
- `Status`
- `Lines`
- `ReservationState`
- `PaymentState`
- `FulfillmentState`
- `Totals`
- `CreatedAt`
- `CancelledAt`

Behavior:

- create from approved quote
- mark backordered
- mark ready for fulfillment
- record payment outcome
- authorize shipment
- cancel
- derive returnable quantities

### OrderLine

Attributes:

- `OrderLineID`
- `SKU`
- `ProductNameSnapshot`
- `OrderedQuantity`
- `ReservedQuantity`
- `ShippedQuantity`
- `ReturnedQuantity`
- `UnitPrice`
- `DiscountAmount`
- `LineTotal`

Behavior:

- reserve quantity
- release reserved quantity
- ship quantity
- register returned quantity
- compute remaining shippable quantity
- compute remaining returnable quantity

### Shipment

Attributes:

- `ShipmentID`
- `OrderID`
- `Status`
- `Lines`
- `ShippedAt`

Behavior:

- allocate lines
- mark shipped
- reject invalid shipment attempt

### ReturnRequest

Attributes:

- `ReturnRequestID`
- `OrderID`
- `Status`
- `Reason`
- `Lines`
- `RequestedAt`
- `ReviewedAt`
- `RefundState`

Behavior:

- validate eligibility
- accept
- reject
- mark refunded

### Refund

Attributes:

- `RefundID`
- `ReturnRequestID`
- `Amount`
- `Status`

Behavior:

- mark pending
- mark completed
- mark failed

### PluginRegistration

Attributes:

- `PluginKey`
- `PluginType`
- `Status`
- `Config`
- `Version`

Behavior:

- enable
- disable
- update configuration

## Value Objects

The following value objects are canonical even if not every implementation models them as explicit types.

### Identity and Classification

- `CustomerID`
- `ProductID`
- `QuoteID`
- `OrderID`
- `ShipmentID`
- `ReturnRequestID`
- `SKU`
- `ProductCategory`
- `CustomerTier`
- `PaymentTerms`

### Quantitative Values

- `Money`
- `Quantity`
- `Percentage`
- `Threshold`

### Workflow Values

- `QuoteStatus`
- `ApprovalStatus`
- `OrderStatus`
- `ShipmentStatus`
- `ReturnStatus`
- `PaymentStatus`
- `RefundStatus`

### Composite Values

- `QuoteTotals`
- `OrderTotals`
- `DiscountBreakdown`
- `PolicyFinding`
- `ReturnPolicy`
- `ReservationState`

## Enumerations

### ProductCategory

- `Standard`
- `CustomBuild`
- `Clearance`

### CustomerTier

- `Standard`
- `Preferred`
- `Enterprise`

### PaymentTerms

- `Prepaid`
- `Invoice30`

### StockShortagePolicy

- `RejectOrder`
- `AllowBackorder`

### QuoteStatus

- `Draft`
- `PendingApproval`
- `Approved`
- `Rejected`
- `Converted`
- `Expired`

### ApprovalStatus

- `NotRequired`
- `Required`
- `Pending`
- `Approved`
- `Rejected`
- `NeedsMoreInfo`

### OrderStatus

- `PendingReservation`
- `Backordered`
- `ReadyForPayment`
- `PaymentReview`
- `ReadyForFulfillment`
- `PartiallyShipped`
- `Shipped`
- `Cancelled`

### PaymentStatus

- `NotRequired`
- `Pending`
- `Accepted`
- `Failed`
- `ManualReview`

### ShipmentStatus

- `Pending`
- `PartiallyShipped`
- `Shipped`

### ReturnStatus

- `Requested`
- `Accepted`
- `Rejected`
- `RefundPending`
- `Refunded`

### RefundStatus

- `NotStarted`
- `Pending`
- `Completed`
- `Failed`

## Policies and Rule Objects

These are canonical business concepts. Some implementations may encode them as domain services, strategy objects, rules, tables, or scripts.

### PricingPolicy

Responsibilities:

- apply customer-tier discounts
- apply quantity-based discounts
- apply promotional adjustments
- apply plugin contributions
- produce deterministic totals and discount breakdowns

### DiscountPolicy

Responsibilities:

- decide whether a line or quote discount is allowed
- detect approval thresholds
- reject prohibited discount levels
- respect category-level discountability constraints

### ApprovalPolicy

Responsibilities:

- determine if quote approval is required
- generate reason codes
- interpret manager decision effects

### InventoryPolicy

Responsibilities:

- decide reservation outcome under stock shortage
- determine whether to reject or backorder

### PaymentPolicy

Responsibilities:

- decide whether manual review is required
- determine whether shipment is blocked by payment state

### ReturnPolicyEvaluator

Responsibilities:

- validate return window
- validate product category returnability
- validate quantity limits against shipped and already returned quantities

### ShipmentPolicy

Responsibilities:

- determine whether order is eligible for shipping
- validate partial shipment rules

### PluginContributionPolicy

Responsibilities:

- select enabled plugins by extension point
- merge plugin output deterministically
- resolve precedence between core rules and plugin-contributed rules

## Domain Services

Canonical domain services where cross-aggregate coordination is meaningful:

### QuotePricingService

- recalculates quote totals from quote lines, customer traits, product rules, and plugins

### QuoteApprovalService

- evaluates approval requirements and produces policy findings or approval requests

### OrderCreationService

- converts approved quote data into an order snapshot

### InventoryReservationService

- coordinates stock reservation for order lines

### PaymentReviewService

- determines payment review requirement and allowed next state

### ShipmentService

- validates shipment eligibility and records shipped quantities

### ReturnService

- validates and processes return requests

### ReportingProjectionService

- updates read models or reporting summaries from domain changes

## Aggregate Invariants

### Customer Invariants

- tier must be one of the supported tiers
- payment terms must be valid
- inactive customers cannot create new quotes

### Product Invariants

- base price cannot be negative
- setup fee cannot be negative
- category must be valid
- a clearance item is not returnable

### StockRecord Invariants

- on-hand quantity cannot be negative
- reserved quantity cannot be negative
- reserved quantity cannot exceed on-hand quantity
- available quantity must always equal on-hand minus reserved

### Quote Invariants

- quote must have at least one line before submission
- quote lines must have positive quantity
- totals must equal the sum of lines plus setup fees, minus discounts, plus taxes
- rejected quotes cannot convert to orders
- unapproved quotes cannot convert when approval is required
- converted quotes cannot be edited

### Order Invariants

- order can only be created from a quote that is approved or otherwise requires no approval
- ordered quantity must be positive
- shipped quantity per line cannot exceed reserved or ordered quantity according to chosen fulfillment rules
- returned quantity cannot exceed shipped quantity minus previously returned quantity
- cancelled orders cannot be shipped
- shipped orders cannot be cancelled

### Shipment Invariants

- shipment must reference an existing order
- shipment lines must map to order lines
- shipped quantity must be positive
- total shipped quantity per line cannot exceed remaining shippable quantity

### ReturnRequest Invariants

- return must reference a shipped order
- return lines must reference shipped order lines
- requested return quantity must be positive
- return quantity cannot exceed remaining returnable quantity

## Lifecycle Rules

### Quote Lifecycle

`Draft -> PendingApproval -> Approved -> Converted`

Alternative paths:

- `Draft -> Approved` when no approval is required
- `PendingApproval -> Rejected`
- `PendingApproval -> NeedsMoreInfo` if modeled explicitly
- `Draft/Approved -> Expired` if expiration is later introduced

Transition rules:

- lines may be edited only in `Draft`
- submission triggers policy evaluation
- approval decision freezes commercial terms for conversion

### Order Lifecycle

`PendingReservation -> Backordered | ReadyForPayment -> PaymentReview | ReadyForFulfillment -> PartiallyShipped -> Shipped`

Alternative path:

- `PendingReservation | ReadyForPayment | PaymentReview | ReadyForFulfillment -> Cancelled`

Transition rules:

- order creation triggers reservation attempt
- payment review is required when threshold or policy says so
- shipment cannot start unless payment policy allows it

### Return Lifecycle

`Requested -> Accepted -> RefundPending -> Refunded`

Alternative path:

- `Requested -> Rejected`

## Domain Events

These are canonical events. Not every architecture must implement them literally, but the business moments should remain visible.

### Quote Events

- `QuoteCreated`
- `QuoteLineAdded`
- `QuoteLineUpdated`
- `QuoteDiscountApplied`
- `QuoteSubmittedForApproval`
- `QuoteApproved`
- `QuoteRejected`
- `QuoteConvertedToOrder`

### Inventory Events

- `StockReserved`
- `StockReservationFailed`
- `StockReleased`
- `StockConsumedForShipment`
- `StockRestockedFromReturn`

### Order Events

- `OrderCreated`
- `OrderBackordered`
- `OrderCancelled`
- `PaymentReviewRequired`
- `PaymentAccepted`
- `PaymentFailed`

### Fulfillment Events

- `ShipmentCreated`
- `OrderPartiallyShipped`
- `OrderShipped`

### Return Events

- `ReturnRequested`
- `ReturnAccepted`
- `ReturnRejected`
- `RefundCompleted`

### Extension Events

- `PluginEnabled`
- `PluginDisabled`
- `PluginContributionApplied`

## Cross-Aggregate References

Use IDs or snapshots rather than direct object graphs across aggregates.

Canonical references:

- `Quote -> CustomerID`
- `QuoteLine -> SKU`
- `Order -> SourceQuoteID`
- `Order -> CustomerID`
- `OrderLine -> SKU`
- `Shipment -> OrderID`
- `ReturnRequest -> OrderID`
- `ReturnRequest Line -> OrderLineID`

Recommended snapshotting:

- product name on quote and order lines
- unit price at time of quote/order
- customer tier or payment terms if needed for auditability

Snapshotting is important because catalog or customer records may change after quote/order creation.

## Consistency Boundaries

Canonical transactional expectations:

- quote edits and recalculation should be atomic within the quote aggregate
- reservation updates must be atomic within a stock record
- order creation and reservation may be a single transaction in simpler implementations or coordinated steps in more decoupled ones
- shipment creation and order shipped-quantity updates must remain consistent
- accepted return and returned-quantity updates must remain consistent

The canonical model does not force distributed transactions. It only defines business consistency expectations so implementations can choose how to satisfy them.

## Read Model Expectations

These read concerns are canonical but do not need to be modeled as aggregates:

- orders awaiting approval
- quotes by status
- low stock products
- quote-to-order conversion rate
- top discounted products
- return rate by category

These can be served by:

- direct queries in simple architectures
- projections in event-aware architectures
- reporting components in component-based or modular monolith variants

## Extension Points

The canonical model defines these plugin extension points:

### Pricing Extension Point

- contribute surcharge or discount adjustments
- example: seasonal discount plugin

### Approval Extension Point

- contribute extra approval reasons
- example: custom-build risk review plugin

### Shipping Extension Point

- contribute shipping eligibility or shipment enrichment
- example: fragile-item handling plugin

Constraints:

- plugin behavior must be deterministic
- core invariants cannot be bypassed by plugins
- precedence between core policies and plugin output must be explicit

## Minimum Canonical Scenarios

Every implementation should be able to express these scenarios using the same domain language:

1. Create a quote with standard products and no approval requirement.
2. Create a quote with a discount that requires manager approval.
3. Reject a quote because discount exceeds hard limit.
4. Convert an approved quote into an order and reserve stock.
5. Backorder or reject an order due to shortage policy.
6. Block shipment because payment is not yet accepted.
7. Ship part of an order and reflect partial fulfillment.
8. Cancel an unshipped order and release reserved stock.
9. Accept a valid return and restock inventory.
10. Reject a return for a clearance item.
11. Enable a pricing plugin and observe changed totals.

## Mapping Guidance For Different Architectures

This canonical model is not prescriptive about code shape.

Examples:

- Transaction Script may treat aggregates as records and keep behavior in procedural services.
- Active Record may attach persistence and simple business methods to model structs.
- Rich Domain Model and DDD may implement aggregates, value objects, and domain services explicitly.
- Hexagonal, Clean, and Onion may preserve the same domain model while changing dependency direction and adapter boundaries.
- Rules Engine variants may externalize parts of `DiscountPolicy`, `ApprovalPolicy`, or `ReturnPolicyEvaluator`.
- Microkernel variants may realize extension points through plugin registries and plugin capability contracts.

## Deliverable Role

This document is the canonical reference model for all future architecture implementations of the application. If later design work conflicts with this document, the conflict should be resolved explicitly rather than drifting silently between implementations.
