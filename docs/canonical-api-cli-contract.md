# Canonical API and CLI Contract

## Purpose

This document defines the canonical external contract for the application described in:

- [prd.md](/c:/Users/buffo/Code/architecture/01.application.architectures/docs/prd.md)
- [canonical-domain-model.md](/c:/Users/buffo/Code/architecture/01.application.architectures/docs/canonical-domain-model.md)
- [canonical-use-cases.md](/c:/Users/buffo/Code/architecture/01.application.architectures/docs/canonical-use-cases.md)

Its role is to keep all architecture implementations externally consistent while allowing freedom in internal design.

The canonical contract includes:

- transport-neutral command and query shapes
- standard identifiers and payload conventions
- error model
- idempotency expectations
- canonical HTTP API mapping
- canonical CLI mapping

## Contract Goals

The external contract should:

- expose the same business behaviors in every implementation
- keep business naming stable across HTTP, CLI, and tests
- be simple enough for a first implementation
- be rich enough to cover all required workflows
- avoid binding the system too tightly to one transport style

## Contract Style

The canonical contract follows these rules:

- commands are intent-based
- queries return snapshots or views
- IDs are opaque strings
- statuses are returned as stable enum-like strings
- money uses explicit amount and currency fields
- timestamps use RFC3339 UTC strings
- error responses use stable machine-readable codes

## Global Conventions

### Identifiers

All aggregate identifiers are opaque strings:

- `customerId`
- `productId`
- `quoteId`
- `quoteLineId`
- `orderId`
- `orderLineId`
- `shipmentId`
- `returnRequestId`
- `pluginKey`

`sku` is also a stable externally supplied identifier.

### Money Shape

```json
{
  "amount": "1250.00",
  "currency": "USD"
}
```

Rules:

- amount is serialized as a string to avoid float ambiguity
- a single currency is sufficient for the toy app, but the field remains explicit

### Quantity Shape

Quantities are integers unless a future variant needs decimals. Canonically:

```json
{
  "quantity": 5
}
```

### Timestamp Shape

```json
{
  "createdAt": "2026-05-17T12:30:00Z"
}
```

### Status Values

Status values must align with the canonical domain model:

- `Draft`
- `PendingApproval`
- `Approved`
- `Rejected`
- `Converted`
- `PendingReservation`
- `Backordered`
- `ReadyForPayment`
- `PaymentReview`
- `ReadyForFulfillment`
- `PartiallyShipped`
- `Shipped`
- `Cancelled`
- `Requested`
- `Accepted`
- `RefundPending`
- `Refunded`

### Actor Field

Commands that represent human action should include an actor field such as:

- `requestedBy`
- `submittedBy`
- `reviewedBy`
- `cancelledBy`
- `shippedBy`
- `processedBy`

The canonical contract does not define authentication. It only preserves the business actor identity in the request.

### Idempotency Key

Commands with retry sensitivity may accept:

```json
{
  "idempotencyKey": "unique-client-command-key"
}
```

Canonical use:

- `convert quote to order`
- `capture payment`
- `create shipment`
- `accept return`
- `complete refund`

## Standard Envelope Guidance

The canonical contract does not require strict envelopes for every implementation, but the following structure is recommended for consistency:

### Success Response

```json
{
  "data": {},
  "meta": {
    "requestId": "req-123"
  }
}
```

### Error Response

```json
{
  "error": {
    "code": "InsufficientStock",
    "message": "Requested quantity exceeds available stock",
    "details": {}
  },
  "meta": {
    "requestId": "req-123"
  }
}
```

## Canonical Error Model

### Error Categories

- `ValidationError`
- `NotFound`
- `Conflict`
- `BusinessRuleViolation`
- `ConcurrencyError`
- `InfrastructureError`

### Canonical Business Error Codes

- `CustomerInactive`
- `ProductUnavailable`
- `QuoteEmpty`
- `QuoteNotApprovable`
- `QuoteAlreadyConverted`
- `ApprovalRequired`
- `InsufficientStock`
- `ShipmentNotAllowedUntilPaymentAccepted`
- `ReturnNotEligible`
- `OrderAlreadyCancelled`
- `OrderAlreadyShipped`
- `PluginNotRegistered`
- `PluginAlreadyEnabled`

