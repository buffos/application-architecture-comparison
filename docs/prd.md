# Product Requirements Document

## Title

Architecture Comparison Toy App: Policy-Driven Order Management

## Purpose

Build the same small-but-deep Go application multiple times using different architectural styles so an experienced software engineer can compare:

- what stays the same across architectures
- what changes in code structure and dependency direction
- where each architecture becomes useful
- what accidental complexity each architecture introduces
- which patterns naturally emerge in each style

The product must be simple enough to finish repeatedly, but rich enough to expose real architectural tradeoffs. The application is therefore intentionally a **single business domain with multiple cross-cutting concerns**, not a CRUD demo.

## Product Summary

The application is a **policy-driven order management system** for a small B2B shop that sells configurable office equipment.

Users create quotes and convert them into orders. The system calculates prices, validates stock, applies discount and approval rules, reserves inventory, accepts payment, handles shipment, and supports cancellation and returns. Some business behavior must be configurable through rules and plugins.

This domain is chosen because it naturally supports:

- transaction-heavy flows
- domain rules and invariants
- state transitions
- orchestration across multiple subsystems
- pluggable behaviors
- read/write separation opportunities
- both anemic and rich modeling styles
- both simple and sophisticated implementations

## Product Goals

The application must:

1. Be implementable in all architectures listed in `docs/architectures.md`.
2. Expose enough complexity to justify architectural patterns beyond basic layering.
3. Keep the core business scope stable across all implementations.
4. Allow side-by-side comparison of similarities, differences, strengths, and weaknesses.
5. Make patterns visible in code, not only in documentation.

## Non-Goals

The application does not need:

- real third-party integrations
- production-grade auth
- distributed deployment
- real payment processing
- real shipping carriers
- advanced UI
- multi-tenant support
- full accounting

These concerns can be simulated with in-memory or local adapters unless a given architecture benefits from a different implementation.

## Core Product Idea

The product manages the lifecycle of a customer order from quote to fulfillment to return.

The system should support both:

- straightforward CRUD-like operations
- deeper business operations with rules, policies, workflows, and invariants

That balance is necessary so the same product can demonstrate:

- Transaction Script
- Active Record
- Rich Domain Model
- DDD
- Layered
- Clean / Onion / Hexagonal
- Modular Monolith
- Component-Based
- Microkernel / Plugin
- Rules Engine

## Users

### 1. Sales Clerk

- creates quotes
- adds products and quantities
- requests discounts
- converts approved quotes into orders

### 2. Warehouse Clerk

- sees reservations
- allocates stock
- ships orders
- records returns

### 3. Manager

- approves high discounts
- reviews policy violations
- configures pricing and approval policies

### 4. System Administrator

- enables plugins
- configures rule sets
- views operational summaries

## Domain Scope

### Core Entities

- Customer
- Product
- Inventory Item / Stock Record
- Quote
- Quote Line
- Order
- Order Line
- Payment
- Shipment
- Return
- Discount Policy
- Approval Request
- Policy Evaluation Result
- Plugin Registration

### Optional Value Objects or Equivalent Concepts

- Money
- SKU
- Quantity
- Address
- Order Status
- Return Reason
- Discount Percentage

Not every architecture must model these as first-class value objects, but the requirements should make them useful where appropriate.

## Core Business Capabilities

### 1. Catalog and Pricing

- The system stores products with SKU, name, base price, category, and availability flag.
- Some products are marked as configurable and may carry extra setup fees.
- Prices may be adjusted by customer tier, quantity breaks, and promotional rules.

### 2. Quote Creation

- A sales clerk can create a quote for a customer.
- A quote can contain multiple lines.
- Each line includes product, quantity, optional configuration notes, and calculated unit price.
- Quote totals must include subtotal, discount amount, tax amount, and grand total.

### 3. Discount Handling

- Discounts can be applied manually or automatically.
- Discounts above a threshold require manager approval.
- Certain product categories are not discountable.
- Multiple discount rules may conflict; the system must define deterministic precedence.

### 4. Quote Approval

- A quote can be submitted for approval.
- Approval outcome can be approved, rejected, or needs more info.
- An approved quote can be converted into an order.
- A rejected quote cannot be converted.

### 5. Inventory Reservation

- Converting a quote into an order attempts to reserve stock.
- If sufficient stock is unavailable, the order is marked backordered or rejected depending on policy.
- Reservation must update available and reserved quantities consistently.

### 6. Payment Capture

- Payment is simulated.
- Payment may succeed, fail, or require manual review.
- Shipment cannot happen before payment is accepted, unless customer terms allow invoicing.

### 7. Shipment

- A warehouse clerk can allocate and ship an order.
- Shipping updates order status.
- Partial shipment must be supported.

### 8. Cancellation

- Orders can be cancelled before shipment.
- Cancellation releases reserved inventory.
- Cancellation after shipment is not allowed; a return must be used instead.

### 9. Returns

- Returned items must reference a prior shipped order.
- Return eligibility depends on time window and product category.
- Accepted returns update inventory and refund status.

### 10. Reporting / Read Models

- The system exposes summaries such as the following:
- quote-to-order conversion rate
- orders awaiting approval
- low stock items
- top discounted products
- return rate by category

These reports are intentionally useful for demonstrating projection/read-model patterns.

## Required Business Rules

The application must include at least the following rules:

1. Products in the `CustomBuild` category require approval before order conversion.
2. Discounts above `15%` require manager approval.
3. Discounts above `25%` are always rejected.
4. `Clearance` products cannot be returned.
5. Orders over a configurable monetary threshold require payment review.
6. Customers with `Preferred` tier receive automatic discount eligibility.
7. If available stock is below requested quantity, the system must either backorder or reject based on product policy.
8. Shipment cannot proceed until payment is accepted, unless the customer is on invoice terms.
9. Cancellation is only allowed before any line has shipped.
10. Returned quantity cannot exceed shipped quantity minus previously returned quantity.

These rules should be represented in a way that different architectures can express differently:

- hardcoded service logic
- domain methods
- rule objects
- policies
- decision tables
- plugin-provided behavior

## Required Workflows

The product must support these end-to-end flows.

### Happy Path

1. Create customer and products.
2. Create quote.
3. Add lines.
4. Apply eligible discounts.
5. Submit for approval if needed.
6. Approve quote.
7. Convert quote to order.
8. Reserve inventory.
9. Capture payment.
10. Ship order.

### Discount Approval Flow

1. Sales clerk applies large discount.
2. System flags approval requirement.
3. Manager approves or rejects.
4. Quote status updates accordingly.

### Stock Shortage Flow

1. Quote converts to order.
2. Reservation detects shortage.
3. System applies backorder or rejection policy.
4. Reason is visible to the user.

### Return Flow

1. User requests return for shipped order.
2. System validates eligibility.
3. Accepted return updates inventory and refund state.
4. Rejected return explains policy failure.

### Plugin / Rules Variation Flow

1. Admin enables a pricing or shipping plugin.
2. Plugin contributes behavior without changing core use cases.
3. Result is visible in quote totals or shipment handling.

This flow is essential for Microkernel / Plugin Architecture.

## Architecture Demonstration Requirements

The same product requirements must allow later implementations to demonstrate the following patterns and concerns.

### Structural Patterns

- controllers / handlers
- application services / use cases
- repositories
- domain services
- ports and adapters
- dependency inversion
- modules / bounded contexts
- components with clear contracts
- plugin extension points

### Behavioral Patterns

- validation
- policy evaluation
- state transitions
- orchestration
- domain events or equivalent notification hooks
- transactional consistency
- command-style operations
- projections / reporting views

### Data and Modeling Patterns

- anemic model style
- active record style
- rich aggregate behavior
- value objects
- repository abstraction
- read model vs write model tension
- rule encapsulation

## Explicit Comparison Constraints

To make architecture comparisons fair, all future implementations should preserve:

- the same business language
- the same primary workflows
- the same business rules
- the same externally visible behaviors
- the same seed scenarios and acceptance tests where practical

They may vary in:

- package structure
- dependency direction
- modeling style
- persistence strategy
- plugin/rule wiring
- internal boundaries

## Functional Requirements

### FR-1 Product Management

- Create and list products.
- Update price, category, and availability.
- Maintain per-product stock quantity and replenishment policy.

### FR-2 Customer Management

- Create and list customers.
- Record customer tier and payment terms.

### FR-3 Quote Management

- Create a quote for a customer.
- Add, update, and remove quote lines.
- Recalculate totals whenever quote contents change.
- Submit quote for approval.
- Approve or reject quote.

### FR-4 Order Management