### HTTP Mapping Guidance

- `400 Bad Request`: malformed input, missing required fields, invalid enums
- `404 Not Found`: missing resource
- `409 Conflict`: state conflict, duplicate conversion, already cancelled, optimistic concurrency conflict
- `422 Unprocessable Entity`: business rule violation
- `500 Internal Server Error`: unexpected infrastructure or system failure

## Concurrency and Versioning

For mutable resources, implementations may expose:

- `version`
- `etag`

Canonical recommendation:

- include a numeric `version` in mutable resource snapshots
- allow conditional update semantics later if needed

This is useful for comparing optimistic concurrency strategies across architectures.

## Core Resource Shapes

### Customer Snapshot

```json
{
  "customerId": "cust-001",
  "name": "Acme Corp",
  "tier": "Preferred",
  "paymentTerms": "Invoice30",
  "status": "Active",
  "version": 1,
  "createdAt": "2026-05-17T12:30:00Z"
}
```

### Product Snapshot

```json
{
  "productId": "prod-001",
  "sku": "CHAIR-001",
  "name": "Office Chair",
  "category": "Standard",
  "basePrice": { "amount": "120.00", "currency": "USD" },
  "availability": "Available",
  "isConfigurable": false,
  "setupFee": { "amount": "0.00", "currency": "USD" },
  "stockShortagePolicy": "RejectOrder",
  "returnPolicy": {
    "returnable": true,
    "returnWindowDays": 30
  },
  "version": 1
}
```

### Stock Snapshot

```json
{
  "sku": "CHAIR-001",
  "onHand": 20,
  "reserved": 4,
  "available": 16,
  "reorderThreshold": 5,
  "version": 3
}
```

### Quote Snapshot

```json
{
  "quoteId": "quote-001",
  "customerId": "cust-001",
  "status": "PendingApproval",
  "approvalStatus": "Pending",
  "lines": [],
  "totals": {
    "subtotal": { "amount": "200.00", "currency": "USD" },
    "discountTotal": { "amount": "20.00", "currency": "USD" },
    "taxTotal": { "amount": "18.00", "currency": "USD" },
    "grandTotal": { "amount": "198.00", "currency": "USD" }
  },
  "policyFindings": [],
  "version": 4,
  "createdAt": "2026-05-17T12:30:00Z",
  "submittedAt": "2026-05-17T12:35:00Z"
}
```

### Quote Line Snapshot

```json
{
  "quoteLineId": "ql-001",
  "sku": "CHAIR-001",
  "productName": "Office Chair",
  "quantity": 2,
  "baseUnitPrice": { "amount": "120.00", "currency": "USD" },
  "adjustedUnitPrice": { "amount": "110.00", "currency": "USD" },
  "setupFee": { "amount": "0.00", "currency": "USD" },
  "discounts": [
    {
      "code": "PREFERRED_TIER",
      "description": "Preferred customer discount",
      "amount": { "amount": "20.00", "currency": "USD" }
    }
  ],
  "lineTotal": { "amount": "220.00", "currency": "USD" },
  "configurationNote": ""
}
```

### Order Snapshot

```json
{
  "orderId": "order-001",
  "sourceQuoteId": "quote-001",
  "customerId": "cust-001",
  "status": "ReadyForFulfillment",
  "paymentStatus": "Accepted",
  "reservationStatus": "Reserved",
  "fulfillmentStatus": "NotShipped",
  "lines": [],
  "totals": {
    "subtotal": { "amount": "200.00", "currency": "USD" },
    "discountTotal": { "amount": "20.00", "currency": "USD" },
    "taxTotal": { "amount": "18.00", "currency": "USD" },
    "grandTotal": { "amount": "198.00", "currency": "USD" }
  },
  "version": 2,
  "createdAt": "2026-05-17T12:40:00Z"
}
```

### Shipment Snapshot

```json
{
  "shipmentId": "ship-001",
  "orderId": "order-001",
  "status": "Shipped",
  "lines": [
    {
      "orderLineId": "ol-001",
      "quantity": 2
    }
  ],
  "shippedAt": "2026-05-17T13:00:00Z"
}
```

### Return Request Snapshot