- Convert approved quote to order.
- Track status transitions.
- Prevent invalid transitions.

### FR-5 Inventory Management

- Reserve stock during order creation.
- Release reservations on cancellation.
- Add stock back on accepted returns.

### FR-6 Payment

- Simulate payment authorization/capture outcome.
- Track payment state per order.

### FR-7 Shipping

- Allocate order lines.
- Support full or partial shipment.
- Track shipment records.

### FR-8 Returns

- Create return request.
- Accept or reject return.
- Track refund status.

### FR-9 Rules and Policies

- Evaluate approval, pricing, shipping, inventory, and return rules.
- Allow at least one implementation approach where rules are configurable.

### FR-10 Reporting

- Provide read endpoints or commands for core summaries.

### FR-11 Extensibility

- Support at least one extension point for pluggable business behavior.
- Extension examples may include pricing calculators, shipping strategies, or approval policies.

## Non-Functional Requirements

### NFR-1 Simplicity of Delivery

- A single engineer should be able to implement one architecture variant without excessive infrastructure.

### NFR-2 Repeatability

- The app should be runnable locally with minimal setup.
- Prefer in-memory or lightweight persistence for early variants.

### NFR-3 Testability

- The product must allow unit, integration, and end-to-end scenario tests.
- Business rules should be testable independently where architecture permits.

### NFR-4 Observability

- Important business actions should be logged or traceable.
- Failures should include clear reasons.

### NFR-5 Determinism

- Given the same inputs and configuration, calculations and rule outcomes must be deterministic.

### NFR-6 Evolution

- The design should leave room to add new rules, new product categories, and new plugins without rewriting the core business language.

## Suggested Delivery Shape

To keep the app implementation manageable while still rich enough for comparison, the first implementation should target:

- CLI, HTTP API, or both
- in-memory repositories first
- optional file-based persistence later

A web UI is not required. API or CLI flows are enough to demonstrate architecture.

## Acceptance Scenarios

### Scenario 1: Standard Order

- Given an in-stock standard product
- When a clerk creates and approves a quote and converts it to an order
- Then inventory is reserved, payment is accepted, and shipment succeeds

### Scenario 2: Approval Required

- Given a quote with a discount above approval threshold
- When the quote is submitted
- Then it cannot convert to an order until approved

### Scenario 3: Approval Rejected

- Given a quote with a prohibited discount
- When validation runs
- Then the quote is rejected with a clear reason

### Scenario 4: Stock Shortage

- Given insufficient stock
- When an approved quote converts to an order
- Then the result follows backorder or reject policy deterministically

### Scenario 5: Invalid Shipment

- Given payment is not accepted and the customer is not on invoice terms
- When shipment is attempted
- Then shipment is rejected

### Scenario 6: Valid Return

- Given a shipped returnable item within return window
- When a return is accepted
- Then inventory is adjusted and refund state is updated

### Scenario 7: Invalid Return

- Given a clearance item
- When a return is requested
- Then the request is rejected by policy

### Scenario 8: Plugin Behavior

- Given a plugin that adds a seasonal surcharge or discount
- When quote totals are recalculated
- Then the plugin contribution is reflected without changing the core workflow

## Why This App Fits The Architecture Comparison

This product is intentionally balanced across multiple axes:

- It has CRUD concerns, so simple architectures remain viable.
- It has workflow and invariants, so richer domain approaches have room to matter.
- It has integration boundaries, so ports/adapters and clean boundaries are meaningful.
- It has configurable policies, so rules engines and policy objects are justified.
- It has pluggable behavior, so microkernel architecture is not artificial.
- It has reporting needs, so separate read concerns can emerge naturally.

If the app were simpler, many architectures would collapse into the same shape and the comparison would be superficial. If it were much larger, the exercise would turn into product development instead of architectural study.

## Implementation Guidance For Future Variants

Each architecture implementation should answer these questions explicitly:

1. Where does business logic live?
2. What owns transaction boundaries?
3. How are rules represented?
4. How are dependencies inverted, if at all?
5. What is the role of repositories, records, or entities?
6. How are modules or components defined?
7. How is extensibility introduced?
8. What becomes easier?
9. What becomes more verbose or indirect?

## Deliverable

The deliverable for each future implementation is the same product with the same business behavior, implemented in Go, using a different architectural style from `docs/architectures.md`.