```json
{
  "returnRequestId": "ret-001",
  "orderId": "order-001",
  "status": "Accepted",
  "refundStatus": "RefundPending",
  "reason": "Damaged",
  "lines": [
    {
      "orderLineId": "ol-001",
      "quantity": 1
    }
  ],
  "policyFindings": [],
  "requestedAt": "2026-05-18T10:00:00Z"
}
```

### Plugin Snapshot

```json
{
  "pluginKey": "seasonal-pricing",
  "pluginType": "Pricing",
  "status": "Enabled",
  "version": "1.0.0",
  "config": {
    "discountPercentage": 5
  }
}
```

## Canonical HTTP API

### API Base

- base path: `/api/v1`
- content type: `application/json`

### Customers

#### `POST /api/v1/customers`

Create customer.

Request:

```json
{
  "name": "Acme Corp",
  "tier": "Preferred",
  "paymentTerms": "Invoice30"
}
```

Response:

- `201 Created` with customer snapshot

#### `GET /api/v1/customers`

List customers.

#### `PATCH /api/v1/customers/{customerId}/commercial-terms`

Update tier or payment terms.

Request:

```json
{
  "tier": "Enterprise",
  "paymentTerms": "Prepaid"
}
```

### Products

#### `POST /api/v1/products`

Create product.

#### `GET /api/v1/products`

List products.

Query params:

- `category`
- `availability`

#### `PATCH /api/v1/products/{productId}`

Update product commercial data.

### Inventory

#### `POST /api/v1/inventory/{sku}/receive`

Receive stock.

Request:

```json
{
  "quantity": 25
}
```

#### `PATCH /api/v1/inventory/{sku}/reorder-threshold`

Request:

```json
{
  "reorderThreshold": 5
}
```

#### `GET /api/v1/inventory/{sku}`

Read stock record.

### Quotes

#### `POST /api/v1/quotes`

Create quote.

Request:

```json
{
  "customerId": "cust-001"
}
```

#### `GET /api/v1/quotes`

List quotes.

Query params:

- `status`
- `customerId`

#### `GET /api/v1/quotes/{quoteId}`

Read quote detail.

#### `POST /api/v1/quotes/{quoteId}/lines`

Add quote line.

Request:

```json
{
  "sku": "CHAIR-001",
  "quantity": 2,
  "configurationNote": "",
  "requestedDiscount": {
    "type": "Percentage",
    "value": 10
  }
}
```

#### `PATCH /api/v1/quotes/{quoteId}/lines/{quoteLineId}`

Update quote line.

#### `DELETE /api/v1/quotes/{quoteId}/lines/{quoteLineId}`

Remove quote line.

#### `POST /api/v1/quotes/{quoteId}/reprice`

Explicitly reprice quote.

#### `POST /api/v1/quotes/{quoteId}/submit`

Submit quote for approval evaluation.

Request:

```json
{
  "submittedBy": "sales-clerk-1"
}
```

Response:

- `200 OK` with quote snapshot and policy findings

#### `POST /api/v1/quotes/{quoteId}/approve`

Approve quote.

Request:

```json
{
  "reviewedBy": "manager-1",
  "decisionComment": "Approved"
}
```

#### `POST /api/v1/quotes/{quoteId}/reject`

Reject quote.

Request:

```json
{
  "reviewedBy": "manager-1",
  "decisionComment": "Discount too high"
}
```

### Orders

#### `POST /api/v1/quotes/{quoteId}/convert-to-order`

Convert quote to order.

Headers:

- `Idempotency-Key: convert-quote-001`

Request:

```json
{
  "requestedBy": "sales-clerk-1"
}
```

Response:

- `201 Created` with order snapshot
- `200 OK` may be used for an idempotent replay returning the same created order

#### `GET /api/v1/orders`

List orders.

Query params:

- `status`
- `customerId`
- `paymentStatus`

#### `GET /api/v1/orders/{orderId}`

Read order detail.

#### `POST /api/v1/orders/{orderId}/cancel`

Cancel order.

Request:

```json
{
  "cancelledBy": "sales-clerk-1",
  "reason": "Customer request"
}
```

### Payment

#### `POST /api/v1/orders/{orderId}/capture-payment`

Capture payment.

Headers:

- `Idempotency-Key: payment-001`

Request:

```json
{
  "scenario": "Accept"
}
```

Allowed canonical scenario values:

- `Accept`
- `Fail`
- `ManualReview`
- `Auto`

#### `POST /api/v1/orders/{orderId}/approve-payment-review`

Request:

```json
{
  "reviewedBy": "manager-1",
  "decision": "Approve"
}
```

#### `GET /api/v1/orders/{orderId}/payment`

Read payment status.

### Fulfillment

#### `POST /api/v1/orders/{orderId}/shipments`

Create shipment.

Headers:

- `Idempotency-Key: shipment-001`

Request:

```json
{
  "shippedBy": "warehouse-clerk-1",
  "lines": [
    {
      "orderLineId": "ol-001",
      "quantity": 1
    }
  ]
}
```

#### `GET /api/v1/shipments`

List shipments.

Query params:

- `orderId`
- `status`

#### `GET /api/v1/shipments/{shipmentId}`

Read shipment detail.

### Returns

#### `POST /api/v1/orders/{orderId}/returns`

Create return request.

Request:

```json
{
  "requestedBy": "warehouse-clerk-1",
  "reason": "Damaged",
  "lines": [
    {
      "orderLineId": "ol-001",
      "quantity": 1
    }
  ]
}
```

#### `GET /api/v1/returns/{returnRequestId}`

Read return detail.

#### `POST /api/v1/returns/{returnRequestId}/accept`

Headers:

- `Idempotency-Key: return-accept-001`

Request:

```json
{
  "reviewedBy": "warehouse-clerk-1"
}
```

#### `POST /api/v1/returns/{returnRequestId}/reject`

Request:

```json
{
  "reviewedBy": "warehouse-clerk-1",
  "reason": "Outside return window"
}
```

#### `POST /api/v1/returns/{returnRequestId}/complete-refund`

Headers:

- `Idempotency-Key: refund-001`

Request:

```json
{
  "processedBy": "manager-1"
}
```

### Plugins

#### `POST /api/v1/plugins`

Register plugin.

#### `GET /api/v1/plugins`

List plugins.

#### `POST /api/v1/plugins/{pluginKey}/enable`

Enable plugin.

#### `POST /api/v1/plugins/{pluginKey}/disable`

Disable plugin.

#### `PATCH /api/v1/plugins/{pluginKey}/configuration`

Update plugin configuration.

### Reporting

#### `GET /api/v1/reports/orders-awaiting-approval`

#### `GET /api/v1/reports/low-stock-items`

#### `GET /api/v1/reports/quote-conversion`

#### `GET /api/v1/reports/top-discounted-products`

#### `GET /api/v1/reports/return-rate-by-category`

## Canonical CLI Contract

### CLI Goals

The CLI should:

- expose the same use cases as the HTTP API
- support scripts and demos easily
- make architecture comparisons easy without requiring a web UI

Canonical executable name examples:

- `policy-order-app`
- `app`

The exact binary name may vary. Command semantics should not.

## CLI Conventions

### Output Modes

Support:

- human-readable text
- JSON output via `--output json`

JSON output should align with the HTTP response `data` payload shape where practical.

### Exit Codes

- `0`: success
- `2`: validation error
- `3`: not found
- `4`: business rule violation or conflict
- `10`: unexpected system failure

### Actor Flag

Commands that represent human decisions should accept actor flags such as:

- `--submitted-by`
- `--requested-by`
- `--reviewed-by`
- `--cancelled-by`
- `--shipped-by`
- `--processed-by`

### Idempotency Flag

Retry-sensitive commands may accept:

- `--idempotency-key`

## Canonical CLI Commands

### Customers

```text
app customers create --name "Acme Corp" --tier Preferred --payment-terms Invoice30
app customers list
app customers update-commercial-terms --customer-id cust-001 --tier Enterprise --payment-terms Prepaid
```

### Products

```text
app products create --sku CHAIR-001 --name "Office Chair" --category Standard --base-price 120.00 --currency USD --availability Available --stock-shortage-policy RejectOrder --return-window-days 30
app products list
app products update --product-id prod-001 --base-price 130.00
```

### Inventory

```text
app inventory receive --sku CHAIR-001 --quantity 25
app inventory set-reorder-threshold --sku CHAIR-001 --reorder-threshold 5
app inventory get --sku CHAIR-001
```

### Quotes

```text
app quotes create --customer-id cust-001
app quotes add-line --quote-id quote-001 --sku CHAIR-001 --quantity 2 --requested-discount-type Percentage --requested-discount-value 10
app quotes update-line --quote-id quote-001 --quote-line-id ql-001 --quantity 3
app quotes remove-line --quote-id quote-001 --quote-line-id ql-001
app quotes reprice --quote-id quote-001
app quotes submit --quote-id quote-001 --submitted-by sales-clerk-1
app quotes approve --quote-id quote-001 --reviewed-by manager-1 --decision-comment "Approved"
app quotes reject --quote-id quote-001 --reviewed-by manager-1 --decision-comment "Too much discount"
app quotes get --quote-id quote-001
app quotes list --status PendingApproval
```

### Orders

```text
app orders convert-from-quote --quote-id quote-001 --requested-by sales-clerk-1 --idempotency-key convert-001
app orders cancel --order-id order-001 --cancelled-by sales-clerk-1 --reason "Customer request"
app orders get --order-id order-001
app orders list --status ReadyForFulfillment
```

### Payment

```text
app payments capture --order-id order-001 --scenario Auto --idempotency-key payment-001
app payments approve-review --order-id order-001 --reviewed-by manager-1 --decision Approve
app payments get --order-id order-001
```

### Fulfillment

```text
app shipments create --order-id order-001 --line ol-001=1 --shipped-by warehouse-clerk-1 --idempotency-key shipment-001
app shipments get --shipment-id ship-001
app shipments list --order-id order-001
```

### Returns

```text
app returns request --order-id order-001 --line ol-001=1 --reason Damaged --requested-by warehouse-clerk-1
app returns accept --return-request-id ret-001 --reviewed-by warehouse-clerk-1 --idempotency-key return-accept-001
app returns reject --return-request-id ret-001 --reviewed-by warehouse-clerk-1 --reason "Outside return window"
app returns complete-refund --return-request-id ret-001 --processed-by manager-1 --idempotency-key refund-001
app returns get --return-request-id ret-001
```

### Plugins

```text
app plugins register --plugin-key seasonal-pricing --plugin-type Pricing --version 1.0.0 --config '{"discountPercentage":5}'
app plugins enable --plugin-key seasonal-pricing
app plugins disable --plugin-key seasonal-pricing
app plugins update-config --plugin-key seasonal-pricing --config '{"discountPercentage":7}'
app plugins list
```

### Reports

```text
app reports orders-awaiting-approval
app reports low-stock-items
app reports quote-conversion
app reports top-discounted-products
app reports return-rate-by-category
```

## Contract Parity Rules

Every architecture implementation should preserve:

- the same business command names or near-equivalent names
- the same business resources and identifiers
- the same status vocabulary
- the same error codes
- the same meaning of idempotent commands
- the same core request and response fields

They may vary in:

- exact envelope structure
- pagination details
- logging format
- internal DTO organization
- whether both HTTP and CLI are implemented initially

## Minimum Required Surface for First Implementation

If the first architecture variant needs a reduced first slice, it should still include:

1. `customers create`
2. `products create`
3. `inventory receive`
4. `quotes create`
5. `quotes add-line`
6. `quotes submit`
7. `quotes approve`
8. `orders convert-from-quote`
9. `payments capture`
10. `shipments create`
11. `returns request`
12. `returns accept`
13. `plugins enable`
14. `reports low-stock-items`

This minimum slice is enough to exercise the main architectural patterns without requiring the full surface immediately.

## Suggested Testing Contract

To compare architectures fairly, the same acceptance tests should be runnable against:

- application services directly
- HTTP endpoints
- CLI commands

Canonical test assets should eventually cover:

- happy path quote-to-order flow
- discount approval flow
- stock shortage flow
- invalid shipment flow
- return acceptance and rejection flow
- plugin pricing variation flow

## Deliverable Role

This document is the canonical external contract reference for future implementations. If two architecture variants expose materially different API or CLI behavior, this document is the baseline used to detect and explain the divergence.
